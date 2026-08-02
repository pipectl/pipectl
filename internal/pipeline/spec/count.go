package spec

type CountStep struct {
	Message string `yaml:"message"`
}

func (s *CountStep) StepType() string {
	return "count"
}

func (s *CountStep) Validate() error { return nil }
