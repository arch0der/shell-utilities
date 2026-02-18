// stat - Display file or filesystem status
// Usage: stat <file>...
package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: stat <file>...")
		os.Exit(1)
	}

	exitCode := 0
	for _, path := range os.Args[1:] {
		if err := statFile(path); err != nil {
			fmt.Fprintln(os.Stderr, "stat:", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func statFile(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	fileType := "regular file"
	if info.IsDir() {
		fileType = "directory"
	} else if info.Mode()&os.ModeSymlink != 0 {
		fileType = "symbolic link"
	} else if info.Mode()&os.ModeNamedPipe != 0 {
		fileType = "fifo"
	} else if info.Mode()&os.ModeDevice != 0 {
		fileType = "device"
	}

	fmt.Printf("  File: %s\n", path)
	fmt.Printf("  Size: %-15d FileType: %s\n", info.Size(), fileType)
	fmt.Printf("  Mode: %s (%04o)\n", info.Mode(), info.Mode().Perm())
	fmt.Printf("Modify: %s\n", info.ModTime().Format("2006-01-02 15:04:05.000000000 -0700"))

	// Linux-specific: Inode, UID, GID via syscall
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		fmt.Printf(" Inode: %-15d Links: %d\n", stat.Ino, stat.Nlink)
		fmt.Printf("   UID: %-10d  GID: %d\n", stat.Uid, stat.Gid)
		atime := stat.Atim
		fmt.Printf("Access: %d.%09d\n", atime.Sec, atime.Nsec)
	}
	return nil
}
