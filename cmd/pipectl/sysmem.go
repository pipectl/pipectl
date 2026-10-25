package main

import "fmt"

// peakRSSMultiplier is a conservative estimate of how much memory pipectl
// may use relative to --max-input-size. The performance benchmark measured
// peak-RSS/input-size ratios across CSV, JSON, and JSONL from 10MB to 1GB
// inputs, ranging ~9.3x (CSV, 1GB) up to ~15.8x (JSON, 10MB) (see
// benchmark/README.md and website/docs/performance.md); this constant
// rounds up from the highest observed ratio for headroom. pipectl reads the
// whole payload into memory and parses it into a record structure alongside
// the raw bytes, so peak usage runs well above the raw file size.
const peakRSSMultiplier = 16

// memoryWarningThreshold is the fraction of total system RAM above which
// pipectl warns that --max-input-size may be too high for this machine.
const memoryWarningThreshold = 0.5

// memoryWarning returns a one-line advisory if maxInputSize, scaled by
// peakRSSMultiplier, could plausibly exceed memoryWarningThreshold of
// totalRAM. It returns "" when no warning is warranted, including when
// totalRAM is unknown (0) — the check is advisory only and never blocks
// execution.
//
// This is intentionally based on total system RAM, not currently-available
// memory: "available" is a racy signal (what's free now may not be free by
// the time the large allocation actually happens) and its semantics differ
// enough between OSes that it would invite false precision. Total RAM is
// stable and gives a conservative, honest baseline.
func memoryWarning(maxInputSize, totalRAM int64) string {
	if totalRAM <= 0 || maxInputSize <= 0 {
		return ""
	}

	estimated := float64(maxInputSize) * peakRSSMultiplier
	if estimated <= float64(totalRAM)*memoryWarningThreshold {
		return ""
	}

	return fmt.Sprintf(
		"warning: --max-input-size %s may use up to ~%s of memory, which exceeds %.0f%% of this machine's %s of RAM (lower --max-input-size or free up memory before running on large files)",
		formatBytes(maxInputSize), formatBytes(int64(estimated)), memoryWarningThreshold*100, formatBytes(totalRAM),
	)
}

// formatBytes renders a byte count using the same units as --max-input-size.
func formatBytes(n int64) string {
	switch {
	case n >= gigabyte:
		return fmt.Sprintf("%.1fGB", float64(n)/gigabyte)
	case n >= megabyte:
		return fmt.Sprintf("%.1fMB", float64(n)/megabyte)
	case n >= kilobyte:
		return fmt.Sprintf("%.1fKB", float64(n)/kilobyte)
	default:
		return fmt.Sprintf("%dB", n)
	}
}
