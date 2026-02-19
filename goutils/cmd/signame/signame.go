// signame - convert between signal numbers and names
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var signals = map[int]string{
	1:"SIGHUP",2:"SIGINT",3:"SIGQUIT",4:"SIGILL",5:"SIGTRAP",6:"SIGABRT",
	7:"SIGBUS",8:"SIGFPE",9:"SIGKILL",10:"SIGUSR1",11:"SIGSEGV",12:"SIGUSR2",
	13:"SIGPIPE",14:"SIGALRM",15:"SIGTERM",16:"SIGSTKFLT",17:"SIGCHLD",
	18:"SIGCONT",19:"SIGSTOP",20:"SIGTSTP",21:"SIGTTIN",22:"SIGTTOU",
	23:"SIGURG",24:"SIGXCPU",25:"SIGXFSZ",26:"SIGVTALRM",27:"SIGPROF",
	28:"SIGWINCH",29:"SIGIO",30:"SIGPWR",31:"SIGSYS",
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("%-6s %-14s %s\n", "Num", "Name", "Description")
		fmt.Println(strings.Repeat("─", 50))
		for n := 1; n <= 31; n++ {
			name := signals[n]
			sig := syscall.Signal(n)
			fmt.Printf("%-6d %-14s %s\n", n, name, sig.String())
		}
		return
	}
	for _, arg := range os.Args[1:] {
		arg = strings.ToUpper(strings.TrimPrefix(strings.TrimPrefix(arg, "SIG"), "sig"))
		// Try number
		if n, err := strconv.Atoi(arg); err == nil {
			if name, ok := signals[n]; ok {
				fmt.Printf("%d = %s (%s)\n", n, name, syscall.Signal(n).String())
			} else {
				fmt.Printf("%d = unknown signal\n", n)
			}
			continue
		}
		// Try name
		name := "SIG" + arg
		for n, sn := range signals {
			if sn == name || sn == "SIG"+arg {
				fmt.Printf("%s = %d (%s)\n", sn, n, syscall.Signal(n).String())
				goto next
			}
		}
		fmt.Printf("%s = unknown signal\n", arg)
		next:
	}
}
