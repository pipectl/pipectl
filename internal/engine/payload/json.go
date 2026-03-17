package payload

const JSONType string = "json"

type RecordPayload interface {
	Payload
	GetRecords() []map[string]interface{}
	EnsureRecords()
	Value() interface{}
}

type JSONShape string

const (
	JSONObjectShape JSONShape = "object"
	JSONArrayShape  JSONShape = "array"
)

type JSON struct {
	Records []map[string]interface{}
	Shape   JSONShape
}

func (p *JSON) Type() string {
	return JSONType
}

func (p *JSON) RecordCount() int {
	return len(p.Records)
}

func (p *JSON) GetRecords() []map[string]interface{} {
	return p.Records
}

func (p *JSON) EnsureRecords() {
	if p.Records == nil {
		p.Records = []map[string]interface{}{{}}
	}
	if p.Shape == "" {
		p.Shape = JSONObjectShape
	}
}

func (p *JSON) Value() interface{} {
	switch p.Shape {
	case JSONArrayShape:
		return p.Records
	case JSONObjectShape, "":
		if len(p.Records) == 0 {
			return map[string]interface{}{}
		}
		return p.Records[0]
	default:
		if len(p.Records) == 0 {
			return map[string]interface{}{}
		}
		return p.Records[0]
	}
}
