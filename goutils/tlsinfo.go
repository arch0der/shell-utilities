// tlsinfo - Show TLS certificate details for a host.
//
// Usage:
//
//	tlsinfo [OPTIONS] HOST[:PORT]
//
// Options:
//
//	-p PORT   Port (default: 443)
//	-t DUR    Timeout (default: 10s)
//	-j        JSON output
//	-k        Skip verification (show cert even if invalid)
//	-c        Check only: exit 0 if valid, 1 if expired/invalid
//	-w N      Warn if cert expires within N days (exit 2)
//	-a        Show all certs in chain
//
// Examples:
//
//	tlsinfo google.com
//	tlsinfo -p 8443 internal.example.com
//	tlsinfo -w 30 api.example.com    # warn if <30 days left
//	tlsinfo -j example.com | jq .expiry
//	tlsinfo -a example.com           # show full chain
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var (
	port    = flag.String("p", "443", "port")
	timeout = flag.Duration("t", 10*time.Second, "timeout")
	asJSON  = flag.Bool("j", false, "JSON output")
	insecure = flag.Bool("k", false, "skip verify")
	check   = flag.Bool("c", false, "check mode")
	warnDays = flag.Int("w", 0, "warn days")
	chain   = flag.Bool("a", false, "show chain")
)

type CertInfo struct {
	Subject    string   `json:"subject"`
	Issuer     string   `json:"issuer"`
	SANs       []string `json:"sans"`
	NotBefore  string   `json:"not_before"`
	NotAfter   string   `json:"not_after"`
	DaysLeft   int      `json:"days_left"`
	IsCA       bool     `json:"is_ca"`
	Serial     string   `json:"serial"`
	Expired    bool     `json:"expired"`
}

func certInfo(cert *x509.Certificate) CertInfo {
	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)
	var sans []string
	sans = append(sans, cert.DNSNames...)
	for _, ip := range cert.IPAddresses {
		sans = append(sans, ip.String())
	}
	return CertInfo{
		Subject:   cert.Subject.CommonName,
		Issuer:    cert.Issuer.CommonName,
		SANs:      sans,
		NotBefore: cert.NotBefore.Format("2006-01-02"),
		NotAfter:  cert.NotAfter.Format("2006-01-02"),
		DaysLeft:  daysLeft,
		IsCA:      cert.IsCA,
		Serial:    cert.SerialNumber.String(),
		Expired:   now.After(cert.NotAfter),
	}
}

func printCert(c CertInfo, indent string) {
	fmt.Printf("%sSubject:    %s\n", indent, c.Subject)
	fmt.Printf("%sIssuer:     %s\n", indent, c.Issuer)
	if len(c.SANs) > 0 {
		fmt.Printf("%sSANs:       %s\n", indent, strings.Join(c.SANs, ", "))
	}
	fmt.Printf("%sNot Before: %s\n", indent, c.NotBefore)
	fmt.Printf("%sNot After:  %s\n", indent, c.NotAfter)
	status := fmt.Sprintf("%d days remaining", c.DaysLeft)
	if c.Expired {
		status = "EXPIRED"
	}
	fmt.Printf("%sExpiry:     %s\n", indent, status)
	fmt.Printf("%sSerial:     %s\n", indent, c.Serial)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: tlsinfo [OPTIONS] HOST[:PORT]")
		os.Exit(1)
	}

	host := args[0]
	p := *port
	if h, po, err := net.SplitHostPort(host); err == nil {
		host, p = h, po
	}
	addr := net.JoinHostPort(host, p)

	dialer := &net.Dialer{Timeout: *timeout}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{
		InsecureSkipVerify: *insecure,
		ServerName:         host,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "tlsinfo: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		fmt.Fprintln(os.Stderr, "tlsinfo: no certificates")
		os.Exit(1)
	}

	target := certs[0]
	info := certInfo(target)

	if *check {
		if info.Expired {
			fmt.Println("EXPIRED")
			os.Exit(1)
		}
		fmt.Println("OK")
		os.Exit(0)
	}

	if *warnDays > 0 && info.DaysLeft <= *warnDays && !info.Expired {
		fmt.Fprintf(os.Stderr, "WARNING: certificate expires in %d days\n", info.DaysLeft)
		defer os.Exit(2)
	}

	if *asJSON {
		var infos []CertInfo
		showCerts := certs
		if !*chain {
			showCerts = certs[:1]
		}
		for _, c := range showCerts {
			infos = append(infos, certInfo(c))
		}
		var out interface{} = infos[0]
		if *chain {
			out = infos
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Printf("Host: %s\n", addr)
	fmt.Println()
	fmt.Println("Leaf Certificate:")
	printCert(info, "  ")

	if *chain && len(certs) > 1 {
		for i, c := range certs[1:] {
			fmt.Printf("\nChain[%d]:\n", i+1)
			printCert(certInfo(c), "  ")
		}
	}
}
