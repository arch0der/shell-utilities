// binwatch - watch a file or directory for changes and run a command
package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func hashFile(path string) string {
	f, err := os.Open(path); if err != nil { return "" }
	defer f.Close()
	h := md5.New(); io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func hashDir(root string, pattern string) string {
	h := md5.New()
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() { return nil }
		if pattern != "" {
			matched, _ := filepath.Match(pattern, info.Name())
			if !matched { return nil }
		}
		fmt.Fprintf(h, "%s:%d:%d:", path, info.Size(), info.ModTime().UnixNano())
		return nil
	})
	return fmt.Sprintf("%x", h.Sum(nil))
}

func runCmd(cmdStr string) {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 { return }
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	fmt.Fprintf(os.Stderr, "\n[binwatch] Running: %s\n", cmdStr)
	if err := cmd.Run(); err != nil { fmt.Fprintf(os.Stderr, "[binwatch] Command failed: %v\n", err) }
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: binwatch [options] <path> <command>
  -i <interval>   poll interval (default: 1s)
  -p <pattern>    glob pattern to filter (for dir watching)
  -q              quiet: no banner messages
  -1              run once immediately then watch`)
	os.Exit(1)
}

func main() {
	interval := time.Second
	pattern := ""
	quiet := false
	runOnce := false
	var watchPath, command string

	args := os.Args[1:]
	rest := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-i": i++; interval, _ = time.ParseDuration(args[i]); if interval == 0 { interval = time.Second }
		case "-p": i++; pattern = args[i]
		case "-q": quiet = true
		case "-1": runOnce = true
		default: rest = append(rest, args[i])
		}
	}
	if len(rest) < 2 { usage() }
	watchPath = rest[0]
	command = strings.Join(rest[1:], " ")

	hash := func() string {
		info, err := os.Stat(watchPath); if err != nil { return "" }
		if info.IsDir() { return hashDir(watchPath, pattern) }
		return hashFile(watchPath)
	}

	if !quiet { fmt.Fprintf(os.Stderr, "[binwatch] Watching %s every %s\n", watchPath, interval) }
	last := hash()
	if runOnce { runCmd(command); last = hash() }

	for {
		time.Sleep(interval)
		cur := hash()
		if cur != last {
			last = cur
			if !quiet { fmt.Fprintf(os.Stderr, "[binwatch] Change detected in %s\n", watchPath) }
			runCmd(command)
		}
	}
}
