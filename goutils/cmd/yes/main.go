// yes - Output a string repeatedly until killed
// Usage: yes [string]
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	s := "y"
	if len(os.Args) > 1 {
		s = strings.Join(os.Args[1:], " ")
	}
	for {
		fmt.Println(s)
	}
}
