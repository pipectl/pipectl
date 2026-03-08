package plan

import (
	"testing"

	"github.com/shanebell/pipectl/internal/engine/steps/rename"
	"github.com/shanebell/pipectl/internal/pipeline/spec"
)

func TestBuildRenameStep(t *testing.T) {
	pipeline := spec.Pipeline{
		Steps: []spec.StepWrapper{
			{
				Step: &spec.RenameStep{
					Fields: map[string]string{
						"firstName": "first_name",
						"lastName":  "last_name",
					},
				},
			},
		},
	}

	executableSteps, err := Build(pipeline)
	if err != nil {
		t.Fatalf("build returned error: %v", err)
	}

	if len(executableSteps) != 1 {
		t.Fatalf("unexpected step count: got %d want %d", len(executableSteps), 1)
	}

	renameStep, ok := executableSteps[0].(*rename.Step)
	if !ok {
		t.Fatalf("expected *rename.Step, got %T", executableSteps[0])
	}

	expected := map[string]string{
		"firstName": "first_name",
		"lastName":  "last_name",
	}
	if len(renameStep.Fields) != len(expected) {
		t.Fatalf("unexpected fields count: got %d want %d", len(renameStep.Fields), len(expected))
	}

	for from, to := range expected {
		if got := renameStep.Fields[from]; got != to {
			t.Fatalf("unexpected field mapping for %q: got %q want %q", from, got, to)
		}
	}
}
