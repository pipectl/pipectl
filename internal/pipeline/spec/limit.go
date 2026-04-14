package spec

import "fmt"

type LimitStep struct {
	Count int `yaml:"count"`
}

func (s *LimitStep) StepType() string {
	return "limit"
}

func (s *LimitStep) String() string {
	return fmt.Sprintf("[%s] count=%d", s.StepType(), s.Count)
}

func (s *LimitStep) Validate() error {
	if s.Count < 1 {
		return fmt.Errorf("limit count must be at least 1")
	}
	return nil
}
