package repositories

import (
	"fmt"

	"github.com/ghostec/Will.IAM/models"
)

// Tokens contract
type Tokens interface {
	Get(string) (*models.Token, error)
	Save(*models.Token) error
}

type tokens struct {
	storage *Storage
}

func (ts tokens) Get(accessToken string) (*models.Token, error) {
	t := new(models.Token)
	if _, err := ts.storage.PG.DB.Query(
		t, "SELECT * FROM tokens WHERE access_token = ?0 OR sso_access_token = ?0",
		accessToken,
	); err != nil {
		return nil, err
	}
	if t.AccessToken == "" {
		return nil, fmt.Errorf("access token not found")
	}
	return t, nil
}

func (ts tokens) Save(token *models.Token) error {
	token.SSOAccessToken = token.AccessToken
	_, err := ts.storage.PG.DB.Exec(`INSERT INTO tokens (access_token,
	refresh_token, sso_access_token, token_type, expiry, email, updated_at)
	VALUES (?access_token, ?refresh_token, ?sso_access_token, ?token_type,
	?expiry, ?email, now()) ON CONFLICT (refresh_token) DO UPDATE SET
	access_token = ?access_token, expiry = ?expiry`, token)
	return err
}

// NewTokens ctor
func NewTokens(storage *Storage) Tokens {
	return &tokens{storage: storage}
}
