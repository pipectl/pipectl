package spec

import (
	"fmt"
	"strings"
)

var validHTTPMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// validateHTTPMethod normalizes method (trim + uppercase) and checks it against
// validHTTPMethods, returning errors prefixed with stepType to match each step's
// own error conventions.
func validateHTTPMethod(stepType, method string) (string, error) {
	method = strings.ToUpper(strings.TrimSpace(method))

	if method == "" {
		return "", fmt.Errorf("%s method is required", stepType)
	}

	for _, m := range validHTTPMethods {
		if method == m {
			return method, nil
		}
	}

	return "", fmt.Errorf("%s method must be one of: %s", stepType, strings.Join(validHTTPMethods, ", "))
}
