package count

import (
	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type Step struct {
	Message string
}

func (s *Step) Name() string {
	return "count"
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
	context.Logger.Log("  records: %d", context.Payload.RecordCount())
	return nil
}
