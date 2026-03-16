package payload

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
)

type Payload interface {
	Type() string
}

func Read(input []byte, format string) (Payload, error) {
	switch format {

	case "json":
		var data map[string]interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return nil, err
		}
		return &JSON{Data: data}, nil

	case "csv":
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
	if format == "json" {

		// TODO: which payload types can be converted to JSON?
		switch payload.Type() {

		case "json":
			jsonPayload, _ := payload.(*JSON)
			output, err := json.MarshalIndent(jsonPayload.Data, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return err
			}
			fmt.Println(string(output))

		case "csv":
			csvPayload, _ := payload.(*CSV)
			// TODO how to convert from CSV to JSON?
			fmt.Println("TODO: convert CSV to JSON")
			fmt.Println(csvPayload.Rows)

		default:
			return fmt.Errorf("Cannot convert to JSON")
		}

	} else if format == "csv" {
		switch payload.Type() {
		case "csv":
			csvPayload, _ := payload.(*CSV)
			buf := new(bytes.Buffer)
			csvWriter := csv.NewWriter(buf)
			if err := csvWriter.WriteAll(csvPayload.Rows); err != nil {
				fmt.Println("Error writing CSV:", err)
				return err
			}
			fmt.Println(buf.String())

		case "json":
			// TODO convert JSON to CSV

		default:
			return fmt.Errorf("Cannot convert to CSV")
		}

	}

	return nil
}
