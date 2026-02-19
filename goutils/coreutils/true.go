package main

import "os"

func init() { register("true", runTrue) }

func runTrue() { os.Exit(0) }
