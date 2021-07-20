package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
)

// AuthorizationRequest represents the constructed request by the client
// described in [section 4.1.1](https://tools.ietf.org/html/rfc6749#section-4.1.1),
// extended by [PKFC](https://tools.ietf.org/html/rfc7636#section-4.3)
//
// "form" tags are used to define that fields will be read using the
// "application / x-www-form-urlencoded" format
type AuthorizationRequest struct {
	ClientID     ClientID `query:"client_id" form:"client_id" validate:"required"`
	RedirectURI  string   `query:"redirect_uri" form:"redirect_uri" validate:"uri|eq="`
	ResponseType string   `query:"response_type" form:"response_type" validate:"required,eq=code"`
	State        string   `query:"state" form:"state"`
	Scope        string   `query:"scope" form:"scope"`

	// PKFC extension
	CodeChallenge       string `query:"code_challenge" form:"code_challenge" validate:"required"`
	CodeChallengeMethod string `query:"code_challenge_method" form:"code_challenge_method" validate:"eq=S256|eq=plain|eq="`
}

func (ar AuthorizationRequest) Valid() error {
	return nil
}

// AuthorizationResponse represents the response sended to client after
// authorization request, described in
// [section 4.1.2](https://tools.ietf.org/html/rfc6749#section-4.1.2)
type AuthorizationResponse struct {
	Code  string
	State string
}

func NewAuthorizationResponse(code, state string) AuthorizationResponse {
	return AuthorizationResponse{
		Code:  code,
		State: state,
	}
}

func (ar AuthorizationResponse) URLQuery() string {
	q := url.Values{}
	q.Set("code", ar.Code)
	q.Set("state", ar.State)

	return q.Encode()
}

func GenerateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
