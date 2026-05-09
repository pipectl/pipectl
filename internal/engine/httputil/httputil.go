package httputil

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pipectl/pipectl/internal/engine/payload"
)

const (
	DefaultTimeoutSeconds = 60
	MaxTimeoutSeconds     = 300
)

// MarshalPayload encodes p to bytes and returns the matching Content-Type.
// Returns (nil, "", nil) when the payload type does not produce a request body.
func MarshalPayload(p payload.Payload) ([]byte, string, error) {
	switch v := p.(type) {
	case *payload.JSON:
		body, err := json.Marshal(v.Value())
		if err != nil {
			return nil, "", fmt.Errorf("could not encode JSON payload: %w", err)
		}
		return body, "application/json", nil
	case *payload.JSONL:
		body, err := marshalJSONL(v)
		if err != nil {
			return nil, "", fmt.Errorf("could not encode JSONL payload: %w", err)
		}
		return body, "application/x-ndjson", nil
	case *payload.CSV:
		body, err := marshalCSV(v)
		if err != nil {
			return nil, "", fmt.Errorf("could not encode CSV payload: %w", err)
		}
		return body, "text/csv", nil
	default:
		return nil, "", fmt.Errorf("unsupported payload type %v", p.Type())
	}
}

// BuildClient returns an *http.Client optionally configured with a proxy.
func BuildClient(proxy string) (*http.Client, error) {
	client := &http.Client{}
	if proxy == "" {
		return client, nil
	}

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL %q: %w", proxy, err)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyURL(proxyURL)
	client.Transport = transport
	return client, nil
}

// ResolveTimeout validates seconds and returns the corresponding duration.
// A value of 0 returns the default timeout.
func ResolveTimeout(seconds int) (time.Duration, error) {
	if seconds == 0 {
		return DefaultTimeoutSeconds * time.Second, nil
	}
	if seconds < 0 {
		return 0, fmt.Errorf("invalid timeout %d: must be a positive number of seconds", seconds)
	}
	if seconds > MaxTimeoutSeconds {
		return 0, fmt.Errorf("invalid timeout %d: maximum is %d seconds", seconds, MaxTimeoutSeconds)
	}
	return time.Duration(seconds) * time.Second, nil
}

func marshalCSV(input *payload.CSV) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.WriteAll(input.Rows); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
