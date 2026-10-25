//go:build !darwin && !linux && !windows

package main

// totalSystemMemory is unavailable on this platform; the memory warning is
// simply skipped rather than failing the run.
func totalSystemMemory() (int64, bool) {
	return 0, false
}
