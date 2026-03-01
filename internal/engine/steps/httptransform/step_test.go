package httptransform

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/engine/payload"
)

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

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true}`))
			}))
			defer target.Close()

			step := &Step{
				URL:    target.URL,
				Method: method,
			}

			ctx := &engine.ExecutionContext{
				Payload: &payload.JSON{Data: map[string]interface{}{"a": "b"}},
			}

			if err := step.Execute(ctx); err != nil {
				t.Fatalf("execute returned error: %v", err)
			}
		})
	}
}

func TestResolveTimeoutDefault(t *testing.T) {
	step := &Step{}

	got, err := step.resolveTimeout()
	if err != nil {
		t.Fatalf("resolveTimeout returned error: %v", err)
	}
	if got != 60*time.Second {
		t.Fatalf("expected default timeout 60s, got %v", got)
	}
}

func TestResolveTimeoutCustom(t *testing.T) {
	step := &Step{Timeout: 2}

	got, err := step.resolveTimeout()
	if err != nil {
		t.Fatalf("resolveTimeout returned error: %v", err)
	}
	if got != 2*time.Second {
		t.Fatalf("expected timeout 2s, got %v", got)
	}
}

func TestResolveTimeoutInvalid(t *testing.T) {
	step := &Step{Timeout: -1}

	_, err := step.resolveTimeout()
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
	if !strings.Contains(err.Error(), "invalid timeout") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveTimeoutAboveMaximum(t *testing.T) {
	step := &Step{Timeout: 301}

	_, err := step.resolveTimeout()
	if err == nil {
		t.Fatal("expected error for timeout above maximum")
	}
	if !strings.Contains(err.Error(), "maximum is 300 seconds") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteWithTimeoutInSeconds(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"from":"slow-target"}`))
	}))
	defer target.Close()

	step := &Step{
		URL:     target.URL,
		Method:  http.MethodPost,
		Timeout: 1,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"a": "b"}},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("unexpected timeout error: %v", err)
	}
}

func TestExecuteWithoutProxy(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"from":"target"}`))
	}))
	defer target.Close()

	step := &Step{
		URL:    target.URL,
		Method: http.MethodPost,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"a": "b"}},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.JSON)
	if !ok {
		t.Fatalf("expected payload.JSON, got %T", ctx.Payload)
	}
	if got, ok := out.Data["from"]; !ok || got != "target" {
		t.Fatalf("unexpected response payload: %#v", out.Data)
	}
}

func TestExecuteWithHeaders(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer {{API_TOKEN}}" {
			t.Fatalf("expected Authorization header to be set, got %q", got)
		}
		if got := r.Header.Get("X-Custom-Header"); got != "custom-value" {
			t.Fatalf("expected X-Custom-Header header to be set, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"from":"target-with-headers"}`))
	}))
	defer target.Close()

	step := &Step{
		URL:    target.URL,
		Method: http.MethodPost,
		Headers: map[string]string{
			"Authorization":   "Bearer {{API_TOKEN}}",
			"X-Custom-Header": "custom-value",
		},
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"a": "b"}},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	out, ok := ctx.Payload.(*payload.JSON)
	if !ok {
		t.Fatalf("expected payload.JSON, got %T", ctx.Payload)
	}
	if got, ok := out.Data["from"]; !ok || got != "target-with-headers" {
		t.Fatalf("unexpected response payload: %#v", out.Data)
	}
}

func TestExecuteWithProxy(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"from":"target-via-proxy"}`))
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

	step := &Step{
		URL:    target.URL,
		Method: http.MethodPost,
		Proxy:  proxy.URL,
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"a": "b"}},
	}

	if err := step.Execute(ctx); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	select {
	case <-proxyCalled:
	default:
		t.Fatalf("expected proxy to be called")
	}

	out, ok := ctx.Payload.(*payload.JSON)
	if !ok {
		t.Fatalf("expected payload.JSON, got %T", ctx.Payload)
	}
	if got, ok := out.Data["from"]; !ok || got != "target-via-proxy" {
		t.Fatalf("unexpected response payload: %#v", out.Data)
	}
}

func TestExecuteWithInvalidProxyURL(t *testing.T) {
	step := &Step{
		URL:    "http://example.com",
		Method: http.MethodGet,
		Proxy:  "://bad proxy",
	}

	ctx := &engine.ExecutionContext{
		Payload: &payload.JSON{Data: map[string]interface{}{"a": "b"}},
	}

	err := step.Execute(ctx)
	if err == nil {
		t.Fatal("expected an error for invalid proxy URL")
	}
	if !strings.Contains(err.Error(), "invalid proxy URL") {
		t.Fatalf("unexpected error: %v", err)
	}
}
