package oauth

import (
	"encoding/base64"
	"os"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/llbarbosas/noodle-account/core/util"
)

type UserInfo struct {
	Subject             ResourceOwnerID     `json:"sub"`
	Name                string              `json:"name,omitempty"`
	GivenName           string              `json:"given_name,omitempty"`
	LastName            string              `json:"last_name,omitempty"`
	Picture             string              `json:"picture,omitempty"`
	Email               string              `json:"email,omitempty"`
	EmailVerified       bool                `json:"email_verified,omitempty"`
	PhoneNumber         string              `json:"phone_number,omitempty"`
	PhoneNumberVerified bool                `json:"phone_number_verified,omitempty"`
	Roles               []ResourceOwnerRole `json:"roles,omitempty"`
	UpdatedAt           int64               `json:"updated_at,omitempty"`
}

type IDToken struct {
	*jwt.StandardClaims
	Issuer   string   `json:"iss"`
	Audience []string `json:"aud"`
	IssuedAt int64    `json:"iat"`
	JTI      string   `json:"jti"`
	UserInfo
}

func (it IDToken) BindUserData(user ResourceOwner, scopes []string) IDToken {
	t := IDToken{
		StandardClaims: it.StandardClaims,
		Issuer:         it.Issuer,
		Audience:       it.Audience,
		IssuedAt:       it.IssuedAt,
		JTI:            it.JTI,
	}

	t.StandardClaims.Subject = it.StandardClaims.Subject
	t.UserInfo.Subject = it.UserInfo.Subject
	t.UpdatedAt = user.UpdatedAt()
	t.Roles = user.Roles()

	// TODO: Do this filter on API return
	if util.Contains(scopes, "profile") {
		t.Name = user.Name()
		t.GivenName = user.GivenName()
		t.LastName = user.Name()
		t.Picture = user.Picture()
	}

	if util.Contains(scopes, "email") {
		t.Email = user.Email()
		t.EmailVerified = user.EmailVerified()
	}

	if util.Contains(scopes, "phone") {
		t.PhoneNumber = user.PhoneNumber()
		t.PhoneNumberVerified = user.PhoneNumberVerified()
	}

	return t
}

func (it IDToken) JWT() (*jwt.Token, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims = &it

	return token, nil
}

func (it IDToken) JWTString() (string, error) {
	token, err := it.JWT()

	if err != nil {
		return "", err
	}

	signBytes, err := base64.RawURLEncoding.DecodeString(os.Getenv("JWT_PRIVATEKEY"))

	if err != nil {
		return "", err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		return "", err
	}

	return token.SignedString(signKey)
}

func IdTokenClaimsFromAccessToken(at AccessToken) IDToken {
	return IDToken{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: at.ExpiresAt,
			Subject:   string(at.Subject),
		},
		Issuer:   at.Issuer,
		Audience: at.Audience,
		IssuedAt: at.IssuedAt,
		JTI:      at.JTI,
		UserInfo: UserInfo{
			Subject: at.Subject,
		},
	}
}
