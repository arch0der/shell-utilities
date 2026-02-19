// hashmap - in-memory key-value store with get/set/del/list via stdin commands
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: hashmap [init.json]
  Interactive key-value store. Commands:
    set <key> <value>   store a value
    get <key>           retrieve a value
    del <key>           delete a key
    has <key>           check existence (exit 0/1)
    list [prefix]       list keys
    keys                list all keys
    vals                list all values
    dump                print as JSON
    load <json>         merge JSON into store
    clear               clear all keys
    count               number of keys
    quit / exit         exit`)
	os.Exit(1)
}

func main() {
	store := map[string]string{}
	if len(os.Args) > 1 {
		data, err := os.ReadFile(os.Args[1])
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		for k, v := range m { store[k] = fmt.Sprintf("%v", v) }
	}

	sc := bufio.NewScanner(os.Stdin)
	isInteractive := false
	if fi, _ := os.Stdin.Stat(); fi.Mode()&os.ModeCharDevice != 0 { isInteractive = true }
	if isInteractive { fmt.Print("hashmap> ") }

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") { if isInteractive { fmt.Print("hashmap> ") }; continue }
		parts := strings.SplitN(line, " ", 3)
		cmd := strings.ToLower(parts[0])
		switch cmd {
		case "set":
			if len(parts) < 3 { fmt.Fprintln(os.Stderr, "set needs key and value"); break }
			store[parts[1]] = parts[2]; fmt.Println("OK")
		case "get":
			if len(parts) < 2 { fmt.Fprintln(os.Stderr, "get needs key"); break }
			if v, ok := store[parts[1]]; ok { fmt.Println(v) } else { fmt.Println("(nil)") }
		case "del":
			if len(parts) < 2 { break }
			delete(store, parts[1]); fmt.Println("OK")
		case "has":
			if len(parts) < 2 { break }
			if _, ok := store[parts[1]]; ok { fmt.Println("true") } else { fmt.Println("false") }
		case "list":
			prefix := ""; if len(parts) > 1 { prefix = parts[1] }
			keys := []string{}
			for k := range store { if strings.HasPrefix(k, prefix) { keys = append(keys, k) } }
			sort.Strings(keys); for _, k := range keys { fmt.Println(k) }
		case "keys":
			keys := make([]string, 0, len(store)); for k := range store { keys = append(keys, k) }
			sort.Strings(keys); for _, k := range keys { fmt.Println(k) }
		case "vals":
			keys := make([]string, 0, len(store)); for k := range store { keys = append(keys, k) }
			sort.Strings(keys); for _, k := range keys { fmt.Println(store[k]) }
		case "dump":
			b, _ := json.MarshalIndent(store, "", "  "); fmt.Println(string(b))
		case "load":
			if len(parts) < 2 { break }
			raw := strings.Join(parts[1:], " ")
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(raw), &m); err != nil { fmt.Fprintln(os.Stderr, err); break }
			for k, v := range m { store[k] = fmt.Sprintf("%v", v) }; fmt.Println("OK")
		case "clear":
			store = map[string]string{}; fmt.Println("OK")
		case "count":
			fmt.Println(len(store))
		case "quit", "exit", "q":
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown command %q\n", cmd)
		}
		if isInteractive { fmt.Print("hashmap> ") }
	}
}
