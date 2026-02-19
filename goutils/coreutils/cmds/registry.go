package cmds

import "sort"

var registry = map[string]func(){}

func register(name string, fn func()) { registry[name] = fn }

func Lookup(name string) (func(), bool) {
	fn, ok := registry[name]
	return fn, ok
}

func List() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
