package reverse

import (
	"fmt"
	"slices"
)

type Reverse struct {
	text string
}

func New(text string) *Reverse {
	return &Reverse{text: text}
}

func (s *Reverse) Name() string {
	return "reverse"
}

func (s *Reverse) Run(input []byte) ([]byte, error) {
	if s.text != "" {
		runes := []rune(s.text)
		slices.Reverse(runes)
		return []byte(fmt.Sprintf("%v\n", string(runes))), nil
	}
	return input, nil
}
