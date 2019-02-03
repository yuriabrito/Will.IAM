package models

import "encoding/json"

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
