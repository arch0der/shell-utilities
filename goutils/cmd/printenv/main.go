// printenv - Print environment variables
// Usage: printenv [variable...]
package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	if len(os.Args) > 1 {
		// Print specific variables
		exitCode := 0
		for _, name := range os.Args[1:] {
			val, ok := os.LookupEnv(name)
			if !ok {
				exitCode = 1
				continue
			}
			fmt.Println(val)
		}
		os.Exit(exitCode)
	}

	// Print all variables, sorted
	env := os.Environ()
	sort.Strings(env)
	for _, e := range env {
		fmt.Println(e)
	}
}
