package payload

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Payload interface {
	Type() string
	RecordCount() int
}

func Read(input []byte, format string) (Payload, error) {
	switch format {

	case JSONType:
		var value interface{}
		if err := json.Unmarshal(input, &value); err != nil {
			return nil, fmt.Errorf("invalid JSON input: %w", err)
		}
		switch data := value.(type) {
		case map[string]interface{}:
			return &JSON{
				Items: []map[string]interface{}{data},
				Shape: JSONObjectShape,
			}, nil
		case []interface{}:
			records := make([]map[string]interface{}, 0, len(data))
			for i, item := range data {
				record, ok := item.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid JSON input: array item %d is %T, expected object", i, item)
				}
				records = append(records, record)
			}
			return &JSON{
				Items: records,
				Shape: JSONArrayShape,
			}, nil
		default:
			return nil, fmt.Errorf("invalid JSON input: expected object or array of objects")
		}

	case JSONLType:
		scanner := bufio.NewScanner(bytes.NewReader(input))
		records := make([]map[string]interface{}, 0)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				return nil, fmt.Errorf("invalid JSONL input on line %d: blank lines are not permitted", lineNumber)
			}

			var value interface{}
			if err := json.Unmarshal([]byte(line), &value); err != nil {
				return nil, fmt.Errorf("invalid JSONL input on line %d: %w", lineNumber, err)
			}
			record, ok := value.(map[string]interface{})
			if !ok || record == nil {
				return nil, fmt.Errorf("invalid JSONL input on line %d: expected object", lineNumber)
			}

			records = append(records, record)
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("invalid JSONL input: %w", err)
		}

		return &JSONL{Items: records}, nil

	default:
		return nil, fmt.Errorf("unsupported input format")
	}
}

func ReadCSV(input []byte, delimiter rune) (Payload, error) {
	r := csv.NewReader(bytes.NewReader(input))
	if delimiter != 0 {
		r.Comma = delimiter
	}
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV input: %w", err)
	}
	return &CSV{Rows: rows}, nil
}

func Convert(payload Payload, format string) (Payload, error) {
	switch format {
	case JSONType:
		switch typed := payload.(type) {
		case *JSON:
			return &JSON{
				Items: typed.Items,
				Shape: typed.Shape,
			}, nil
		case *JSONL:
			return &JSON{
				Items: typed.Items,
				Shape: JSONArrayShape,
			}, nil
		case *CSV:
			records, err := csvRowsToRecords(typed.Rows)
			if err != nil {
				return nil, err
			}
			return &JSON{
				Items: records,
				Shape: JSONArrayShape,
			}, nil
		default:
			return nil, fmt.Errorf("cannot convert %s payload to JSON", payload.Type())
		}

	case JSONLType:
		switch typed := payload.(type) {
		case *JSONL:
			return &JSONL{Items: typed.Items}, nil
		case *JSON:
			return &JSONL{Items: typed.Items}, nil
		case *CSV:
			records, err := csvRowsToRecords(typed.Rows)
			if err != nil {
				return nil, err
			}
			return &JSONL{Items: records}, nil
		default:
			return nil, fmt.Errorf("cannot convert %s payload to JSONL", payload.Type())
		}

	case CSVType:
		switch typed := payload.(type) {
		case *CSV:
			return &CSV{Rows: typed.Rows}, nil
		case *JSON:
			rows, err := jsonRecordsToCSVRows(typed.Items)
			if err != nil {
				return nil, err
			}
			return &CSV{Rows: rows}, nil
		case *JSONL:
			rows, err := jsonRecordsToCSVRows(typed.Items)
			if err != nil {
				return nil, err
			}
			return &CSV{Rows: rows}, nil
		default:
			return nil, fmt.Errorf("cannot convert %s payload to CSV", payload.Type())
		}

	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

func Write(payload Payload, format string, writer io.Writer) error {
	converted, err := Convert(payload, format)
	if err != nil {
		return err
	}

	if writer == nil {
		writer = os.Stdout
	}

	switch typed := converted.(type) {
	case *JSON:
		output, err := json.MarshalIndent(typed.Value(), "", "  ")
		if err != nil {
			return fmt.Errorf("marshal JSON output: %w", err)
		}
		_, err = fmt.Fprintln(writer, string(output))
		return err
	case *JSONL:
		for _, record := range typed.Items {
			output, err := json.Marshal(record)
			if err != nil {
				return fmt.Errorf("marshal JSONL output: %w", err)
			}
			if _, err := fmt.Fprintln(writer, string(output)); err != nil {
				return err
			}
		}
		return nil
	case *CSV:
		buf := new(bytes.Buffer)
		csvWriter := csv.NewWriter(buf)
		if err := csvWriter.WriteAll(typed.Rows); err != nil {
			return fmt.Errorf("write CSV output: %w", err)
		}
		_, err := fmt.Fprintln(writer, buf.String())
		return err
	default:
		return fmt.Errorf("unsupported converted payload type: %T", converted)
	}
}

func csvRowsToRecords(rows [][]string) ([]map[string]interface{}, error) {
	if len(rows) == 0 {
		return []map[string]interface{}{}, nil
	}

	headers := rows[0]
	records := make([]map[string]interface{}, 0, len(rows)-1)

	for rowIndex, row := range rows[1:] {
		if len(row) != len(headers) {
			return nil, fmt.Errorf(
				"invalid CSV payload: row %d has %d columns, expected %d",
				rowIndex+2,
				len(row),
				len(headers),
			)
		}

		record := make(map[string]interface{}, len(headers))
		for columnIndex, header := range headers {
			if err := assignCSVField(record, header, row[columnIndex]); err != nil {
				return nil, fmt.Errorf("invalid CSV payload: row %d column %q: %w", rowIndex+2, header, err)
			}
		}
		records = append(records, record)
	}

	return records, nil
}

func assignCSVField(record map[string]interface{}, header, rawValue string) error {
	parts := strings.Split(header, ".")
	current := record

	for i, part := range parts {
		last := i == len(parts)-1
		if last {
			if _, exists := current[part]; exists {
				return fmt.Errorf("duplicate field %q", headerForError(parts[:i+1]))
			}

			current[part] = parseCSVFieldValue(rawValue)
			return nil
		}

		existing, exists := current[part]
		if !exists {
			child := make(map[string]interface{})
			current[part] = child
			current = child
			continue
		}

		child, ok := existing.(map[string]interface{})
		if !ok {
			return fmt.Errorf("field %q conflicts with nested field %q", headerForError(parts[:i+1]), header)
		}
		current = child
	}

	return nil
}

func parseCSVFieldValue(raw string) interface{} {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return raw
	}

	if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
		var decoded interface{}
		if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
			switch decoded.(type) {
			case []interface{}, map[string]interface{}:
				return decoded
			}
		}
	}

	return raw
}

func headerForError(parts []string) string {
	return strings.Join(parts, ".")
}

func jsonRecordsToCSVRows(records []map[string]interface{}) ([][]string, error) {
	if len(records) == 0 {
		return [][]string{}, nil
	}

	headerSet := make(map[string]struct{})
	for _, record := range records {
		flattened, err := flattenJSONRecord(record)
		if err != nil {
			return nil, err
		}
		for key := range flattened {
			headerSet[key] = struct{}{}
		}
	}

	headers := make([]string, 0, len(headerSet))
	for key := range headerSet {
		headers = append(headers, key)
	}
	sort.Strings(headers)

	rows := make([][]string, 0, len(records)+1)
	rows = append(rows, headers)

	for _, record := range records {
		flattened, err := flattenJSONRecord(record)
		if err != nil {
			return nil, err
		}

		row := make([]string, len(headers))
		for index, header := range headers {
			row[index] = flattened[header]
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func flattenJSONRecord(record map[string]interface{}) (map[string]string, error) {
	flattened := make(map[string]string)
	for key, value := range record {
		if err := flattenJSONValue(flattened, key, value); err != nil {
			return nil, err
		}
	}

	return flattened, nil
}

func flattenJSONValue(flattened map[string]string, key string, value interface{}) error {
	switch typed := value.(type) {
	case map[string]interface{}:
		for nestedKey, nestedValue := range typed {
			childKey := nestedKey
			if key != "" {
				childKey = key + "." + nestedKey
			}
			if err := flattenJSONValue(flattened, childKey, nestedValue); err != nil {
				return err
			}
		}
		return nil
	case []interface{}:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return fmt.Errorf("marshal JSON array field %q: %w", key, err)
		}
		flattened[key] = string(encoded)
		return nil
	case nil:
		flattened[key] = ""
		return nil
	case string:
		flattened[key] = typed
		return nil
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return fmt.Errorf("marshal JSON field %q: %w", key, err)
		}
		flattened[key] = strings.Trim(string(encoded), "\"")
		return nil
	}
}
