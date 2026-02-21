package pipeline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/shanebell/pipectl/internal/steps"
	"github.com/shanebell/pipectl/internal/steps/normalize"
	"github.com/shanebell/pipectl/internal/steps/redact"
	validate_json "github.com/shanebell/pipectl/internal/steps/validate-json"
)

type Pipeline struct {
	ID    string        `yaml:"id"`
	Steps []StepWrapper `yaml:"steps"`
}

type Step interface {
	StepType() string
	BuildExecutor() (steps.ExecutableStep, error)
}

type StepWrapper struct {
	Step Step
}

type ValidateJSONStep struct {
	Schema string `yaml:"schema"`
}

func (s *ValidateJSONStep) StepType() string {
	return "validate-json"
}

func (s *ValidateJSONStep) BuildExecutor() (steps.ExecutableStep, error) {
	return &validate_json.ValidateJSONStep{
		Schema: s.Schema,
	}, nil
}

func (s *ValidateJSONStep) String() string {
	return fmt.Sprintf("[%s] schema: %v", s.StepType(), s.Schema)
}

type NormalizeStep struct {
	Fields map[string]string `yaml:"fields"`
}

func (s *NormalizeStep) StepType() string {
	return "normalize"
}

func (s *NormalizeStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *NormalizeStep) BuildExecutor() (steps.ExecutableStep, error) {
	return &normalize.NormalizeStep{
		Fields: s.Fields,
	}, nil
}

type RedactStep struct {
	Fields []string `yaml:"fields"`
}

func (s *RedactStep) StepType() string {
	return "redact"
}

func (s *RedactStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *RedactStep) BuildExecutor() (steps.ExecutableStep, error) {
	return &redact.RedactStep{
		Fields: s.Fields,
	}, nil
}

// custom unmarshal for different steps
func (w *StepWrapper) UnmarshalYAML(b []byte) error {
	var raw map[string]yaml.RawMessage
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return err
	}

	if len(raw) != 1 {
		return fmt.Errorf("step must contain exactly one key")
	}

	for key, value := range raw {
		switch key {

		case "validate-json":
			var step ValidateJSONStep
			if err := yaml.Unmarshal(value, &step); err != nil {
				return err
			}
			w.Step = &step

		case "normalize":
			var step NormalizeStep
			if err := yaml.Unmarshal(value, &step); err != nil {
				return err
			}
			w.Step = &step

		case "redact":
			var step RedactStep
			if err := yaml.Unmarshal(value, &step); err != nil {
				return err
			}
			w.Step = &step

		default:
			return fmt.Errorf("unknown step type: %s", key)
		}
	}

	return nil
}

func LoadPayload(input []byte, format string) (steps.Payload, error) {
	switch format {

	case "json":
		var data map[string]interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return nil, err
		}
		return &steps.JSONPayload{Data: data}, nil

	case "text":
		return &steps.TextPayload{Text: string(input)}, nil

	default:
		return nil, fmt.Errorf("unsupported input format")
	}
}

func RunFromFile(path string, input []byte) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// parse the raw yaml
	var raw map[string]yaml.RawMessage
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}

	var pipeline Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return err
	}

	// DEBUG INFO
	fmt.Println("----------------")
	fmt.Printf("Pipeline: %s\n", pipeline.ID)
	for _, s := range pipeline.Steps {
		fmt.Printf("- %v\n", s.Step)
	}
	fmt.Println("----------------")
	fmt.Println()
	// END DEBUG INFO

	var executableSteps []steps.ExecutableStep

	// build a list of steps and validate along the way
	for _, s := range pipeline.Steps {
		executor, err := s.Step.BuildExecutor()
		if err != nil {
			return err
		}
		executableSteps = append(executableSteps, executor)
	}

	// TODO how to determine input type?
	payload, err := LoadPayload(input, "json")
	if err != nil {
		return err
	}

	context := &steps.ExecutionContext{
		Payload: payload,
	}

	// execute each step
	for _, executableStep := range executableSteps {
		fmt.Printf("\nExecuting step [%s]\n", executableStep.Name())
		if err := executableStep.Execute(context); err != nil {
			return err
		}
	}

	// TODO pipeline should definte output type
	output, err := json.MarshalIndent(context.Payload, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}
	fmt.Println(string(output))

	return nil
}
