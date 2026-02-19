 1048576, 0}, "megabyte": {"B", 1048576, 0}, "megabytes": {"B", 1048576, 0},
	"gb":  {"B", 1073741824, 0}, "gigabyte": {"B", 1073741824, 0}, "gigabytes": {"B", 1073741824, 0},
	"tb":  {"B", 1099511627776, 0}, "terabyte": {"B", 1099511627776, 0},
	"kib": {"B", 1024, 0}, "mib": {"B", 1048576, 0}, "gib": {"B", 1073741824, 0},
	"bit": {"B", 0.125, 0}, "bits": {"B", 0.125, 0},
	"kbit": {"B", 128, 0}, "mbit": {"B", 131072, 0}, "gbit": {"B", 134217728, 0},

	// Angle (base: degree)
	"deg":     {"deg", 1, 0}, "degree": {"deg", 1, 0}, "degrees": {"deg", 1, 0},
	"rad":     {"deg", 180 / math.Pi, 0}, "radian": {"deg", 180 / math.Pi, 0}, "radians": {"deg", 180 / math.Pi, 0},
	"grad":    {"deg", 0.9, 0}, "gradian": {"deg", 0.9, 0},
	"arcmin":  {"deg", 1.0 / 60.0, 0},
	"arcsec":  {"deg", 1.0 / 3600.0, 0},
}

func convert(val float64, from, to string) (float64, error) {
	fromUnit, ok1 := unitTable[strings.ToLower(from)]
	toUnit, ok2 := unitTable[strings.ToLower(to)]
	if !ok1 {
		return 0, fmt.Errorf("unknown unit: %s", from)
	}
	if !ok2 {
		return 0, fmt.Errorf("unknown unit: %s", to)
	}
	if fromUnit.base != toUnit.base {
		return 0, fmt.Errorf("incompatible units: %s (%s) vs %s (%s)", from, fromUnit.base, to, toUnit.base)
	}

	// Temperature special case: val*factor + offset gives base (°C)
	// Then from base: (base - offset_to) / factor_to
	base := val*fromUnit.factor + fromUnit.offset
	result := (base - toUnit.offset) / toUnit.factor
	return result, nil
}

func formatResult(v float64) string {
	if math.Abs(v) >= 1e10 || (math.Abs(v) < 1e-4 && v != 0) {
		return fmt.Sprintf("%g", v)
	}
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%g", v)
}

func doConvert(line string) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: <value> <from> <to>")
		return
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid number: %s\n", fields[0])
		return
	}
	result, err := convert(val, fields[1], fields[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	from := unitTable[strings.ToLower(fields[1])]
	to := unitTable[strings.ToLower(fields[2])]
	fmt.Printf("%s %s = %s %s\n", formatResult(val), fields[1], formatResult(result), fields[2])
	_ = from
	_ = to
}

func main() {
	if len(os.Args) >= 4 {
		doConvert(strings.Join(os.Args[1:], " "))
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--list" {
		fmt.Println("Available units (grouped by dimension):")
		fmt.Println("  Length:      m km cm mm mi ft in yd nm ly au")
		fmt.Println("  Mass:        kg g mg lb oz t tonne stone")
		fmt.Println("  Time:        s ms min h d wk yr")
		fmt.Println("  Temperature: c f k (celsius fahrenheit kelvin)")
		fmt.Println("  Area:        m2 km2 ft2 mi2 ac ha")
		fmt.Println("  Volume:      l ml m3 gal qt pt cup floz tbsp tsp")
		fmt.Println("  Speed:       m/s kph mph kn fps")
		fmt.Println("  Energy:      j kj cal kcal wh kwh btu ev")
		fmt.Println("  Power:       w kw mw hp")
		fmt.Println("  Pressure:    pa kpa bar atm psi mmhg")
		fmt.Println("  Storage:     b kb mb gb tb bit kbit mbit gbit")
		fmt.Println("  Angle:       deg rad grad arcmin arcsec")
		return
	}

	// Interactive mode
	fmt.Println("units - unit converter (type 'quit' to exit, '--list' for units)")
	fmt.Println("Usage: <value> <from_unit> <to_unit>")
	fmt.Println("Example: 100 f c")
	fmt.Println()
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("units> ")
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "quit" || line == "q" || line == "exit" {
			break
		}
		if line == "--list" || line == "list" {
			os.Args = []string{"units", "--list"}
			// rerun
			fmt.Println("Length: m km cm mm mi ft in yd | Mass: kg g lb oz | Time: s min h d | Temp: c f k | Speed: kph mph | Energy: j kj cal kwh | Storage: b kb mb gb")
			continue
		}
		if line == "" {
			continue
		}
		doConvert(line)
	}
}
