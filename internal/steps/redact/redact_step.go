package redact

import (
	"fmt"
	"slices"

	"github.com/shanebell/pipectl/internal/steps"
)

type RedactStep struct {
	Fields []string
}

func (s *RedactStep) Name() string {
	return "redact"
}

func (s *RedactStep) Execute(context *steps.ExecutionContext) error {
	jsonPayload, ok := context.Payload.(*steps.JSONPayload)
	if !ok {
		return fmt.Errorf("%v requires JSON payload, got %s", s.Name(), context.Payload.Type())
	}

	for _, f := range s.Fields {
		fmt.Printf("- field: %v\n", f)
	}
	fmt.Printf("- payload: %v\n", jsonPayload.Data)

	for k, v := range jsonPayload.Data {
		if slices.Contains(s.Fields, k) {
			fmt.Printf("- redacting field: %v, value: %v\n", k, v)
			jsonPayload.Data[k] = "REDACTED"
		}
	}

	return nil
}
