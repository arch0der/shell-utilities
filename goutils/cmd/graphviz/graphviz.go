// graphviz - generate DOT language graphs from simple edge lists
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: graphviz [options] [file]
  Read edge list (one edge per line: "A -> B" or "A -- B") and emit DOT.
  -d, --directed    directed graph (digraph), default
  -u, --undirected  undirected graph
  -l <label>        graph label
  -r <rankdir>      layout direction: LR|TB|RL|BT (default: TB)
  --attrs           pass through attribute lines (key=value)`)
	os.Exit(1)
}

func main() {
	directed := true
	label := ""
	rankdir := "TB"
	var file string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d", "--directed": directed = true
		case "-u", "--undirected": directed = false
		case "-l": i++; label = args[i]
		case "-r": i++; rankdir = args[i]
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			file = args[i]
		}
	}

	var r *os.File
	if file != "" {
		var err error; r, err = os.Open(file)
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer r.Close()
	} else { r = os.Stdin }

	graphType := "digraph"; edgeOp := "->"
	if !directed { graphType = "graph"; edgeOp = "--" }

	fmt.Printf("%s G {\n", graphType)
	fmt.Printf("  rankdir=%s;\n", rankdir)
	fmt.Printf("  node [shape=box, fontname=\"Helvetica\"];\n")
	fmt.Printf("  edge [fontname=\"Helvetica\", fontsize=10];\n")
	if label != "" { fmt.Printf("  label=%q;\n", label) }
	fmt.Println()

	sc := bufio.NewScanner(r)
	nodes := map[string]bool{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") { continue }
		// Attribute line
		if strings.Contains(line, "=") && !strings.Contains(line, "->") && !strings.Contains(line, "--") {
			fmt.Printf("  %s;\n", line); continue
		}
		// Edge line
		sep := "->"; if strings.Contains(line, "--") { sep = "--" }
		parts := strings.SplitN(line, sep, 2)
		if len(parts) != 2 {
			// Solo node
			node := strings.Trim(strings.TrimSpace(line), `"`)
			if !nodes[node] { nodes[node] = true; fmt.Printf("  %q;\n", node) }
			continue
		}
		from := strings.Trim(strings.TrimSpace(parts[0]), `"`)
		rest := strings.TrimSpace(parts[1])
		// Extract label if present: B [label="x"]
		toNode := rest; attrs := ""
		if idx := strings.Index(rest, "["); idx >= 0 { toNode = strings.TrimSpace(rest[:idx]); attrs = " " + rest[idx:] }
		to := strings.Trim(toNode, `"`)
		nodes[from] = true; nodes[to] = true
		fmt.Printf("  %q %s %q%s;\n", from, edgeOp, to, attrs)
	}
	fmt.Println("}")
}
