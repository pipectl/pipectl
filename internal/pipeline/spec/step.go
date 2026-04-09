package spec

type Step interface {
	StepType() string
	// Validate checks that the step is correctly configured. It is called by
	// StepWrapper after unmarshaling, including when the YAML value is null,
	// so every step must implement it even if it has no required configuration.
	Validate() error
}

type StepWrapper struct {
	Step Step
}
