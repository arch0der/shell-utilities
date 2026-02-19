// script - Record a terminal session to a typescript file
// Usage: script [-a] [-q] [-t timingfile] [file]
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	appendMode = flag.Bool("a", false, "Append to typescript file")
	quiet      = flag.Bool("q", false, "Quiet mode")
	timing     = flag.String("t", "", "Write timing data to file")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: script [-a] [-q] [-t timingfile] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()

	outFile := "typescript"
	if flag.NArg() > 0 {
		outFile = flag.Arg(0)
	}

	flags := os.O_CREATE | os.O_WRONLY
	if *appendMode {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	f, err := os.OpenFile(outFile, flags, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "script:", err)
		os.Exit(1)
	}
	defer f.Close()

	var timingFile *os.File
	if *timing != "" {
		timingFile, err = os.Create(*timing)
		if err != nil {
			fmt.Fprintln(os.Stderr, "script:", err)
			os.Exit(1)
		}
		defer timingFile.Close()
	}

	if !*quiet {
		msg := fmt.Sprintf("Script started, output log file is '%s'\n", outFile)
		fmt.Print(msg)
		f.WriteString(msg)
	}

	header := fmt.Sprintf("Script started on %s\n", time.Now().Format("Mon Jan  2 15:04:05 2006"))
	f.WriteString(header)

	// Run shell and capture output
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell)
	cmd.Stdin = os.Stdin

	// Tee stdout/stderr to both terminal and file
	start := time.Now()
	pr, pw, _ := os.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	// Copy from pipe to both stdout and file
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := pr.Read(buf)
			if n > 0 {
				data := buf[:n]
				os.Stdout.Write(data)
				f.Write(data)
				if timingFile != nil {
					elapsed := time.Since(start).Seconds()
					fmt.Fprintf(timingFile, "%.6f %d\n", elapsed, n)
				}
			}
			if err != nil {
				break
			}
		}
	}()

	// Handle SIGWINCH (terminal resize)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	go func() {
		for range sigCh {
			// Terminal resized - in a real pty-based impl we'd forward this
		}
	}()

	cmd.Run()
	pw.Close()
	time.Sleep(100 * time.Millisecond) // Let goroutine flush

	footer := fmt.Sprintf("\nScript done on %s\n", time.Now().Format("Mon Jan  2 15:04:05 2006"))
	f.WriteString(footer)

	if !*quiet {
		fmt.Printf("Script done, output log file is '%s'\n", outFile)
	}
}
