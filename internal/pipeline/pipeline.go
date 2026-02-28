package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
	"github.com/shanebell/pipectl/internal/pipeline/plan"
	"github.com/shanebell/pipectl/internal/pipeline/spec"
)

func log(p spec.Pipeline) {
	fmt.Printf("Executing pipeline\n")
	fmt.Printf("- ID: %s\n", p.ID)
	fmt.Printf("- Input: %s\n", p.Input.Format)
	fmt.Println("- Steps:")
	for i, step := range p.Steps {
		fmt.Printf("  %v. %s\n", i+1, step.Step.StepType())
	}
	fmt.Printf("- Output: %s\n", p.Output.Format)
}

func Run(path string, input []byte) error {
	p, err := spec.Load(path)
	if err != nil {
		return err
	}

	log(p)

	executableSteps, err := plan.Build(p)
	if err != nil {
		return err
	}

	pipelineEngine := engine.New(executableSteps)

	inputPayload, err := payload.Read(input, p.Input.Format)
	if err != nil {
		return err
	}

	ctx := &engine.ExecutionContext{Payload: inputPayload}

	fmt.Printf("\nRunning steps...\n")
	if err := pipelineEngine.Run(ctx); err != nil {
		return err
	}

	fmt.Printf("\nOutput:\n")
	if err := payload.Write(ctx.Payload, p.Output.Format); err != nil {
		return err
	}

	return nil
}
