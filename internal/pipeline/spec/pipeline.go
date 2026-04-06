package spec

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

var validFormats = []string{"json", "jsonl", "csv"}

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

	if !isValidFormat(p.Input.Format) {
		return Pipeline{}, fmt.Errorf("input format must be one of: json, jsonl, csv")
	}

	if p.Input.Delimiter != "" && len([]rune(p.Input.Delimiter)) != 1 {
		return Pipeline{}, fmt.Errorf("input delimiter must be a single character")
	}

	if !isValidFormat(p.Output.Format) {
		return Pipeline{}, fmt.Errorf("output format must be one of: json, jsonl, csv")
	}

	return p, nil
}

func isValidFormat(f string) bool {
	for _, v := range validFormats {
		if f == v {
			return true
		}
	}
	return false
}
