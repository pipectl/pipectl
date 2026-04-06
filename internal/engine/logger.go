package engine

import (
	"fmt"
	"io"
	"os"
)

// Logger writes diagnostic output to an io.Writer.
// Log writes unconditionally; Debug writes only when verbose is enabled.
// All methods are nil-safe.
type Logger struct {
	w       io.Writer
	verbose bool
}

func NewLogger(verbose bool) *Logger {
	return &Logger{w: os.Stderr, verbose: verbose}
}

// NewLoggerWithWriter creates a Logger that writes to w. Intended for tests.
func NewLoggerWithWriter(w io.Writer, verbose bool) *Logger {
	return &Logger{w: w, verbose: verbose}
}

// Log writes a message unconditionally.
func (l *Logger) Log(format string, args ...interface{}) {
	if l == nil {
		return
	}
	fmt.Fprintf(l.w, format+"\n", args...)
}

// Debug writes a message only when verbose mode is enabled.
func (l *Logger) Debug(format string, args ...interface{}) {
	if l == nil || !l.verbose {
		return
	}
	fmt.Fprintf(l.w, format+"\n", args...)
}
