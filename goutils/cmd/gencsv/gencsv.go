// gencsv - generate test/mock CSV data with realistic fake values
package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	firstNames = []string{"Alice","Bob","Carol","David","Eve","Frank","Grace","Hank","Iris","Jack","Kate","Liam","Mia","Noah","Olivia","Paul","Quinn","Rose","Sam","Tara"}
	lastNames  = []string{"Smith","Jones","Brown","Taylor","Wilson","Davis","Miller","Moore","Anderson","Jackson","White","Harris","Martin","Garcia","Martinez","Robinson","Clark","Lewis","Walker","Hall"}
	domains    = []string{"gmail.com","yahoo.com","outlook.com","example.com","test.org","company.io"}
	cities     = []string{"New York","Los Angeles","Chicago","Houston","Phoenix","Philadelphia","San Antonio","San Diego","Dallas","San Jose","Austin","Jacksonville","Denver","Seattle","Nashville"}
	countries  = []string{"US","UK","CA","AU","DE","FR","JP","BR","IN","MX"}
	statuses   = []string{"active","inactive","pending","suspended"}
	depts      = []string{"Engineering","Marketing","Sales","Finance","HR","Legal","Operations","Design"}
)

type colDef struct{ name, kind string }

func parseSchema(spec string) []colDef {
	var cols []colDef
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, ":"); idx > 0 {
			cols = append(cols, colDef{part[:idx], part[idx+1:]})
		} else {
			cols = append(cols, colDef{part, "string"})
		}
	}
	return cols
}

var defaultSchema = []colDef{
	{"id","id"},{"first_name","first"},{"last_name","last"},{"email","email"},
	{"city","city"},{"country","country"},{"age","age"},{"status","status"},{"created_at","date"},
}

func genValue(rng *rand.Rand, kind string, rowN int) string {
	switch kind {
	case "id": return strconv.Itoa(rowN + 1)
	case "first": return firstNames[rng.Intn(len(firstNames))]
	case "last": return lastNames[rng.Intn(len(lastNames))]
	case "email":
		f := strings.ToLower(firstNames[rng.Intn(len(firstNames))])
		l := strings.ToLower(lastNames[rng.Intn(len(lastNames))])
		return fmt.Sprintf("%s.%s@%s", f, l, domains[rng.Intn(len(domains))])
	case "name": return firstNames[rng.Intn(len(firstNames))] + " " + lastNames[rng.Intn(len(lastNames))]
	case "city": return cities[rng.Intn(len(cities))]
	case "country": return countries[rng.Intn(len(countries))]
	case "age": return strconv.Itoa(18 + rng.Intn(62))
	case "status": return statuses[rng.Intn(len(statuses))]
	case "dept": return depts[rng.Intn(len(depts))]
	case "date":
		days := rng.Intn(3*365)
		return time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	case "datetime":
		days := rng.Intn(3*365)
		return time.Now().AddDate(0, 0, -days).Format("2006-01-02T15:04:05Z")
	case "int": return strconv.Itoa(rng.Intn(10000))
	case "float": return fmt.Sprintf("%.2f", rng.Float64()*10000)
	case "bool": if rng.Intn(2) == 0 { return "true" }; return "false"
	case "uuid":
		b := make([]byte, 16); rng.Read(b); b[6] = (b[6]&0x0f)|0x40; b[8] = (b[8]&0x3f)|0x80
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4],b[4:6],b[6:8],b[8:10],b[10:])
	case "phone": return fmt.Sprintf("+1-%03d-%03d-%04d", 200+rng.Intn(800), rng.Intn(1000), rng.Intn(10000))
	case "price": return fmt.Sprintf("%.2f", 0.99+float64(rng.Intn(999))+rng.Float64())
	default: return fmt.Sprintf("%s_%d", kind, rng.Intn(1000))
	}
}

func main() {
	n := 10
	schema := ""
	noHeader := false
	seed := time.Now().UnixNano()
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n": i++; n, _ = strconv.Atoi(args[i])
		case "-s", "--schema": i++; schema = args[i]
		case "-H": noHeader = true
		case "--seed": i++; seed, _ = strconv.ParseInt(args[i], 10, 64)
		default: if v, err := strconv.Atoi(args[i]); err == nil { n = v }
		}
	}

	var cols []colDef
	if schema != "" { cols = parseSchema(schema) } else { cols = defaultSchema }
	rng := rand.New(rand.NewSource(seed))
	w := csv.NewWriter(os.Stdout)
	if !noHeader {
		headers := make([]string, len(cols))
		for i, c := range cols { headers[i] = c.name }
		w.Write(headers)
	}
	for row := 0; row < n; row++ {
		record := make([]string, len(cols))
		for i, c := range cols { record[i] = genValue(rng, c.kind, row) }
		w.Write(record)
	}
	w.Flush()
}
