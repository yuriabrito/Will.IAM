package errors

import (
	"encoding/json"
	"fmt"
)

// InvalidPageError happens when page sent is not an integer
type InvalidPageError struct {
	str string
}

// NewInvalidPageError ctor
func NewInvalidPageError(str string) *InvalidPageError {
	return &InvalidPageError{
		str: str,
	}
}

func (e *InvalidPageError) Error() string {
	return fmt.Sprintf("%s is not a valid page number", e.str)
}

// Serialize returns the error serialized
func (e *InvalidPageError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-007",
		"error":       "InvalidPageError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}

// InvalidPageSizeError happens when page sent is not an integer
type InvalidPageSizeError struct {
	str string
}

// NewInvalidPageSizeError ctor
func NewInvalidPageSizeError(str string) *InvalidPageSizeError {
	return &InvalidPageSizeError{
		str: str,
	}
}

func (e *InvalidPageSizeError) Error() string {
	return fmt.Sprintf("%s is not a valid page size", e.str)
}

// Serialize returns the error serialized
func (e *InvalidPageSizeError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-008",
		"error":       "InvalidPageSizeError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}
