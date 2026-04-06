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
	Delimiter string `yaml:"delimiter,omitempty"`
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

	if p.Input.Delimiter != "" && len([]rune(p.Input.Delimiter)) != 1 {
		return Pipeline{}, fmt.Errorf("input delimiter must be a single character")
	}

	return p, nil
}
