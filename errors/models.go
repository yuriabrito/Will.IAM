package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// EntityNotFoundError happens when validation of a struct fails
type EntityNotFoundError struct {
	typ reflect.Type
	ref interface{}
}

// NewEntityNotFoundError ctor
func NewEntityNotFoundError(typ interface{}, ref interface{}) *EntityNotFoundError {
	return &EntityNotFoundError{
		typ: reflect.TypeOf(typ),
		ref: ref,
	}
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s %#v not found", e.typ.Name(), e.ref)
}

// Serialize returns the error serialized
func (e *EntityNotFoundError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-001",
		"error":       "EntityNotFoundError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}
