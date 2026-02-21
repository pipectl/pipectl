package pipeline

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/shanebell/pipectl/internal/steps/reverse"

	"github.com/shanebell/pipectl/internal/steps/echo"
)

type Schema struct {
	Id    string       `yaml:"id"`
	Steps []StepSchema `yaml:"steps,omitempty"`
}

// QUESTION: how to handle different types of steps?
// Different struct for each one?
// One giant one?
type StepSchema struct {
	Type string `yaml:"type"`
	Text string `yaml:"text"`
}

func (s StepSchema) String() string {
	return fmt.Sprintf("Step [%v]: %v", s.Type, s.Text)
}

func RunFromFile(path string, input []byte) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var spec Schema
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return err
	}

	// DEBUG INFO
	fmt.Println("----------------")
	fmt.Printf("Pipeline: %s\n", spec.Id)
	for _, s := range spec.Steps {
		fmt.Printf("- %v\n", s)
	}
	fmt.Println("----------------")
	fmt.Println()
	// END DEBUG INFO

	var steps []Step

	// build a list of steps and validate along the way
	for _, s := range spec.Steps {
		switch s.Type {
		case "echo":
			steps = append(steps, echo.New(s.Text))

		case "reverse":
			steps = append(steps, reverse.New(s.Text))

		default:
			// invalid step type, exit

			return fmt.Errorf("unknown step type: %s", s.Type)
		}
	}

	// execute each step
	for _, step := range steps {
		fmt.Printf("Running step [%s]\n", step.Name())
		output, err := step.Run(input)
		if err != nil {
			return fmt.Errorf("step %s failed: %w", step.Name(), err)
		}
		fmt.Printf("Output: %s\n", output)

		// pass the output to the next step
		input = output
	}

	// Write final output to stdout
	//os.Stdout.Write(input)
	//os.Stdout.Write([]byte("\n"))

	return nil
}
