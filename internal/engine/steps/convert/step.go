package convert

import (
	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type Step struct {
	Format string
}

func (s *Step) Name() string {
	return "convert"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case *payload.JSON, *payload.JSONL, *payload.CSV:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(ctx *engine.ExecutionContext) error {
	ctx.Logger.Debug("  format: %s", s.Format)

	converted, err := payload.Convert(ctx.Payload, s.Format)
	if err != nil {
		return err
	}

	ctx.Payload = converted
	return nil
}
