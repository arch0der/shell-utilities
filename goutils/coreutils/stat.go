package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func init() { register("stat", runStat) }

func runStat() {
	args := os.Args[1:]
	format := ""
	dereference := false
	fileSystem := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-L" || a == "--dereference":
			dereference = true
		case a == "-f" || a == "--file-system":
			fileSystem = true
		case a == "--format" && i+1 < len(args):
			i++
			format = args[i]
		case strings.HasPrefix(a, "--format="):
			format = a[9:]
		case a == "-c" && i+1 < len(args):
			i++
			format = args[i]
		case strings.HasPrefix(a, "-c"):
			format = a[2:]
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	_ = fileSystem

	statFn := os.Lstat
	if dereference {
		statFn = os.Stat
	}

	exitCode := 0
	for _, f := range files {
		info, err := statFn(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "stat: %s: %v\n", f, err)
			exitCode = 1
			continue
		}
		stat := info.Sys().(*syscall.Stat_t)
		if format != "" {
			result := formatStat(format, f, info, stat)
			fmt.Println(result)
		} else {
			printStat(f, info, stat)
		}
	}
	os.Exit(exitCode)
}

func printStat(name string, info os.FileInfo, stat *syscall.Stat_t) {
	fmt.Printf("  File: %s\n", name)
	fmt.Printf("  Size: %-10d\tBlocks: %-10d IO Block: %-6d %s\n",
		info.Size(), stat.Blocks, stat.Blksize, fileTypeStr(info.Mode()))
	fmt.Printf("Device: %xh/%dd\tInode: %-10d  Links: %d\n",
		stat.Dev, stat.Dev, stat.Ino, stat.Nlink)
	fmt.Printf("Access: (%04o/%s)  Uid: (%5d/%8s)   Gid: (%5d/%8s)\n",
		uint32(info.Mode()), info.Mode().String(), stat.Uid, "?", stat.Gid, "?")
	fmt.Printf("Access: %s\n", time.Unix(stat.Atim.Sec, stat.Atim.Nsec).Format("2006-01-02 15:04:05.000000000 -0700"))
	fmt.Printf("Modify: %s\n", info.ModTime().Format("2006-01-02 15:04:05.000000000 -0700"))
	fmt.Printf("Change: %s\n", time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec).Format("2006-01-02 15:04:05.000000000 -0700"))
}

func fileTypeStr(m os.FileMode) string {
	switch {
	case m.IsDir():
		return "directory"
	case m&os.ModeSymlink != 0:
		return "symbolic link"
	case m&os.ModeNamedPipe != 0:
		return "fifo"
	case m&os.ModeSocket != 0:
		return "socket"
	case m&os.ModeDevice != 0:
		return "block special file"
	case m&os.ModeCharDevice != 0:
		return "character special file"
	default:
		return "regular file"
	}
}

func formatStat(format string, name string, info os.FileInfo, stat *syscall.Stat_t) string {
	var sb strings.Builder
	for i := 0; i < len(format); i++ {
		if format[i] != '%' || i+1 >= len(format) {
			sb.WriteByte(format[i])
			continue
		}
		i++
		switch format[i] {
		case 'n':
			sb.WriteString(name)
		case 's':
			sb.WriteString(fmt.Sprintf("%d", info.Size()))
		case 'i':
			sb.WriteString(fmt.Sprintf("%d", stat.Ino))
		case 'a':
			sb.WriteString(fmt.Sprintf("%o", uint32(info.Mode())&07777))
		case 'A':
			sb.WriteString(info.Mode().String())
		case 'u':
			sb.WriteString(fmt.Sprintf("%d", stat.Uid))
		case 'g':
			sb.WriteString(fmt.Sprintf("%d", stat.Gid))
		case 'h':
			sb.WriteString(fmt.Sprintf("%d", stat.Nlink))
		case 'Y':
			sb.WriteString(fmt.Sprintf("%d", info.ModTime().Unix()))
		case 'y':
			sb.WriteString(info.ModTime().Format("2006-01-02 15:04:05"))
		case 'F':
			sb.WriteString(fileTypeStr(info.Mode()))
		default:
			sb.WriteByte('%')
			sb.WriteByte(format[i])
		}
	}
	return sb.String()
}
