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

	output := captureWriteOutput(t, csvPayload, JSONType)

	assertContains(t, output, "[\n")
	assertContains(t, output, `"id": "1"`)
	assertContains(t, output, `"name": "alice"`)
	assertContains(t, output, `"id": "2"`)
	assertContains(t, output, `"name": "bob"`)
}

func TestWriteJSONConvertsNestedCSVFieldsAndArrays(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"name", "nested.name", "values"},
			{" John Smith ", "nested value", `["quick","brown","fox"]`},
		},
	}

	output := captureWriteOutput(t, csvPayload, JSONType)

	assertContains(t, output, `"name": " John Smith "`)
	assertContains(t, output, `"nested": {`)
	assertContains(t, output, `"name": "nested value"`)
	assertContains(t, output, `"values": [`)
	assertContains(t, output, `"quick"`)
	assertContains(t, output, `"brown"`)
	assertContains(t, output, `"fox"`)
}

func TestWriteJSONConvertsHeaderOnlyCSVToEmptyJSONArray(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
		},
	}

	output := captureWriteOutput(t, csvPayload, JSONType)

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

	err := Write(csvPayload, JSONType, nil)
	if err == nil {
		t.Fatal("expected error for row length mismatch")
	}

	assertContains(t, err.Error(), "row 2")
	assertContains(t, err.Error(), "expected 2")
}

func TestWriteJSONReturnsErrorForConflictingNestedCSVFields(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"nested", "nested.name"},
			{"value", "nested value"},
		},
	}

	err := Write(csvPayload, JSONType, nil)
	if err == nil {
		t.Fatal("expected error for conflicting nested fields")
	}

	assertContains(t, err.Error(), `field "nested" conflicts with nested field "nested.name"`)
}
