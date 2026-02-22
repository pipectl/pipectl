package steps

type Payload interface {
	Type() string
}

type JSONPayload struct {
	Data map[string]interface{}
}

func (j *JSONPayload) Type() string {
	return "json"
}

type CSVPayload struct {
	Rows [][]string
}

func (c *CSVPayload) Type() string {
	return "csv"
}

type TextPayload struct {
	Text string
}

func (t *TextPayload) Type() string {
	return "text"
}

type ExecutionContext struct {
	Payload Payload
}

type ExecutableStep interface {
	Execute(*ExecutionContext) error
	Supports(payload Payload) bool
	Name() string
}
