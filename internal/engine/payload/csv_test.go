package payload

import (
	"strings"
	"testing"
)

func TestCSVType(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "alice"},
		},
	}

	if got := csvPayload.Type(); got != CSVType {
		t.Fatalf("expected payload type %q, got %q", CSVType, got)
	}
}

func TestCSVRecordCount(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "alice"},
			{"2", "bob"},
		},
	}

	if got := csvPayload.RecordCount(); got != 2 {
		t.Fatalf("expected record count 2, got %d", got)
	}
}

func TestWriteJSONConvertsCSVRowsToJSONArray(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "alice"},
			{"2", "bob"},
		},
	}

	output := captureStdout(t, func() {
		if err := Write(csvPayload, JSONType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	assertContains(t, output, "[\n")
	assertContains(t, output, `"id": "1"`)
	assertContains(t, output, `"name": "alice"`)
	assertContains(t, output, `"id": "2"`)
	assertContains(t, output, `"name": "bob"`)
}

func TestWriteJSONConvertsHeaderOnlyCSVToEmptyJSONArray(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
		},
	}

	output := captureStdout(t, func() {
		if err := Write(csvPayload, JSONType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	if strings.TrimSpace(output) != "[]" {
		t.Fatalf("expected empty JSON array, got %q", output)
	}
}

func TestWriteJSONReturnsErrorForCSVRowLengthMismatch(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1"},
		},
	}

	err := Write(csvPayload, JSONType)
	if err == nil {
		t.Fatal("expected error for row length mismatch")
	}

	assertContains(t, err.Error(), "row 2")
	assertContains(t, err.Error(), "expected 2")
}
