// cron - Simple cron-like task scheduler (runs in foreground)
// Usage: cron <crontab_file>
// Format: minute hour day month weekday command
//         * * * * * command
//         */5 * * * * command  (every 5 minutes)
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Job struct {
	min     string
	hour    string
	day     string
	month   string
	weekday string
	cmd     string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: cron <crontab_file>")
		os.Exit(1)
	}
	jobs := loadJobs(os.Args[1])
	fmt.Printf("cron: loaded %d job(s), running...\n", len(jobs))
	for {
		now := time.Now()
		for _, job := range jobs {
			if matches(job, now) {
				go runJob(job.cmd)
			}
		}
		// Sleep until next minute
		next := now.Truncate(time.Minute).Add(time.Minute)
		time.Sleep(time.Until(next))
	}
}

func loadJobs(path string) []Job {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cron:", err)
		os.Exit(1)
	}
	defer f.Close()
	var jobs []Job
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		jobs = append(jobs, Job{
			min:     fields[0],
			hour:    fields[1],
			day:     fields[2],
			month:   fields[3],
			weekday: fields[4],
			cmd:     strings.Join(fields[5:], " "),
		})
	}
	return jobs
}

func matches(j Job, t time.Time) bool {
	return matchField(j.min, t.Minute(), 0, 59) &&
		matchField(j.hour, t.Hour(), 0, 23) &&
		matchField(j.day, t.Day(), 1, 31) &&
		matchField(j.month, int(t.Month()), 1, 12) &&
		matchField(j.weekday, int(t.Weekday()), 0, 6)
}

func matchField(field string, val, min, max int) bool {
	if field == "*" {
		return true
	}
	// */n - every n
	if strings.HasPrefix(field, "*/") {
		n, err := strconv.Atoi(field[2:])
		if err != nil || n == 0 {
			return false
		}
		return (val-min)%n == 0
	}
	// List: 1,2,5
	for _, part := range strings.Split(field, ",") {
		// Range: 1-5
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, _ := strconv.Atoi(bounds[0])
			hi, _ := strconv.Atoi(bounds[1])
			if val >= lo && val <= hi {
				return true
			}
		} else {
			n, _ := strconv.Atoi(part)
			if n == val {
				return true
			}
		}
	}
	return false
}

func runJob(cmd string) {
	fmt.Printf("cron: running: %s\n", cmd)
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cron: job failed: %v\n", err)
	}
}
