package pipeline

import (
	"io"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
	"github.com/shanebell/pipectl/internal/pipeline/plan"
	"github.com/shanebell/pipectl/internal/pipeline/spec"
)

func Run(path string, input []byte, output io.Writer, verbose bool) error {
	p, err := spec.Load(path)
	if err != nil {
		return err
	}

	logger := engine.NewLogger(verbose)

	logger.Log("pipeline: %s [%s → %s, %d steps]", p.ID, p.Input.Format, p.Output.Format, len(p.Steps))
	for i, step := range p.Steps {
		logger.Debug("  %d. %s", i+1, step.Step.StepType())
	}

	executableSteps, err := plan.Build(p)
	if err != nil {
		return err
	}

	pipelineEngine := engine.New(executableSteps)

	inputPayload, err := payload.Read(input, p.Input.Format)
	if err != nil {
		return err
	}

	ctx := &engine.ExecutionContext{Payload: inputPayload, Logger: logger}

	if err := pipelineEngine.Run(ctx); err != nil {
		return err
	}

	return payload.Write(ctx.Payload, p.Output.Format, output)
}
