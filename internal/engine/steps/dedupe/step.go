package dedupe

import (
	"fmt"
	"strings"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type Step struct {
	Fields        []string
	CaseSensitive bool
}

func (s *Step) Name() string {
	return "dedupe"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case *payload.CSV, payload.JSONRecordPayload:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(ctx *engine.ExecutionContext) error {
	switch p := ctx.Payload.(type) {
	case payload.JSONRecordPayload:
		return s.dedupeJSON(p, ctx.Logger)
	case *payload.CSV:
		return s.dedupeCsv(p, ctx.Logger)
	default:
		return fmt.Errorf("unsupported payload type %T", ctx.Payload)
	}
}

func (s *Step) recordKey(record map[string]interface{}) string {
	var b strings.Builder
	for _, f := range s.Fields {
		v := fmt.Sprintf("%v", record[f])
		if !s.CaseSensitive {
			v = strings.ToLower(v)
		}
		b.WriteString(f + "=" + v + "\n")
	}
	return b.String()
}

func (s *Step) dedupeJSON(p payload.JSONRecordPayload, logger *engine.Logger) error {
	seen := make(map[string]struct{})
	records := p.Records()
	out := records[:0]
	removed := 0

	for _, rec := range records {
		key := s.recordKey(rec)
		if _, exists := seen[key]; exists {
			removed++
			continue
		}
		seen[key] = struct{}{}
		out = append(out, rec)
	}

	if removed > 0 {
		logger.Debug("  removed %d duplicate records", removed)
	}

	switch p := p.(type) {
	case *payload.JSON:
		p.Items = out
	case *payload.JSONL:
		p.Items = out
	}

	return nil
}

func (s *Step) dedupeCsv(csvPayload *payload.CSV, logger *engine.Logger) error {
	if len(csvPayload.Rows) == 0 {
		return nil
	}

	headerRow := csvPayload.Rows[0]
	seen := make(map[string]struct{})
	out := [][]string{headerRow}
	removed := 0

	for _, row := range csvPayload.Rows[1:] {
		record := payload.CSVRowToRecord(headerRow, row)
		key := s.recordKey(record)
		if _, exists := seen[key]; exists {
			removed++
			continue
		}
		seen[key] = struct{}{}
		out = append(out, row)
	}

	if removed > 0 {
		logger.Debug("  removed %d duplicate rows", removed)
	}

	csvPayload.Rows = out
	return nil
}
