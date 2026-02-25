package http_transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
)

type Step struct {
	URL    string
	Method string
}

func (s *Step) Name() string {
	return "http-transform"
}

func (s *Step) Supports(p payload.Payload) bool {
	return p.Type() == payload.JSONType
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	fmt.Printf("- transforming via HTTP %v to URL: %v\n", s.Method, s.URL)

	var bodyReader io.Reader
	if s.Method == "POST" {
		jsonBody, _ := json.Marshal(context.Payload)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(s.Method, s.URL, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("X-Pipectl-Step", "http-transform")

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %d\n", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %s\n", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("Error parsing JSON response: %s\n", err)
	}
	context.Payload = &payload.JSON{Data: data}

	return nil
}
