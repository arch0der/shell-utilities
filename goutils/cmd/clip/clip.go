// clip - copy stdin to clipboard (xclip/xsel/pbcopy) and optionally echo it
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	paste := len(os.Args) > 1 && (os.Args[1] == "-p" || os.Args[1] == "--paste")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		if paste { cmd = exec.Command("pbpaste") } else { cmd = exec.Command("pbcopy") }
	default:
		if paste {
			if _, err := exec.LookPath("xclip"); err == nil {
				cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
			} else {
				cmd = exec.Command("xsel", "--clipboard", "--output")
			}
		} else {
			if _, err := exec.LookPath("xclip"); err == nil {
				cmd = exec.Command("xclip", "-selection", "clipboard")
			} else {
				cmd = exec.Command("xsel", "--clipboard", "--input")
			}
		}
	}

	if paste {
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "clip: clipboard read failed:", err); os.Exit(1)
		}
		return
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	cmd.Stdin = os.Stdin
	in, _ := cmd.StdinPipe()
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "clip: clipboard tool not found:", err); os.Exit(1)
	}
	in.Write(data); in.Close()
	cmd.Wait()
	os.Stdout.Write(data)
}
