package payload

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
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
				Records: []map[string]interface{}{data},
				Shape:   JSONObjectShape,
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
				Records: records,
				Shape:   JSONArrayShape,
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
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
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

		return &JSONL{Records: records}, nil

	case CSVType:
		rows, err := csv.NewReader(bytes.NewReader(input)).ReadAll()
		if err != nil {
			panic(err)
		}
		return &CSV{Rows: rows}, nil

	default:
		return nil, fmt.Errorf("unsupported input format")
	}
}

func Write(payload Payload, format string) error {
	if format == JSONType {

		// TODO: which payload types can be converted to JSON?
		switch payload.Type() {

		case JSONType:
			jsonPayload, _ := payload.(*JSON)
			output, err := json.MarshalIndent(jsonPayload.Value(), "", "  ")
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return err
			}
			fmt.Println(string(output))

		case JSONLType:
			jsonlPayload, _ := payload.(*JSONL)
			output, err := json.MarshalIndent(jsonlPayload.Value(), "", "  ")
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return err
			}
			fmt.Println(string(output))

		case CSVType:
			csvPayload, _ := payload.(*CSV)
			// TODO how to convert from CSV to JSON?
			fmt.Println("TODO: convert CSV to JSON")
			fmt.Println(csvPayload.Rows)

		default:
			return fmt.Errorf("Cannot convert to JSON")
		}

	} else if format == JSONLType {
		switch payload.Type() {
		case JSONLType:
			jsonlPayload, _ := payload.(*JSONL)
			for _, record := range jsonlPayload.Records {
				output, err := json.Marshal(record)
				if err != nil {
					fmt.Println("Error marshalling JSONL:", err)
					return err
				}
				fmt.Println(string(output))
			}

		case JSONType:
			jsonPayload, _ := payload.(*JSON)
			for _, record := range jsonPayload.Records {
				output, err := json.Marshal(record)
				if err != nil {
					fmt.Println("Error marshalling JSONL:", err)
					return err
				}
				fmt.Println(string(output))
			}

		default:
			return fmt.Errorf("Cannot convert to JSONL")
		}

	} else if format == CSVType {
		switch payload.Type() {
		case CSVType:
			csvPayload, _ := payload.(*CSV)
			buf := new(bytes.Buffer)
			csvWriter := csv.NewWriter(buf)
			if err := csvWriter.WriteAll(csvPayload.Rows); err != nil {
				fmt.Println("Error writing CSV:", err)
				return err
			}
			fmt.Println(buf.String())

		case JSONType:
			// TODO convert JSON to CSV

		default:
			return fmt.Errorf("Cannot convert to CSV")
		}

	}

	return nil
}
