package constants

// Actions are all Will.IAM supported permission verbs
var Actions = struct {
	Roles           []string
	ServiceAccounts []string
}{
	Roles: []string{
		"CreateRole",
		"ListRoles",
	},
	ServiceAccounts: []string{
		"CreateServiceAccount",
		"ListServiceAccounts",
	},
}
