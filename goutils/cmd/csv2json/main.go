// csv2json - Convert CSV to JSON
// Usage: csv2json [-d delim] [-no-header] [file]
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	delim    = flag.String("d", ",", "Field delimiter")
	noHeader = flag.Bool("no-header", false, "Treat first row as data, not header")
	pretty   = flag.Bool("p", false, "Pretty-print output")
	array    = flag.Bool("a", false, "Output as array of arrays (not objects)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: csv2json [-d delim] [-no-header] [-p] [-a] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var r io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "csv2json:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	csvr := csv.NewReader(r)
	if len(*delim) > 0 {
		csvr.Comma = []rune(*delim)[0]
	}
	csvr.LazyQuotes = true
	csvr.TrimLeadingSpace = true

	records, err := csvr.ReadAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, "csv2json:", err)
		os.Exit(1)
	}

	if len(records) == 0 {
		fmt.Println("[]")
		return
	}

	var result interface{}

	if *array || *noHeader {
		rows := make([]interface{}, len(records))
		for i, row := range records {
			cols := make([]interface{}, len(row))
			for j, v := range row {
				cols[j] = coerceValue(v)
			}
			rows[i] = cols
		}
		result = rows
	} else {
		headers := records[0]
		rows := make([]interface{}, len(records)-1)
		for i, row := range records[1:] {
			obj := map[string]interface{}{}
			for j, h := range headers {
				if j < len(row) {
					obj[h] = coerceValue(row[j])
				} else {
					obj[h] = nil
				}
			}
			rows[i] = obj
		}
		result = rows
	}

	var b []byte
	if *pretty {
		b, err = json.MarshalIndent(result, "", "  ")
	} else {
		b, err = json.Marshal(result)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "csv2json:", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

func coerceValue(s string) interface{} {
	if s == "" {
		return nil
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return s
}
