package httptransform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
)

type Step struct {
	URL    string
	Method string
	Proxy  string
}

func (s *Step) Name() string {
	return "http-transform"
}

func (s *Step) Supports(p payload.Payload) bool {
	return p.Type() == payload.JSONType
}

func (s *Step) Execute(context *engine.ExecutionContext) error {
	fmt.Printf("- transforming via HTTP %v to URL: %v\n", s.Method, s.URL)
	if s.Proxy != "" {
		fmt.Printf("- using proxy: %v\n", s.Proxy)
	}

	transformedPayload, err := s.transformPayload(context.Payload)
	if err != nil {
		return err
	}

	context.Payload = transformedPayload
	return nil
}

func (s *Step) transformPayload(inputPayload payload.Payload) (*payload.JSON, error) {
	var bodyReader io.Reader
	if s.Method == "POST" {
		jsonBody, _ := json.Marshal(inputPayload)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(s.Method, s.URL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Pipectl-Step", "http-transform")

	client := &http.Client{}
	if s.Proxy != "" {
		proxyURL, err := url.Parse(s.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL %q: %w", s.Proxy, err)
		}

		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.Proxy = http.ProxyURL(proxyURL)
		client.Transport = transport
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error calling HTTP service: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d\n", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s\n", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %s\n", err)
	}

	return &payload.JSON{Data: data}, nil
}
