package validate_json

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/steps"
)

type ValidateJSONStep struct {
	Schema string
}

func (s *ValidateJSONStep) Name() string {
	return "validate-json"
}

func (s *ValidateJSONStep) Supports(payload steps.Payload) bool {
	return payload.Type() == "json"
}

func (s *ValidateJSONStep) Execute(context *steps.ExecutionContext) error {
	jsonPayload, ok := context.Payload.(*steps.JSONPayload)
	if !ok {
		return fmt.Errorf("%v requires JSON payload, got %s", s.Name(), context.Payload.Type())
	}

	fmt.Printf("- schema: %v\n", s.Schema)
	fmt.Printf("- payload: %v\n", jsonPayload.Data)

	// TODO implement JSON schema validation

	return nil
}
