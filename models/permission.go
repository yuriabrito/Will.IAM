package models

import (
	"fmt"
	"strings"
)

// OwnershipLevel defines the holding rights about a resource
type OwnershipLevel string

// OwnershipLevels are all possible holding rights
// Owner: can exercise the action over the resource and provide the exact same
// ownership:action:resource rights to other parties
// Lender: can only exercise the action over the resource
var OwnershipLevels = struct {
	Owner  OwnershipLevel
	Lender OwnershipLevel
}{
	Owner:  "RO",
	Lender: "RL",
}

// Less returns true if o < oo; Lender < Owner
func (o OwnershipLevel) Less(oo OwnershipLevel) bool {
	if o == OwnershipLevels.Lender && oo == OwnershipLevels.Owner {
		return true
	}
	return false
}

// Action is defined by IAM clients
type Action string

// BuildAction from string
func BuildAction(str string) Action {
	return Action(str)
}

// Service type
type Service string

// BuildService from string
func BuildService(str string) Service {
	return Service(str)
}

// ResourceHierarchy is either a complete or an open hierarchy to something
// Eg:
// Complete: maestro::sniper-3d::na::sniper3d-red
// Open: maestro::sniper-3d::stag::*
type ResourceHierarchy struct {
	Open      bool
	Size      int
	Hierarchy []string `json:"resourceHierarchy" pg:"resource_hierarchy,array"`
}

// BuildResourceHierarchy from string
func BuildResourceHierarchy(parts []string) ResourceHierarchy {
	open := false
	size := len(parts)
	if parts[len(parts)-1] == "*" {
		open = true
	}
	return ResourceHierarchy{Open: open, Size: size, Hierarchy: parts}
}

// Contains checks whether rh.Hierarchy contains orh.Hierarchy
func (rh ResourceHierarchy) Contains(orh ResourceHierarchy) bool {
	if rh.Size > orh.Size {
		return false
	}
	for i := range orh.Hierarchy {
		if rh.Hierarchy[i] == "*" {
			return true
		}
		if orh.Hierarchy[i] != rh.Hierarchy[i] {
			return false
		}
	}
	return true
}

// ToString will join the Hierarchy slice to return a
// resource hierarchy in the format {1}::{2}...
func (rh ResourceHierarchy) ToString() string {
	return strings.Join(rh.Hierarchy, "::")
}

// Permission is bound to a role and
// defines the onwership level of an action over a resource
type Permission struct {
	RoleID            string         `json:"roleId" pg:"role_id"`
	Service           Service        `json:"service" pg:"service"`
	OwnershipLevel    OwnershipLevel `json:"ownershipLevel" pg:"ownership_level"`
	Action            Action         `json:"action" pg:"action"`
	ResourceHierarchy ResourceHierarchy
}

// ValidatePermission validates a permission in string format
func ValidatePermission(str string) (bool, error) {
	// Format: OwnershipLevel::Action::Service::{ResourceHierarchy}
	parts := strings.Split(str, "::")
	if len(parts) < 4 {
		return false, fmt.Errorf(
			"Incomplete permission. Expected format: " +
				"Service::OwnershipLevel::Action::{ResourceHierarchy}",
		)
	}
	ol := OwnershipLevel(parts[1])
	if ol != OwnershipLevels.Owner && ol != OwnershipLevels.Lender {
		return false, fmt.Errorf("OwnershipLevel needs to be RO or RL")
	}
	for _, part := range parts {
		if part == "" {
			return false, fmt.Errorf("No parts can be empty")
		}
	}
	return true, nil
}

// BuildPermission transforms a string into struct
func BuildPermission(str string) (Permission, error) {
	if valid, err := ValidatePermission(str); !valid {
		return Permission{}, err
	}
	parts := strings.Split(str, "::")
	service := BuildService(parts[0])
	ol := OwnershipLevel(parts[1])
	action := BuildAction(parts[2])
	rh := BuildResourceHierarchy(parts[3:])
	return Permission{
		Service:           service,
		OwnershipLevel:    ol,
		Action:            action,
		ResourceHierarchy: rh,
	}, nil
}

// IsPresent checks if a permission is satisfied in a slice
func (p Permission) IsPresent(permissions []Permission) bool {
	for _, pp := range permissions {
		if pp.Action != p.Action || pp.Service != p.Service ||
			pp.OwnershipLevel.Less(p.OwnershipLevel) {
			continue
		}
		if pp.ResourceHierarchy.Contains(p.ResourceHierarchy) {
			return true
		}
	}
	return false
}

// ToString converts a permission to it's equivalent string format
func (p Permission) ToString() string {
	return fmt.Sprintf(
		"%s::%s::%s::%s", p.Service, p.OwnershipLevel, p.Action,
		strings.Join(p.ResourceHierarchy.Hierarchy, "::"),
	)
}
