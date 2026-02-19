//go:build !linux

package main

import "fmt"

func setHostname(name string) error {
	return fmt.Errorf("setting hostname not supported on this platform")
}
