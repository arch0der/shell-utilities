// confirm - Interactive yes/no confirmation prompt for shell scripts.
//
// Usage:
//
//	confirm [OPTIONS] [PROMPT]
//
// Options:
//
//	-d Y|N    Default answer if Enter pressed (default: N)
//	-t DUR    Auto-answer with default after timeout
//	-y        Non-interactive: always answer yes (for CI)
//	-n        Non-interactive: always answer no
//	-q        Quiet: no prompt, just read answer
//
// Exit codes:
//
//	0 = yes
//	1 = no
//
// Examples:
//
//	confirm "Deploy to production?" && ./deploy.sh
//	confirm -d Y "Continue?" || exit 1
//	confirm -t 10s "Auto-yes in 10s?" -d Y
//	confirm -y   # always yes (CI mode)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	defaultAns = flag.String("d", "N", "default answer (Y/N)")
	timeout    = flag.Duration("t", 0, "auto-answer timeout")
	alwaysYes  = flag.Bool("y", false, "always yes")
	alwaysNo   = flag.Bool("n", false, "always no")
	quiet      = flag.Bool("q", false, "quiet mode")
)

func main() {
	flag.Parse()
	prompt := strings.Join(flag.Args(), " ")

	if *alwaysYes {
		os.Exit(0)
	}
	if *alwaysNo {
		os.Exit(1)
	}

	def := strings.ToUpper(*defaultAns)
	hint := "[y/N]"
	if def == "Y" {
		hint = "[Y/n]"
	}

	if !*quiet {
		if prompt != "" {
			fmt.Fprintf(os.Stderr, "%s %s ", prompt, hint)
		} else {
			fmt.Fprintf(os.Stderr, "Confirm %s ", hint)
		}
	}

	if *timeout > 0 {
		// Show countdown
		done := make(chan string, 1)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			done <- strings.TrimSpace(line)
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		remaining := int(*timeout / time.Second)

		for {
			select {
			case ans := <-done:
				ticker.Stop()
				fmt.Fprintf(os.Stderr, "\n")
				if isYes(ans, def) {
					os.Exit(0)
				}
				os.Exit(1)
			case <-ticker.C:
				remaining--
				fmt.Fprintf(os.Stderr, "\r%s %s [%ds] ", prompt, hint, remaining)
				if remaining <= 0 {
					fmt.Fprintf(os.Stderr, "\n")
					if def == "Y" {
						os.Exit(0)
					}
					os.Exit(1)
				}
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		// EOF or non-interactive: use default
		if def == "Y" {
			os.Exit(0)
		}
		os.Exit(1)
	}
	ans := strings.TrimSpace(line)
	if isYes(ans, def) {
		os.Exit(0)
	}
	os.Exit(1)
}

func isYes(ans, def string) bool {
	if ans == "" {
		return def == "Y"
	}
	return strings.ToUpper(ans) == "Y" || strings.ToLower(ans) == "yes"
}
