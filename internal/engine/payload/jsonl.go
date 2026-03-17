package payload

const JSONLType string = "jsonl"

type JSONL struct {
	Items []map[string]interface{}
}

func (p *JSONL) Type() string {
	return JSONLType
}

func (p *JSONL) RecordCount() int {
	return len(p.Items)
}

func (p *JSONL) Records() []map[string]interface{} {
	return p.Items
}

func (p *JSONL) Value() interface{} {
	return p.Items
}
