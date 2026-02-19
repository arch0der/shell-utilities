// bench - Benchmark a command by running it N times.
//
// Usage:
//
//	bench [OPTIONS] -- COMMAND [ARGS...]
//
// Options:
//
//	-n N      Number of runs (default: 10)
//	-w N      Warmup runs not counted (default: 1)
//	-t        Show per-run times
//	-j        JSON output
//	-q        Suppress command output
//	-s        Shell mode: run via sh -c
//
// Output:
//
//	min, max, mean, median, stddev, total
//
// Examples:
//
//	bench -- sleep 0.1
//	bench -n 100 -q -- grep -r "pattern" /etc
//	bench -s -n 20 -- "echo hello | wc -c"
//	bench -n 5 -t -- make build
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"sort"
	"time"
)

var (
	runs    = flag.Int("n", 10, "number of runs")
	warmup  = flag.Int("w", 1, "warmup runs")
	verbose = flag.Bool("t", false, "per-run times")
	asJSON  = flag.Bool("j", false, "JSON output")
	quiet   = flag.Bool("q", false, "suppress output")
	shell   = flag.Bool("s", false, "shell mode")
)

type Stats struct {
	Runs   int     `json:"runs"`
	Min    float64 `json:"min_ms"`
	Max    float64 `json:"max_ms"`
	Mean   float64 `json:"mean_ms"`
	Median float64 `json:"median_ms"`
	Stddev float64 `json:"stddev_ms"`
	Total  float64 `json:"total_ms"`
}

func runCmd(args []string) (time.Duration, error) {
	var cmd *exec.Cmd
	if *shell {
		cmd = exec.Command("sh", append([]string{"-c"}, args...)...)
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	if *quiet {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	} else {
		cmd.Stdout = os.Stderr // route to stderr so timing is clean
		cmd.Stderr = os.Stderr
	}
	start := time.Now()
	err := cmd.Run()
	return time.Since(start), err
}

func stddev(vals []float64, mean float64) float64 {
	if len(vals) < 2 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		d := v - mean
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(vals)-1))
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: bench [OPTIONS] -- COMMAND [ARGS...]")
		os.Exit(1)
	}

	// Warmup
	for i := 0; i < *warmup; i++ {
		if !*quiet {
			fmt.Fprintf(os.Stderr, "warmup %d/%d\n", i+1, *warmup)
		}
		runCmd(args)
	}

	var times []float64
	for i := 0; i < *runs; i++ {
		d, err := runCmd(args)
		ms := float64(d.Microseconds()) / 1000.0
		times = append(times, ms)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bench: run %d failed: %v\n", i+1, err)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "run %d: %.3fms\n", i+1, ms)
		}
	}

	sorted := make([]float64, len(times))
	copy(sorted, times)
	sort.Float64s(sorted)

	total := 0.0
	for _, t := range times {
		total += t
	}
	mean := total / float64(len(times))
	median := sorted[len(sorted)/2]
	if len(sorted)%2 == 0 {
		median = (sorted[len(sorted)/2-1] + sorted[len(sorted)/2]) / 2
	}

	stats := Stats{
		Runs:   *runs,
		Min:    sorted[0],
		Max:    sorted[len(sorted)-1],
		Mean:   mean,
		Median: median,
		Stddev: stddev(times, mean),
		Total:  total,
	}

	if *asJSON {
		b, _ := json.MarshalIndent(stats, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Printf("runs:    %d\n", stats.Runs)
	fmt.Printf("min:     %.3fms\n", stats.Min)
	fmt.Printf("max:     %.3fms\n", stats.Max)
	fmt.Printf("mean:    %.3fms\n", stats.Mean)
	fmt.Printf("median:  %.3fms\n", stats.Median)
	fmt.Printf("stddev:  %.3fms\n", stats.Stddev)
	fmt.Printf("total:   %.3fms\n", stats.Total)
}
