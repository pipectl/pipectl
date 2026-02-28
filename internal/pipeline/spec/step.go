package spec

type Step interface {
	StepType() string
}

type StepWrapper struct {
	Step Step
}
