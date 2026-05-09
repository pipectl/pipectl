package engine

import (
	"time"

	"github.com/pipectl/pipectl/internal/engine/payload"
)

type StepTiming struct {
	Name       string
	Duration   time.Duration
	RecordsIn  int
	RecordsOut int
}

type ExecutionContext struct {
	Payload       payload.Payload
	Logger        *Logger
	CollectTiming bool
	TimingResults []StepTiming
}
