package spec

import "fmt"

// Defaults holds pipeline-level default values, grouped by concern, that are
// inherited by steps which don't set the corresponding option themselves.
type Defaults struct {
	HTTP *HTTPDefaults `yaml:"http,omitempty"`
	Text *TextDefaults `yaml:"text,omitempty"`
}

// HTTPDefaults are inherited by http-request and http-transform steps.
type HTTPDefaults struct {
	Proxy   string            `yaml:"proxy,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
	Timeout int               `yaml:"timeout,omitempty"`
}

// TextDefaults are inherited by assert, filter, and dedupe steps.
type TextDefaults struct {
	CaseSensitive *bool `yaml:"case-sensitive,omitempty"`
}

func (d *Defaults) Validate() error {
	if d.HTTP != nil {
		if d.HTTP.Timeout < minTimeoutSeconds {
			return fmt.Errorf("defaults.http timeout must be >= %d", minTimeoutSeconds)
		}
		if d.HTTP.Timeout > maxTimeoutSeconds {
			return fmt.Errorf("defaults.http timeout must be <= %d seconds", maxTimeoutSeconds)
		}
	}
	return nil
}

// applyDefaults merges d onto each step's unset fields. An explicit step-level
// value always takes precedence over the matching default.
func applyDefaults(steps []StepWrapper, d *Defaults) {
	for _, sw := range steps {
		switch s := sw.Step.(type) {
		case *HTTPRequestStep:
			s.Proxy, s.Headers, s.Timeout = mergeHTTPDefaults(s.Proxy, s.Headers, s.Timeout, d)
		case *HTTPTransformStep:
			s.Proxy, s.Headers, s.Timeout = mergeHTTPDefaults(s.Proxy, s.Headers, s.Timeout, d)
		case *AssertStep:
			s.CaseSensitive = resolveCaseSensitive(s.CaseSensitive, d)
		case *FilterStep:
			s.CaseSensitive = resolveCaseSensitive(s.CaseSensitive, d)
		case *DedupeStep:
			s.CaseSensitive = resolveCaseSensitive(s.CaseSensitive, d)
		}
	}
}

func mergeHTTPDefaults(proxy string, headers map[string]string, timeout int, d *Defaults) (string, map[string]string, int) {
	if d == nil || d.HTTP == nil {
		return proxy, headers, timeout
	}

	if proxy == "" {
		proxy = d.HTTP.Proxy
	}

	if timeout == 0 {
		timeout = d.HTTP.Timeout
	}

	if len(d.HTTP.Headers) > 0 {
		merged := make(map[string]string, len(d.HTTP.Headers)+len(headers))
		for k, v := range d.HTTP.Headers {
			merged[k] = v
		}
		for k, v := range headers {
			merged[k] = v
		}
		headers = merged
	}

	return proxy, headers, timeout
}

// resolveCaseSensitive returns stepValue if set, else d's text default if set,
// else true. The result is always non-nil.
func resolveCaseSensitive(stepValue *bool, d *Defaults) *bool {
	effective := true
	if d != nil && d.Text != nil && d.Text.CaseSensitive != nil {
		effective = *d.Text.CaseSensitive
	}
	if stepValue != nil {
		effective = *stepValue
	}
	return &effective
}
