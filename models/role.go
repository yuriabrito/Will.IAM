package models

// Role type
type Role struct {
	ID   string `json:"id" pg:"id"`
	Name string `json:"name" pg:"name"`
	CreatedUpdatedAt
}
