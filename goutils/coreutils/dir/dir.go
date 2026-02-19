package main

import "os"

func main() {
	// dir is like ls with -C -b
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-C", "-b")
	} else {
		os.Args = append([]string{os.Args[0], "-C", "-b"}, os.Args[1:]...)
	}
	runLs()
}
