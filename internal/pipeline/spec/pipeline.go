package spec

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

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

var varToken = regexp.MustCompile(`\$\{([^}]+)\}`)

func substituteVars(data []byte, vars map[string]string) ([]byte, error) {
	if len(vars) > 0 {
		args := make([]string, 0, len(vars)*2)
		for k, v := range vars {
			args = append(args, "${"+k+"}", v)
		}
		data = []byte(strings.NewReplacer(args...).Replace(string(data)))
	}

	matches := varToken.FindAllStringSubmatch(string(data), -1)
	if len(matches) > 0 {
		seen := map[string]bool{}
		var names []string
		for _, m := range matches {
			if !seen[m[1]] {
				seen[m[1]] = true
				names = append(names, m[1])
			}
		}
		sort.Strings(names)
		return nil, fmt.Errorf("unresolved variables: %s", strings.Join(names, ", "))
	}

	return data, nil
}

func Load(path string, vars map[string]string) (Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Pipeline{}, err
	}

	data, err = substituteVars(data, vars)
	if err != nil {
		return Pipeline{}, err
	}

	var p Pipeline
	if err := yaml.NewDecoder(bytes.NewReader(data), yaml.DisallowUnknownField()).Decode(&p); err != nil {
		return Pipeline{}, fmt.Errorf("decode pipeline: %w", err)
	}

	if p.ID == "" {
		return Pipeline{}, fmt.Errorf("pipeline id must be specified")
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

	if len(p.Steps) == 0 {
		return Pipeline{}, fmt.Errorf("pipeline must have at least one step")
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
