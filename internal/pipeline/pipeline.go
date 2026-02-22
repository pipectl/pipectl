package pipeline

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/shanebell/pipectl/internal/steps"
	"github.com/shanebell/pipectl/internal/steps/normalize"
	"github.com/shanebell/pipectl/internal/steps/redact"
	"github.com/shanebell/pipectl/internal/steps/validate_json"
)

type Pipeline struct {
	ID     string        `yaml:"id"`
	Input  Input         `yaml:"input"`
	Steps  []StepWrapper `yaml:"steps"`
	Output Output        `yaml:"output"`
}

type Input struct {
	Format    string `yaml:"format"`
	Encoding  string `yaml:"encoding:omitempty"`
	Schema    string `yaml:"schema:omitempty"`
	Delimiter string `yaml:"delimiter:omitempty"`
	HasHeader bool   `yaml:"has_header:omitempty"`
	MaxSize   int    `yaml:"max_size:omitempty"`
}

type Output struct {
	Format string `yaml:"format"`
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
	Strategy string   `yaml:"strategy"`
	Fields   []string `yaml:"fields"`
}

func (s *RedactStep) StepType() string {
	return "redact"
}

func (s *RedactStep) String() string {
	return fmt.Sprintf("[%s] fields: %v", s.StepType(), s.Fields)
}

func (s *RedactStep) BuildExecutor() (steps.ExecutableStep, error) {
	return &redact.RedactStep{
		Fields:   s.Fields,
		Strategy: s.Strategy,
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

	case "csv":
		rows, err := csv.NewReader(bytes.NewReader(input)).ReadAll()
		if err != nil {
			panic(err)
		}
		return &steps.CSVPayload{Rows: rows}, nil

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

	var pipeline Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return err
	}

	// DEBUG INFO
	fmt.Println("----------------")
	fmt.Printf("Pipeline: %s\n", pipeline.ID)
	fmt.Printf("Input: %s\n", pipeline.Input.Format)
	fmt.Println("Steps:")
	for _, s := range pipeline.Steps {
		fmt.Printf("- %v\n", s.Step)
	}
	fmt.Println("----------------")
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

	payload, err := LoadPayload(input, pipeline.Input.Format)
	if err != nil {
		return err
	}

	context := &steps.ExecutionContext{
		Payload: payload,
	}

	// execute each step
	for _, executableStep := range executableSteps {
		fmt.Printf("\nExecuting step [%s]\n", executableStep.Name())

		if !executableStep.Supports(context.Payload) {
			return fmt.Errorf("step [%s] does not support payload type [%s]", executableStep.Name(), context.Payload.Type())
		}

		if err := executableStep.Execute(context); err != nil {
			return err
		}
	}

	// produce final output
	// TODO move this into a separate function and convert to a switch statement
	fmt.Println("\nOutput:")
	if pipeline.Output.Format == "json" {

		// TODO: which payload types can be converted to JSON?
		switch context.Payload.Type() {

		case "json":
			jsonPayload, _ := context.Payload.(*steps.JSONPayload)
			output, err := json.MarshalIndent(jsonPayload.Data, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return nil
			}
			fmt.Println(string(output))

		case "csv":
			csvPayload, _ := context.Payload.(*steps.CSVPayload)
			// TODO how do we convert CSV to JSON?
			fmt.Println("TODO: convert CSV to JSON")
			fmt.Println(csvPayload.Rows)

		default:
			return fmt.Errorf("Cannot convert to JSON")
		}

	} else if pipeline.Output.Format == "csv" {
		switch context.Payload.Type() {
		case "csv":
			csvPayload, _ := context.Payload.(*steps.CSVPayload)
			buf := new(bytes.Buffer)
			csvWriter := csv.NewWriter(buf)
			if err := csvWriter.WriteAll(csvPayload.Rows); err != nil {
				fmt.Println("Error writing CSV:", err)
				return nil
			}
			fmt.Println(buf.String())

		case "json":
			// TODO convert JSON to CSV

		default:
			return fmt.Errorf("Cannot convert to CSV")
		}

	}

	return nil
}
