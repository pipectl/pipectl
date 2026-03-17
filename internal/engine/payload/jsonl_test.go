package payload

import "testing"

func TestJSONLType(t *testing.T) {
	jsonlPayload := &JSONL{
		Records: []map[string]interface{}{
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
	if len(jsonlPayload.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(jsonlPayload.Records))
	}
	if jsonlPayload.Records[1]["id"] != float64(2) {
		t.Fatalf("unexpected records: %#v", jsonlPayload.Records)
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
		Records: []map[string]interface{}{
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
