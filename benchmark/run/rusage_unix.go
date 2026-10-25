//go:build unix

package main

import (
	"os"
	"runtime"
	"syscall"
)

// peakRSSBytes returns a subprocess's peak resident set size in bytes, using
// the rusage stats the OS already tracked for it (no dependency on the
// external `time` binary). syscall.Rusage.Maxrss is reported in different
// units per OS: bytes on macOS, kilobytes on Linux.
func peakRSSBytes(ps *os.ProcessState) (int64, bool) {
	ru, ok := ps.SysUsage().(*syscall.Rusage)
	if !ok || ru == nil {
		return 0, false
	}
	if runtime.GOOS == "darwin" {
		return ru.Maxrss, true
	}
	return ru.Maxrss * 1024, true
}
