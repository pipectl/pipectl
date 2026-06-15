package wizard

import (
	"fmt"
	"strings"
)

// Result holds the answers collected by the wizard form.
type Result struct {
	ID           string
	InputFormat  string
	Steps        []string
	OutputFormat string
	OutputFile   string
}

// stepTemplate is a YAML snippet for a single step, indented to sit inside the
// `steps:` list (two-space indent, dash on the first line).
var stepTemplates = map[string]string{
	"select": `  - select:
      fields: [field1, field2]  # fields to keep; all others are dropped`,

	"rename": `  - rename:
      fields:
        old_name: new_name  # add more entries as needed`,

	"cast": `  - cast:
      fields:
        field_name:
          type: string  # int | float | bool | time | string`,

	"default": `  - default:
      fields:
        field_name: default_value  # applied when field is missing or empty`,

	"normalize": `  - normalize:
      fields:
        field_name: lower  # lower | upper | trim | trim-left | trim-right | collapse-spaces | capitalize`,

	"redact": `  - redact:
      strategy: mask  # mask | sha256 | partial-first | partial-last
      fields: [field_name]`,

	"convert": `  - convert:
      format: json  # json | jsonl | csv`,

	"filter": `  - filter:
      field: field_name
      equals: value  # equals | not-equals | contains | starts-with | ends-with | greater-than | less-than`,

	"sort": `  - sort:
      field: field_name
      direction: asc  # asc | desc`,

	"limit": `  - limit:
      count: 100`,

	"dedupe": `  - dedupe:
      fields: [field1, field2]
      case-sensitive: true`,

	"assert": `  - assert:
      min-records: 1  # also: max-records, records-equal, field-exists`,

	"validate-json": `  - validate-json:
      schema: ./schema.json`,

	"count": `  - count:
      message: "Record count"`,

	"log": `  - log:
      message: "Debug checkpoint"
      count: true   # print record count
      sample: 3     # print N sample records`,

	"http-request": `  - http-request:
      url: https://example.com/endpoint
      method: POST  # GET | POST | PUT | PATCH | DELETE | HEAD | OPTIONS
      timeout: 30   # seconds (0-300)`,

	"http-transform": `  - http-transform:
      url: https://example.com/transform
      method: POST  # GET | POST | PUT | PATCH | DELETE | HEAD | OPTIONS
      timeout: 30
      expect-format: json  # json | jsonl | csv`,
}

// Render converts a wizard Result into a pipeline YAML string.
func Render(r Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "id: %s\n", r.ID)
	b.WriteString("input:\n")
	fmt.Fprintf(&b, "  format: %s\n", r.InputFormat)
	b.WriteString("steps:\n")
	for _, step := range r.Steps {
		if tpl, ok := stepTemplates[step]; ok {
			b.WriteString(tpl)
			b.WriteByte('\n')
		}
	}
	b.WriteString("output:\n")
	fmt.Fprintf(&b, "  format: %s\n", r.OutputFormat)
	return b.String()
}
