package steps

import (
	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
)

type ExecutableStep interface {
	Execute(ctx *engine.ExecutionContext) error
	Supports(payload payload.Payload) bool
	Name() string
}
