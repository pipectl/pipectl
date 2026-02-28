package engine

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/payload"
)

type ExecutableStep interface {
	Execute(ctx *ExecutionContext) error
	Supports(payload payload.Payload) bool
	Name() string
}

type Engine struct {
	steps []ExecutableStep
}

func New(steps []ExecutableStep) *Engine {
	return &Engine{
		steps: steps,
	}
}

func (e *Engine) Run(ctx *ExecutionContext) error {
	for _, step := range e.steps {

		fmt.Printf("\nExecuting step [%s]\n", step.Name())

		if !step.Supports(ctx.Payload) {
			return fmt.Errorf("step [%s] does not support payload type [%s]", step.Name(), ctx.Payload.Type())
		}

		if err := step.Execute(ctx); err != nil {
			return fmt.Errorf("step [%s] failed: %w", step.Name(), err)
		}
	}

	return nil
}
