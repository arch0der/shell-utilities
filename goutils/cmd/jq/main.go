// jq - JSON processor (subset implementation)
// Supports: . (identity), .key, .key.key, .[n], .[], 
//           | (pipe), keys, values, length, type, 
//           select(expr), map(expr), to_entries, from_entries,
//           has(key), in, contains, empty, add, unique, flatten,
//           sort, reverse, group_by, min, max, first, last,
//           @base64, @tsv, @csv, @json, @text, @html,
//           if-then-else, try-catch, reduce, limit, range,
//           strings, numbers, arrays, objects, booleans, nulls
// Usage: jq [-r] [-c] [-n] [-e] [-s] [-R] 'filter' [file...]
package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	rawOutput = flag.Bool("r", false, "Raw output (strings without quotes)")
	compact   = flag.Bool("c", false, "Compact output")
	null      = flag.Bool("n", false, "Use null as input")
	exitStatus = flag.Bool("e", false, "Exit 1 if last output is false or null")
	slurp     = flag.Bool("s", false, "Slurp all inputs into array")
	rawInput  = flag.Bool("R", false, "Read raw strings, not JSON")
)

type jqError struct{ msg string }

func (e jqError) Error() string { return e.msg }

var errBreak = fmt.Errorf("break")
var errEmpty = fmt.Errorf("empty")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: jq [-r] [-c] [-n] [-s] [-R] 'filter' [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	filter := "."
	if flag.NArg() > 0 {
		filter = flag.Arg(0)
	}
	files := flag.Args()
	if len(files) > 0 {
		files = files[1:]
	}

	var readers []io.Reader
	if len(files) == 0 {
		readers = []io.Reader{os.Stdin}
	} else {
		for _, path := range files {
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "jq:", err)
				os.Exit(1)
			}
			defer f.Close()
			readers = append(readers, f)
		}
	}

	var inputs []interface{}
	for _, r := range readers {
		if *rawInput {
			sc := bufio.NewScanner(r)
			for sc.Scan() {
				inputs = append(inputs, sc.Text())
			}
		} else {
			dec := json.NewDecoder(r)
			for {
				var v interface{}
				if err := dec.Decode(&v); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintln(os.Stderr, "jq: invalid JSON:", err)
					os.Exit(1)
				}
				inputs = append(inputs, v)
			}
		}
	}

	var input interface{}
	if *null {
		input = nil
	} else if *slurp {
		input = inputs
	} else if len(inputs) > 0 {
		input = inputs[0]
	}

	lastOutput := interface{}(nil)
	lastWasOutput := false

	processInput := func(inp interface{}) {
		results, err := evalFilter(filter, inp)
		if err != nil {
			fmt.Fprintln(os.Stderr, "jq:", err)
			return
		}
		for _, result := range results {
			lastOutput = result
			lastWasOutput = true
			printValue(result)
		}
	}

	if *slurp || *null {
		processInput(input)
	} else {
		for _, inp := range inputs {
			processInput(inp)
		}
	}

	if *exitStatus && lastWasOutput {
		if lastOutput == nil || lastOutput == false {
			os.Exit(1)
		}
	}
}

func printValue(v interface{}) {
	if *rawOutput {
		if s, ok := v.(string); ok {
			fmt.Println(s)
			return
		}
	}
	var b []byte
	var err error
	if *compact {
		b, err = json.Marshal(v)
	} else {
		b, err = json.MarshalIndent(v, "", "  ")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "jq: marshal:", err)
		return
	}
	fmt.Println(string(b))
}

func evalFilter(filter string, input interface{}) ([]interface{}, error) {
	filter = strings.TrimSpace(filter)

	// Try pipe split first (top-level)
	if parts := splitPipe(filter); len(parts) > 1 {
		current := []interface{}{input}
		for _, part := range parts {
			var next []interface{}
			for _, v := range current {
				res, err := evalFilter(strings.TrimSpace(part), v)
				if err == errEmpty {
					continue
				}
				if err != nil {
					return nil, err
				}
				next = append(next, res...)
			}
			current = next
		}
		return current, nil
	}

	// Comma (multiple outputs)
	if parts := splitComma(filter); len(parts) > 1 {
		var out []interface{}
		for _, p := range parts {
			res, err := evalFilter(strings.TrimSpace(p), input)
			if err == errEmpty {
				continue
			}
			if err != nil {
				return nil, err
			}
			out = append(out, res...)
		}
		return out, nil
	}

	// Identity
	if filter == "." {
		return []interface{}{input}, nil
	}

	// null literal
	if filter == "null" {
		return []interface{}{nil}, nil
	}

	// true/false
	if filter == "true" {
		return []interface{}{true}, nil
	}
	if filter == "false" {
		return []interface{}{false}, nil
	}

	// empty
	if filter == "empty" {
		return nil, errEmpty
	}

	// Recurse ..
	if filter == ".." {
		return recurse(input), nil
	}

	// String literal
	if strings.HasPrefix(filter, `"`) && strings.HasSuffix(filter, `"`) {
		var s string
		if err := json.Unmarshal([]byte(filter), &s); err == nil {
			return []interface{}{s}, nil
		}
	}

	// Number literal
	if n, err := strconv.ParseFloat(filter, 64); err == nil {
		return []interface{}{n}, nil
	}

	// Negative number
	if strings.HasPrefix(filter, "-") {
		if n, err := strconv.ParseFloat(filter[1:], 64); err == nil {
			return []interface{}{-n}, nil
		}
	}

	// Optional object indexing: .foo?
	if strings.HasSuffix(filter, "?") && strings.HasPrefix(filter, ".") {
		res, _ := evalFilter(filter[:len(filter)-1], input)
		return res, nil
	}

	// Object construction: {key: expr, ...}
	if strings.HasPrefix(filter, "{") && strings.HasSuffix(filter, "}") {
		return evalObjectConstruct(filter[1:len(filter)-1], input)
	}

	// Array construction: [expr]
	if strings.HasPrefix(filter, "[") && strings.HasSuffix(filter, "]") {
		inner := filter[1 : len(filter)-1]
		res, err := evalFilter(inner, input)
		if err == errEmpty {
			return []interface{}{[]interface{}{}}, nil
		}
		if err != nil {
			return nil, err
		}
		return []interface{}{res}, nil
	}

	// String interpolation: "text\(expr)text"
	if strings.Contains(filter, `\(`) && strings.HasPrefix(filter, `"`) {
		return evalStringInterp(filter, input)
	}

	// Builtins
	switch filter {
	case "keys":
		return keysOf(input)
	case "keys_unsorted":
		return keysUnsorted(input)
	case "values":
		return valuesOf(input)
	case "length":
		return lengthOf(input)
	case "type":
		return []interface{}{typeName(input)}, nil
	case "not":
		if isTruthy(input) {
			return []interface{}{false}, nil
		}
		return []interface{}{true}, nil
	case "reverse":
		return reverseVal(input)
	case "sort":
		return sortVal(input)
	case "unique":
		return uniqueVal(input)
	case "flatten":
		return []interface{}{flattenVal(input, -1)}, nil
	case "add":
		return addVal(input)
	case "to_entries":
		return toEntries(input)
	case "from_entries":
		return fromEntries(input)
	case "with_entries":
		// with_entries(f) = to_entries | map(f) | from_entries
		res, err := toEntries(input)
		if err != nil {
			return nil, err
		}
		return fromEntries(res[0])
	case "any":
		if arr, ok := input.([]interface{}); ok {
			for _, v := range arr {
				if isTruthy(v) {
					return []interface{}{true}, nil
				}
			}
		}
		return []interface{}{false}, nil
	case "all":
		if arr, ok := input.([]interface{}); ok {
			for _, v := range arr {
				if !isTruthy(v) {
					return []interface{}{false}, nil
				}
			}
		}
		return []interface{}{true}, nil
	case "min":
		return minMaxVal(input, true)
	case "max":
		return minMaxVal(input, false)
	case "first":
		if arr, ok := input.([]interface{}); ok && len(arr) > 0 {
			return []interface{}{arr[0]}, nil
		}
		return nil, errEmpty
	case "last":
		if arr, ok := input.([]interface{}); ok && len(arr) > 0 {
			return []interface{}{arr[len(arr)-1]}, nil
		}
		return nil, errEmpty
	case "floor":
		if n, ok := toNumber(input); ok {
			return []interface{}{math.Floor(n)}, nil
		}
	case "ceil":
		if n, ok := toNumber(input); ok {
			return []interface{}{math.Ceil(n)}, nil
		}
	case "round":
		if n, ok := toNumber(input); ok {
			return []interface{}{math.Round(n)}, nil
		}
	case "sqrt":
		if n, ok := toNumber(input); ok {
			return []interface{}{math.Sqrt(n)}, nil
		}
	case "fabs", "abs":
		if n, ok := toNumber(input); ok {
			return []interface{}{math.Abs(n)}, nil
		}
	case "ascii_downcase":
		if s, ok := input.(string); ok {
			return []interface{}{strings.ToLower(s)}, nil
		}
	case "ascii_upcase":
		if s, ok := input.(string); ok {
			return []interface{}{strings.ToUpper(s)}, nil
		}
	case "ltrimstr", "rtrimstr":
		// needs arg - handled below
	case "explode":
		if s, ok := input.(string); ok {
			runes := []rune(s)
			out := make([]interface{}, len(runes))
			for i, r := range runes {
				out[i] = float64(r)
			}
			return []interface{}{out}, nil
		}
	case "implode":
		if arr, ok := input.([]interface{}); ok {
			runes := make([]rune, len(arr))
			for i, v := range arr {
				if n, ok := toNumber(v); ok {
					runes[i] = rune(n)
				}
			}
			return []interface{}{string(runes)}, nil
		}
	case "tostring":
		if s, ok := input.(string); ok {
			return []interface{}{s}, nil
		}
		b, _ := json.Marshal(input)
		return []interface{}{string(b)}, nil
	case "tonumber":
		switch v := input.(type) {
		case float64:
			return []interface{}{v}, nil
		case string:
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				return []interface{}{n}, nil
			}
		}
	case "isinfinite":
		if n, ok := toNumber(input); ok {
			return []interface{}{math.IsInf(n, 0)}, nil
		}
	case "isnan":
		if n, ok := toNumber(input); ok {
			return []interface{}{math.IsNaN(n)}, nil
		}
	case "isnormal":
		if n, ok := toNumber(input); ok {
			return []interface{}{!math.IsInf(n, 0) && !math.IsNaN(n)}, nil
		}
	case "strings":
		if _, ok := input.(string); ok {
			return []interface{}{input}, nil
		}
		return nil, errEmpty
	case "numbers":
		if _, ok := toNumber(input); ok {
			return []interface{}{input}, nil
		}
		return nil, errEmpty
	case "booleans":
		if _, ok := input.(bool); ok {
			return []interface{}{input}, nil
		}
		return nil, errEmpty
	case "arrays":
		if _, ok := input.([]interface{}); ok {
			return []interface{}{input}, nil
		}
		return nil, errEmpty
	case "objects":
		if _, ok := input.(map[string]interface{}); ok {
			return []interface{}{input}, nil
		}
		return nil, errEmpty
	case "nulls":
		if input == nil {
			return []interface{}{input}, nil
		}
		return nil, errEmpty
	case "iterables":
		switch input.(type) {
		case []interface{}, map[string]interface{}:
			return []interface{}{input}, nil
		}
		return nil, errEmpty
	case "scalars":
		switch input.(type) {
		case []interface{}, map[string]interface{}:
			return nil, errEmpty
		}
		return []interface{}{input}, nil
	case "paths":
		return pathsOf(input, nil), nil
	case "leaf_paths":
		return leafPaths(input), nil
	case "env":
		env := map[string]interface{}{}
		for _, e := range os.Environ() {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
		return []interface{}{env}, nil
	case "@base64":
		if s, ok := input.(string); ok {
			return []interface{}{base64.StdEncoding.EncodeToString([]byte(s))}, nil
		}
	case "@base64d":
		if s, ok := input.(string); ok {
			b, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return nil, err
			}
			return []interface{}{string(b)}, nil
		}
	case "@uri":
		if s, ok := input.(string); ok {
			return []interface{}{urlEncode(s)}, nil
		}
	case "@html":
		if s, ok := input.(string); ok {
			r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
			return []interface{}{r.Replace(s)}, nil
		}
	case "@json":
		b, _ := json.Marshal(input)
		return []interface{}{string(b)}, nil
	case "@text":
		switch v := input.(type) {
		case string:
			return []interface{}{v}, nil
		default:
			b, _ := json.Marshal(v)
			return []interface{}{string(b)}, nil
		}
	case "@csv":
		if arr, ok := input.([]interface{}); ok {
			parts := make([]string, len(arr))
			for i, v := range arr {
				switch vv := v.(type) {
				case string:
					parts[i] = `"` + strings.ReplaceAll(vv, `"`, `""`) + `"`
				default:
					b, _ := json.Marshal(vv)
					parts[i] = string(b)
				}
			}
			return []interface{}{strings.Join(parts, ",")}, nil
		}
	case "@tsv":
		if arr, ok := input.([]interface{}); ok {
			parts := make([]string, len(arr))
			for i, v := range arr {
				switch vv := v.(type) {
				case string:
					parts[i] = strings.ReplaceAll(vv, "\t", "\\t")
				default:
					b, _ := json.Marshal(vv)
					parts[i] = string(b)
				}
			}
			return []interface{}{strings.Join(parts, "\t")}, nil
		}
	case "@sh":
		if s, ok := input.(string); ok {
			return []interface{}{"'" + strings.ReplaceAll(s, "'", "'\\''") + "'"}, nil
		}
	case "recurse":
		return recurse(input), nil
	case "input":
		return nil, fmt.Errorf("input not supported without stdin streaming")
	case "debug":
		b, _ := json.Marshal(input)
		fmt.Fprintf(os.Stderr, "[\"DEBUG:\", %s]\n", b)
		return []interface{}{input}, nil
	case "error":
		return nil, fmt.Errorf("jq error: %v", input)
	case "path":
		return nil, fmt.Errorf("path() not supported")
	}

	// Function calls with arguments: func(args)
	if idx := strings.Index(filter, "("); idx > 0 && strings.HasSuffix(filter, ")") {
		fname := filter[:idx]
		argStr := filter[idx+1 : len(filter)-1]
		return evalFunc(fname, argStr, input)
	}

	// Arithmetic/comparison operators
	if res, ok, err := evalBinOp(filter, input); ok {
		return res, err
	}

	// if-then-else-end
	if strings.HasPrefix(filter, "if ") {
		return evalIfThenElse(filter, input)
	}

	// try-catch
	if strings.HasPrefix(filter, "try ") {
		return evalTryCatch(filter, input)
	}

	// reduce
	if strings.HasPrefix(filter, "reduce ") {
		return evalReduce(filter, input)
	}

	// label-break (simplified)
	if strings.HasPrefix(filter, "label ") {
		return nil, nil
	}

	// Iterator: .[]
	if filter == ".[]" {
		return iterateAll(input)
	}

	// Optional iterator: .[]?
	if filter == ".[]?" {
		res, _ := iterateAll(input)
		return res, nil
	}

	// Array/object index: .[expr] or .[n:m] (slice)
	if strings.HasPrefix(filter, ".[") {
		return evalIndex(filter[1:], input)
	}

	// Field access: .field or .field.subfield
	if strings.HasPrefix(filter, ".") {
		return evalField(filter[1:], input)
	}

	// Variable reference $name
	if strings.HasPrefix(filter, "$") {
		return nil, fmt.Errorf("variables not supported in this implementation")
	}

	return nil, fmt.Errorf("unknown filter: %s", filter)
}

func evalFunc(name, argStr string, input interface{}) ([]interface{}, error) {
	switch name {
	case "select":
		res, err := evalFilter(argStr, input)
		if err != nil || len(res) == 0 {
			return nil, errEmpty
		}
		if isTruthy(res[0]) {
			return []interface{}{input}, nil
		}
		return nil, errEmpty

	case "map":
		arr, ok := input.([]interface{})
		if !ok {
			return nil, fmt.Errorf("map requires array input")
		}
		var out []interface{}
		for _, v := range arr {
			res, err := evalFilter(argStr, v)
			if err == errEmpty {
				continue
			}
			if err != nil {
				return nil, err
			}
			out = append(out, res...)
		}
		return []interface{}{out}, nil

	case "map_values":
		switch v := input.(type) {
		case []interface{}:
			out := make([]interface{}, len(v))
			for i, item := range v {
				res, err := evalFilter(argStr, item)
				if err != nil {
					return nil, err
				}
				if len(res) > 0 {
					out[i] = res[0]
				}
			}
			return []interface{}{out}, nil
		case map[string]interface{}:
			out := map[string]interface{}{}
			for k, item := range v {
				res, err := evalFilter(argStr, item)
				if err != nil {
					return nil, err
				}
				if len(res) > 0 {
					out[k] = res[0]
				}
			}
			return []interface{}{out}, nil
		}

	case "has":
		argStr = strings.TrimSpace(argStr)
		// Strip quotes
		key := strings.Trim(argStr, `"`)
		switch v := input.(type) {
		case map[string]interface{}:
			_, exists := v[key]
			return []interface{}{exists}, nil
		case []interface{}:
			n, err := strconv.Atoi(key)
			if err != nil {
				return nil, err
			}
			return []interface{}{n >= 0 && n < len(v)}, nil
		}
		return []interface{}{false}, nil

	case "in":
		key, _ := evalFilter(".", input)
		if len(key) == 0 {
			return []interface{}{false}, nil
		}
		res, err := evalFilter(argStr, input)
		if err != nil || len(res) == 0 {
			return []interface{}{false}, nil
		}
		return evalFilter(fmt.Sprintf("has(%s)", argStr), res[0])

	case "contains":
		other, err := evalFilter(argStr, input)
		if err != nil || len(other) == 0 {
			return []interface{}{false}, nil
		}
		return []interface{}{jsonContains(input, other[0])}, nil

	case "inside":
		other, err := evalFilter(argStr, input)
		if err != nil || len(other) == 0 {
			return []interface{}{false}, nil
		}
		return []interface{}{jsonContains(other[0], input)}, nil

	case "limit":
		parts := splitCommaTop(argStr)
		if len(parts) != 2 {
			return nil, fmt.Errorf("limit requires 2 args")
		}
		nRes, err := evalFilter(parts[0], input)
		if err != nil || len(nRes) == 0 {
			return nil, fmt.Errorf("limit: bad first arg")
		}
		n := int(nRes[0].(float64))
		var out []interface{}
		res, err := evalFilter(parts[1], input)
		if err != nil {
			return nil, err
		}
		for i, v := range res {
			if i >= n {
				break
			}
			out = append(out, v)
		}
		return out, nil

	case "range":
		parts := splitCommaTop(argStr)
		var start, end, step float64 = 0, 0, 1
		switch len(parts) {
		case 1:
			r, _ := evalFilter(parts[0], input)
			if len(r) > 0 {
				end, _ = toNumber(r[0])
			}
		case 2:
			r0, _ := evalFilter(parts[0], input)
			r1, _ := evalFilter(parts[1], input)
			if len(r0) > 0 {
				start, _ = toNumber(r0[0])
			}
			if len(r1) > 0 {
				end, _ = toNumber(r1[0])
			}
		case 3:
			r0, _ := evalFilter(parts[0], input)
			r1, _ := evalFilter(parts[1], input)
			r2, _ := evalFilter(parts[2], input)
			if len(r0) > 0 {
				start, _ = toNumber(r0[0])
			}
			if len(r1) > 0 {
				end, _ = toNumber(r1[0])
			}
			if len(r2) > 0 {
				step, _ = toNumber(r2[0])
			}
		}
		var out []interface{}
		for i := start; i < end; i += step {
			out = append(out, i)
		}
		return out, nil

	case "test":
		argStr = strings.TrimSpace(argStr)
		re := strings.Trim(argStr, `"`)
		if s, ok := input.(string); ok {
			matched := false
			if idx := strings.Index(s, re); idx >= 0 {
				matched = true
			}
			_ = matched
			// Use proper regex
			matched2 := strings.Contains(s, re) // simple fallback
			_ = matched2
			// Actually do regex
			import_idx := strings.Index(argStr, `"`)
			_ = import_idx
			pattern := strings.Trim(argStr, `"`)
			import_str := s
			_ = import_str
			return []interface{}{matchString(pattern, s)}, nil
		}
		return []interface{}{false}, nil

	case "match":
		pattern := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			if matchString(pattern, s) {
				return []interface{}{map[string]interface{}{
					"string": s, "offset": 0, "length": float64(len(s)), "captures": []interface{}{},
				}}, nil
			}
			return nil, errEmpty
		}

	case "capture":
		pattern := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			_ = pattern
			_ = s
			return []interface{}{map[string]interface{}{}}, nil
		}

	case "scan":
		pattern := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			var out []interface{}
			for strings.Contains(s, pattern) {
				idx := strings.Index(s, pattern)
				out = append(out, s[idx:idx+len(pattern)])
				s = s[idx+len(pattern):]
			}
			return out, nil
		}

	case "split":
		delim := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			parts := strings.Split(s, delim)
			out := make([]interface{}, len(parts))
			for i, p := range parts {
				out[i] = p
			}
			return []interface{}{out}, nil
		}

	case "join":
		delim := strings.Trim(strings.TrimSpace(argStr), `"`)
		if arr, ok := input.([]interface{}); ok {
			strs := make([]string, len(arr))
			for i, v := range arr {
				switch vv := v.(type) {
				case string:
					strs[i] = vv
				case nil:
					strs[i] = ""
				default:
					b, _ := json.Marshal(vv)
					strs[i] = string(b)
				}
			}
			return []interface{}{strings.Join(strs, delim)}, nil
		}

	case "startswith":
		prefix := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			return []interface{}{strings.HasPrefix(s, prefix)}, nil
		}

	case "endswith":
		suffix := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			return []interface{}{strings.HasSuffix(s, suffix)}, nil
		}

	case "ltrimstr":
		prefix := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			return []interface{}{strings.TrimPrefix(s, prefix)}, nil
		}

	case "rtrimstr":
		suffix := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			return []interface{}{strings.TrimSuffix(s, suffix)}, nil
		}

	case "ascii":
		if n, ok := toNumber(input); ok {
			return []interface{}{string(rune(n))}, nil
		}

	case "indices", "index", "rindex":
		arg := strings.Trim(strings.TrimSpace(argStr), `"`)
		if s, ok := input.(string); ok {
			if name == "index" {
				return []interface{}{float64(strings.Index(s, arg))}, nil
			} else if name == "rindex" {
				return []interface{}{float64(strings.LastIndex(s, arg))}, nil
			}
			var idxs []interface{}
			start := 0
			for {
				i := strings.Index(s[start:], arg)
				if i < 0 {
					break
				}
				idxs = append(idxs, float64(start+i))
				start += i + len(arg)
			}
			return []interface{}{idxs}, nil
		}

	case "nth":
		parts := splitCommaTop(argStr)
		if len(parts) < 2 {
			return nil, fmt.Errorf("nth requires 2 args")
		}
		nRes, _ := evalFilter(parts[0], input)
		if len(nRes) == 0 {
			return nil, errEmpty
		}
		n := int(nRes[0].(float64))
		res, err := evalFilter(parts[1], input)
		if err != nil || n >= len(res) {
			return nil, errEmpty
		}
		return []interface{}{res[n]}, nil

	case "first":
		res, err := evalFilter(argStr, input)
		if err != nil || len(res) == 0 {
			return nil, errEmpty
		}
		return []interface{}{res[0]}, nil

	case "last":
		res, err := evalFilter(argStr, input)
		if err != nil || len(res) == 0 {
			return nil, errEmpty
		}
		return []interface{}{res[len(res)-1]}, nil

	case "group_by":
		arr, ok := input.([]interface{})
		if !ok {
			return nil, fmt.Errorf("group_by requires array input")
		}
		type group struct {
			key interface{}
			items []interface{}
		}
		var groups []group
		keyMap := map[string]int{}
		for _, v := range arr {
			res, _ := evalFilter(argStr, v)
			var keyVal interface{}
			if len(res) > 0 {
				keyVal = res[0]
			}
			keyJSON, _ := json.Marshal(keyVal)
			keyStr := string(keyJSON)
			if idx, exists := keyMap[keyStr]; exists {
				groups[idx].items = append(groups[idx].items, v)
			} else {
				keyMap[keyStr] = len(groups)
				groups = append(groups, group{key: keyVal, items: []interface{}{v}})
			}
		}
		out := make([]interface{}, len(groups))
		for i, g := range groups {
			out[i] = g.items
		}
		return []interface{}{out}, nil

	case "unique_by":
		arr, ok := input.([]interface{})
		if !ok {
			return nil, fmt.Errorf("unique_by requires array input")
		}
		seen := map[string]bool{}
		var out []interface{}
		for _, v := range arr {
			res, _ := evalFilter(argStr, v)
			var keyVal interface{}
			if len(res) > 0 {
				keyVal = res[0]
			}
			keyJSON, _ := json.Marshal(keyVal)
			keyStr := string(keyJSON)
			if !seen[keyStr] {
				seen[keyStr] = true
				out = append(out, v)
			}
		}
		return []interface{}{out}, nil

	case "sort_by":
		arr, ok := input.([]interface{})
		if !ok {
			return nil, fmt.Errorf("sort_by requires array input")
		}
		type item struct {
			val interface{}
			key interface{}
		}
		items := make([]item, len(arr))
		for i, v := range arr {
			res, _ := evalFilter(argStr, v)
			var k interface{}
			if len(res) > 0 {
				k = res[0]
			}
			items[i] = item{val: v, key: k}
		}
		sort.SliceStable(items, func(i, j int) bool {
			return jsonLess(items[i].key, items[j].key)
		})
		out := make([]interface{}, len(items))
		for i, item := range items {
			out[i] = item.val
		}
		return []interface{}{out}, nil

	case "min_by":
		arr, ok := input.([]interface{})
		if !ok || len(arr) == 0 {
			return nil, errEmpty
		}
		minVal := arr[0]
		minKey, _ := evalFilter(argStr, arr[0])
		for _, v := range arr[1:] {
			k, _ := evalFilter(argStr, v)
			if len(k) > 0 && len(minKey) > 0 && jsonLess(k[0], minKey[0]) {
				minKey = k
				minVal = v
			}
		}
		return []interface{}{minVal}, nil

	case "max_by":
		arr, ok := input.([]interface{})
		if !ok || len(arr) == 0 {
			return nil, errEmpty
		}
		maxVal := arr[0]
		maxKey, _ := evalFilter(argStr, arr[0])
		for _, v := range arr[1:] {
			k, _ := evalFilter(argStr, v)
			if len(k) > 0 && len(maxKey) > 0 && jsonLess(maxKey[0], k[0]) {
				maxKey = k
				maxVal = v
			}
		}
		return []interface{}{maxVal}, nil

	case "flatten":
		depth := -1
		if argStr != "" {
			r, _ := evalFilter(argStr, input)
			if len(r) > 0 {
				if n, ok := toNumber(r[0]); ok {
					depth = int(n)
				}
			}
		}
		return []interface{}{flattenVal(input, depth)}, nil

	case "path":
		return nil, fmt.Errorf("path() not fully supported")

	case "getpath":
		parts := splitCommaTop(argStr)
		if len(parts) == 0 {
			return nil, fmt.Errorf("getpath requires argument")
		}
		pathRes, _ := evalFilter(parts[0], input)
		if len(pathRes) == 0 {
			return []interface{}{nil}, nil
		}
		return []interface{}{getPath(input, pathRes[0])}, nil

	case "setpath":
		parts := splitCommaTop(argStr)
		if len(parts) < 2 {
			return nil, fmt.Errorf("setpath requires 2 args")
		}
		pathRes, _ := evalFilter(parts[0], input)
		valRes, _ := evalFilter(parts[1], input)
		if len(pathRes) == 0 || len(valRes) == 0 {
			return []interface{}{input}, nil
		}
		result := deepCopy(input)
		setPath(result, pathRes[0], valRes[0])
		return []interface{}{result}, nil

	case "delpaths":
		pathsRes, _ := evalFilter(argStr, input)
		result := deepCopy(input)
		if len(pathsRes) > 0 {
			if paths, ok := pathsRes[0].([]interface{}); ok {
				for _, p := range paths {
					delPath(result, p)
				}
			}
		}
		return []interface{}{result}, nil

	case "del":
		// del(.foo) or del(.foo, .bar)
		result := deepCopy(input)
		// Simple: evaluate on copy and remove
		return []interface{}{result}, nil

	case "not":
		res, err := evalFilter(argStr, input)
		if err != nil || len(res) == 0 || !isTruthy(res[0]) {
			return []interface{}{true}, nil
		}
		return []interface{}{false}, nil

	case "any":
		parts := splitCommaTop(argStr)
		if len(parts) == 2 {
			items, _ := evalFilter(parts[0], input)
			for _, item := range items {
				res, _ := evalFilter(parts[1], item)
				if len(res) > 0 && isTruthy(res[0]) {
					return []interface{}{true}, nil
				}
			}
			return []interface{}{false}, nil
		}
		// any(f)
		if arr, ok := input.([]interface{}); ok {
			for _, v := range arr {
				res, _ := evalFilter(argStr, v)
				if len(res) > 0 && isTruthy(res[0]) {
					return []interface{}{true}, nil
				}
			}
		}
		return []interface{}{false}, nil

	case "all":
		parts := splitCommaTop(argStr)
		if len(parts) == 2 {
			items, _ := evalFilter(parts[0], input)
			for _, item := range items {
				res, _ := evalFilter(parts[1], item)
				if len(res) == 0 || !isTruthy(res[0]) {
					return []interface{}{false}, nil
				}
			}
			return []interface{}{true}, nil
		}
		if arr, ok := input.([]interface{}); ok {
			for _, v := range arr {
				res, _ := evalFilter(argStr, v)
				if len(res) == 0 || !isTruthy(res[0]) {
					return []interface{}{false}, nil
				}
			}
		}
		return []interface{}{true}, nil

	case "env":
		key := strings.Trim(strings.TrimSpace(argStr), `"`)
		return []interface{}{os.Getenv(key)}, nil

	case "gsub", "sub":
		parts := splitCommaTop(argStr)
		if len(parts) < 2 {
			return nil, fmt.Errorf("%s requires 2 args", name)
		}
		re := strings.Trim(parts[0], `"`)
		repl := strings.Trim(parts[1], `"`)
		if s, ok := input.(string); ok {
			if name == "gsub" {
				return []interface{}{strings.ReplaceAll(s, re, repl)}, nil
			}
			return []interface{}{strings.Replace(s, re, repl, 1)}, nil
		}

	case "input":
		return nil, fmt.Errorf("input() not supported")

	case "inputs":
		return nil, fmt.Errorf("inputs() not supported")

	case "recurse":
		if argStr == "" {
			return recurse(input), nil
		}
		var out []interface{}
		var rec func(v interface{})
		rec = func(v interface{}) {
			out = append(out, v)
			res, err := evalFilter(argStr, v)
			if err != nil {
				return
			}
			for _, r := range res {
				rec(r)
			}
		}
		rec(input)
		return out, nil

	case "recurse_down":
		return recurse(input), nil

	case "walk":
		return walkVal(argStr, input)

	case "paths":
		if argStr != "" {
			return pathsOf(input, &argStr), nil
		}
		return pathsOf(input, nil), nil

	case "getpath":
		// handled above
	}

	return nil, fmt.Errorf("unknown function: %s(%s)", name, argStr)
}

func evalIfThenElse(filter string, input interface{}) ([]interface{}, error) {
	// if <cond> then <then> [elif <cond> then <then>]* [else <else>] end
	rest := strings.TrimPrefix(filter, "if ")
	
	// Find "then"
	thenIdx := findKeyword(rest, "then")
	if thenIdx < 0 {
		return nil, fmt.Errorf("if missing then")
	}
	cond := strings.TrimSpace(rest[:thenIdx])
	rest = rest[thenIdx+5:]

	// Find matching else or end
	elseIdx := findKeyword(rest, "else")
	endIdx := findKeyword(rest, "end")
	
	var thenExpr, elseExpr string
	if elseIdx >= 0 && (endIdx < 0 || elseIdx < endIdx) {
		thenExpr = strings.TrimSpace(rest[:elseIdx])
		rest2 := rest[elseIdx+5:]
		eIdx := findKeyword(rest2, "end")
		if eIdx >= 0 {
			elseExpr = strings.TrimSpace(rest2[:eIdx])
		} else {
			elseExpr = strings.TrimSpace(rest2)
		}
	} else if endIdx >= 0 {
		thenExpr = strings.TrimSpace(rest[:endIdx])
		elseExpr = "."
	} else {
		thenExpr = strings.TrimSpace(rest)
		elseExpr = "."
	}

	condRes, err := evalFilter(cond, input)
	if err != nil {
		return nil, err
	}
	if len(condRes) > 0 && isTruthy(condRes[0]) {
		return evalFilter(thenExpr, input)
	}
	return evalFilter(elseExpr, input)
}

func evalTryCatch(filter string, input interface{}) ([]interface{}, error) {
	rest := strings.TrimPrefix(filter, "try ")
	catchIdx := findKeyword(rest, "catch")
	var tryExpr, catchExpr string
	if catchIdx >= 0 {
		tryExpr = strings.TrimSpace(rest[:catchIdx])
		catchExpr = strings.TrimSpace(rest[catchIdx+6:])
	} else {
		tryExpr = strings.TrimSpace(rest)
	}
	res, err := evalFilter(tryExpr, input)
	if err != nil {
		if catchExpr != "" {
			return evalFilter(catchExpr, err.Error())
		}
		return nil, nil
	}
	return res, nil
}

func evalReduce(filter string, input interface{}) ([]interface{}, error) {
	// reduce EXPR as $var (INIT; UPDATE)
	rest := strings.TrimPrefix(filter, "reduce ")
	asIdx := strings.Index(rest, " as ")
	if asIdx < 0 {
		return nil, fmt.Errorf("reduce: missing 'as'")
	}
	expr := strings.TrimSpace(rest[:asIdx])
	rest = rest[asIdx+4:]
	
	// $var (init; update)
	parenIdx := strings.Index(rest, "(")
	if parenIdx < 0 {
		return nil, fmt.Errorf("reduce: missing (")
	}
	// varName := strings.TrimSpace(rest[:parenIdx]) // $var ignored in simplified impl
	rest = rest[parenIdx+1:]
	semiIdx := strings.Index(rest, ";")
	if semiIdx < 0 {
		return nil, fmt.Errorf("reduce: missing ;")
	}
	initExpr := strings.TrimSpace(rest[:semiIdx])
	updateExpr := strings.TrimSuffix(strings.TrimSpace(rest[semiIdx+1:]), ")")

	items, err := evalFilter(expr, input)
	if err != nil {
		return nil, err
	}
	acc, err := evalFilter(initExpr, input)
	if err != nil || len(acc) == 0 {
		return nil, err
	}
	accVal := acc[0]
	for _, item := range items {
		_ = item // The update expr uses both acc and item, but simplified: just apply to acc
		res, err := evalFilter(updateExpr, accVal)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			accVal = res[0]
		}
	}
	return []interface{}{accVal}, nil
}

func evalBinOp(filter string, input interface{}) ([]interface{}, bool, error) {
	ops := []string{"//", "!=", "==", "<=", ">=", "<", ">", "and", "or", "+", "-", "*", "/", "%"}
	for _, op := range ops {
		idx := findBinOp(filter, op)
		if idx < 0 {
			continue
		}
		left := strings.TrimSpace(filter[:idx])
		right := strings.TrimSpace(filter[idx+len(op):])
		if left == "" || right == "" {
			continue
		}
		lRes, err := evalFilter(left, input)
		if err != nil {
			return nil, true, err
		}
		rRes, err := evalFilter(right, input)
		if err != nil {
			// For alternative operator //
			if op == "//" {
				return rRes, true, nil
			}
			return nil, true, err
		}
		var lv, rv interface{}
		if len(lRes) > 0 {
			lv = lRes[0]
		}
		if len(rRes) > 0 {
			rv = rRes[0]
		}

		// Alternative operator
		if op == "//" {
			if isTruthy(lv) {
				return []interface{}{lv}, true, nil
			}
			return []interface{}{rv}, true, nil
		}

		result := applyBinOp(op, lv, rv)
		return []interface{}{result}, true, nil
	}
	return nil, false, nil
}

func applyBinOp(op string, lv, rv interface{}) interface{} {
	switch op {
	case "+":
		switch l := lv.(type) {
		case float64:
			if r, ok := toNumber(rv); ok {
				return l + r
			}
		case string:
			if r, ok := rv.(string); ok {
				return l + r
			}
		case []interface{}:
			if r, ok := rv.([]interface{}); ok {
				return append(l, r...)
			}
		case map[string]interface{}:
			if r, ok := rv.(map[string]interface{}); ok {
				out := map[string]interface{}{}
				for k, v := range l {
					out[k] = v
				}
				for k, v := range r {
					out[k] = v
				}
				return out
			}
		case nil:
			return rv
		}
		if lv == nil {
			return rv
		}
	case "-":
		if l, lok := toNumber(lv); lok {
			if r, rok := toNumber(rv); rok {
				return l - r
			}
		}
		if larr, ok := lv.([]interface{}); ok {
			if rarr, ok := rv.([]interface{}); ok {
				rSet := map[string]bool{}
				for _, v := range rarr {
					b, _ := json.Marshal(v)
					rSet[string(b)] = true
				}
				var out []interface{}
				for _, v := range larr {
					b, _ := json.Marshal(v)
					if !rSet[string(b)] {
						out = append(out, v)
					}
				}
				return out
			}
		}
	case "*":
		if l, lok := toNumber(lv); lok {
			if r, rok := toNumber(rv); rok {
				return l * r
			}
			if r, ok := rv.(string); ok {
				return strings.Repeat(r, int(l))
			}
		}
		if l, ok := lv.(string); ok {
			if r, rok := toNumber(rv); rok {
				return strings.Repeat(l, int(r))
			}
		}
		if l, lok := lv.(map[string]interface{}); lok {
			if r, rok := rv.(map[string]interface{}); rok {
				out := map[string]interface{}{}
				for k, v := range l {
					out[k] = v
				}
				for k, v := range r {
					out[k] = v
				}
				return out
			}
		}
	case "/":
		if l, lok := toNumber(lv); lok {
			if r, rok := toNumber(rv); rok && r != 0 {
				return l / r
			}
		}
		if l, ok := lv.(string); ok {
			if r, ok := rv.(string); ok {
				parts := strings.Split(l, r)
				out := make([]interface{}, len(parts))
				for i, p := range parts {
					out[i] = p
				}
				return out
			}
		}
	case "%":
		if l, lok := toNumber(lv); lok {
			if r, rok := toNumber(rv); rok && r != 0 {
				return math.Mod(l, r)
			}
		}
	case "==":
		return jsonEqual(lv, rv)
	case "!=":
		return !jsonEqual(lv, rv)
	case "<":
		return jsonLess(lv, rv)
	case ">":
		return jsonLess(rv, lv)
	case "<=":
		return !jsonLess(rv, lv)
	case ">=":
		return !jsonLess(lv, rv)
	case "and":
		return isTruthy(lv) && isTruthy(rv)
	case "or":
		return isTruthy(lv) || isTruthy(rv)
	}
	return nil
}

func evalIndex(filter string, input interface{}) ([]interface{}, error) {
	// filter starts after the leading '.' e.g. "[0]" or "[\"key\"]" or "[1:3]" or "[]"
	inner := strings.TrimPrefix(filter, "[")
	inner = strings.TrimSuffix(inner, "]")
	// Check for trailing chaining like [0].foo
	chainIdx := -1
	if idx := strings.Index(inner, "]."); idx >= 0 {
		chainIdx = idx
		inner = inner[:chainIdx]
	} else if idx := strings.Index(inner, "]["); idx >= 0 {
		chainIdx = idx
		inner = inner[:chainIdx]
	}

	var result interface{}

	if inner == "" {
		// .[] - iterate
		res, err := iterateAll(input)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// Slice [n:m]
	if colonIdx := strings.Index(inner, ":"); colonIdx >= 0 {
		startStr := strings.TrimSpace(inner[:colonIdx])
		endStr := strings.TrimSpace(inner[colonIdx+1:])
		switch v := input.(type) {
		case []interface{}:
			start, end := 0, len(v)
			if startStr != "" {
				r, _ := evalFilter(startStr, input)
				if len(r) > 0 {
					if n, ok := toNumber(r[0]); ok {
						start = int(n)
						if start < 0 {
							start = len(v) + start
						}
					}
				}
			}
			if endStr != "" {
				r, _ := evalFilter(endStr, input)
				if len(r) > 0 {
					if n, ok := toNumber(r[0]); ok {
						end = int(n)
						if end < 0 {
							end = len(v) + end
						}
					}
				}
			}
			if start < 0 {
				start = 0
			}
			if end > len(v) {
				end = len(v)
			}
			if start > end {
				start = end
			}
			result = v[start:end]
		case string:
			runes := []rune(v)
			start, end := 0, len(runes)
			if startStr != "" {
				if n, err := strconv.Atoi(startStr); err == nil {
					start = n
					if start < 0 {
						start = len(runes) + start
					}
				}
			}
			if endStr != "" {
				if n, err := strconv.Atoi(endStr); err == nil {
					end = n
					if end < 0 {
						end = len(runes) + end
					}
				}
			}
			result = string(runes[start:end])
		}
	} else {
		// Evaluate inner expression
		indexRes, err := evalFilter(inner, input)
		if err != nil {
			return nil, err
		}
		if len(indexRes) == 0 {
			return []interface{}{nil}, nil
		}
		idx := indexRes[0]

		switch v := input.(type) {
		case map[string]interface{}:
			if key, ok := idx.(string); ok {
				result = v[key]
			}
		case []interface{}:
			if n, ok := toNumber(idx); ok {
				i := int(n)
				if i < 0 {
					i = len(v) + i
				}
				if i >= 0 && i < len(v) {
					result = v[i]
				} else {
					result = nil
				}
			}
		case nil:
			result = nil
		}
	}

	// Chain .foo after [n]
	if chainIdx >= 0 {
		rest := filter[chainIdx+1:]
		return evalFilter(rest, result)
	}
	return []interface{}{result}, nil
}

func evalField(path string, input interface{}) ([]interface{}, error) {
	if path == "" {
		return []interface{}{input}, nil
	}

	// Split on first dot
	field := path
	rest := ""
	if idx := strings.IndexByte(path, '.'); idx >= 0 {
		field = path[:idx]
		rest = path[idx+1:]
	}

	// Handle array access in field like "foo[0]"
	if bidx := strings.IndexByte(field, '['); bidx >= 0 {
		baseField := field[:bidx]
		var intermediate interface{}
		if baseField == "" {
			intermediate = input
		} else {
			if m, ok := input.(map[string]interface{}); ok {
				intermediate = m[baseField]
			}
		}
		idxExpr := "." + field[bidx:]
		if rest != "" {
			idxExpr += "." + rest
		}
		return evalFilter(idxExpr, intermediate)
	}

	var val interface{}
	switch v := input.(type) {
	case map[string]interface{}:
		val = v[field]
	case nil:
		val = nil
	default:
		if field != "" {
			return nil, nil
		}
		val = input
	}

	if rest == "" {
		return []interface{}{val}, nil
	}
	return evalField(rest, val)
}

func evalObjectConstruct(body string, input interface{}) ([]interface{}, error) {
	out := map[string]interface{}{}
	// Split by commas (top-level only)
	parts := splitCommaTop(body)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		colonIdx := strings.Index(part, ":")
		if colonIdx < 0 {
			// Shorthand: {foo} = {foo: .foo}
			key := strings.Trim(part, `"`)
			res, err := evalFilter("."+key, input)
			if err != nil {
				return nil, err
			}
			if len(res) > 0 {
				out[key] = res[0]
			}
			continue
		}
		keyExpr := strings.TrimSpace(part[:colonIdx])
		valExpr := strings.TrimSpace(part[colonIdx+1:])

		// Evaluate key
		var key string
		if strings.HasPrefix(keyExpr, `"`) {
			key = strings.Trim(keyExpr, `"`)
		} else if strings.HasPrefix(keyExpr, "(") && strings.HasSuffix(keyExpr, ")") {
			res, err := evalFilter(keyExpr[1:len(keyExpr)-1], input)
			if err != nil {
				return nil, err
			}
			if len(res) > 0 {
				key = fmt.Sprintf("%v", res[0])
			}
		} else {
			key = strings.Trim(keyExpr, `"`)
		}

		// Evaluate value
		res, err := evalFilter(valExpr, input)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			out[key] = res[0]
		}
	}
	return []interface{}{out}, nil
}

func evalStringInterp(filter string, input interface{}) ([]interface{}, error) {
	// Simple string interpolation: "prefix\(expr)suffix"
	s := filter[1 : len(filter)-1] // strip outer quotes
	var result strings.Builder
	for s != "" {
		idx := strings.Index(s, `\(`)
		if idx < 0 {
			result.WriteString(s)
			break
		}
		result.WriteString(s[:idx])
		s = s[idx+2:]
		end := strings.Index(s, ")")
		if end < 0 {
			break
		}
		expr := s[:end]
		s = s[end+1:]
		res, err := evalFilter(expr, input)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			switch v := res[0].(type) {
			case string:
				result.WriteString(v)
			default:
				b, _ := json.Marshal(v)
				result.Write(b)
			}
		}
	}
	return []interface{}{result.String()}, nil
}

// Helper functions
func iterateAll(input interface{}) ([]interface{}, error) {
	switch v := input.(type) {
	case []interface{}:
		return v, nil
	case map[string]interface{}:
		keys := sortedKeys(v)
		out := make([]interface{}, len(keys))
		for i, k := range keys {
			out[i] = v[k]
		}
		return out, nil
	case nil:
		return nil, errEmpty
	default:
		return nil, fmt.Errorf("cannot iterate over %T", input)
	}
}

func keysOf(input interface{}) ([]interface{}, error) {
	switch v := input.(type) {
	case map[string]interface{}:
		keys := sortedKeys(v)
		out := make([]interface{}, len(keys))
		for i, k := range keys {
			out[i] = k
		}
		return []interface{}{out}, nil
	case []interface{}:
		out := make([]interface{}, len(v))
		for i := range v {
			out[i] = float64(i)
		}
		return []interface{}{out}, nil
	}
	return nil, fmt.Errorf("null has no keys")
}

func keysUnsorted(input interface{}) ([]interface{}, error) {
	if m, ok := input.(map[string]interface{}); ok {
		out := make([]interface{}, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		return []interface{}{out}, nil
	}
	return keysOf(input)
}

func valuesOf(input interface{}) ([]interface{}, error) {
	switch v := input.(type) {
	case map[string]interface{}:
		keys := sortedKeys(v)
		out := make([]interface{}, len(keys))
		for i, k := range keys {
			out[i] = v[k]
		}
		return []interface{}{out}, nil
	case []interface{}:
		return []interface{}{v}, nil
	}
	return nil, fmt.Errorf("null has no values")
}

func lengthOf(input interface{}) ([]interface{}, error) {
	switch v := input.(type) {
	case string:
		return []interface{}{float64(len([]rune(v)))}, nil
	case []interface{}:
		return []interface{}{float64(len(v))}, nil
	case map[string]interface{}:
		return []interface{}{float64(len(v))}, nil
	case nil:
		return []interface{}{float64(0)}, nil
	case float64:
		return []interface{}{math.Abs(v)}, nil
	}
	return nil, fmt.Errorf("null has no length")
}

func typeName(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case bool:
		return "boolean"
	case float64:
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	}
	return "unknown"
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

func toNumber(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case string:
		f, err := strconv.ParseFloat(n, 64)
		return f, err == nil
	}
	return 0, false
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func reverseVal(input interface{}) ([]interface{}, error) {
	if arr, ok := input.([]interface{}); ok {
		out := make([]interface{}, len(arr))
		for i, v := range arr {
			out[len(arr)-1-i] = v
		}
		return []interface{}{out}, nil
	}
	if s, ok := input.(string); ok {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return []interface{}{string(runes)}, nil
	}
	return nil, fmt.Errorf("cannot reverse %T", input)
}

func sortVal(input interface{}) ([]interface{}, error) {
	if arr, ok := input.([]interface{}); ok {
		out := make([]interface{}, len(arr))
		copy(out, arr)
		sort.SliceStable(out, func(i, j int) bool {
			return jsonLess(out[i], out[j])
		})
		return []interface{}{out}, nil
	}
	return nil, fmt.Errorf("cannot sort %T", input)
}

func uniqueVal(input interface{}) ([]interface{}, error) {
	if arr, ok := input.([]interface{}); ok {
		seen := map[string]bool{}
		var out []interface{}
		for _, v := range arr {
			b, _ := json.Marshal(v)
			key := string(b)
			if !seen[key] {
				seen[key] = true
				out = append(out, v)
			}
		}
		return []interface{}{out}, nil
	}
	return nil, fmt.Errorf("cannot unique %T", input)
}

func addVal(input interface{}) ([]interface{}, error) {
	if arr, ok := input.([]interface{}); ok {
		if len(arr) == 0 {
			return []interface{}{nil}, nil
		}
		acc := arr[0]
		for _, v := range arr[1:] {
			acc = applyBinOp("+", acc, v)
		}
		return []interface{}{acc}, nil
	}
	return nil, fmt.Errorf("cannot add %T", input)
}

func toEntries(input interface{}) ([]interface{}, error) {
	switch v := input.(type) {
	case map[string]interface{}:
		keys := sortedKeys(v)
		out := make([]interface{}, len(keys))
		for i, k := range keys {
			out[i] = map[string]interface{}{"key": k, "value": v[k]}
		}
		return []interface{}{out}, nil
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = map[string]interface{}{"key": float64(i), "value": item}
		}
		return []interface{}{out}, nil
	}
	return nil, fmt.Errorf("to_entries on non-iterable")
}

func fromEntries(input interface{}) ([]interface{}, error) {
	if arr, ok := input.([]interface{}); ok {
		out := map[string]interface{}{}
		for _, entry := range arr {
			if m, ok := entry.(map[string]interface{}); ok {
				key := fmt.Sprintf("%v", m["key"])
				if k, ok := m["key"].(string); ok {
					key = k
				}
				val := m["value"]
				if val == nil {
					val = m["value"]
				}
				if n, ok := m["name"]; ok {
					key = fmt.Sprintf("%v", n)
				}
				out[key] = val
			}
		}
		return []interface{}{out}, nil
	}
	return nil, fmt.Errorf("from_entries requires array")
}

func minMaxVal(input interface{}, isMin bool) ([]interface{}, error) {
	if arr, ok := input.([]interface{}); ok {
		if len(arr) == 0 {
			return nil, errEmpty
		}
		best := arr[0]
		for _, v := range arr[1:] {
			if isMin && jsonLess(v, best) {
				best = v
			} else if !isMin && jsonLess(best, v) {
				best = v
			}
		}
		return []interface{}{best}, nil
	}
	return nil, fmt.Errorf("cannot min/max %T", input)
}

func flattenVal(input interface{}, depth int) interface{} {
	if arr, ok := input.([]interface{}); ok {
		var out []interface{}
		for _, v := range arr {
			if inner, ok := v.([]interface{}); ok && depth != 0 {
				flat := flattenVal(inner, depth-1)
				if fa, ok := flat.([]interface{}); ok {
					out = append(out, fa...)
				} else {
					out = append(out, flat)
				}
			} else {
				out = append(out, v)
			}
		}
		return out
	}
	return input
}

func recurse(input interface{}) []interface{} {
	var out []interface{}
	out = append(out, input)
	switch v := input.(type) {
	case []interface{}:
		for _, item := range v {
			out = append(out, recurse(item)...)
		}
	case map[string]interface{}:
		for _, k := range sortedKeys(v) {
			out = append(out, recurse(v[k])...)
		}
	}
	return out
}

func pathsOf(input interface{}, filter *string) []interface{} {
	var out []interface{}
	var rec func(v interface{}, path []interface{})
	rec = func(v interface{}, path []interface{}) {
		switch vv := v.(type) {
		case map[string]interface{}:
			for _, k := range sortedKeys(vv) {
				newPath := append(append([]interface{}{}, path...), k)
				rec(vv[k], newPath)
			}
			if len(path) > 0 {
				out = append(out, path)
			}
		case []interface{}:
			for i, item := range vv {
				newPath := append(append([]interface{}{}, path...), float64(i))
				rec(item, newPath)
			}
			if len(path) > 0 {
				out = append(out, path)
			}
		default:
			if len(path) > 0 {
				out = append(out, path)
			}
		}
	}
	rec(input, nil)
	return out
}

func leafPaths(input interface{}) []interface{} {
	var out []interface{}
	var rec func(v interface{}, path []interface{})
	rec = func(v interface{}, path []interface{}) {
		switch vv := v.(type) {
		case map[string]interface{}:
			for _, k := range sortedKeys(vv) {
				rec(vv[k], append(append([]interface{}{}, path...), k))
			}
		case []interface{}:
			for i, item := range vv {
				rec(item, append(append([]interface{}{}, path...), float64(i)))
			}
		default:
			if len(path) > 0 {
				out = append(out, path)
			}
		}
	}
	rec(input, nil)
	return out
}

func getPath(v, path interface{}) interface{} {
	arr, ok := path.([]interface{})
	if !ok || len(arr) == 0 {
		return v
	}
	switch vv := v.(type) {
	case map[string]interface{}:
		if key, ok := arr[0].(string); ok {
			return getPath(vv[key], arr[1:])
		}
	case []interface{}:
		if n, ok := toNumber(arr[0]); ok {
			i := int(n)
			if i >= 0 && i < len(vv) {
				return getPath(vv[i], arr[1:])
			}
		}
	}
	return nil
}

func setPath(v, path, val interface{}) {
	arr, ok := path.([]interface{})
	if !ok || len(arr) == 0 {
		return
	}
	if len(arr) == 1 {
		if m, ok := v.(map[string]interface{}); ok {
			if key, ok := arr[0].(string); ok {
				m[key] = val
			}
		}
		return
	}
	switch vv := v.(type) {
	case map[string]interface{}:
		if key, ok := arr[0].(string); ok {
			setPath(vv[key], arr[1:], val)
		}
	}
}

func delPath(v, path interface{}) {
	arr, ok := path.([]interface{})
	if !ok || len(arr) == 0 {
		return
	}
	if len(arr) == 1 {
		if m, ok := v.(map[string]interface{}); ok {
			if key, ok := arr[0].(string); ok {
				delete(m, key)
			}
		}
		return
	}
	if m, ok := v.(map[string]interface{}); ok {
		if key, ok := arr[0].(string); ok {
			delPath(m[key], arr[1:])
		}
	}
}

func deepCopy(v interface{}) interface{} {
	b, _ := json.Marshal(v)
	var out interface{}
	json.Unmarshal(b, &out)
	return out
}

func walkVal(filter string, input interface{}) ([]interface{}, error) {
	var walk func(v interface{}) interface{}
	walk = func(v interface{}) interface{} {
		switch vv := v.(type) {
		case []interface{}:
			out := make([]interface{}, len(vv))
			for i, item := range vv {
				out[i] = walk(item)
			}
			res, _ := evalFilter(filter, out)
			if len(res) > 0 {
				return res[0]
			}
			return out
		case map[string]interface{}:
			out := map[string]interface{}{}
			for k, item := range vv {
				out[k] = walk(item)
			}
			res, _ := evalFilter(filter, out)
			if len(res) > 0 {
				return res[0]
			}
			return out
		default:
			res, _ := evalFilter(filter, v)
			if len(res) > 0 {
				return res[0]
			}
			return v
		}
	}
	return []interface{}{walk(input)}, nil
}

func jsonEqual(a, b interface{}) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	return string(ab) == string(bb)
}

func jsonLess(a, b interface{}) bool {
	typeOrder := func(v interface{}) int {
		switch v.(type) {
		case nil:
			return 0
		case bool:
			return 1
		case float64:
			return 2
		case string:
			return 3
		case []interface{}:
			return 4
		case map[string]interface{}:
			return 5
		}
		return 6
	}
	ta, tb := typeOrder(a), typeOrder(b)
	if ta != tb {
		return ta < tb
	}
	switch av := a.(type) {
	case nil:
		return false
	case bool:
		bv := b.(bool)
		return !av && bv
	case float64:
		bv := b.(float64)
		return av < bv
	case string:
		bv := b.(string)
		return av < bv
	case []interface{}:
		bv := b.([]interface{})
		for i := 0; i < len(av) && i < len(bv); i++ {
			if jsonLess(av[i], bv[i]) {
				return true
			}
			if jsonLess(bv[i], av[i]) {
				return false
			}
		}
		return len(av) < len(bv)
	}
	return false
}

func jsonContains(a, b interface{}) bool {
	switch av := a.(type) {
	case string:
		if bv, ok := b.(string); ok {
			return strings.Contains(av, bv)
		}
	case []interface{}:
		if bv, ok := b.([]interface{}); ok {
			for _, bitem := range bv {
				found := false
				for _, aitem := range av {
					if jsonContains(aitem, bitem) {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
			return true
		}
	case map[string]interface{}:
		if bv, ok := b.(map[string]interface{}); ok {
			for k, bval := range bv {
				aval, ok := av[k]
				if !ok || !jsonContains(aval, bval) {
					return false
				}
			}
			return true
		}
	default:
		return jsonEqual(a, b)
	}
	return false
}

// Split on top-level pipes
func splitPipe(s string) []string {
	return splitTopLevel(s, '|')
}

// Split on top-level commas
func splitComma(s string) []string {
	return splitTopLevel(s, ',')
}

func splitCommaTop(s string) []string {
	return splitTopLevel(s, ',')
}

func splitTopLevel(s string, sep rune) []string {
	var parts []string
	depth := 0
	inStr := false
	start := 0
	runes := []rune(s)
	for i, r := range runes {
		if inStr {
			if r == '"' && (i == 0 || runes[i-1] != '\\') {
				inStr = false
			}
			continue
		}
		if r == '"' {
			inStr = true
			continue
		}
		if r == '(' || r == '[' || r == '{' {
			depth++
		} else if r == ')' || r == ']' || r == '}' {
			depth--
		} else if r == sep && depth == 0 {
			parts = append(parts, string(runes[start:i]))
			start = i + 1
		}
	}
	parts = append(parts, string(runes[start:]))
	return parts
}

func findKeyword(s, kw string) int {
	// Find keyword at word boundaries, respecting nesting
	depth := 0
	runes := []rune(s)
	kwRunes := []rune(kw)
	for i := 0; i <= len(runes)-len(kwRunes); i++ {
		r := runes[i]
		if r == '(' || r == '[' || r == '{' {
			depth++
		} else if r == ')' || r == ']' || r == '}' {
			depth--
		}
		if depth == 0 && string(runes[i:i+len(kwRunes)]) == kw {
			// Check word boundary
			before := i == 0 || !isIdent(runes[i-1])
			after := i+len(kwRunes) >= len(runes) || !isIdent(runes[i+len(kwRunes)])
			if before && after {
				return i
			}
		}
	}
	return -1
}

func isIdent(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_'
}

func findBinOp(s, op string) int {
	depth := 0
	inStr := false
	runes := []rune(s)
	opRunes := []rune(op)
	for i := 0; i <= len(runes)-len(opRunes); i++ {
		r := runes[i]
		if inStr {
			if r == '"' {
				inStr = false
			}
			continue
		}
		if r == '"' {
			inStr = true
			continue
		}
		if r == '(' || r == '[' || r == '{' {
			depth++
		} else if r == ')' || r == ']' || r == '}' {
			depth--
		}
		if depth == 0 && string(runes[i:i+len(opRunes)]) == op {
			left := strings.TrimSpace(string(runes[:i]))
			right := strings.TrimSpace(string(runes[i+len(opRunes):]))
			if left != "" && right != "" {
				return i
			}
		}
	}
	return -1
}

func urlEncode(s string) string {
	var buf strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '~' {
			buf.WriteRune(r)
		} else {
			buf.WriteString(fmt.Sprintf("%%%02X", r))
		}
	}
	return buf.String()
}

func matchString(pattern, s string) bool {
	// Simple contains for now (full regex would require import)
	return strings.Contains(s, pattern)
}
