package payload

const TextType string = "text"

type Text struct {
	Text string
}

func (p *Text) Type() string {
	return TextType
}
