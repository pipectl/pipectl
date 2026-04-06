package payload

import (
	"strings"
	"testing"
)

func TestReadCSVWithDefaultDelimiter(t *testing.T) {
	input := []byte("id,name\n1,alice\n2,bob\n")
	p, err := ReadCSV(input, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	csv, ok := p.(*CSV)
	if !ok {
		t.Fatalf("expected *CSV, got %T", p)
	}
	if got := csv.RecordCount(); got != 2 {
		t.Fatalf("expected 2 records, got %d", got)
	}
	if csv.Rows[1][1] != "alice" {
		t.Fatalf("unexpected value: got %q want %q", csv.Rows[1][1], "alice")
	}
}

func TestReadCSVWithCustomDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter rune
	}{
		{name: "pipe", input: "id|name\n1|alice\n", delimiter: '|'},
		{name: "tab", input: "id\tname\n1\talice\n", delimiter: '\t'},
		{name: "semicolon", input: "id;name\n1;alice\n", delimiter: ';'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := ReadCSV([]byte(tt.input), tt.delimiter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			csv, ok := p.(*CSV)
			if !ok {
				t.Fatalf("expected *CSV, got %T", p)
			}
			if csv.Rows[0][0] != "id" || csv.Rows[0][1] != "name" {
				t.Fatalf("unexpected headers: %v", csv.Rows[0])
			}
			if csv.Rows[1][0] != "1" || csv.Rows[1][1] != "alice" {
				t.Fatalf("unexpected row: %v", csv.Rows[1])
			}
		})
	}
}

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
