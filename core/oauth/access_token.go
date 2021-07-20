package oauth

import (
	"os"
	"strings"
	"time"

	"github.com/llbarbosas/noodle-account/core/util"
)

var (
	tokenSignKey                  = []byte(os.Getenv("JWT_SECRET"))
	defaultAccessTokenExpiration  = time.Hour * 24
	defaultRefreshTokenExpiration = time.Hour * 72
)

// AccessTokenRequest ...
// described in [section 4.1.3](https://tools.ietf.org/html/rfc6749#section-4.1.3)
// extended by [PKFC](https://tools.ietf.org/html/rfc7636#section-4.5)
type AccessTokenRequest struct {
	GrantType string `query:"grant_type" form:"grant_type"`

	// "authorization_code" grant_type
	Code        string `query:"code" form:"code"`
	RedirectURI string `query:"redirect_uri" form:"redirect_uri"`
	ClientID    string `query:"client_id" form:"client_id"`

	// "refresh_token" grant_type
	RefreshToken string `query:"refresh_token" form:"refresh_token"`
	Scope        string `query:"scope" form:"scope"`

	// "client_credentials" grant_type
	// ClientID     string `query:"client_id" form:"client_id"`
	ClientSecret string `query:"client_secret" form:"client_secret"`

	// PKFC extention
	CodeVerifier string `query:"code_verifier" form:"code_verifier"`
}

// AccessTokenResponse ...
// described in [section 5.1](https://tools.ietf.org/html/rfc6749#section-5.1)
type AccessTokenResponse struct {
	AccessToken           AccessToken
	AccessTokenCode       string `json:"access_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshTokenCode      string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
	IDToken               string `json:"id_token,omitempty"`
}

type RefreshToken struct {
	Code      string
	ExpiresAt int64
}

type AccessToken struct {
	Code      string
	Issuer    string
	Audience  []string
	Subject   ResourceOwnerID
	ClientID  ClientID
	IssuedAt  int64
	JTI       string
	ExpiresAt int64
	Scope     []string
	RefreshToken
}

func newAccessToken(subject ResourceOwnerID, clientID ClientID, scopes []string) (*AccessToken, error) {
	code, err := util.RandomString(100)

	if err != nil {
		return nil, err
	}

	refreshTokenCode, err := util.RandomString(100)

	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	refreshToken := RefreshToken{
		Code:      refreshTokenCode,
		ExpiresAt: currentTime.Add(defaultRefreshTokenExpiration).Unix(),
	}

	accessToken := AccessToken{
		Code:         code,
		ExpiresAt:    currentTime.Add(defaultAccessTokenExpiration).Unix(),
		Issuer:       "https://auth.noodle.com",
		Audience:     []string{},
		Subject:      subject,
		IssuedAt:     currentTime.Unix(),
		JTI:          util.NewUUIDStr(),
		ClientID:     clientID,
		RefreshToken: refreshToken,
	}

	return &accessToken, nil
}

func newAccessTokenResponse(subject ResourceOwner, clientID ClientID, scope string) (*AccessTokenResponse, error) {
	scopes := ProcessScopes(scope)
	scopesStr := strings.Join(scopes, " ")

	accessToken, err := newAccessToken(subject.ID(), clientID, scopes)

	if err != nil {
		return nil, err
	}

	var (
		idTokenString string
	)

	if util.Contains(scopes, "openid") {
		idToken := IdTokenClaimsFromAccessToken(*accessToken).BindUserData(subject, scopes)
		idTokenString, err = idToken.JWTString()

		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	accessTokenResponse := AccessTokenResponse{
		AccessToken:           *accessToken,
		AccessTokenCode:       accessToken.Code,
		TokenType:             "bearer",
		ExpiresIn:             defaultAccessTokenExpiration.Milliseconds(),
		RefreshTokenCode:      accessToken.RefreshToken.Code,
		RefreshTokenExpiresIn: defaultRefreshTokenExpiration.Milliseconds(),
		Scope:                 scopesStr,
		IDToken:               idTokenString,
	}

	return &accessTokenResponse, nil
}

func newAccessTokenResponseFromToken(accessToken *AccessToken, scope string, resourceOwner ResourceOwner) (*AccessTokenResponse, error) {
	return newAccessTokenResponse(resourceOwner, accessToken.ClientID, scope)
}
