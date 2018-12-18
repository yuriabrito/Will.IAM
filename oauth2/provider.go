package oauth2

// AuthResult is the result of a successful authentication
type AuthResult struct {
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
}

// Provider is the contract any OAuth2 implementation must follow
type Provider interface {
	BuildAuthURL(string) string
	ExchangeCode(string) (*AuthResult, error)
	Authenticate(string) (*AuthResult, error)
}

// ProviderBlankMock is a Provider mock will all dummy implementations
type ProviderBlankMock struct{}

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
	return &AuthResult{
		AccessToken: "any",
		Email:       "any",
	}, nil
}
