// logger - Log messages to syslog or stdout
// Usage: logger [-p priority] [-t tag] [-s] [message]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/syslog"
	"os"
	"strings"
	"time"
)

var (
	priority = flag.String("p", "user.notice", "Priority: facility.level (e.g. user.info, local0.err)")
	tag      = flag.String("t", "", "Tag/program name")
	stderr   = flag.Bool("s", false, "Also log to stderr")
)

var facilityMap = map[string]syslog.Priority{
	"kern": syslog.LOG_KERN, "user": syslog.LOG_USER, "mail": syslog.LOG_MAIL,
	"daemon": syslog.LOG_DAEMON, "auth": syslog.LOG_AUTH, "syslog": syslog.LOG_SYSLOG,
	"lpr": syslog.LOG_LPR, "news": syslog.LOG_NEWS, "local0": syslog.LOG_LOCAL0,
	"local1": syslog.LOG_LOCAL1, "local2": syslog.LOG_LOCAL2, "local3": syslog.LOG_LOCAL3,
	"local4": syslog.LOG_LOCAL4, "local5": syslog.LOG_LOCAL5, "local6": syslog.LOG_LOCAL6,
	"local7": syslog.LOG_LOCAL7,
}

var levelMap = map[string]syslog.Priority{
	"emerg": syslog.LOG_EMERG, "alert": syslog.LOG_ALERT, "crit": syslog.LOG_CRIT,
	"err": syslog.LOG_ERR, "warning": syslog.LOG_WARNING, "notice": syslog.LOG_NOTICE,
	"info": syslog.LOG_INFO, "debug": syslog.LOG_DEBUG,
}

func parsePriority(s string) syslog.Priority {
	parts := strings.SplitN(s, ".", 2)
	p := syslog.LOG_USER | syslog.LOG_NOTICE
	if len(parts) == 2 {
		if f, ok := facilityMap[parts[0]]; ok {
			p = f
		}
		if l, ok := levelMap[parts[1]]; ok {
			p |= l
		}
	} else {
		if l, ok := levelMap[parts[0]]; ok {
			p = syslog.LOG_USER | l
		}
	}
	return p
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: logger [-p priority] [-t tag] [-s] [message]")
		flag.PrintDefaults()
	}
	flag.Parse()

	prio := parsePriority(*priority)
	tagStr := *tag
	if tagStr == "" {
		tagStr = os.Args[0]
	}

	logMsg := func(msg string) {
		// Try syslog first
		w, err := syslog.New(prio, tagStr)
		if err == nil {
			w.Write([]byte(msg))
			w.Close()
		} else {
			// Fallback: print to stderr with timestamp
			now := time.Now().Format("Jan  2 15:04:05")
			fmt.Fprintf(os.Stderr, "%s %s: %s\n", now, tagStr, msg)
		}
		if *stderr {
			now := time.Now().Format("Jan  2 15:04:05")
			fmt.Fprintf(os.Stderr, "%s %s: %s\n", now, tagStr, msg)
		}
	}

	if flag.NArg() > 0 {
		logMsg(strings.Join(flag.Args(), " "))
		return
	}

	// Read from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		logMsg(sc.Text())
	}
}
