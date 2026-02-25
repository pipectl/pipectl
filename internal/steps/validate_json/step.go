package validate_json

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
)

type Step struct {
	Schema string
}

func (s *Step) Name() string {
	return "validate-json"
}

func (s *Step) Supports(p payload.Payload) bool {
	return p.Type() == payload.JSONType
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	jsonPayload, ok := context.Payload.(*payload.JSON)
	if !ok {
		return fmt.Errorf("%v requires JSON payload, got %s", s.Name(), context.Payload.Type())
	}

	fmt.Printf("- schema: %v\n", s.Schema)
	fmt.Printf("- payload: %v\n", jsonPayload.Data)

	// TODO implement JSON schema validation

	return nil
}
