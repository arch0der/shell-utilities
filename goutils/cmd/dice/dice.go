// dice - roll dice using standard notation (e.g. 2d6, 1d20+5, 4d6kh3)
package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var diceRe = regexp.MustCompile(`(?i)^(\d*)d(\d+)(kh(\d+)|kl(\d+))?([-+]\d+)?$`)

func roll(expr string) {
	m := diceRe.FindStringSubmatch(strings.TrimSpace(expr))
	if m == nil { fmt.Fprintf(os.Stderr, "dice: invalid expression %q\n", expr); return }

	numDice := 1
	if m[1] != "" { numDice, _ = strconv.Atoi(m[1]) }
	sides, _ := strconv.Atoi(m[2])
	if sides < 1 || numDice < 1 { fmt.Fprintln(os.Stderr, "dice: invalid dice spec"); return }

	rolls := make([]int, numDice)
	for i := range rolls { rolls[i] = rand.Intn(sides) + 1 }

	kept := rolls
	keepLabel := ""
	if m[4] != "" { // kh = keep highest
		n, _ := strconv.Atoi(m[4])
		sorted := append([]int{}, rolls...)
		sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
		kept = sorted[:n]
		keepLabel = fmt.Sprintf(" keep-high-%d", n)
	} else if m[5] != "" { // kl = keep lowest
		n, _ := strconv.Atoi(m[5])
		sorted := append([]int{}, rolls...)
		sort.Ints(sorted)
		kept = sorted[:n]
		keepLabel = fmt.Sprintf(" keep-low-%d", n)
	}

	mod := 0
	if m[6] != "" { mod, _ = strconv.Atoi(m[6]) }

	sum := 0
	for _, v := range kept { sum += v }
	sum += mod

	rollStrs := make([]string, len(rolls))
	for i, v := range rolls { rollStrs[i] = strconv.Itoa(v) }

	modStr := ""
	if mod != 0 { modStr = fmt.Sprintf("%+d", mod) }
	fmt.Printf("%s%s: [%s]%s = %d\n", expr, keepLabel, strings.Join(rollStrs, ", "), modStr, sum)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		roll("1d6"); return
	}
	for _, arg := range os.Args[1:] { roll(arg) }
}
