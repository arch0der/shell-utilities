// ini2json - Convert INI/config files to JSON.
//
// Usage:
//
//	ini2json [OPTIONS] [FILE...]
//	cat config.ini | ini2json
//
// Options:
//
//	-c        Compact output
//	-s        Flatten sections (no nesting, key becomes "section.key")
//	-t        Try to parse values as numbers/booleans
//
// Examples:
//
//	ini2json config.ini
//	cat /etc/php.ini | ini2json -t
//	ini2json -s app.conf | jq '.database.host'
//
// Supports:
//   [section] headers
//   key = value pairs
//   key: value pairs
//   # and ; comments
//   Global (no section) keys go to top level
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	compact = flag.Bool("c", false, "compact")
	flat    = flag.Bool("s", false, "flatten with section.key")
	typed   = flag.Bool("t", false, "parse types")
)

func parseValue(s string) interface{} {
	if !*typed {
		return s
	}
	if s == "true" || s == "yes" || s == "on" {
		return true
	}
	if s == "false" || s == "no" || s == "off" {
		return false
	}
	if s == "null" || s == "~" {
		return nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func process(r *os.File) map[string]interface{} {
	result := make(map[string]interface{})
	sections := make(map[string]map[string]interface{})
	currentSection := ""

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			if _, ok := sections[currentSection]; !ok {
				sections[currentSection] = make(map[string]interface{})
			}
			continue
		}
		// Key = value or Key: value
		var key, val string
		for _, sep := range []string{"=", ":"} {
			if idx := strings.Index(line, sep); idx > 0 {
				key = strings.TrimSpace(line[:idx])
				val = strings.TrimSpace(line[idx+1:])
				// Strip inline comments
				for _, c := range []string{" #", " ;"} {
					if i := strings.Index(val, c); i >= 0 {
						val = strings.TrimSpace(val[:i])
					}
				}
				// Strip quotes
				val = strings.Trim(val, `"'`)
				break
			}
		}
		if key == "" {
			continue
		}
		parsed := parseValue(val)
		if currentSection == "" {
			result[key] = parsed
		} else {
			if *flat {
				result[currentSection+"."+key] = parsed
			} else {
				if sections[currentSection] == nil {
					sections[currentSection] = make(map[string]interface{})
				}
				sections[currentSection][key] = parsed
			}
		}
	}

	if !*flat {
		for sec, vals := range sections {
			result[sec] = vals
		}
	}
	return result
}

func main() {
	flag.Parse()
	files := flag.Args()

	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ini2json: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	// Merge all results
	merged := make(map[string]interface{})
	for _, r := range readers {
		m := process(r)
		for k, v := range m {
			merged[k] = v
		}
	}

	var b []byte
	if *compact {
		b, _ = json.Marshal(merged)
	} else {
		b, _ = json.MarshalIndent(merged, "", "  ")
	}
	fmt.Println(string(b))
}
