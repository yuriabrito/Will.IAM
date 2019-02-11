package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cespare/xxhash"
	"github.com/ghostec/Will.IAM/constants"
	"github.com/ghostec/Will.IAM/errors"
	"github.com/ghostec/Will.IAM/models"
	"github.com/go-redis/redis"
)

// Tokens contract
type Tokens interface {
	FromCache(string) (*models.AuthResult, error)
	ToCache(*models.AuthResult) error
	Get(string) (*models.Token, error)
	Save(*models.Token) error
	Clone() Tokens
	setStorage(*Storage)
}

type tokens struct {
	*withStorage
}

func (ts *tokens) Clone() Tokens {
	return NewTokens(ts.storage.Clone())
}

func buildTokenCacheKey(accessToken string) string {
	return fmt.Sprintf("accessToken-%d", xxhash.Sum64String(accessToken))
}

func (ts tokens) ToCache(auth *models.AuthResult) error {
	if !constants.TokensCacheEnabled {
		return nil
	}
	bts, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	return ts.storage.Redis.Client.Set(
		buildTokenCacheKey(auth.AccessToken), bts,
		time.Duration(constants.TokensCacheTTL)*time.Second,
	).Err()
}

func (ts tokens) FromCache(
	accessToken string,
) (*models.AuthResult, error) {
	if !constants.TokensCacheEnabled {
		return nil, nil
	}
	res := ts.storage.Redis.Client.Get(buildTokenCacheKey(accessToken))
	bts, err := res.Bytes()
	if err != nil && err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	auth := &models.AuthResult{}
	if err := json.Unmarshal(bts, &auth); err != nil {
		return nil, err
	}
	return auth, nil
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
		return nil, errors.NewEntityNotFoundError(models.Token{}, accessToken)
	}
	return t, nil
}

func (ts tokens) Save(token *models.Token) error {
	token.SSOAccessToken = token.AccessToken
	_, err := ts.storage.PG.DB.Exec(`INSERT INTO tokens (access_token,
	refresh_token, sso_access_token, token_type, expiry, email, updated_at)
	VALUES (?access_token, ?refresh_token, ?sso_access_token, ?token_type,
	?expiry, ?email, now()) ON CONFLICT (refresh_token) DO UPDATE SET
	access_token = ?access_token, expiry = ?expiry, updated_at = now()`, token)
	return err
}

// NewTokens ctor
func NewTokens(storage *Storage) Tokens {
	return &tokens{&withStorage{storage: storage}}
}
