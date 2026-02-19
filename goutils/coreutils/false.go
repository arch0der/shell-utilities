package main

import "os"

func init() { register("false", runFalse) }

func runFalse() { os.Exit(1) }
