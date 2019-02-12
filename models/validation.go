package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Validation holds errors and whether a model is valid or not
type Validation struct {
	jsonErrors map[string]string
}

// Valid returns wheter a model is valid or not
func (v Validation) Valid() bool {
	return v.jsonErrors == nil
}

// Errors of validations as a json []byte { "errors": { ... } }
func (v Validation) Errors() []byte {
	bts, _ := json.Marshal(map[string]interface{}{
		"errors": v.jsonErrors,
	})
	return bts
}

// AddError in key,value format
func (v *Validation) AddError(key, value string) {
	if v.jsonErrors == nil {
		v.jsonErrors = map[string]string{}
	}
	v.jsonErrors[key] = value
}

func (v Validation) Error() error {
	if v.jsonErrors == nil {
		return nil
	}
	strs := []string{}
	for k, v := range v.jsonErrors {
		strs = append(strs, fmt.Sprintf("%s: %s", k, v))
	}
	return fmt.Errorf(strings.Join(strs, "\t"))
}
