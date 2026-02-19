// strace - Trace system calls (wraps system strace or uses /proc/pid/syscall)
// Usage: strace [-c] [-e syscall] <command> [args...]
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	summary = flag.Bool("c", false, "Count and summarize syscalls")
	filter  = flag.String("e", "", "Trace only these syscalls (comma-separated)")
	output  = flag.String("o", "", "Write output to file")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: strace [-c] [-e syscall,...] [-o file] <command> [args...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Use system strace if available
	if sysstrace, err := exec.LookPath("strace"); err == nil {
		args := []string{}
		if *summary {
			args = append(args, "-c")
		}
		if *filter != "" {
			args = append(args, "-e", "trace="+*filter)
		}
		if *output != "" {
			args = append(args, "-o", *output)
		}
		args = append(args, flag.Args()...)
		cmd := exec.Command(sysstrace, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				os.Exit(e.ExitCode())
			}
		}
		return
	}

	// Fallback: run command and poll /proc/pid/syscall
	filterSet := map[string]bool{}
	if *filter != "" {
		for _, s := range strings.Split(*filter, ",") {
			filterSet[strings.TrimSpace(s)] = true
		}
	}

	var logOut *os.File = os.Stderr
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintln(os.Stderr, "strace:", err)
			os.Exit(1)
		}
		defer f.Close()
		logOut = f
	}

	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	counts := map[string]int{}
	start := time.Now()

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "strace:", err)
		os.Exit(1)
	}

	pid := cmd.Process.Pid
	done := make(chan struct{})

	go func() {
		defer close(done)
		prev := ""
		for {
			data, err := os.ReadFile(fmt.Sprintf("/proc/%d/syscall", pid))
			if err != nil {
				return
			}
			cur := strings.TrimSpace(string(data))
			if cur != prev && cur != "" {
				fields := strings.Fields(cur)
				if len(fields) > 0 {
					name := syscallNumberToName(fields[0])
					if len(filterSet) == 0 || filterSet[name] {
						counts[name]++
						if !*summary {
							args := ""
							if len(fields) > 1 {
								args = strings.Join(fields[1:], ", ")
							}
							fmt.Fprintf(logOut, "%s(%s) = ?\n", name, args)
						}
					}
				}
				prev = cur
			}
			time.Sleep(500 * time.Microsecond)
		}
	}()

	cmd.Wait()
	elapsed := time.Since(start)
	<-done

	if *summary && len(counts) > 0 {
		total := 0
		for _, c := range counts {
			total += c
		}
		fmt.Fprintf(logOut, "%% time     seconds  calls   syscall\n")
		fmt.Fprintf(logOut, "------- ---------- ------- ------------------\n")
		for name, c := range counts {
			pct := float64(c) / float64(total) * 100
			secs := elapsed.Seconds() * float64(c) / float64(total)
			fmt.Fprintf(logOut, "%7.2f   %.6f %7d %s\n", pct, secs, c, name)
		}
		fmt.Fprintf(logOut, "------- ---------- ------- ------------------\n")
		fmt.Fprintf(logOut, "100.00   %.6f %7d total\n", elapsed.Seconds(), total)
	}
}

// x86_64 syscall number map (subset)
var syscallNums = map[string]string{
	"0": "read", "1": "write", "2": "open", "3": "close",
	"4": "stat", "5": "fstat", "6": "lstat", "7": "poll",
	"8": "lseek", "9": "mmap", "10": "mprotect", "11": "munmap",
	"12": "brk", "21": "access", "22": "pipe", "23": "select",
	"32": "dup", "33": "dup2", "39": "getpid", "41": "socket",
	"42": "connect", "43": "accept", "44": "sendto", "45": "recvfrom",
	"56": "clone", "57": "fork", "58": "vfork", "59": "execve",
	"60": "exit", "61": "wait4", "62": "kill", "63": "uname",
	"72": "fcntl", "78": "getdents", "79": "getcwd", "80": "chdir",
	"81": "fchdir", "82": "rename", "83": "mkdir", "84": "rmdir",
	"85": "creat", "86": "link", "87": "unlink", "88": "symlink",
	"89": "readlink", "90": "chmod", "91": "fchmod", "92": "chown",
	"257": "openat", "262": "newfstatat", "295": "openat",
}

func syscallNumberToName(num string) string {
	if name, ok := syscallNums[num]; ok {
		return name
	}
	return "sys_" + num
}
