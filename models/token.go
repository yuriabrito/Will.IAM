package models

// Token type
type Token struct {
	AccessToken  string `json:"access_token" pg:"access_token"`
	RefreshToken string `json:"refresh_token" pg:"refresh_token"`
	TokenType    string `json:"token_type" pg:"token_type"`
	Expiry       string `json:"expiry" pg:"expiry"`
	Email        string `json:"email" pg:"email"`
}
