package _log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Message string
	Count   bool
	Sample  int
}

func (s *Step) Name() string {
	return "log"
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

	if s.Count {
		fmt.Printf("- records: %d\n", s.recordCount(context.Payload))
	}

	s.printSample(context.Payload)
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

func (s *Step) printSample(p payload.Payload) {
	limit := s.Sample
	if limit < 0 {
		limit = 0
	}
	if limit == 0 {
		return
	}

	switch v := p.(type) {
	case *payload.CSV:
		if len(v.Rows) == 0 {
			return
		}
		if len(v.Rows) <= 1 {
			fmt.Printf("- sample (%d):\n", 0)
			fmt.Println(strings.Join(v.Rows[0], ","))
			return
		}
		rows := v.Rows[1:]
		if len(rows) > limit {
			rows = rows[:limit]
		}
		fmt.Printf("- sample (%d):\n", len(rows))
		fmt.Println(strings.Join(v.Rows[0], ","))
		for _, row := range rows {
			fmt.Println(strings.Join(row, ","))
		}
	case *payload.JSON:
		if len(v.Data) == 0 {
			return
		}
		fmt.Printf("- sample (1):\n")
		raw, err := json.Marshal(v.Data)
		if err != nil {
			fmt.Printf("%v\n", v.Data)
			return
		}
		fmt.Println(string(raw))
	case *payload.Text:
		lines := nonEmptyLines(v.Text)
		if len(lines) == 0 {
			return
		}
		if len(lines) > limit {
			lines = lines[:limit]
		}
		fmt.Printf("- sample (%d):\n", len(lines))
		for _, line := range lines {
			fmt.Println(line)
		}
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
