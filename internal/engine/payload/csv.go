package payload

const CSVType string = "csv"

type CSV struct {
	Rows [][]string
}

func (p *CSV) Type() string {
	return CSVType
}
