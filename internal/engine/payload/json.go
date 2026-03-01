package payload

const JSONType string = "json"

type JSON struct {
	Data map[string]interface{}
}

func (p *JSON) Type() string {
	return JSONType
}
