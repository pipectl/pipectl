package main

import (
	"strings"
	"testing"
)

func TestMemoryWarning(t *testing.T) {
	const eightGB = 8 * gigabyte
	const thirtyTwoGB = 32 * gigabyte

	tests := []struct {
		name         string
		maxInputSize int64
		totalRAM     int64
		wantWarning  bool
		wantContains string
	}{
		{
			name:         "default max-input-size on a common 8GB machine, no warning",
			maxInputSize: 256 * megabyte,
			totalRAM:     eightGB,
			wantWarning:  false,
		},
		{
			name:         "large max-input-size on modest machine, warns",
			maxInputSize: 2 * gigabyte,
			totalRAM:     eightGB,
			wantWarning:  true,
			wantContains: "exceeds 50% of this machine's 8.0GB of RAM",
		},
		{
			name:         "large max-input-size on ample machine, no warning",
			maxInputSize: 1 * gigabyte,
			totalRAM:     thirtyTwoGB,
			wantWarning:  false,
		},
		{
			name:         "default max-input-size on a constrained 2GB machine, warns",
			maxInputSize: 256 * megabyte,
			totalRAM:     2 * gigabyte,
			wantWarning:  true,
		},
		{
			name:         "unknown total RAM never warns",
			maxInputSize: 100 * gigabyte,
			totalRAM:     0,
			wantWarning:  false,
		},
		{
			name:         "zero or negative max-input-size never warns",
			maxInputSize: 0,
			totalRAM:     eightGB,
			wantWarning:  false,
		},
		{
			name:         "boundary just above threshold warns",
			maxInputSize: eightGB/2/peakRSSMultiplier + 1,
			totalRAM:     eightGB,
			wantWarning:  true,
		},
		{
			name:         "boundary just below threshold does not warn",
			maxInputSize: eightGB / 2 / peakRSSMultiplier,
			totalRAM:     eightGB,
			wantWarning:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := memoryWarning(tt.maxInputSize, tt.totalRAM)
			if tt.wantWarning && got == "" {
				t.Fatalf("memoryWarning(%d, %d) = %q, want a warning", tt.maxInputSize, tt.totalRAM, got)
			}
			if !tt.wantWarning && got != "" {
				t.Fatalf("memoryWarning(%d, %d) = %q, want no warning", tt.maxInputSize, tt.totalRAM, got)
			}
			if tt.wantContains != "" && !strings.Contains(got, tt.wantContains) {
				t.Fatalf("memoryWarning(%d, %d) = %q, want it to contain %q", tt.maxInputSize, tt.totalRAM, got, tt.wantContains)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{name: "bytes", bytes: 500, want: "500B"},
		{name: "kilobytes", bytes: 2 * kilobyte, want: "2.0KB"},
		{name: "megabytes", bytes: 5 * megabyte, want: "5.0MB"},
		{name: "gigabytes", bytes: 2 * gigabyte, want: "2.0GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatBytes(tt.bytes); got != tt.want {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}
