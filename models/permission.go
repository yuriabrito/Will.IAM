package models

// Ownership defines the holding rights about a resource
type Ownership string

// Ownerships are all possible holding rights
// Owner: can exercise the action over the resource and provide the exact same
// ownership:action:resource rights to other parties
// Lender: can only exercise the action over the resource
var Ownerships = struct {
	Owner  Ownership
	Lender Ownership
}{
	Owner:  "owner",
	Lender: "lender",
}

// Action is defined by IAM clients
type Action string

// Resource is either a complete or an open hierarchy to something
// Eg:
// Complete: maestro::sniper-3d::na::sniper3d-red
// Open: maestro::sniper-3d::stag::*
type Resource struct {
	Hierarchy string
}

// Rule defines the onwership level of an action over a resource
type Rule struct {
	Ownership Ownership
	Action    Action
	Resource  Resource
}
