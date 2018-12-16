package repositories

import "github.com/ghostec/Will.IAM/models"

// Tokens contract
type Tokens interface {
	Get(string) (*models.Token, error)
	Save(*models.Token) error
}

type tokens struct {
	storage *Storage
}

func (ts tokens) Get(accessToken string) (*models.Token, error) {
	t := &models.Token{}
	if _, err := ts.storage.PG.DB.Query(
		t, "SELECT * FROM tokens WHERE access_token = ?", accessToken,
	); err != nil {
		return nil, err
	}
	return t, nil
}

func (ts tokens) Save(token *models.Token) error {
	_, err := ts.storage.PG.DB.Exec(`INSERT INTO tokens (access_token, refresh_token,
	token_type, expiry, email, updated_at) VALUES (?access_token, ?refresh_token,
	?token_type, ?expiry, ?email, now())`, token)
	return err
}

// NewTokens ctor
func NewTokens(storage *Storage) Tokens {
	return &tokens{storage: storage}
}
