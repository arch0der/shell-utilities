// kill - Send signals to processes
// Usage: kill [-s signal] <pid>...
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

var sigName = flag.String("s", "TERM", "Signal name (TERM, KILL, HUP, INT, etc.)")

var signals = map[string]syscall.Signal{
	"HUP":  syscall.SIGHUP,
	"INT":  syscall.SIGINT,
	"QUIT": syscall.SIGQUIT,
	"KILL": syscall.SIGKILL,
	"TERM": syscall.SIGTERM,
	"USR1": syscall.SIGUSR1,
	"USR2": syscall.SIGUSR2,
	"STOP": syscall.SIGSTOP,
	"CONT": syscall.SIGCONT,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: kill [-s signal] <pid>...")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	sig, ok := signals[*sigName]
	if !ok {
		fmt.Fprintf(os.Stderr, "kill: unknown signal %s\n", *sigName)
		os.Exit(1)
	}

	exitCode := 0
	for _, pidStr := range flag.Args() {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: invalid pid: %s\n", pidStr)
			exitCode = 1
			continue
		}
		p, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintln(os.Stderr, "kill:", err)
			exitCode = 1
			continue
		}
		if err := p.Signal(sig); err != nil {
			fmt.Fprintln(os.Stderr, "kill:", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
