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
	switch p.(type) {
	case payload.RecordPayload, *payload.CSV:
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
		fmt.Printf("- records: %d\n", context.Payload.RecordCount())
	}

	s.printSample(context.Payload)
	return nil
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
	case payload.RecordPayload:
		records := v.GetRecords()
		if len(records) == 0 {
			return
		}
		if len(records) > limit {
			records = records[:limit]
		}
		fmt.Printf("- sample (%d):\n", len(records))
		for _, record := range records {
			raw, err := json.Marshal(record)
			if err != nil {
				fmt.Printf("%v\n", record)
				continue
			}
			fmt.Println(string(raw))
		}
	}
}
