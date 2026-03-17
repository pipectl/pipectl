package payload

import "testing"

func TestJSONLType(t *testing.T) {
	jsonlPayload := &JSONL{
		Items: []map[string]interface{}{
			{"id": 1, "name": "alice"},
		},
	}

	if got := jsonlPayload.Type(); got != JSONLType {
		t.Fatalf("expected payload type %q, got %q", JSONLType, got)
	}
}

func TestReadJSONLParsesEachLineAsRecord(t *testing.T) {
	got, err := Read([]byte("{\"id\":1}\n{\"id\":2}\n"), JSONLType)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	jsonlPayload, ok := got.(*JSONL)
	if !ok {
		t.Fatalf("expected payload.JSONL, got %T", got)
	}
	if len(jsonlPayload.Items) != 2 {
		t.Fatalf("expected 2 records, got %d", len(jsonlPayload.Items))
	}
	if jsonlPayload.Items[1]["id"] != float64(2) {
		t.Fatalf("unexpected records: %#v", jsonlPayload.Items)
	}
}

func TestReadJSONLRejectsNonObjectLine(t *testing.T) {
	_, err := Read([]byte("{\"id\":1}\n[1,2,3]\n"), JSONLType)
	if err == nil {
		t.Fatal("expected error for non-object JSONL line")
	}

	assertContains(t, err.Error(), "expected object")
}

func TestWriteJSONLPreservesLineDelimitedOutput(t *testing.T) {
	jsonlPayload := &JSONL{
		Items: []map[string]interface{}{
			{"id": 1},
			{"id": 2},
		},
	}

	output := captureStdout(t, func() {
		if err := Write(jsonlPayload, JSONLType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	assertContains(t, output, "{\"id\":1}\n")
	assertContains(t, output, "{\"id\":2}\n")
}

func TestWriteJSONLConvertsCSVRowsToLineDelimitedObjects(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1", "alice"},
			{"2", "bob"},
		},
	}

	output := captureStdout(t, func() {
		if err := Write(csvPayload, JSONLType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	assertContains(t, output, "{\"id\":\"1\",\"name\":\"alice\"}\n")
	assertContains(t, output, "{\"id\":\"2\",\"name\":\"bob\"}\n")
}

func TestWriteJSONLReturnsErrorForCSVRowLengthMismatch(t *testing.T) {
	csvPayload := &CSV{
		Rows: [][]string{
			{"id", "name"},
			{"1"},
		},
	}

	err := Write(csvPayload, JSONLType)
	if err == nil {
		t.Fatal("expected error for row length mismatch")
	}

	assertContains(t, err.Error(), "row 2")
	assertContains(t, err.Error(), "expected 2")
}
