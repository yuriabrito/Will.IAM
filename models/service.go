package models

// Service type
type Service struct {
	ID                      string `json:"id" pg:"id"`
	Name                    string `json:"name" pg:"name"`
	PermissionName          string `json:"permissionName" pg:"permission_name"`
	ServiceAccountID        string `json:"serviceAccountID" pg:"service_account_id"`
	CreatorServiceAccountID string `json:"creatorServiceAccountID" pg:"creator_service_account_id"`
	AMURL                   string `json:"amUrl" sql:"am_url"`
	CreatedUpdatedAt
}

// Validate Service model
func (s Service) Validate(fields ...string) Validation {
	v := &Validation{}
	if s.Name == "" {
		v.AddError("name", "required")
	}
	if s.PermissionName == "" {
		v.AddError("permissionName", "required")
	}
	return *v
}
