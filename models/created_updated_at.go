package models

// CreatedUpdatedAt define common timestamp fields to several models
type CreatedUpdatedAt struct {
	CreatedAt string `json:"createdAt" pg:"created_at"`
	UpdatedAt string `json:"updatedAt" pg:"updated_at"`
}
