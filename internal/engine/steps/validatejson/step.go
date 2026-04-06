package validatejson

import (
	"fmt"
	"os"
	"strings"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
	"github.com/xeipuuv/gojsonschema"
)

type Step struct {
	Schema string
}

func (s *Step) Name() string {
	return "validate-json"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case payload.JSONRecordPayload:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	context.Logger.Debug("  schema: %s", s.Schema)

	schemaLoader, err := s.schemaLoader()
	if err != nil {
		return err
	}

	switch value := context.Payload.(type) {
	case *payload.JSONL:
		return s.validateJSONLRecords(schemaLoader, value.Records())
	case payload.JSONRecordPayload:
		return s.validateJSONPayload(schemaLoader, value.Value())
	default:
		return fmt.Errorf("%v received invalid payload type %v", s.Name(), context.Payload.Type())
	}
}

func (s *Step) validateJSONPayload(schemaLoader gojsonschema.JSONLoader, data interface{}) error {
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

func (s *Step) validateJSONLRecords(schemaLoader gojsonschema.JSONLoader, records []map[string]interface{}) error {
	for i, record := range records {
		if err := s.validateJSONPayload(schemaLoader, record); err != nil {
			return fmt.Errorf("JSONL record %d: %w", i+1, err)
		}
	}

	return nil
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
