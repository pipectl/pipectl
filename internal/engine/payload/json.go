package payload

const JSONType string = "json"

type JSONRecordPayload interface {
	Payload
	Records() []map[string]interface{}
	Value() interface{}
}

type JSONShape string

const (
	JSONObjectShape JSONShape = "object"
	JSONArrayShape  JSONShape = "array"
)

type JSON struct {
	Items []map[string]interface{}
	Shape JSONShape
}

func (p *JSON) Type() string {
	return JSONType
}

func (p *JSON) RecordCount() int {
	return len(p.Items)
}

func (p *JSON) Records() []map[string]interface{} {
	return p.Items
}

func (p *JSON) Value() interface{} {
	switch p.Shape {
	case JSONArrayShape:
		return p.Items
	case JSONObjectShape, "":
		if len(p.Items) == 0 {
			return map[string]interface{}{}
		}
		return p.Items[0]
	default:
		if len(p.Items) == 0 {
			return map[string]interface{}{}
		}
		return p.Items[0]
	}
}
