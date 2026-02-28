package pipeline

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
	"github.com/shanebell/pipectl/internal/steps"
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

var stepRegistry = map[string]Step{
	"validate-json":  &ValidateJSONStep{},
	"normalize":      &NormalizeStep{},
	"redact":         &RedactStep{},
	"filter":         &FilterStep{},
	"http-transform": &HTTPTransformStep{},
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
		var step, ok = stepRegistry[key]
		if !ok {
			return fmt.Errorf("unknown step type: %s", key)
		}
		if err := yaml.Unmarshal(value, step); err != nil {
			return err
		}
		w.Step = step
	}

	return nil
}

func LoadPayload(input []byte, format string) (payload.Payload, error) {
	switch format {

	case "json":
		var data map[string]interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return nil, err
		}
		return &payload.JSON{Data: data}, nil

	case "csv":
		rows, err := csv.NewReader(bytes.NewReader(input)).ReadAll()
		if err != nil {
			panic(err)
		}
		return &payload.CSV{Rows: rows}, nil

	case "text":
		return &payload.Text{Text: string(input)}, nil

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

	stepPayload, err := LoadPayload(input, pipeline.Input.Format)
	if err != nil {
		return err
	}

	context := &engine.ExecutionContext{
		Payload: stepPayload,
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
			jsonPayload, _ := context.Payload.(*payload.JSON)
			output, err := json.MarshalIndent(jsonPayload.Data, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return nil
			}
			fmt.Println(string(output))

		case "csv":
			csvPayload, _ := context.Payload.(*payload.CSV)
			// TODO how do we convert CSV to JSON?
			fmt.Println("TODO: convert CSV to JSON")
			fmt.Println(csvPayload.Rows)

		default:
			return fmt.Errorf("Cannot convert to JSON")
		}

	} else if pipeline.Output.Format == "csv" {
		switch context.Payload.Type() {
		case "csv":
			csvPayload, _ := context.Payload.(*payload.CSV)
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
