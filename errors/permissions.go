package errors

import (
	"encoding/json"
	"fmt"
)

// UserDoesntHavePermissionError happens when user doesn't have a permission
type UserDoesntHavePermissionError struct {
	permission string
}

// NewUserDoesntHavePermissionError ctor
func NewUserDoesntHavePermissionError(
	permission string,
) *UserDoesntHavePermissionError {
	return &UserDoesntHavePermissionError{permission: permission}
}

func (e *UserDoesntHavePermissionError) Error() string {
	return fmt.Sprintf("user doesn't have permission %s", e.permission)
}

// Serialize returns the error serialized
func (e *UserDoesntHavePermissionError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-003",
		"error":       "UserDoesntHavePermissionError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}

// StatusCode implements ErrorWithStatusCode
func (e *UserDoesntHavePermissionError) StatusCode() int {
	return 403
}

// UserDoesntHavePermissionsError happens when user doesn't have [] permissions
type UserDoesntHavePermissionsError struct {
	permissions []string
}

// NewUserDoesntHavePermissionsError ctor
func NewUserDoesntHavePermissionsError(
	permissions []string,
) *UserDoesntHavePermissionsError {
	return &UserDoesntHavePermissionsError{permissions: permissions}
}

func (e *UserDoesntHavePermissionsError) Error() string {
	str := ""
	for i := range e.permissions {
		str = fmt.Sprintf("%s %s", str, e.permissions[i])
	}
	return fmt.Sprintf("user doesn't have permissions%s", str)
}

// Serialize returns the error serialized
func (e *UserDoesntHavePermissionsError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-004",
		"error":       "UserDoesntHavePermissionsError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}

// StatusCode implements ErrorWithStatusCode
func (e *UserDoesntHavePermissionsError) StatusCode() int {
	return 403
}

// UserDoesntHaveAllPermissionsError happens when HasAll{X}Permissions is false
type UserDoesntHaveAllPermissionsError struct {
}

// NewUserDoesntHaveAllPermissionsError ctor
func NewUserDoesntHaveAllPermissionsError() *UserDoesntHaveAllPermissionsError {
	return &UserDoesntHaveAllPermissionsError{}
}

func (e *UserDoesntHaveAllPermissionsError) Error() string {
	return "user doesn't have all permissions"
}

// Serialize returns the error serialized
func (e *UserDoesntHaveAllPermissionsError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "ERR-005",
		"error":       "UserDoesntHaveAllPermissionsError",
		"description": e.Error(),
		"success":     false,
	})

	return g
}

// StatusCode implements ErrorWithStatusCode
func (e *UserDoesntHaveAllPermissionsError) StatusCode() int {
	return 403
}
