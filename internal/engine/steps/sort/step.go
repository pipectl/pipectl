package sort

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

const (
	DirectionAsc  = "asc"
	DirectionDesc = "desc"
)

type Step struct {
	Field     string
	Direction string
}

func (s *Step) Name() string {
	return "sort"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case *payload.JSON, *payload.JSONL, *payload.CSV:
		return true
	default:
		return false
	}
}

func (s *Step) Execute(ctx *engine.ExecutionContext) error {
	switch p := ctx.Payload.(type) {
	case *payload.JSON:
		if p.Shape != payload.JSONArrayShape {
			return fmt.Errorf("cannot sort a JSON object payload, only JSON arrays are supported")
		}
		sort.SliceStable(p.Items, s.jsonLess(p.Items))
		ctx.Logger.Debug("  sorted %d records by %s %s", len(p.Items), s.Field, s.Direction)
		return nil
	case *payload.JSONL:
		sort.SliceStable(p.Items, s.jsonLess(p.Items))
		ctx.Logger.Debug("  sorted %d records by %s %s", len(p.Items), s.Field, s.Direction)
		return nil
	case *payload.CSV:
		return s.sortCSV(p, ctx.Logger)
	default:
		return fmt.Errorf("unsupported payload type %T", ctx.Payload)
	}
}

func (s *Step) jsonLess(records []map[string]interface{}) func(i, j int) bool {
	return func(i, j int) bool {
		aVal, aExists := records[i][s.Field]
		bVal, bExists := records[j][s.Field]

		aMissing := !aExists || aVal == nil
		bMissing := !bExists || bVal == nil

		return s.less(fmt.Sprintf("%v", aVal), fmt.Sprintf("%v", bVal), aMissing, bMissing)
	}
}

func (s *Step) sortCSV(p *payload.CSV, logger *engine.Logger) error {
	if len(p.Rows) <= 1 {
		return nil
	}

	header := p.Rows[0]
	colIndex := -1
	for i, h := range header {
		if h == s.Field {
			colIndex = i
			break
		}
	}

	if colIndex == -1 {
		return fmt.Errorf("sort: field %q not found in CSV headers", s.Field)
	}

	rows := make([][]string, len(p.Rows)-1)
	copy(rows, p.Rows[1:])

	sort.SliceStable(rows, func(i, j int) bool {
		a := rows[i][colIndex]
		b := rows[j][colIndex]
		return s.less(a, b, a == "", b == "")
	})

	p.Rows = append([][]string{header}, rows...)
	logger.Debug("  sorted %d rows by %s %s", len(rows), s.Field, s.Direction)
	return nil
}

func (s *Step) less(a, b string, aMissing, bMissing bool) bool {
	if aMissing && bMissing {
		return false
	}
	if aMissing {
		return false // nulls last
	}
	if bMissing {
		return true // nulls last
	}

	fa, errA := strconv.ParseFloat(strings.TrimSpace(a), 64)
	fb, errB := strconv.ParseFloat(strings.TrimSpace(b), 64)
	if errA == nil && errB == nil {
		if s.Direction == DirectionDesc {
			return fa > fb
		}
		return fa < fb
	}

	if s.Direction == DirectionDesc {
		return a > b
	}
	return a < b
}
