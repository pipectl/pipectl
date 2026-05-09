package httprequest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pipectl/pipectl/internal/engine"
	"github.com/pipectl/pipectl/internal/engine/httputil"
	"github.com/pipectl/pipectl/internal/engine/payload"
)

type Step struct {
	payload.AllFormatsSupport
	URL     string
	Method  string
	Proxy   string
	Headers map[string]string
	Timeout int
}

func (s *Step) Name() string {
	return "http-request"
}

func (s *Step) Execute(ctx *engine.ExecutionContext) error {
	ctx.Logger.Debug("  %s %s", strings.ToUpper(s.Method), s.URL)
	if s.Proxy != "" {
		ctx.Logger.Debug("  proxy: %s", s.Proxy)
	}

	if err := s.sendRequest(ctx.Payload); err != nil {
		return err
	}

	return nil
}

func (s *Step) sendRequest(inputPayload payload.Payload) error {
	method := strings.ToUpper(s.Method)

	var bodyBytes []byte
	var autoContentType string
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete {
		var err error
		bodyBytes, autoContentType, err = httputil.MarshalPayload(inputPayload)
		if err != nil {
			return fmt.Errorf("http-request could not encode payload: %w", err)
		}
	}

	var bodyReader io.Reader
	if bodyBytes != nil {
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(method, s.URL, bodyReader)
	if err != nil {
		return err
	}

	requestTimeout, err := httputil.ResolveTimeout(s.Timeout)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(req.Context(), requestTimeout)
	defer cancel()
	req = req.WithContext(timeoutCtx)

	for key, value := range s.Headers {
		req.Header.Set(key, value)
	}
	if autoContentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", autoContentType)
	}

	client, err := httputil.BuildClient(s.Proxy)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
