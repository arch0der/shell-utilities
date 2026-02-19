package main

import (
	"fmt"
	"os"
	"strconv"
)

func init() {
	register("test", runTest)
	register("[", runBracket)
}

func runBracket() {
	args := os.Args[1:]
	if len(args) > 0 && args[len(args)-1] == "]" {
		args = args[:len(args)-1]
	}
	os.Args = append([]string{"test"}, args...)
	runTest()
}

func runTest() {
	args := os.Args[1:]
	if evalTest(args) {
		os.Exit(0)
	}
	os.Exit(1)
}

func evalTest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	if args[0] == "!" {
		return !evalTest(args[1:])
	}
	for i, a := range args {
		if a == "-a" && i > 0 {
			return evalTest(args[:i]) && evalTest(args[i+1:])
		}
		if a == "-o" && i > 0 {
			return evalTest(args[:i]) || evalTest(args[i+1:])
		}
	}
	if len(args) == 1 {
		return args[0] != ""
	}
	if len(args) == 2 {
		op, val := args[0], args[1]
		switch op {
		case "-z":
			return val == ""
		case "-n":
			return val != ""
		case "-e":
			_, err := os.Stat(val)
			return err == nil
		case "-f":
			info, err := os.Stat(val)
			return err == nil && !info.IsDir()
		case "-d":
			info, err := os.Stat(val)
			return err == nil && info.IsDir()
		case "-r":
			f, err := os.Open(val)
			if err == nil { f.Close() }
			return err == nil
		case "-w":
			f, err := os.OpenFile(val, os.O_WRONLY, 0)
			if err == nil { f.Close() }
			return err == nil
		case "-x":
			info, err := os.Stat(val)
			return err == nil && info.Mode()&0111 != 0
		case "-s":
			info, err := os.Stat(val)
			return err == nil && info.Size() > 0
		case "-L", "-h":
			_, err := os.Lstat(val)
			return err == nil
		case "-p":
			info, err := os.Lstat(val)
			return err == nil && info.Mode()&os.ModeNamedPipe != 0
		case "-S":
			info, err := os.Lstat(val)
			return err == nil && info.Mode()&os.ModeSocket != 0
		}
		return false
	}
	if len(args) == 3 {
		left, op, right := args[0], args[1], args[2]
		switch op {
		case "=", "==":
			return left == right
		case "!=":
			return left != right
		case "<":
			return left < right
		case ">":
			return left > right
		}
		a2, aerr := strconv.ParseInt(left, 10, 64)
		b2, berr := strconv.ParseInt(right, 10, 64)
		if aerr != nil || berr != nil {
			fmt.Fprintln(os.Stderr, "test: integer expression expected")
			return false
		}
		switch op {
		case "-eq":
			return a2 == b2
		case "-ne":
			return a2 != b2
		case "-lt":
			return a2 < b2
		case "-le":
			return a2 <= b2
		case "-gt":
			return a2 > b2
		case "-ge":
			return a2 >= b2
		}
	}
	return false
}
