package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func init() { register("kill", runKill) }

var signalMap = map[string]syscall.Signal{
	"HUP": syscall.SIGHUP, "INT": syscall.SIGINT, "QUIT": syscall.SIGQUIT,
	"ILL": syscall.SIGILL, "TRAP": syscall.SIGTRAP, "ABRT": syscall.SIGABRT,
	"KILL": syscall.SIGKILL, "SEGV": syscall.SIGSEGV, "PIPE": syscall.SIGPIPE,
	"ALRM": syscall.SIGALRM, "TERM": syscall.SIGTERM, "USR1": syscall.SIGUSR1,
	"USR2": syscall.SIGUSR2, "CONT": syscall.SIGCONT, "STOP": syscall.SIGSTOP,
	"TSTP": syscall.SIGTSTP, "TTIN": syscall.SIGTTIN, "TTOU": syscall.SIGTTOU,
}

func runKill() {
	args := os.Args[1:]
	sig := syscall.SIGTERM
	listSignals := false
	pids := []int{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-l" || a == "--list" {
			listSignals = true
		} else if strings.HasPrefix(a, "-") {
			sigName := strings.TrimPrefix(a, "-")
			sigName = strings.TrimPrefix(sigName, "SIG")
			if s, ok := signalMap[strings.ToUpper(sigName)]; ok {
				sig = s
			} else if n, err := strconv.Atoi(sigName); err == nil {
				sig = syscall.Signal(n)
			}
		} else {
			pid, _ := strconv.Atoi(a)
			pids = append(pids, pid)
		}
	}

	if listSignals {
		for name := range signalMap {
			fmt.Println(name)
		}
		return
	}

	exitCode := 0
	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: %d: %v\n", pid, err)
			exitCode = 1
			continue
		}
		if err := proc.Signal(sig); err != nil {
			fmt.Fprintf(os.Stderr, "kill: %d: %v\n", pid, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
