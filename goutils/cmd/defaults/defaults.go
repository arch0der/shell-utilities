// defaults - print a value or a default if empty/missing (env var helper)
package main

import (
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: defaults <value> <default>  |  defaults -e VAR <default>
  Print <value> if non-empty, otherwise print <default>.
  -e VAR    use environment variable VAR as the value
  -E VAR    also export the result as VAR=... to stdout for eval
  Supports multiple pairs: defaults VAL1 DEF1 VAL2 DEF2 ...`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 { usage() }
	args := os.Args[1:]

	useEnv := false
	exportMode := false

	if args[0] == "-e" || args[0] == "-E" {
		exportMode = args[0] == "-E"
		useEnv = true
		args = args[1:]
	}

	if useEnv {
		if len(args) < 2 { usage() }
		varName := args[0]
		def := strings.Join(args[1:], " ")
		val := os.Getenv(varName)
		if val == "" { val = def }
		if exportMode { fmt.Printf("export %s=%q\n", varName, val) } else { fmt.Println(val) }
		return
	}

	// Process pairs: val default val default ...
	if len(args)%2 != 0 { usage() }
	for i := 0; i < len(args); i += 2 {
		val, def := args[i], args[i+1]
		if val == "" || val == "null" || val == "nil" || val == "None" { val = def }
		fmt.Println(val)
	}
}
