package models

import "time"

// Token type
type Token struct {
	ID             string    `json:"-" pg:"id"`
	AccessToken    string    `json:"access_token" pg:"access_token"`
	RefreshToken   string    `json:"refresh_token" pg:"refresh_token"`
	SSOAccessToken string    `json:"sso_access_token" pg:"sso_access_token"`
	TokenType      string    `json:"token_type" pg:"token_type"`
	Expiry         time.Time `json:"expiry" pg:"expiry"`
	Email          string    `json:"email" pg:"email"`
	CreatedUpdatedAt
}

// AccessTokenAuth stores a ServiceAccountID and the (maybe refreshed)
// AccessToken
type AccessTokenAuth struct {
	ServiceAccountID string
	AccessToken      string
	Email            string
}
