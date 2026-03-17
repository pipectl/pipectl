package payload

const CSVType string = "csv"

type CSV struct {
	Rows [][]string
}

func (p *CSV) Type() string {
	return CSVType
}

func (p *CSV) RecordCount() int {
	if len(p.Rows) == 0 {
		return 0
	}

	return len(p.Rows) - 1
}
