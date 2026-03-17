package payload

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestJSONType(t *testing.T) {
	jsonPayload := &JSON{
		Items: []map[string]interface{}{
			{"id": 1, "name": "alice"},
		},
		Shape: JSONObjectShape,
	}

	if got := jsonPayload.Type(); got != JSONType {
		t.Fatalf("expected payload type %q, got %q", JSONType, got)
	}
}

func TestJSONRecordCount(t *testing.T) {
	jsonPayload := &JSON{
		Items: []map[string]interface{}{
			{"id": 1},
			{"id": 2},
		},
		Shape: JSONArrayShape,
	}

	if got := jsonPayload.RecordCount(); got != 2 {
		t.Fatalf("expected record count 2, got %d", got)
	}
}

func TestReadJSONTreatsObjectAsSingleRecord(t *testing.T) {
	got, err := Read([]byte(`{"id":1,"name":"alice"}`), JSONType)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	jsonPayload, ok := got.(*JSON)
	if !ok {
		t.Fatalf("expected payload.JSON, got %T", got)
	}
	if jsonPayload.Shape != JSONObjectShape {
		t.Fatalf("expected shape %q, got %q", JSONObjectShape, jsonPayload.Shape)
	}
	if len(jsonPayload.Items) != 1 {
		t.Fatalf("expected 1 record, got %d", len(jsonPayload.Items))
	}
	if jsonPayload.Items[0]["name"] != "alice" {
		t.Fatalf("unexpected records: %#v", jsonPayload.Items)
	}
}

func TestReadJSONTreatsArrayAsRecords(t *testing.T) {
	got, err := Read([]byte(`[{"id":1},{"id":2}]`), JSONType)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	jsonPayload := got.(*JSON)
	if jsonPayload.Shape != JSONArrayShape {
		t.Fatalf("expected shape %q, got %q", JSONArrayShape, jsonPayload.Shape)
	}
	if len(jsonPayload.Items) != 2 {
		t.Fatalf("expected 2 records, got %d", len(jsonPayload.Items))
	}
}

func TestReadJSONRejectsPrimitiveArray(t *testing.T) {
	_, err := Read([]byte(`["quick","brown","fox"]`), JSONType)
	if err == nil {
		t.Fatal("expected error for primitive JSON array")
	}

	assertContains(t, err.Error(), "expected object")
}

func TestReadJSONRejectsMixedArray(t *testing.T) {
	_, err := Read([]byte(`[{"id":1},"brown"]`), JSONType)
	if err == nil {
		t.Fatal("expected error for mixed JSON array")
	}

	assertContains(t, err.Error(), "expected object")
}

func TestWriteJSONPreservesObjectShape(t *testing.T) {
	jsonPayload := &JSON{
		Items: []map[string]interface{}{
			{"id": 1, "name": "alice"},
		},
		Shape: JSONObjectShape,
	}

	output := captureStdout(t, func() {
		if err := Write(jsonPayload, JSONType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	if strings.Contains(output, "[") {
		t.Fatalf("expected object output, got %q", output)
	}
	assertContains(t, output, `"id": 1`)
	assertContains(t, output, `"name": "alice"`)
}

func TestWriteJSONPreservesArrayShape(t *testing.T) {
	jsonPayload := &JSON{
		Items: []map[string]interface{}{
			{"id": 1},
			{"id": 2},
		},
		Shape: JSONArrayShape,
	}

	output := captureStdout(t, func() {
		if err := Write(jsonPayload, JSONType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	assertContains(t, output, "[\n")
	assertContains(t, output, `"id": 1`)
	assertContains(t, output, `"id": 2`)
}

func TestWriteCSVConvertsJSONPayloadWithNestedObjects(t *testing.T) {
	jsonPayload := &JSON{
		Items: []map[string]interface{}{
			{
				"name": "John Smith",
				"address": map[string]interface{}{
					"city":    "Sydney",
					"country": "AU",
				},
			},
		},
		Shape: JSONObjectShape,
	}

	output := captureStdout(t, func() {
		if err := Write(jsonPayload, CSVType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	assertContains(t, output, "address.city,address.country,name\n")
	assertContains(t, output, "Sydney,AU,John Smith\n")
}

func TestWriteCSVConvertsJSONArrayFieldsToJSONString(t *testing.T) {
	jsonPayload := &JSON{
		Items: []map[string]interface{}{
			{
				"id":   float64(1),
				"tags": []interface{}{"new", "vip"},
			},
		},
		Shape: JSONObjectShape,
	}

	output := captureStdout(t, func() {
		if err := Write(jsonPayload, CSVType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	assertContains(t, output, "id,tags\n")
	assertContains(t, output, "1,\"[\"\"new\"\",\"\"vip\"\"]\"\n")
}

func TestWriteCSVIncludesFlattenedKeysFromAllJSONRecords(t *testing.T) {
	jsonPayload := &JSON{
		Items: []map[string]interface{}{
			{"id": float64(1), "name": "alice"},
			{"id": float64(2), "address": map[string]interface{}{"city": "Sydney"}},
		},
		Shape: JSONArrayShape,
	}

	output := captureStdout(t, func() {
		if err := Write(jsonPayload, CSVType); err != nil {
			t.Fatalf("Write returned error: %v", err)
		}
	})

	assertContains(t, output, "address.city,id,name\n")
	assertContains(t, output, "Sydney,2,\n")
	assertContains(t, output, ",1,alice\n")
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe returned error: %v", err)
	}
	defer reader.Close()

	os.Stdout = writer
	defer func() {
		os.Stdout = original
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("closing writer returned error: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		t.Fatalf("reading stdout returned error: %v", err)
	}

	return buf.String()
}

func assertContains(t *testing.T, value, expected string) {
	t.Helper()
	if !strings.Contains(value, expected) {
		t.Fatalf("expected output to contain %q, got %q", expected, value)
	}
}
