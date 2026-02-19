package main

import (
	"fmt"
	"runtime"
)

func main() {
	arch := runtime.GOARCH
	// Map Go arch names to uname-style names
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i686"
	case "arm64":
		arch = "aarch64"
	case "arm":
		arch = "armv7l"
	}
	fmt.Println(arch)
}
