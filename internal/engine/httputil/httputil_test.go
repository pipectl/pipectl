package httputil

import (
	"strings"
	"testing"
	"time"

	"github.com/pipectl/pipectl/internal/engine/payload"
)

func TestMarshalPayloadJSON(t *testing.T) {
	p := &payload.JSON{
		Items: []map[string]any{{"key": "value"}},
		Shape: payload.JSONArrayShape,
	}

	body, ct, err := MarshalPayload(p)
	if err != nil {
		t.Fatalf("MarshalPayload: %v", err)
	}
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
	if !strings.Contains(string(body), "value") {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestMarshalPayloadJSONL(t *testing.T) {
	p := &payload.JSONL{
		Items: []map[string]any{
			{"a": 1},
			{"b": 2},
		},
	}

	body, ct, err := MarshalPayload(p)
	if err != nil {
		t.Fatalf("MarshalPayload: %v", err)
	}
	if ct != "application/x-ndjson" {
		t.Errorf("expected application/x-ndjson, got %q", ct)
	}
	lines := strings.Split(strings.TrimRight(string(body), "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestMarshalPayloadCSV(t *testing.T) {
	p := &payload.CSV{
		Rows: [][]string{
			{"name", "age"},
			{"alice", "30"},
		},
	}

	body, ct, err := MarshalPayload(p)
	if err != nil {
		t.Fatalf("MarshalPayload: %v", err)
	}
	if ct != "text/csv" {
		t.Errorf("expected text/csv, got %q", ct)
	}
	if !strings.Contains(string(body), "alice") {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestResolveTimeoutDefault(t *testing.T) {
	got, err := ResolveTimeout(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != DefaultTimeoutSeconds*time.Second {
		t.Errorf("expected %v, got %v", DefaultTimeoutSeconds*time.Second, got)
	}
}

func TestResolveTimeoutCustom(t *testing.T) {
	got, err := ResolveTimeout(30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 30*time.Second {
		t.Errorf("expected 30s, got %v", got)
	}
}

func TestResolveTimeoutNegative(t *testing.T) {
	_, err := ResolveTimeout(-1)
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
	if !strings.Contains(err.Error(), "invalid timeout") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveTimeoutAboveMax(t *testing.T) {
	_, err := ResolveTimeout(MaxTimeoutSeconds + 1)
	if err == nil {
		t.Fatal("expected error for timeout above max")
	}
	if !strings.Contains(err.Error(), "maximum is 300 seconds") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuildClientNoProxy(t *testing.T) {
	client, err := BuildClient("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestBuildClientInvalidProxy(t *testing.T) {
	_, err := BuildClient("://bad-url")
	if err == nil {
		t.Fatal("expected error for invalid proxy URL")
	}
	if !strings.Contains(err.Error(), "invalid proxy URL") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuildClientWithProxy(t *testing.T) {
	client, err := BuildClient("http://proxy.example.com:8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Transport == nil {
		t.Fatal("expected transport to be set for proxy")
	}
}
