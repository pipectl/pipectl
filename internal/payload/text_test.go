package payload

import "testing"

func TestTextType(t *testing.T) {
	textPayload := &Text{Text: "hello"}

	if got := textPayload.Type(); got != TextType {
		t.Fatalf("expected payload type %q, got %q", TextType, got)
	}
}
