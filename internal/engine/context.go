package engine

import (
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type ExecutionContext struct {
	Payload payload.Payload
	Logger  *Logger
}
