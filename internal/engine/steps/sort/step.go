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
	NullsFirst    = "first"
	NullsLast     = "last"
)

// Key is one sort key. Keys are compared in order; later keys only break ties
// left by earlier ones.
type Key struct {
	Field     string
	Direction string
	Nulls     string
}

type Step struct {
	payload.AllFormatsSupport
	Keys []Key
}

func (s *Step) Name() string {
	return "sort"
}

func (s *Step) Execute(ctx *engine.ExecutionContext) error {
	switch p := ctx.Payload.(type) {
	case *payload.JSON:
		if p.Shape != payload.JSONArrayShape {
			return fmt.Errorf("cannot sort a JSON object payload, only JSON arrays are supported")
		}
		sort.SliceStable(p.Items, s.jsonLess(p.Items))
		ctx.Logger.Debug("  sorted %d records by %s", len(p.Items), s.describeKeys())
		return nil
	case *payload.JSONL:
		sort.SliceStable(p.Items, s.jsonLess(p.Items))
		ctx.Logger.Debug("  sorted %d records by %s", len(p.Items), s.describeKeys())
		return nil
	case *payload.CSV:
		return s.sortCSV(p, ctx.Logger)
	default:
		return fmt.Errorf("unsupported payload type %T", ctx.Payload)
	}
}

func (s *Step) jsonLess(records []map[string]interface{}) func(i, j int) bool {
	return func(i, j int) bool {
		for _, k := range s.Keys {
			aVal, aExists := records[i][k.Field]
			bVal, bExists := records[j][k.Field]

			aMissing := !aExists || aVal == nil
			bMissing := !bExists || bVal == nil

			if c := compare(k, fmt.Sprintf("%v", aVal), fmt.Sprintf("%v", bVal), aMissing, bMissing); c != 0 {
				return c < 0
			}
		}
		return false
	}
}

func (s *Step) sortCSV(p *payload.CSV, logger *engine.Logger) error {
	if len(p.Rows) <= 1 {
		return nil
	}

	fields := make([]string, len(s.Keys))
	for i, k := range s.Keys {
		fields[i] = k.Field
	}

	header := p.Rows[0]
	colIndex, err := payload.FindColumnIndices(header, fields)
	if err != nil {
		return fmt.Errorf("sort: %w", err)
	}

	rows := make([][]string, len(p.Rows)-1)
	copy(rows, p.Rows[1:])

	sort.SliceStable(rows, func(i, j int) bool {
		for _, k := range s.Keys {
			idx := colIndex[k.Field]
			a := rows[i][idx]
			b := rows[j][idx]
			if c := compare(k, a, b, a == "", b == ""); c != 0 {
				return c < 0
			}
		}
		return false
	})

	p.Rows = append([][]string{header}, rows...)
	logger.Debug("  sorted %d rows by %s", len(rows), s.describeKeys())
	return nil
}

func (s *Step) describeKeys() string {
	parts := make([]string, len(s.Keys))
	for i, k := range s.Keys {
		parts[i] = k.Field + " " + k.Direction
	}
	return strings.Join(parts, ", ")
}

// compare returns -1 if a sorts before b, 1 if after, and 0 if they tie on k.
func compare(k Key, a, b string, aMissing, bMissing bool) int {
	if aMissing && bMissing {
		return 0
	}
	// Null placement is independent of direction; an empty Nulls means last.
	if aMissing {
		if k.Nulls == NullsFirst {
			return -1
		}
		return 1
	}
	if bMissing {
		if k.Nulls == NullsFirst {
			return 1
		}
		return -1
	}

	var c int
	fa, errA := strconv.ParseFloat(strings.TrimSpace(a), 64)
	fb, errB := strconv.ParseFloat(strings.TrimSpace(b), 64)
	if errA == nil && errB == nil {
		switch {
		case fa < fb:
			c = -1
		case fa > fb:
			c = 1
		}
	} else {
		c = strings.Compare(a, b)
	}

	if k.Direction == DirectionDesc {
		return -c
	}
	return c
}
