package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("tsort", runTsort) }

func runTsort() {
	args := os.Args[1:]
	var r io.Reader = os.Stdin
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		fh, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "tsort:", err)
			os.Exit(1)
		}
		defer fh.Close()
		r = fh
	}

	graph := map[string][]string{}
	inDegree := map[string]int{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		for i := 0; i+1 < len(fields); i += 2 {
			from, to := fields[i], fields[i+1]
			if _, ok := graph[from]; !ok {
				graph[from] = nil
				inDegree[from] = 0
			}
			if _, ok := graph[to]; !ok {
				graph[to] = nil
				inDegree[to] = 0
			}
			graph[from] = append(graph[from], to)
			inDegree[to]++
		}
	}

	// Kahn's algorithm
	var queue []string
	for node := range graph {
		if inDegree[node] == 0 {
			queue = append(queue, node)
		}
	}

	exitCode := 0
	count := 0
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fmt.Println(node)
		count++
		for _, next := range graph[node] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	if count != len(graph) {
		fmt.Fprintln(os.Stderr, "tsort: cycle detected")
		exitCode = 1
	}
	os.Exit(exitCode)
}
