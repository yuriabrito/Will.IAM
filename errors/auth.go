package errors

import (
	"encoding/json"
	"fmt"
)

// InvalidAuthorizationTypeError happens when authorization sent is
// neither Bearer nor KeyPair
type InvalidAuthorizationTypeError struct{}

// NewInvalidAuthorizationTypeError ctor
func NewInvalidAuthorizationTypeError() *InvalidAuthorizationTypeError {
	return &InvalidAuthorizationTypeError{}
}

func (e *InvalidAuthorizationTypeError) Error() string {
	return fmt.Sprintf("Invalid authorization header")
}

// Serialize returns the error serialized
func (e *InvalidAuthorizationTypeError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-006",
		"error":       "InvalidAuthorizationTypeError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}
