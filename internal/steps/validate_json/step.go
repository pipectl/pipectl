package validate_json

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/steps"
)

type Step struct {
	Schema string
}

func (s *Step) Name() string {
	return "validate-json"
}

func (s *Step) Supports(payload steps.Payload) bool {
	return payload.Type() == "json"
}

func (s *Step) Execute(context *steps.ExecutionContext) error {
	jsonPayload, ok := context.Payload.(*steps.JSONPayload)
	if !ok {
		return fmt.Errorf("%v requires JSON payload, got %s", s.Name(), context.Payload.Type())
	}

	fmt.Printf("- schema: %v\n", s.Schema)
	fmt.Printf("- payload: %v\n", jsonPayload.Data)

	// TODO implement JSON schema validation

	return nil
}
