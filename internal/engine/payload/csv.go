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

// HeaderIndex returns a map from header name to column index.
func HeaderIndex(headers []string) map[string]int {
	idx := make(map[string]int, len(headers))
	for i, h := range headers {
		idx[h] = i
	}
	return idx
}

// CSVRowToRecord maps a single CSV row to a flat string-keyed record using the provided headers.
// Columns beyond the length of row are omitted.
func CSVRowToRecord(headers []string, row []string) map[string]interface{} {
	record := make(map[string]interface{}, len(headers))
	for i, header := range headers {
		if i < len(row) {
			record[header] = row[i]
		}
	}
	return record
}
