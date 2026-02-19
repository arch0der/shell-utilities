// scp - Secure file copy (wraps ssh for remote, copies locally otherwise)
// Usage: scp [-r] [-P port] source destination
// Remote format: [user@]host:path
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	recursive = flag.Bool("r", false, "Recursive copy")
	port      = flag.String("P", "22", "Remote port")
	preserve  = flag.Bool("p", false, "Preserve modification times and modes")
	quiet     = flag.Bool("q", false, "Quiet mode")
	verbose   = flag.Bool("v", false, "Verbose")
)

func isRemote(s string) bool {
	// remote if contains : and the part before : looks like host[:port]
	idx := strings.Index(s, ":")
	if idx < 0 {
		return false
	}
	host := s[:idx]
	return !strings.HasPrefix(host, "/") && !strings.HasPrefix(host, ".")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: scp [-r] [-P port] [-p] [-q] source destination")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	src := flag.Arg(0)
	dst := flag.Arg(1)

	// If either side is remote, use system ssh/scp
	if isRemote(src) || isRemote(dst) {
		// Try system scp
		scpPath, err := exec.LookPath("scp")
		if err == nil {
			args := []string{}
			if *recursive {
				args = append(args, "-r")
			}
			if *port != "22" {
				args = append(args, "-P", *port)
			}
			if *preserve {
				args = append(args, "-p")
			}
			if *quiet {
				args = append(args, "-q")
			}
			if *verbose {
				args = append(args, "-v")
			}
			args = append(args, src, dst)
			cmd := exec.Command(scpPath, args...)
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
		fmt.Fprintln(os.Stderr, "scp: remote copy requires ssh/scp to be installed")
		os.Exit(1)
	}

	// Local copy
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scp:", err)
		os.Exit(1)
	}

	if srcInfo.IsDir() {
		if !*recursive {
			fmt.Fprintln(os.Stderr, "scp: omitting directory (use -r)")
			os.Exit(1)
		}
		if err := copyDir(src, dst); err != nil {
			fmt.Fprintln(os.Stderr, "scp:", err)
			os.Exit(1)
		}
	} else {
		if err := copyFile(src, dst); err != nil {
			fmt.Fprintln(os.Stderr, "scp:", err)
			os.Exit(1)
		}
	}
}

func copyDir(src, dst string) error {
	srcInfo, _ := os.Stat(src)
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
		} else {
			if err := copyFile(s, d); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	// If dst is a directory, put file inside it
	if info, err := os.Stat(dst); err == nil && info.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	srcInfo, _ := sf.Stat()
	df, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer df.Close()

	n, err := io.Copy(df, sf)
	if !*quiet {
		fmt.Printf("%s -> %s (%d bytes)\n", src, dst, n)
	}
	return err
}
