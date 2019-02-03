package errors

import (
	"encoding/json"
	"fmt"
)

// NonAllowedEmailDomainError happens when validation of a struct fails
type NonAllowedEmailDomainError struct {
	domain string
}

// NewNonAllowedEmailDomainError ctor
func NewNonAllowedEmailDomainError(domain string) *NonAllowedEmailDomainError {
	return &NonAllowedEmailDomainError{domain: domain}
}

func (e *NonAllowedEmailDomainError) Error() string {
	return fmt.Sprintf("email from non-allowed hosted domain: %s", e.domain)
}

// Serialize returns the error serialized
func (e *NonAllowedEmailDomainError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-002",
		"error":       "NonAllowedEmailDomainError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}
