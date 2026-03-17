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
