package payload

type JSON struct {
	Data map[string]interface{}
}

func (p *JSON) Type() string {
	return "json"
}
