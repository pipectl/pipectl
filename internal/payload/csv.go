package payload

type CSV struct {
	Rows [][]string
}

func (p *CSV) Type() string {
	return "csv"
}
