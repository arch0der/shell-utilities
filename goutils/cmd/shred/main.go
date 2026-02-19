// shred - Overwrite a file to hide its contents, then optionally delete it
// Usage: shred [-n passes] [-z] [-u] [-v] file...
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"os"
)

var (
	passes  = flag.Int("n", 3, "Number of overwrite passes")
	zero    = flag.Bool("z", false, "Add final zero-fill pass")
	unlink  = flag.Bool("u", false, "Remove file after shredding")
	verbose = flag.Bool("v", false, "Show progress")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: shred [-n passes] [-z] [-u] [-v] file...")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	exitCode := 0
	for _, path := range flag.Args() {
		if err := shredFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "shred: %s: %v\n", path, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func shredFile(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	size := info.Size()
	buf := make([]byte, 65536)

	for pass := 0; pass < *passes; pass++ {
		if *verbose {
			fmt.Fprintf(os.Stderr, "shred: %s: pass %d/%d (random)...\n", path, pass+1, *passes)
		}
		if _, err := f.Seek(0, 0); err != nil {
			f.Close()
			return err
		}
		written := int64(0)
		for written < size {
			n := int64(len(buf))
			if n > size-written {
				n = size - written
			}
			rand.Read(buf[:n])
			if _, err := f.Write(buf[:n]); err != nil {
				f.Close()
				return err
			}
			written += n
		}
		f.Sync()
	}

	if *zero {
		if *verbose {
			fmt.Fprintf(os.Stderr, "shred: %s: pass %d/%d (zeros)...\n", path, *passes+1, *passes+1)
		}
		f.Seek(0, 0)
		for i := range buf {
			buf[i] = 0
		}
		written := int64(0)
		for written < size {
			n := int64(len(buf))
			if n > size-written {
				n = size - written
			}
			f.Write(buf[:n])
			written += n
		}
		f.Sync()
	}

	f.Close()

	if *unlink {
		if *verbose {
			fmt.Fprintf(os.Stderr, "shred: %s: removing\n", path)
		}
		return os.Remove(path)
	}
	return nil
}
