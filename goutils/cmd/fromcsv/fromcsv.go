// fromcsv - convert CSV to various formats: JSON, YAML, TSV, Markdown table, SQL INSERT
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: fromcsv [options] [file]
  -f <format>   output format: json|jsonl|tsv|md|sql|kv (default: json)
  -t <table>    table name for SQL output (default: data)
  -H            first row is NOT a header (auto-generate col1,col2...)
  -d <delim>    input delimiter (default: ,)`)
	os.Exit(1)
}

func main() {
	format := "json"
	table := "data"
	noHeader := false
	delim := ','
	var file string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f": i++; format = args[i]
		case "-t": i++; table = args[i]
		case "-H": noHeader = true
		case "-d": i++; if len(args[i]) > 0 { delim = rune(args[i][0]) }
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			file = args[i]
		}
	}

	var r io.Reader = os.Stdin
	if file != "" {
		f, err := os.Open(file); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer f.Close(); r = f
	}

	cr := csv.NewReader(r)
	cr.Comma = delim; cr.LazyQuotes = true; cr.TrimLeadingSpace = true
	records, err := cr.ReadAll()
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	if len(records) == 0 { return }

	var headers []string
	var dataRows [][]string
	if noHeader {
		for i := range records[0] { headers = append(headers, fmt.Sprintf("col%d", i+1)) }
		dataRows = records
	} else {
		headers = records[0]; dataRows = records[1:]
	}

	switch format {
	case "json":
		var out []map[string]interface{}
		for _, row := range dataRows {
			m := map[string]interface{}{}
			for i, h := range headers { if i < len(row) { m[h] = row[i] } else { m[h] = "" } }
			out = append(out, m)
		}
		b, _ := json.MarshalIndent(out, "", "  "); fmt.Println(string(b))
	case "jsonl":
		for _, row := range dataRows {
			m := map[string]interface{}{}
			for i, h := range headers { if i < len(row) { m[h] = row[i] } else { m[h] = "" } }
			b, _ := json.Marshal(m); fmt.Println(string(b))
		}
	case "tsv":
		fmt.Println(strings.Join(headers, "\t"))
		for _, row := range dataRows { fmt.Println(strings.Join(row, "\t")) }
	case "md":
		fmt.Printf("| %s |\n", strings.Join(headers, " | "))
		seps := make([]string, len(headers)); for i := range seps { seps[i] = "---" }
		fmt.Printf("| %s |\n", strings.Join(seps, " | "))
		for _, row := range dataRows {
			padded := make([]string, len(headers))
			for i, h := range headers { _ = h; if i < len(row) { padded[i] = row[i] } }
			fmt.Printf("| %s |\n", strings.Join(padded, " | "))
		}
	case "sql":
		cols := make([]string, len(headers))
		for i, h := range headers { cols[i] = "`" + h + "`" }
		for _, row := range dataRows {
			vals := make([]string, len(headers))
			for i := range headers {
				v := ""; if i < len(row) { v = row[i] }
				vals[i] = "'" + strings.ReplaceAll(v, "'", "''") + "'"
			}
			fmt.Printf("INSERT INTO `%s` (%s) VALUES (%s);\n", table, strings.Join(cols, ", "), strings.Join(vals, ", "))
		}
	case "kv":
		for _, row := range dataRows {
			for i, h := range headers { if i < len(row) { fmt.Printf("%s=%s\n", h, row[i]) } }
			fmt.Println()
		}
	default:
		fmt.Fprintf(os.Stderr, "fromcsv: unknown format %q\n", format); usage()
	}
}
