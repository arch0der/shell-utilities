package main

import (
	"fmt"
	"strings"
	"time"
)

// humanSize converts bytes to human-readable string (IEC units)
func humanSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// parseDuration parses a duration string like "5s", "2m", "1h", "1d" or plain seconds
func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}
	// Go's time.ParseDuration handles ns, us, ms, s, m, h
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}
	// Handle 'd' for days and bare numbers (seconds)
	if strings.HasSuffix(s, "d") {
		n := 0.0
		if _, err := fmt.Sscan(s[:len(s)-1], &n); err == nil {
			return time.Duration(n * 24 * float64(time.Hour)), nil
		}
	}
	n := 0.0
	if _, err := fmt.Sscan(s, &n); err == nil {
		return time.Duration(n * float64(time.Second)), nil
	}
	return 0, fmt.Errorf("invalid duration: %s", s)
}

// parseSize parses a size string like "1k", "512b", "2M" or plain bytes
func parseSize(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	mult := int64(1)
	lower := strings.ToLower(s)
	switch {
	case strings.HasSuffix(lower, "kib"):
		mult = 1024
		s = s[:len(s)-3]
	case strings.HasSuffix(lower, "mib"):
		mult = 1024 * 1024
		s = s[:len(s)-3]
	case strings.HasSuffix(lower, "gib"):
		mult = 1024 * 1024 * 1024
		s = s[:len(s)-3]
	case strings.HasSuffix(lower, "kb"):
		mult = 1000
		s = s[:len(s)-2]
	case strings.HasSuffix(lower, "mb"):
		mult = 1000000
		s = s[:len(s)-2]
	case strings.HasSuffix(lower, "gb"):
		mult = 1000000000
		s = s[:len(s)-2]
	case strings.HasSuffix(lower, "k"):
		mult = 1024
		s = s[:len(s)-1]
	case strings.HasSuffix(lower, "m"):
		mult = 1024 * 1024
		s = s[:len(s)-1]
	case strings.HasSuffix(lower, "g"):
		mult = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case strings.HasSuffix(lower, "t"):
		mult = 1024 * 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case strings.HasSuffix(lower, "b"):
		mult = 512
		s = s[:len(s)-1]
	case strings.HasSuffix(lower, "c"):
		mult = 1
		s = s[:len(s)-1]
	case strings.HasSuffix(lower, "w"):
		mult = 2
		s = s[:len(s)-1]
	}
	n := int64(0)
	fmt.Sscan(s, &n)
	return n * mult
}

// pluralS returns "s" for pluralization if n != 1
func pluralS(n int64) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// echoUnescape processes backslash escape sequences in a string
func echoUnescape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		switch s[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case '\\':
			b.WriteByte('\\')
		case 'a':
			b.WriteByte('\a')
		case 'b':
			b.WriteByte('\b')
		case 'f':
			b.WriteByte('\f')
		case 'v':
			b.WriteByte('\v')
		case '0':
			b.WriteByte(0)
		default:
			b.WriteByte('\\')
			b.WriteByte(s[i])
		}
	}
	return b.String()
}
