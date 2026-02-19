package main

import "os"

func init() { register("vdir", runVdir) }

func runVdir() {
	// vdir is like ls with -l -b
	os.Args = append([]string{os.Args[0], "-l", "-b"}, os.Args[1:]...)
	runLs()
}
