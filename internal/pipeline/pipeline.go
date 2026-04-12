package pipeline

import (
	"io"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
	"github.com/pipectl/pipectl/internal/pipeline/plan"
	"github.com/pipectl/pipectl/internal/pipeline/spec"
)

func Run(path string, input []byte, output io.Writer, verbose bool, dryRun bool) error {
	p, err := spec.Load(path)
	if err != nil {
		return err
	}

	logger := engine.NewLogger(verbose)

	logger.Debug("pipeline: %s [%s → %s, %d steps]", p.ID, p.Input.Format, p.Output.Format, len(p.Steps))
	for i, step := range p.Steps {
		logger.Debug("  %d. %s", i+1, step.Step.StepType())
	}

	executableSteps, err := plan.Build(p)
	if err != nil {
		return err
	}

	if dryRun {
		logger.Log("dry run: %d steps would execute", len(executableSteps))
		for i, step := range executableSteps {
			logger.Log("  %d. %s", i+1, step.Name())
		}
		return nil
	}

	pipelineEngine := engine.New(executableSteps)

	var inputPayload payload.Payload
	if p.Input.Format == payload.CSVType {
		var delimiter rune
		if p.Input.Delimiter != "" {
			delimiter = []rune(p.Input.Delimiter)[0]
		}
		inputPayload, err = payload.ReadCSV(input, delimiter)
	} else {
		inputPayload, err = payload.Read(input, p.Input.Format)
	}
	if err != nil {
		return err
	}

	ctx := &engine.ExecutionContext{Payload: inputPayload, Logger: logger}

	if err := pipelineEngine.Run(ctx); err != nil {
		return err
	}

	return payload.Write(ctx.Payload, p.Output.Format, output)
}
