package engine

import (
	"bytes"
	"testing"
)

func TestLogger_quiet_suppressesLog(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{w: &buf, verbose: true, quiet: true}

	l.Log("should not appear")
	l.Debug("should not appear either")

	if buf.Len() != 0 {
		t.Errorf("expected no output from quiet logger, got %q", buf.String())
	}
}

func TestLogger_log(t *testing.T) {
	var buf bytes.Buffer
	l := NewLoggerWithWriter(&buf, false)

	l.Log("hello %s", "world")
	if buf.String() != "hello world\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestLogger_debug_suppressed(t *testing.T) {
	var buf bytes.Buffer
	l := NewLoggerWithWriter(&buf, false)

	l.Debug("should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected no debug output when verbose=false, got %q", buf.String())
	}
}

func TestLogger_debug_verbose(t *testing.T) {
	var buf bytes.Buffer
	l := NewLoggerWithWriter(&buf, true)

	l.Debug("detail")
	if buf.String() != "detail\n" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}
