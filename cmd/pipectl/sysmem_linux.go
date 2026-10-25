package main

import "syscall"

// totalSystemMemory returns total physical RAM in bytes, or false if it
// could not be determined.
func totalSystemMemory() (int64, bool) {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		return 0, false
	}
	return int64(uint64(info.Totalram) * uint64(info.Unit)), true
}
