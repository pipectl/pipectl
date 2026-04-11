package engine

import (
	"fmt"

	"github.com/shanebell/pipectl/internal/engine/payload"
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
	for i, step := range e.steps {
		ctx.Logger.Debug("\n%d. [%s]", i+1, step.Name())

		if !step.Supports(ctx.Payload) {
			return fmt.Errorf("step %d [%s] does not support payload type [%s]", i+1, step.Name(), ctx.Payload.Type())
		}

		if err := step.Execute(ctx); err != nil {
			return fmt.Errorf("step %d [%s] failed: %w", i+1, step.Name(), err)
		}
	}

	return nil
}
