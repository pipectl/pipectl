package payload

import "testing"

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
