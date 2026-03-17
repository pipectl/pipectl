package httptransform

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

type Step struct {
	URL          string
	Method       string
	Proxy        string
	Headers      map[string]string
	Timeout      int
	ExpectFormat string
}

const (
	defaultTimeoutSeconds = 60
	maxTimeoutSeconds     = 300
)

func (s *Step) Name() string {
	return "http-transform"
}

func (s *Step) Supports(p payload.Payload) bool {
	switch p.(type) {
	case payload.JSONRecordPayload:
		return true
	default:
		return false
	}
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

func (s *Step) transformPayload(inputPayload payload.Payload) (payload.Payload, error) {
	var bodyReader io.Reader
	method := strings.ToUpper(s.Method)

	// Set the request body from the step payload
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete {
		switch v := inputPayload.(type) {
		case *payload.JSON:
			jsonBody, _ := json.Marshal(v.Value())
			bodyReader = bytes.NewBuffer(jsonBody)
		case *payload.JSONL:
			body, err := marshalJSONL(v)
			if err != nil {
				return nil, fmt.Errorf("http-transform could not encode JSONL payload: %w", err)
			}
			bodyReader = bytes.NewBuffer(body)
		default:
			return nil, fmt.Errorf("http-transform received invalid payload type %v", inputPayload.Type())
		}
	}

	req, err := http.NewRequest(method, s.URL, bodyReader)
	if err != nil {
		return nil, err
	}
	requestTimeout, err := s.resolveTimeout()
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(req.Context(), requestTimeout)
	defer cancel()
	req = req.WithContext(timeoutCtx)

	// add headers
	for key, value := range s.Headers {
		req.Header.Set(key, value)
	}
	if bodyReader != nil && req.Header.Get("Content-Type") == "" && inputPayload.Type() == payload.JSONLType {
		req.Header.Set("Content-Type", "application/x-ndjson")
	}

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

	expectedFormat, err := s.resolveExpectedFormat()
	if err != nil {
		return nil, err
	}
	if err := validateResponseContentType(resp.Header.Get("Content-Type"), expectedFormat); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s\n", err)
	}

	switch expectedFormat {
	case payload.JSONType:
		return payload.Read(body, payload.JSONType)
	case payload.JSONLType:
		return payload.Read(body, payload.JSONLType)

	case payload.CSVType:
		rows, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
		if err != nil {
			return nil, fmt.Errorf("Error parsing CSV response: %s\n", err)
		}
		return &payload.CSV{Rows: rows}, nil
	}

	return nil, fmt.Errorf("invalid expect-format %q", expectedFormat)
}

func (s *Step) resolveExpectedFormat() (string, error) {
	expectedFormat := strings.ToLower(strings.TrimSpace(s.ExpectFormat))
	if expectedFormat == "" {
		return payload.JSONType, nil
	}

	if expectedFormat != payload.JSONType && expectedFormat != payload.JSONLType && expectedFormat != payload.CSVType {
		return "", fmt.Errorf("invalid expect-format %q: must be json, jsonl or csv", s.ExpectFormat)
	}

	return expectedFormat, nil
}

func validateResponseContentType(contentType string, expectedFormat string) error {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("invalid response Content-Type %q: %w", contentType, err)
	}

	if !contentTypeMatchesFormat(mediaType, expectedFormat) {
		return fmt.Errorf("response Content-Type %q does not match expect-format %q", mediaType, expectedFormat)
	}

	return nil
}

func contentTypeMatchesFormat(mediaType string, expectedFormat string) bool {
	normalizedMediaType := strings.ToLower(strings.TrimSpace(mediaType))
	switch expectedFormat {
	case payload.JSONType:
		return strings.HasSuffix(normalizedMediaType, "/json") || strings.HasSuffix(normalizedMediaType, "+json")
	case payload.JSONLType:
		return normalizedMediaType == "application/x-ndjson" || normalizedMediaType == "application/ndjson" || normalizedMediaType == "application/jsonl"
	case payload.CSVType:
		return strings.HasSuffix(normalizedMediaType, "/csv") || strings.Contains(normalizedMediaType, "csv")
	default:
		return false
	}
}

func marshalJSONL(input *payload.JSONL) ([]byte, error) {
	var body bytes.Buffer
	for _, record := range input.Items {
		raw, err := json.Marshal(record)
		if err != nil {
			return nil, err
		}
		body.Write(raw)
		body.WriteByte('\n')
	}
	return body.Bytes(), nil
}

func (s *Step) resolveTimeout() (time.Duration, error) {
	if s.Timeout == 0 {
		return defaultTimeoutSeconds * time.Second, nil
	}

	if s.Timeout < 0 {
		return 0, fmt.Errorf("invalid timeout %d: must be a positive number of seconds", s.Timeout)
	}

	if s.Timeout > maxTimeoutSeconds {
		return 0, fmt.Errorf("invalid timeout %d: maximum is %d seconds", s.Timeout, maxTimeoutSeconds)
	}

	return time.Duration(s.Timeout) * time.Second, nil
}
