package main

import "testing"

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "single sentence paragraph",
			content: "# step-name\n\nDoes a thing.\n\nMore detail.\n",
			want:    "Does a thing.",
		},
		{
			name:    "truncates to first sentence of a multi-sentence paragraph",
			content: "# assert\n\nChecks record-count, field-existence, and field-value conditions. The pipeline fails if any assertion is not met.\n",
			want:    "Checks record-count, field-existence, and field-value conditions.",
		},
		{
			name:    "no body",
			content: "# step-name\n",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDescription(tt.content)
			if got != tt.want {
				t.Errorf("extractDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFirstSentence(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"single sentence", "Converts field values to a specified type.", "Converts field values to a specified type."},
		{"multiple sentences", "Sends the payload. Continues unchanged.", "Sends the payload."},
		{"no terminating period", "Truncates the payload to at most N records", "Truncates the payload to at most N records"},
		{"period without trailing space is not a boundary", "Sends to example.com and continues.", "Sends to example.com and continues."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstSentence(tt.in)
			if got != tt.want {
				t.Errorf("firstSentence(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
