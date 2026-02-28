package payload

import "testing"

func TestJSONType(t *testing.T) {
	jsonPayload := &JSON{
		Data: map[string]interface{}{
			"id":   1,
			"name": "alice",
		},
	}

	if got := jsonPayload.Type(); got != JSONType {
		t.Fatalf("expected payload type %q, got %q", JSONType, got)
	}
}
