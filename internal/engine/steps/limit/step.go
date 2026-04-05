package limit

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	Count int
}

func (s *Step) Name() string {
	return "limit"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case payload.JSONRecordPayload, *payload.CSV:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	before := context.Payload.RecordCount()

	switch p := context.Payload.(type) {
	case *payload.JSON:
		if len(p.Items) > s.Count {
			p.Items = p.Items[:s.Count]
		}
	case *payload.JSONL:
		if len(p.Items) > s.Count {
			p.Items = p.Items[:s.Count]
		}
	case *payload.CSV:
		// Rows[0] is the header row; data rows start at index 1
		if p.RecordCount() > s.Count {
			p.Rows = p.Rows[:s.Count+1]
		}
	}

	after := context.Payload.RecordCount()
	if before > after {
		fmt.Printf("- limited %d records to %d\n", before, after)
	} else {
		fmt.Printf("- %d records (limit of %d not reached)\n", after, s.Count)
	}

	return nil
}
