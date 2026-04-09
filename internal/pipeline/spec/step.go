package spec

type Step interface {
	StepType() string
}

// Validator is implemented by steps that require validation beyond YAML
// unmarshaling. StepWrapper calls Validate after unmarshaling so that steps
// with required configuration are caught even when the YAML value is null
// (goccy/go-yaml skips custom UnmarshalYAML methods for null values).
type Validator interface {
	Validate() error
}

type StepWrapper struct {
	Step Step
}
