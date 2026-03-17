package payload

const JSONLType string = "jsonl"

type JSONL struct {
	Records []map[string]interface{}
}

func (p *JSONL) Type() string {
	return JSONLType
}

func (p *JSONL) RecordCount() int {
	return len(p.Records)
}

func (p *JSONL) GetRecords() []map[string]interface{} {
	return p.Records
}

func (p *JSONL) EnsureRecords() {
	if p.Records == nil {
		p.Records = []map[string]interface{}{{}}
	}
}

func (p *JSONL) Value() interface{} {
	return p.Records
}
