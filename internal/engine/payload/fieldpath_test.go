package payload

import (
	"reflect"
	"strings"
	"testing"
)

func TestResolveFieldPath(t *testing.T) {
	newRecord := func() map[string]interface{} {
		return map[string]interface{}{
			"name": "Alice",
			"user": map[string]interface{}{
				"email": "alice@example.com",
			},
			"contacts": []interface{}{
				map[string]interface{}{"phone": "0400 000 000"},
			},
		}
	}

	testCases := []struct {
		name     string
		path     string
		expected interface{}
		err      string
	}{
		{name: "top-level key", path: "name", expected: "Alice"},
		{name: "nested key", path: "user.email", expected: "alice@example.com"},
		{name: "array index", path: "contacts[0].phone", expected: "0400 000 000"},
		{name: "index out of range", path: "contacts[1].phone", err: "index 1 out of range"},
		{name: "missing key", path: "user.phone", err: "not found in record"},
		{name: "key through non-object", path: "name.first", err: "expected object before key"},
		{name: "index through non-array", path: "user[0]", err: "expected array before index"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := newRecord()
			value, assign, err := ResolveFieldPath(record, tc.path)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("expected error containing %q, got %v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(value, tc.expected) {
				t.Fatalf("expected %#v, got %#v", tc.expected, value)
			}

			if err := assign("updated"); err != nil {
				t.Fatalf("assign returned error: %v", err)
			}
			if value, _, _ := ResolveFieldPath(record, tc.path); value != "updated" {
				t.Fatalf("expected assigned value %q, got %#v", "updated", value)
			}
		})
	}
}

func TestParseFieldPath(t *testing.T) {
	testCases := []struct {
		path  string
		valid bool
	}{
		{path: "name", valid: true},
		{path: "user.address.street", valid: true},
		{path: "contacts[0].email", valid: true},
		{path: "matrix[0][1]", valid: true},
		{path: "", valid: false},
		{path: "a.", valid: false},
		{path: ".a", valid: false},
		{path: "a..b", valid: false},
		{path: "a[", valid: false},
		{path: "a[]", valid: false},
		{path: "a[-1]", valid: false},
		{path: "a[x]", valid: false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			err := ParseFieldPath(tc.path)
			if tc.valid && err != nil {
				t.Fatalf("expected %q to be valid, got %v", tc.path, err)
			}
			if !tc.valid && err == nil {
				t.Fatalf("expected %q to be invalid", tc.path)
			}
		})
	}
}
