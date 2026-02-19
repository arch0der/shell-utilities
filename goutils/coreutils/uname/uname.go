package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]
	all := false
	kernel := false
	nodename := false
	release := false
	version := false
	machine := false
	processor := false
	hardware := false
	os2 := false

	if len(args) == 0 {
		kernel = true
	}

	for _, a := range args {
		switch a {
		case "-a", "--all":
			all = true
		case "-s", "--kernel-name":
			kernel = true
		case "-n", "--nodename":
			nodename = true
		case "-r", "--kernel-release":
			release = true
		case "-v", "--kernel-version":
			version = true
		case "-m", "--machine":
			machine = true
		case "-p", "--processor":
			processor = true
		case "-i", "--hardware-platform":
			hardware = true
		case "-o", "--operating-system":
			os2 = true
		}
	}

	var utsname syscall.Utsname
	syscall.Uname(&utsname)

	int8ToStr := func(arr [65]int8) string {
		b := make([]byte, 0, 65)
		for _, v := range arr {
			if v == 0 {
				break
			}
			b = append(b, byte(v))
		}
		return string(b)
	}

	var parts []string
	add := func(cond bool, val string) {
		if cond || all {
			parts = append(parts, val)
		}
	}

	add(kernel, int8ToStr(utsname.Sysname))
	add(nodename, int8ToStr(utsname.Nodename))
	add(release, int8ToStr(utsname.Release))
	add(version, int8ToStr(utsname.Version))

	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	case "386":
		arch = "i686"
	}
	add(machine, arch)
	add(processor, arch)
	add(hardware, arch)
	add(os2, runtime.GOOS)

	fmt.Println(strings.Join(parts, " "))
}
