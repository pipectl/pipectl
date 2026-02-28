package pipeline

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
	"github.com/shanebell/pipectl/internal/pipeline/plan"
	"github.com/shanebell/pipectl/internal/pipeline/spec"
)

func log(p spec.Pipeline) {
	fmt.Println("----------------")
	fmt.Printf("Pipeline: %s\n", p.ID)
	fmt.Printf("Input: %s\n", p.Input.Format)
	fmt.Println("Steps:")
	for _, step := range p.Steps {
		if step.Step == nil {
			fmt.Println("- <nil>")
			continue
		}
		fmt.Printf("- %s\n", step.Step.StepType())
	}
	fmt.Printf("Output: %s\n", p.Output.Format)
	fmt.Println("----------------")
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

	if err := pipelineEngine.Run(ctx); err != nil {
		return err
	}

	if err := payload.Write(ctx.Payload, p.Output.Format); err != nil {
		return err
	}

	return nil
}
