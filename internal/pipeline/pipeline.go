package pipeline

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
)

type Pipeline struct {
	ID     string        `yaml:"id"`
	Input  Input         `yaml:"input"`
	Steps  []StepWrapper `yaml:"steps"`
	Output Output        `yaml:"output"`
}

type Input struct {
	Format    string `yaml:"format"`
	Encoding  string `yaml:"encoding,omitempty"`
	Schema    string `yaml:"schema,omitempty"`
	Delimiter string `yaml:"delimiter,omitempty"`
	HasHeader bool   `yaml:"has_header,omitempty"`
	MaxSize   int    `yaml:"max_size,omitempty"`
}

type Output struct {
	Format string `yaml:"format"`
}

type Step interface {
	StepType() string
	BuildExecutor() (engine.ExecutableStep, error)
}

type StepWrapper struct {
	Step Step
}

func log(pipeline Pipeline) {
	fmt.Println("----------------")
	fmt.Printf("Pipeline: %s\n", pipeline.ID)
	fmt.Printf("Input: %s\n", pipeline.Input.Format)
	fmt.Println("Steps:")
	for _, s := range pipeline.Steps {
		fmt.Printf("- %v\n", s.Step)
	}
	fmt.Printf("Output: %s\n", pipeline.Output.Format)
	fmt.Println("----------------")
}

func Run(path string, input []byte) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// read the pipline from YAML
	var pipeline Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return err
	}

	log(pipeline)

	// build the pipeline steps
	var executableSteps []engine.ExecutableStep
	for _, s := range pipeline.Steps {
		executor, err := s.Step.BuildExecutor()
		if err != nil {
			return err
		}
		executableSteps = append(executableSteps, executor)
	}

	pipelineEngine := engine.New(executableSteps)

	inputPayload, err := payload.Read(input, pipeline.Input.Format)
	if err != nil {
		return err
	}

	context := &engine.ExecutionContext{Payload: inputPayload}

	if err := pipelineEngine.Run(context); err != nil {
		return err
	}

	if err := payload.Write(context.Payload, pipeline.Output.Format); err != nil {
		return err
	}

	return nil
}
