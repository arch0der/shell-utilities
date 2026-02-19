// dotenv - load a .env file and print export statements, or run a command with env vars injected
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func parseEnv(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()
	var pairs []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }
		// Strip inline comments
		if idx := strings.Index(line, " #"); idx > 0 { line = strings.TrimSpace(line[:idx]) }
		// KEY=VALUE or KEY="VALUE"
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 { continue }
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		pairs = append(pairs, key+"="+val)
	}
	return pairs, nil
}

func main() {
	envFile := ".env"
	args := os.Args[1:]

	if len(args) > 0 && strings.HasSuffix(args[0], ".env") {
		envFile = args[0]; args = args[1:]
	}
	if len(args) > 0 && args[0] == "-f" && len(args) > 1 {
		envFile = args[1]; args = args[2:]
	}

	pairs, err := parseEnv(envFile)
	if err != nil { fmt.Fprintln(os.Stderr, "dotenv:", err); os.Exit(1) }

	if len(args) == 0 {
		// Print export statements
		for _, p := range pairs { fmt.Printf("export %s\n", p) }
		return
	}

	// Run command with env vars
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), pairs...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok { os.Exit(exit.ExitCode()) }
		fmt.Fprintln(os.Stderr, "dotenv:", err); os.Exit(1)
	}
}
