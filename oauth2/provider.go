package oauth2

import "context"

// AuthResult is the result of a successful authentication
type AuthResult struct {
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
	Picture     string `json:"picture"`
}

// Provider is the contract any OAuth2 implementation must follow
type Provider interface {
	BuildAuthURL(string) string
	ExchangeCode(string) (*AuthResult, error)
	Authenticate(string) (*AuthResult, error)
	WithContext(context.Context) Provider
}

// ProviderBlankMock is a Provider mock will all dummy implementations
type ProviderBlankMock struct {
	Email string
}

// NewProviderBlankMock ctor
func NewProviderBlankMock() *ProviderBlankMock {
	return &ProviderBlankMock{}
}

// BuildAuthURL dummy
func (p *ProviderBlankMock) BuildAuthURL(any string) string {
	return "any"
}

// ExchangeCode dummy
func (p *ProviderBlankMock) ExchangeCode(any string) (*AuthResult, error) {
	return &AuthResult{
		AccessToken: "any",
		Email:       "any",
	}, nil
}

// Authenticate dummy
func (p *ProviderBlankMock) Authenticate(any string) (*AuthResult, error) {
	email := "any@email.com"
	if p.Email != "" {
		email = p.Email
	}
	return &AuthResult{
		AccessToken: any,
		Email:       email,
	}, nil
}

// WithContext does nothing
func (p *ProviderBlankMock) WithContext(ctx context.Context) Provider {
	return p
}
