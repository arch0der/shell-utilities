package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "printf: missing operand")
		os.Exit(1)
	}
	format := args[0]
	argList := args[1:]

	// Process format string with shell-style format specs
	result := sprintfShell(format, argList)
	fmt.Print(result)
}

func sprintfShell(format string, args []string) string {
	var sb strings.Builder
	argIdx := 0
	i := 0
	for i < len(format) {
		if format[i] == '\\' && i+1 < len(format) {
			i++
			switch format[i] {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '\\':
				sb.WriteByte('\\')
			case 'a':
				sb.WriteByte('\a')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'v':
				sb.WriteByte('\v')
			default:
				sb.WriteByte('\\')
				sb.WriteByte(format[i])
			}
			i++
			continue
		}
		if format[i] == '%' && i+1 < len(format) {
			i++
			// Parse flags/width/precision
			start := i - 1
			for i < len(format) && (format[i] == '-' || format[i] == '+' || format[i] == ' ' || format[i] == '#' || format[i] == '0') {
				i++
			}
			for i < len(format) && format[i] >= '0' && format[i] <= '9' {
				i++
			}
			if i < len(format) && format[i] == '.' {
				i++
				for i < len(format) && format[i] >= '0' && format[i] <= '9' {
					i++
				}
			}
			if i >= len(format) {
				break
			}
			spec := format[start : i+1]
			conv := format[i]
			i++
			arg := ""
			if argIdx < len(args) {
				arg = args[argIdx]
				argIdx++
			}
			switch conv {
			case 's':
				// replace %...s with proper format
				goFmt := strings.Replace(spec, "%", "", 1)
				goFmt = "%" + goFmt
				sb.WriteString(fmt.Sprintf(goFmt, arg))
			case 'd', 'i':
				n, _ := strconv.ParseInt(arg, 0, 64)
				goFmt := strings.Replace(spec, string(conv), "d", 1)
				sb.WriteString(fmt.Sprintf(goFmt, n))
			case 'u':
				n, _ := strconv.ParseUint(arg, 0, 64)
				goFmt := strings.Replace(spec, "u", "d", 1)
				sb.WriteString(fmt.Sprintf(goFmt, n))
			case 'o':
				n, _ := strconv.ParseInt(arg, 0, 64)
				sb.WriteString(fmt.Sprintf(spec, n))
			case 'x':
				n, _ := strconv.ParseInt(arg, 0, 64)
				sb.WriteString(fmt.Sprintf(spec, n))
			case 'X':
				n, _ := strconv.ParseInt(arg, 0, 64)
				sb.WriteString(fmt.Sprintf(spec, n))
			case 'f', 'e', 'E', 'g', 'G':
				n, _ := strconv.ParseFloat(arg, 64)
				sb.WriteString(fmt.Sprintf(spec, n))
			case 'c':
				if len(arg) > 0 {
					sb.WriteByte(arg[0])
				}
			case '%':
				sb.WriteByte('%')
				argIdx-- // don't consume arg
			default:
				sb.WriteString(spec)
			}
			continue
		}
		sb.WriteByte(format[i])
		i++
	}
	return sb.String()
}
