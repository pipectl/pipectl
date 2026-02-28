package validatejson

import (
	"fmt"
	"os"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
	"github.com/xeipuuv/gojsonschema"
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

	fmt.Printf("Validating JSON payload against schema %v\n", s.Schema)
	return s.validateJSONPayload(jsonPayload.Data)
}

func (s *Step) validateJSONPayload(data map[string]interface{}) error {
	schemaLoader, err := s.schemaLoader()
	if err != nil {
		return err
	}

	result, err := gojsonschema.Validate(schemaLoader, gojsonschema.NewGoLoader(data))
	if err != nil {
		return fmt.Errorf("unable to validate JSON payload: %w", err)
	}

	if result.Valid() {
		return nil
	}

	validationErrors := make([]string, 0, len(result.Errors()))
	for _, validationErr := range result.Errors() {
		validationErrors = append(validationErrors, validationErr.String())
	}

	return fmt.Errorf("JSON schema validation failed: %s", strings.Join(validationErrors, "; "))
}

func (s *Step) schemaLoader() (gojsonschema.JSONLoader, error) {
	if strings.TrimSpace(s.Schema) == "" {
		return nil, fmt.Errorf("schema is required")
	}

	trimmedSchema := strings.TrimSpace(s.Schema)
	if strings.HasPrefix(trimmedSchema, "{") || strings.HasPrefix(trimmedSchema, "[") {
		return gojsonschema.NewStringLoader(trimmedSchema), nil
	}

	schemaBytes, err := os.ReadFile(trimmedSchema)
	if err != nil {
		return nil, fmt.Errorf("unable to read schema file %q: %w", trimmedSchema, err)
	}

	return gojsonschema.NewBytesLoader(schemaBytes), nil
}
