package httptransform

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shanebell/pipectl/internal/engine"
	"github.com/shanebell/pipectl/internal/payload"
)

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
