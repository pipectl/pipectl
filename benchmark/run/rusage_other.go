//go:build !unix

package main

import "os"

// peakRSSBytes is unavailable outside the syscall.Rusage-based Unix path;
// the benchmark still runs, it just omits peak-RSS/peak-RSS-ratio columns.
func peakRSSBytes(_ *os.ProcessState) (int64, bool) {
	return 0, false
}
