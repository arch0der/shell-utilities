package main

import (
	"fmt"
	"net"
	"os"
)

func init() { register("hostid", runHostid) }

func runHostid() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Fprintln(os.Stderr, "hostid:", err)
		os.Exit(1)
	}
	var id uint32
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) >= 4 {
			id = uint32(iface.HardwareAddr[0])<<24 |
				uint32(iface.HardwareAddr[1])<<16 |
				uint32(iface.HardwareAddr[2])<<8 |
				uint32(iface.HardwareAddr[3])
			break
		}
	}
	fmt.Printf("%08x\n", id)
}
