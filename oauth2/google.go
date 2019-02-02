package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ghostec/Will.IAM/models"
	"github.com/ghostec/Will.IAM/repositories"
	extensionsHttp "github.com/topfreegames/extensions/http"
)

const tokenEndpoint = "https://www.googleapis.com/oauth2/v4/token"
const userEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo"

// GoogleConfig are the basic required informations to use Google
// as oauth2 provider
type GoogleConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	HostedDomains []string
}

var googleConfig GoogleConfig

func buildURL(endpoint, queryStrings string) string {
	return fmt.Sprintf("%s?%s", endpoint, queryStrings)
}

func mapToQueryStrings(m map[string]string) string {
	s := []string{}
	for k, v := range m {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(s, "&")
}

// Google implements Provider
type Google struct {
	config GoogleConfig
	repo   *repositories.All
	client *http.Client
}

// BuildAuthURL returns an URL authenticate with Google
func (g *Google) BuildAuthURL(state string) string {
	qs := mapToQueryStrings(map[string]string{
		"state":        state,
		"redirect_uri": g.config.RedirectURL,
		"client_id":    g.config.ClientID,
		"scope": strings.Join([]string{
			url.QueryEscape("https://www.googleapis.com/auth/userinfo.profile"),
			url.QueryEscape("https://www.googleapis.com/auth/userinfo.email"),
		}, "+"),
		"access_type":            "offline",
		"include_granted_scopes": "true",
		"response_type":          "code",
		"prompt":                 "consent",
	})
	return buildURL("https://accounts.google.com/o/oauth2/v2/auth", qs)
}

func (g *Google) buildExchangeCodeForm(code string) string {
	v := url.Values{}
	v.Add("code", code)
	v.Add("client_id", g.config.ClientID)
	v.Add("client_secret", g.config.ClientSecret)
	v.Add("redirect_uri", g.config.RedirectURL)
	v.Add("grant_type", "authorization_code")
	return v.Encode()
}

// ExchangeCode will trade code for full token with Google
func (g *Google) ExchangeCode(code string) (*AuthResult, error) {
	t, err := g.tokenFromCode(code)
	if err != nil {
		return nil, err
	}
	userInfo, err := g.getUserInfo(t.AccessToken)
	if err != nil {
		return nil, err
	}
	allowed := g.checkHostedDomain(userInfo.HostedDomain)
	if !allowed {
		return nil, fmt.Errorf(
			"email from non-allowed hosted domain %s", userInfo.HostedDomain,
		)
	}
	t.Email = userInfo.Email
	if err := g.repo.Tokens.Save(t); err != nil {
		return nil, err
	}
	return &AuthResult{
		AccessToken: t.AccessToken,
		Email:       t.Email,
		Picture:     userInfo.Picture,
	}, nil
}

func (g *Google) postToTokenEndpoint(
	urlencoded string,
) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(urlencoded))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	tmap := map[string]interface{}{}
	err = json.Unmarshal(body, &tmap)
	if err != nil {
		return nil, err
	}
	return tmap, nil
}

func (g *Google) tokenFromCode(code string) (*models.Token, error) {
	ecf := g.buildExchangeCodeForm(code)
	tmap, err := g.postToTokenEndpoint(ecf)
	if err != nil {
		return nil, err
	}
	return &models.Token{
		AccessToken:  tmap["access_token"].(string),
		RefreshToken: tmap["refresh_token"].(string),
		TokenType:    tmap["token_type"].(string),
		Expiry: time.Now().UTC().Add(
			time.Second * time.Duration(tmap["expires_in"].(float64)),
		),
	}, nil
}

type userInfo struct {
	Email        string `json:"email"`
	HostedDomain string `json:"hd"`
	Picture      string `json:"picture"`
}

func (g *Google) getUserInfo(accessToken string) (*userInfo, error) {
	req, err := http.NewRequest("GET", userEndpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	if err != nil {
		return nil, err
	}
	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	ui := &userInfo{}
	err = json.Unmarshal(body, ui)
	if err != nil {
		return nil, err
	}
	return ui, nil
}

func (g *Google) checkHostedDomain(hd string) bool {
	if g.config.HostedDomains == nil || len(g.config.HostedDomains) == 0 {
		return true
	}
	for _, allowed := range g.config.HostedDomains {
		if hd == allowed {
			return true
		}
	}
	return false
}

func (g *Google) buildRefreshTokenForm(refreshToken string) string {
	v := url.Values{}
	v.Add("refresh_token", refreshToken)
	v.Add("client_id", g.config.ClientID)
	v.Add("client_secret", g.config.ClientSecret)
	v.Add("grant_type", "refresh_token")
	return v.Encode()
}

func (g *Google) maybeRefresh(t *models.Token) (*userInfo, error) {
	if t.Expiry.After(time.Now().UTC()) {
		return nil, nil
	}
	rtf := g.buildRefreshTokenForm(t.RefreshToken)
	tmap, err := g.postToTokenEndpoint(rtf)
	if err != nil {
		return nil, err
	}
	t.AccessToken = tmap["access_token"].(string)
	t.Expiry = time.Now().UTC().Add(
		time.Second * time.Duration(tmap["expires_in"].(float64)),
	)
	userInfo, err := g.getUserInfo(t.AccessToken)
	if err != nil {
		return nil, err
	}
	if err = g.repo.Tokens.Save(t); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// Authenticate verifies if an accessToken is valid and maybe refresh it
func (g *Google) Authenticate(accessToken string) (*AuthResult, error) {
	t, err := g.repo.Tokens.Get(accessToken)
	if err != nil {
		return nil, err
	}
	var userInfo *userInfo
	if userInfo, err = g.maybeRefresh(t); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	authResult := &AuthResult{
		AccessToken: t.AccessToken,
		Email:       t.Email,
	}
	if userInfo != nil {
		authResult.Picture = userInfo.Picture
	}
	return authResult, nil
}

// WithContext returns a new instance of *Google using ctx
func (g Google) WithContext(ctx context.Context) Provider {
	return NewGoogle(g.config, g.repo.WithContext(ctx))
}

// NewGoogle ctor
func NewGoogle(
	config GoogleConfig, repo *repositories.All,
) *Google {
	return &Google{
		config: config,
		repo:   repo,
		client: extensionsHttp.New(),
	}
}
