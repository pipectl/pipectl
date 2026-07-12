package payload

import "fmt"

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

// FindColumnIndices returns a map from each field name to its column index in headers.
// It returns an error naming the first field not found in headers.
func FindColumnIndices(headers []string, fields []string) (map[string]int, error) {
	all := HeaderIndex(headers)
	indices := make(map[string]int, len(fields))
	for _, field := range fields {
		i, ok := all[field]
		if !ok {
			return nil, fmt.Errorf("field %q not found in CSV headers", field)
		}
		indices[field] = i
	}
	return indices, nil
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
