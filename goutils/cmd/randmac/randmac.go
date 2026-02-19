// randmac - Generate random MAC addresses.
//
// Usage:
//
//	randmac [OPTIONS] [N]
//
// Options:
//
//	-n N       Generate N addresses (default: 1)
//	-u         Uppercase
//	-s SEP     Separator (default: :, use - or . or "")
//	-l         Locally administered bit set (default: true for random)
//	-oui OUI   Use specific OUI prefix (e.g. 00:11:22)
//	-v VENDOR  Use OUI for known vendor (apple, cisco, intel, samsung, vmware)
//
// Examples:
//
//	randmac                       # aa:bb:cc:dd:ee:ff
//	randmac -n 5                  # 5 random MACs
//	randmac -s -                  # aa-bb-cc-dd-ee-ff
//	randmac -s ""                 # aabbccddeeff
//	randmac -oui 00:50:56         # VMware-style
//	randmac -v apple              # Apple OUI
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"strings"
)

var knownOUIs = map[string]string{
	"apple":   "00:17:f2",
	"cisco":   "00:17:5a",
	"intel":   "8c:8d:28",
	"samsung": "00:12:fb",
	"vmware":  "00:50:56",
	"qemu":    "52:54:00",
	"dell":    "00:14:22",
	"hp":      "3c:d9:2b",
}

var (
	count  = flag.Int("n", 1, "count")
	upper  = flag.Bool("u", false, "uppercase")
	sep    = flag.String("s", ":", "separator")
	oui    = flag.String("oui", "", "OUI prefix")
	vendor = flag.String("v", "", "vendor name")
)

func genMAC() string {
	b := make([]byte, 6)
	rand.Read(b)

	// Set locally administered and unicast bits (bits 1 and 0 of first byte)
	b[0] = (b[0] | 0x02) & 0xfe

	prefix := *oui
	if *vendor != "" {
		if o, ok := knownOUIs[strings.ToLower(*vendor)]; ok {
			prefix = o
		} else {
			fmt.Fprintf(os.Stderr, "randmac: unknown vendor %q\n", *vendor)
			os.Exit(1)
		}
	}

	if prefix != "" {
		parts := strings.FieldsFunc(prefix, func(r rune) bool {
			return r == ':' || r == '-' || r == '.'
		})
		for i, p := range parts {
			if i >= 3 {
				break
			}
			var v int
			fmt.Sscanf(p, "%x", &v)
			b[i] = byte(v)
		}
		// OUI first byte: clear locally administered if vendor-specified
		if *vendor == "" {
			b[0] &^= 0x02 // clear LA bit for specified OUI
		}
	}

	parts := make([]string, 6)
	for i, x := range b {
		parts[i] = fmt.Sprintf("%02x", x)
	}
	result := strings.Join(parts, *sep)
	if *upper {
		result = strings.ToUpper(result)
	}
	return result
}

func main() {
	flag.Parse()
	for i := 0; i < *count; i++ {
		fmt.Println(genMAC())
	}
}
