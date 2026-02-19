// numconv - convert numbers between units (length, weight, temp, area, volume, speed)
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type unit struct{ name, abbr string; toBase float64; category string; offset float64 }

var units = []unit{
	// Length (base: meter)
	{"meter","m",1,"length",0}, {"kilometer","km",1000,"length",0},
	{"centimeter","cm",0.01,"length",0}, {"millimeter","mm",0.001,"length",0},
	{"inch","in",0.0254,"length",0}, {"foot","ft",0.3048,"length",0},
	{"yard","yd",0.9144,"length",0}, {"mile","mi",1609.344,"length",0},
	{"nautical_mile","nmi",1852,"length",0},
	// Weight (base: kilogram)
	{"kilogram","kg",1,"weight",0}, {"gram","g",0.001,"weight",0},
	{"pound","lb",0.453592,"weight",0}, {"ounce","oz",0.0283495,"weight",0},
	{"ton","t",1000,"weight",0}, {"stone","st",6.35029,"weight",0},
	// Temperature (base: celsius)
	{"celsius","c",1,"temp",0}, {"fahrenheit","f",5.0/9.0,"temp",-32},
	{"kelvin","k",1,"temp",-273.15},
	// Area (base: square meter)
	{"sqmeter","m2",1,"area",0}, {"sqkm","km2",1e6,"area",0},
	{"sqfoot","ft2",0.092903,"area",0}, {"sqmile","mi2",2.59e6,"area",0},
	{"acre","ac",4046.86,"area",0}, {"hectare","ha",10000,"area",0},
	// Volume (base: liter)
	{"liter","l",1,"volume",0}, {"milliliter","ml",0.001,"volume",0},
	{"gallon","gal",3.78541,"volume",0}, {"quart","qt",0.946353,"volume",0},
	{"pint","pt",0.473176,"volume",0}, {"cup","cup",0.236588,"volume",0},
	{"fluid_ounce","floz",0.0295735,"volume",0}, {"cubic_meter","m3",1000,"volume",0},
	// Speed (base: m/s)
	{"mps","mps",1,"speed",0}, {"kph","kph",0.277778,"speed",0},
	{"mph","mph",0.44704,"speed",0}, {"knot","knot",0.514444,"speed",0},
}

func findUnit(s string) *unit {
	s = strings.ToLower(strings.TrimSpace(s))
	for i, u := range units {
		if strings.ToLower(u.name) == s || strings.ToLower(u.abbr) == s { return &units[i] }
	}
	return nil
}

func convert(val float64, from, to *unit) (float64, error) {
	if from.category != to.category { return 0, fmt.Errorf("cannot convert %s to %s (different categories)", from.category, to.category) }
	if from.category == "temp" {
		celsius := (val + from.offset) * from.toBase
		return celsius/to.toBase - to.offset, nil
	}
	return val * from.toBase / to.toBase, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: numconv <value> <from_unit> <to_unit>")
		fmt.Fprintln(os.Stderr, "  categories: length weight temp area volume speed")
		fmt.Fprintln(os.Stderr, "  example: numconv 100 km mi  |  numconv 32 f c")
		os.Exit(1)
	}
	val, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil { fmt.Fprintln(os.Stderr, "numconv: invalid number"); os.Exit(1) }
	from := findUnit(os.Args[2])
	to := findUnit(os.Args[3])
	if from == nil { fmt.Fprintf(os.Stderr, "numconv: unknown unit %q\n", os.Args[2]); os.Exit(1) }
	if to == nil { fmt.Fprintf(os.Stderr, "numconv: unknown unit %q\n", os.Args[3]); os.Exit(1) }
	result, err := convert(val, from, to)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	fmt.Printf("%g %s = %g %s\n", val, from.abbr, result, to.abbr)
}
