// certinfo - show TLS certificate info for a host or PEM file
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func printCert(cert *x509.Certificate, idx int) {
	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)
	status := "✓ valid"
	if now.After(cert.NotAfter) { status = "✗ EXPIRED" }
	if now.Before(cert.NotBefore) { status = "⚠ not yet valid" }
	expireStr := ""
	if daysLeft < 30 { expireStr = fmt.Sprintf(" ⚠ EXPIRING IN %d DAYS", daysLeft) }

	fmt.Printf("─── Certificate #%d ────────────────────────────\n", idx+1)
	fmt.Printf("Subject      : %s\n", cert.Subject.CommonName)
	fmt.Printf("Issuer       : %s\n", cert.Issuer.CommonName)
	fmt.Printf("Status       : %s%s\n", status, expireStr)
	fmt.Printf("Valid From   : %s\n", cert.NotBefore.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Valid Until  : %s\n", cert.NotAfter.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Days Left    : %d\n", daysLeft)
	fmt.Printf("Serial       : %s\n", cert.SerialNumber)
	fmt.Printf("Sig Algo     : %s\n", cert.SignatureAlgorithm)
	if len(cert.DNSNames) > 0 { fmt.Printf("SANs (DNS)   : %s\n", strings.Join(cert.DNSNames, ", ")) }
	if len(cert.IPAddresses) > 0 {
		ips := make([]string, len(cert.IPAddresses))
		for i, ip := range cert.IPAddresses { ips[i] = ip.String() }
		fmt.Printf("SANs (IP)    : %s\n", strings.Join(ips, ", "))
	}
	fmt.Println()
}

func fromHost(host string) {
	if !strings.Contains(host, ":") { host += ":443" }
	conn, err := tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
	if err != nil { fmt.Fprintln(os.Stderr, "certinfo:", err); os.Exit(1) }
	defer conn.Close()
	for i, cert := range conn.ConnectionState().PeerCertificates { printCert(cert, i) }
}

func fromFile(path string) {
	data, err := os.ReadFile(path); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	var idx int
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil { break }
		if block.Type != "CERTIFICATE" { continue }
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		printCert(cert, idx); idx++
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: certinfo <host[:port]|cert.pem>")
		os.Exit(1)
	}
	target := os.Args[1]
	if strings.HasSuffix(target, ".pem") || strings.HasSuffix(target, ".crt") || strings.HasSuffix(target, ".cer") {
		fromFile(target)
	} else {
		// Check if it's a valid hostname/IP (not a file)
		if _, err := os.Stat(target); err == nil { fromFile(target); return }
		// Validate as host
		host, _, err := net.SplitHostPort(target)
		if err != nil { host = target }
		_ = host
		fromHost(target)
	}
}
