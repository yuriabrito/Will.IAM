package models

import "time"

// Token type
type Token struct {
	ID           string    `json:"-" pg:"id"`
	AccessToken  string    `json:"access_token" pg:"access_token"`
	RefreshToken string    `json:"refresh_token" pg:"refresh_token"`
	TokenType    string    `json:"token_type" pg:"token_type"`
	Expiry       time.Time `json:"expiry" pg:"expiry"`
	Email        string    `json:"email" pg:"email"`
	CreatedUpdatedAt
}
