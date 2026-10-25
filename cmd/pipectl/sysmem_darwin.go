package main

import "golang.org/x/sys/unix"

// totalSystemMemory returns total physical RAM in bytes, or false if it
// could not be determined.
func totalSystemMemory() (int64, bool) {
	v, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return 0, false
	}
	return int64(v), true
}
