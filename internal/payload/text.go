package payload

type Text struct {
	Text string
}

func (p *Text) Type() string {
	return "text"
}
