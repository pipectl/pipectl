package spec

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Pipeline struct {
	ID     string        `yaml:"id"`
	Input  Input         `yaml:"input"`
	Steps  []StepWrapper `yaml:"steps"`
	Output Output        `yaml:"output"`
}

type Input struct {
	Format    string `yaml:"format"`
	Schema    string `yaml:"schema,omitempty"`
	Delimiter string `yaml:"delimiter,omitempty"`
	MaxSize   int    `yaml:"max_size,omitempty"`
}

type Output struct {
	Format string `yaml:"format"`
}

func Load(path string) (Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Pipeline{}, err
	}

	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Pipeline{}, fmt.Errorf("decode pipeline: %w", err)
	}

	return p, nil
}
