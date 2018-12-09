package repositories

import "github.com/ghostec/Will.IAM/models"

// Tokens contract
type Tokens interface {
	Get(string) (*models.Token, error)
}

type tokens struct{}

func (t tokens) Get(accessToken string) (*models.Token, error) {
	return nil, nil
}

// NewTokens ctor
func NewTokens() Tokens {
	return &tokens{}
}
