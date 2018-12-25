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

// All checks if action matches any
func (a Action) All() bool {
	return string(a) == "*"
}

// ResourceHierarchy is either a complete or an open hierarchy to something
// Eg:
// Complete: maestro::sniper-3d::na::sniper3d-red
// Open: maestro::sniper-3d::stag::*
type ResourceHierarchy string

// All checks if rh matches any
func (rh ResourceHierarchy) All() bool {
	return string(rh) == "*"
}

type resourceHierarchy struct {
	size      int
	hierarchy []string
}

func buildResourceHierarchy(rh ResourceHierarchy) resourceHierarchy {
	parts := strings.Split(string(rh), "::")
	size := len(parts)
	return resourceHierarchy{size: size, hierarchy: parts}
}

// Contains checks whether rh.Hierarchy contains orh.Hierarchy
func (rh ResourceHierarchy) Contains(orh ResourceHierarchy) bool {
	rhh := buildResourceHierarchy(rh)
	orhh := buildResourceHierarchy(orh)
	if rhh.size > orhh.size {
		return false
	}
	for i := range orhh.hierarchy {
		if rhh.hierarchy[i] == "*" {
			return true
		}
		if orhh.hierarchy[i] != rhh.hierarchy[i] {
			return false
		}
	}
	return true
}

// Permission is bound to a role and
// defines the onwership level of an action over a resource
type Permission struct {
	RoleID            string            `json:"roleId" pg:"role_id"`
	Service           string            `json:"service" pg:"service"`
	OwnershipLevel    OwnershipLevel    `json:"ownershipLevel" pg:"ownership_level"`
	Action            Action            `json:"action" pg:"action"`
	ResourceHierarchy ResourceHierarchy `json:"resourceHierarchy" pg:"resource_hierarchy"`
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
	service := parts[0]
	ol := OwnershipLevel(parts[1])
	action := BuildAction(parts[2])
	rh := ResourceHierarchy(strings.Join(parts[3:], "::"))
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
		if (pp.Service != "*" && pp.Service != p.Service) ||
			(pp.Action != "*" && pp.Action != p.Action) ||
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
		string(p.ResourceHierarchy),
	)
}

// HasServiceFullAccess checks if permission allows it's role
// to execute any action over any resourch hierarchy under it's service
func (p Permission) HasServiceFullAccess() bool {
	if !p.Action.All() {
		return false
	}
	return p.ResourceHierarchy.All()
}

// HasServiceFullOwnership checks if permission allows it's role
// to execute any action over any resourch hierarchy under it's service
// and it's RO
func (p Permission) HasServiceFullOwnership() bool {
	return p.HasServiceFullAccess() && p.OwnershipLevel == OwnershipLevels.Owner
}

// BuildWillIAMPermissionStr builds a permission in the format expected
// by WillIAM handlers
func BuildWillIAMPermissionStr(ro OwnershipLevel, action, rh string) string {
	return fmt.Sprintf("WillIAM::%s::%s::%s", string(ro), action, rh)
}
