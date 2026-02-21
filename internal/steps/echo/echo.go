package echo

import "fmt"

type Echo struct {
	text string
}

func New(text string) *Echo {
	return &Echo{text: text}
}

func (s *Echo) Name() string {
	return "echo"
}

func (s *Echo) Run(input []byte) ([]byte, error) {
	if s.text != "" {
		return []byte(fmt.Sprintf("%v\n", s.text)), nil
	}
	return input, nil
}
