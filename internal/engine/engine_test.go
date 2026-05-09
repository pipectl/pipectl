package engine

import (
	"testing"

	"github.com/pipectl/pipectl/internal/engine/payload"
)

type mockStep struct {
	name string
}

func (m *mockStep) Name() string                      { return m.name }
func (m *mockStep) Supports(_ payload.Payload) bool   { return true }
func (m *mockStep) Execute(_ *ExecutionContext) error { return nil }

func TestEngine_timing_collected(t *testing.T) {
	steps := []ExecutableStep{
		&mockStep{name: "step-a"},
		&mockStep{name: "step-b"},
	}
	eng := New(steps)

	p, err := payload.Read([]byte(`[{"id":1},{"id":2}]`), "json")
	if err != nil {
		t.Fatalf("read payload: %v", err)
	}

	ctx := &ExecutionContext{Payload: p, CollectTiming: true}
	if err := eng.Run(ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(ctx.TimingResults) != 2 {
		t.Fatalf("expected 2 timing results, got %d", len(ctx.TimingResults))
	}
	if ctx.TimingResults[0].Name != "step-a" {
		t.Errorf("result[0].Name = %q, want step-a", ctx.TimingResults[0].Name)
	}
	if ctx.TimingResults[1].Name != "step-b" {
		t.Errorf("result[1].Name = %q, want step-b", ctx.TimingResults[1].Name)
	}
	for i, tr := range ctx.TimingResults {
		if tr.Duration < 0 {
			t.Errorf("result[%d].Duration = %v, want >= 0", i, tr.Duration)
		}
		if tr.RecordsIn != 2 {
			t.Errorf("result[%d].RecordsIn = %d, want 2", i, tr.RecordsIn)
		}
		if tr.RecordsOut != 2 {
			t.Errorf("result[%d].RecordsOut = %d, want 2", i, tr.RecordsOut)
		}
	}
}

func TestEngine_timing_not_collected_when_disabled(t *testing.T) {
	eng := New([]ExecutableStep{&mockStep{name: "step-a"}})

	p, err := payload.Read([]byte(`[{"id":1}]`), "json")
	if err != nil {
		t.Fatalf("read payload: %v", err)
	}

	ctx := &ExecutionContext{Payload: p, CollectTiming: false}
	if err := eng.Run(ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(ctx.TimingResults) != 0 {
		t.Errorf("expected no timing results when CollectTiming=false, got %d", len(ctx.TimingResults))
	}
}
