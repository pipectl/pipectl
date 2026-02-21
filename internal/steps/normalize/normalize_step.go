package normalize

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/steps"
)

type NormalizeStep struct {
	Fields map[string]string
}

func (s *NormalizeStep) Name() string {
	return "normalize"
}

func (s *NormalizeStep) Execute(context *steps.ExecutionContext) error {
	jsonPayload, ok := context.Payload.(*steps.JSONPayload)
	if !ok {
		return fmt.Errorf("%v requires JSON payload, got %s", s.Name(), context.Payload.Type())
	}

	for k, v := range s.Fields {
		fmt.Printf("- field: %v, action: %v\n", k, v)
	}
	fmt.Printf("- payload: %v\n", jsonPayload.Data)

	// TODO implement normalization logic

	return nil
}
