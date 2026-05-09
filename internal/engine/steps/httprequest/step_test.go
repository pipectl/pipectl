package httprequest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

func jsonPayload() *payload.JSON {
	return &payload.JSON{
		Items: []map[string]interface{}{{"a": "b"}},
		Shape: payload.JSONObjectShape,
	}
}

func TestExecuteDoesNotModifyPayload(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	original := jsonPayload()
	ctx := &engine.ExecutionContext{Payload: original}

	step := &Step{URL: target.URL, Method: http.MethodPost}
	if err := step.Execute(ctx); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if ctx.Payload != original {
		t.Fatal("expected payload to be unchanged after http-request")
	}
}

func TestExecuteWithBodyMethods(t *testing.T) {
	methods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					t.Fatalf("expected method %q, got %q", method, r.Method)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("failed reading request body: %v", err)
				}
				if len(body) == 0 {
					t.Fatalf("expected non-empty request body for %s", method)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer target.Close()

			ctx := &engine.ExecutionContext{Payload: jsonPayload()}
			step := &Step{URL: target.URL, Method: method}
			if err := step.Execute(ctx); err != nil {
				t.Fatalf("Execute returned error: %v", err)
			}
		})
	}
}

func TestExecuteWithNonBodyMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					t.Fatalf("expected method %q, got %q", method, r.Method)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("failed reading request body: %v", err)
				}
				if len(body) != 0 {
					t.Fatalf("expected empty request body for %s, got %q", method, body)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer target.Close()

			ctx := &engine.ExecutionContext{Payload: jsonPayload()}
			step := &Step{URL: target.URL, Method: method}
			if err := step.Execute(ctx); err != nil {
				t.Fatalf("Execute returned error: %v", err)
			}
		})
	}
}

func TestExecute2xxStatusCodesSucceed(t *testing.T) {
	codes := []int{http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent}
	for _, code := range codes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))
			defer target.Close()

			ctx := &engine.ExecutionContext{Payload: jsonPayload()}
			step := &Step{URL: target.URL, Method: http.MethodPost}
			if err := step.Execute(ctx); err != nil {
				t.Fatalf("expected success for %d, got: %v", code, err)
			}
		})
	}
}

func TestExecuteNon2xxStatusReturnsError(t *testing.T) {
	codes := []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError}
	for _, code := range codes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))
			defer target.Close()

			ctx := &engine.ExecutionContext{Payload: jsonPayload()}
			step := &Step{URL: target.URL, Method: http.MethodPost}
			err := step.Execute(ctx)
			if err == nil {
				t.Fatalf("expected error for %d response", code)
			}
			if !strings.Contains(err.Error(), "unexpected status code") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestExecuteWithCustomHeaders(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("expected Authorization header, got %q", got)
		}
		if got := r.Header.Get("X-Custom"); got != "value" {
			t.Fatalf("expected X-Custom header, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	ctx := &engine.ExecutionContext{Payload: jsonPayload()}
	step := &Step{
		URL:    target.URL,
		Method: http.MethodPost,
		Headers: map[string]string{
			"Authorization": "Bearer token",
			"X-Custom":      "value",
		},
	}
	if err := step.Execute(ctx); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
}

func TestExecuteWithInvalidProxyURL(t *testing.T) {
	ctx := &engine.ExecutionContext{Payload: jsonPayload()}
	step := &Step{
		URL:    "http://example.com",
		Method: http.MethodPost,
		Proxy:  "://bad proxy",
	}
	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for invalid proxy URL")
	}
	if !strings.Contains(err.Error(), "invalid proxy URL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteWithProxy(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	proxyCalled := make(chan struct{}, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case proxyCalled <- struct{}{}:
		default:
		}
		req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("proxy request creation failed: %v", err), http.StatusInternalServerError)
			return
		}
		req.Header = r.Header.Clone()
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("proxy round trip failed: %v", err), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		for k, vals := range resp.Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	}))
	defer proxy.Close()

	ctx := &engine.ExecutionContext{Payload: jsonPayload()}
	step := &Step{
		URL:    target.URL,
		Method: http.MethodPost,
		Proxy:  proxy.URL,
	}
	if err := step.Execute(ctx); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	select {
	case <-proxyCalled:
	default:
		t.Fatal("expected proxy to be called")
	}
}

func TestExecuteWithJSONPayloadUsesJSONContentType(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected Content-Type application/json, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	ctx := &engine.ExecutionContext{Payload: jsonPayload()}
	step := &Step{URL: target.URL, Method: http.MethodPost}
	if err := step.Execute(ctx); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
}

func TestExecuteWithJSONLPayloadUsesNDJSONContentType(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/x-ndjson" {
			t.Fatalf("expected Content-Type application/x-ndjson, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSONL{
			Items: []map[string]interface{}{{"a": 1}, {"b": 2}},
		},
	}
	step := &Step{URL: target.URL, Method: http.MethodPost}
	if err := step.Execute(ctx); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
}

func TestExecuteWithCSVPayloadUsesCSVContentType(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "text/csv" {
			t.Fatalf("expected Content-Type text/csv, got %q", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed reading request body: %v", err)
		}
		expectedBody := "id,name\n1,Alice\n"
		if string(body) != expectedBody {
			t.Fatalf("unexpected request body: got %q want %q", string(body), expectedBody)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	ctx := &engine.ExecutionContext{
		Payload: &payload.CSV{
			Rows: [][]string{
				{"id", "name"},
				{"1", "Alice"},
			},
		},
	}
	step := &Step{URL: target.URL, Method: http.MethodPost}
	if err := step.Execute(ctx); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
}

func TestExecuteWithInvalidTimeout(t *testing.T) {
	ctx := &engine.ExecutionContext{Payload: jsonPayload()}
	step := &Step{
		URL:     "http://example.com",
		Method:  http.MethodPost,
		Timeout: -1,
	}
	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
	if !strings.Contains(err.Error(), "invalid timeout") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSupportsAllFormats(t *testing.T) {
	step := &Step{}
	if !step.Supports(&payload.JSON{}) {
		t.Fatal("expected Supports to return true for JSON")
	}
	if !step.Supports(&payload.JSONL{}) {
		t.Fatal("expected Supports to return true for JSONL")
	}
	if !step.Supports(&payload.CSV{}) {
		t.Fatal("expected Supports to return true for CSV")
	}
}
