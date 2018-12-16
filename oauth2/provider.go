package oauth2

// AuthResult is the result of a successful authentication
type AuthResult struct {
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
}

// Provider is the contract any OAuth2 implementation must follow
type Provider interface {
	BuildAuthURL() string
	ExchangeCode(string) (*AuthResult, error)
	Authenticate(string) (*AuthResult, error)
}
