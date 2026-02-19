// logfmt - parse and format logfmt-style log lines (key=value pairs)
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func parseLogfmt(line string) map[string]string {
	result := map[string]string{}
	line = strings.TrimSpace(line)
	for line != "" {
		line = strings.TrimLeft(line, " ")
		if line == "" { break }
		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 { break }
		key := line[:eqIdx]
		line = line[eqIdx+1:]
		var val string
		if strings.HasPrefix(line, `"`) {
			end := strings.Index(line[1:], `"`) + 1
			if end <= 0 { val = line[1:]; line = "" } else { val = line[1:end]; line = strings.TrimLeft(line[end+1:], " ") }
		} else {
			spaceIdx := strings.Index(line, " ")
			if spaceIdx < 0 { val = line; line = "" } else { val = line[:spaceIdx]; line = line[spaceIdx+1:] }
		}
		result[key] = val
	}
	return result
}

func formatLogfmt(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m { keys = append(keys, k) }
	// Priority keys first
	priority := []string{"time","ts","level","msg","message","error","err"}
	seen := map[string]bool{}
	var ordered []string
	for _, k := range priority { if _, ok := m[k]; ok { ordered = append(ordered, k); seen[k] = true } }
	sort.Strings(keys)
	for _, k := range keys { if !seen[k] { ordered = append(ordered, k) } }
	parts := make([]string, len(ordered))
	for i, k := range ordered {
		v := m[k]
		if strings.Contains(v, " ") { v = `"` + v + `"` }
		parts[i] = k + "=" + v
	}
	return strings.Join(parts, " ")
}

func main() {
	format := "pretty"
	filter := ""
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-j", "--json": format = "json"
		case "-r", "--raw": format = "raw"
		case "-f": i++; filter = args[i]
		}
	}

	colorLevel := map[string]string{
		"error":"ERROR","err":"ERROR","fatal":"FATAL","crit":"FATAL",
		"warn":"WARN","warning":"WARN",
		"info":"INFO","debug":"DEBUG","trace":"TRACE",
	}
	levelColor := map[string]string{
		"ERROR":"\033[31m","FATAL":"\033[1;31m","WARN":"\033[33m",
		"INFO":"\033[32m","DEBUG":"\033[36m","TRACE":"\033[2m",
	}
	reset := "\033[0m"

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		m := parseLogfmt(line)
		if len(m) == 0 { fmt.Println(line); continue }
		if filter != "" {
			match := false
			for _, v := range m { if strings.Contains(v, filter) { match = true; break } }
			if !match { continue }
		}
		switch format {
		case "json":
			b, _ := json.Marshal(m); fmt.Println(string(b))
		case "raw":
			fmt.Println(formatLogfmt(m))
		default: // pretty
			ts := m["time"]; if ts == "" { ts = m["ts"] }
			if ts != "" {
				for _, f := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
					if t, err := time.Parse(f, ts); err == nil { ts = t.Format("15:04:05"); break }
				}
			}
			level := strings.ToUpper(m["level"])
			if cl, ok := colorLevel[strings.ToLower(m["level"])]; ok { level = cl }
			msg := m["msg"]; if msg == "" { msg = m["message"] }
			col := levelColor[level]
			fmt.Printf("%s%s %-5s%s %s", col, ts, level, reset, msg)
			for k, v := range m {
				if k == "time" || k == "ts" || k == "level" || k == "msg" || k == "message" { continue }
				fmt.Printf("  %s\033[2m=%s%s", k, v, reset)
			}
			fmt.Println()
		}
	}
}
