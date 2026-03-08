package count

import (
	"fmt"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Message string
}

func (s *Step) Name() string {
	return "count"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.Type() {
	case payload.JSONType, payload.CSVType, payload.TextType:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	if s.Message != "" {
		fmt.Printf("- message: %s\n", s.Message)
	}

	fmt.Printf("- records: %d\n", s.recordCount(context.Payload))
	return nil
}

func (s *Step) recordCount(p payload.Payload) int {
	switch v := p.(type) {
	case *payload.CSV:
		if len(v.Rows) == 0 {
			return 0
		}
		return len(v.Rows) - 1
	case *payload.JSON:
		if len(v.Data) == 0 {
			return 0
		}
		return 1
	case *payload.Text:
		return len(nonEmptyLines(v.Text))
	default:
		return 0
	}
}

func nonEmptyLines(text string) []string {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.TrimSuffix(normalized, "\n")
	if normalized == "" {
		return nil
	}

	lines := strings.Split(normalized, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		filtered = append(filtered, line)
	}

	return filtered
}
