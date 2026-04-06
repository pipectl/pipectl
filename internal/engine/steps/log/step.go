package _log

import (
	"encoding/json"
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
	case payload.JSONRecordPayload, *payload.CSV:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	if s.Message != "" {
		context.Logger.Log("  message: %s", s.Message)
	}

	if s.Count {
		context.Logger.Log("  records: %d", context.Payload.RecordCount())
	}

	s.printSample(context.Logger, context.Payload)
	return nil
}

func (s *Step) printSample(logger *engine.Logger, p payload.Payload) {
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
			logger.Log("  sample (0):")
			logger.Log("    %s", strings.Join(v.Rows[0], ","))
			return
		}
		rows := v.Rows[1:]
		if len(rows) > limit {
			rows = rows[:limit]
		}
		logger.Log("  sample (%d):", len(rows))
		logger.Log("    %s", strings.Join(v.Rows[0], ","))
		for _, row := range rows {
			logger.Log("    %s", strings.Join(row, ","))
		}
	case payload.JSONRecordPayload:
		records := v.Records()
		if len(records) == 0 {
			return
		}
		if len(records) > limit {
			records = records[:limit]
		}
		logger.Log("  sample (%d):", len(records))
		for _, record := range records {
			raw, err := json.Marshal(record)
			if err != nil {
				logger.Log("    %v", record)
				continue
			}
			logger.Log("    %s", string(raw))
		}
	}
}
