package oauth

import (
	"fmt"
	"strings"
	"time"

	http "github.com/llbarbosas/noodle-account/core/infra"
	core "github.com/llbarbosas/noodle-account/core/oauth"
)

var (
	errInvalidAuthorizationHeader = http.ResponseFromWWWAuthenticate(fmt.Sprintf(`Bearer realm="%s"`, http.MainRealm))
	errExpiratedToken             = http.ResponseFromWWWAuthenticate(fmt.Sprintf(`Bearer realm="%s", error="invalid_token", error_description="The access token expired"`, http.MainRealm))
	errInvalidToken               = http.ResponseFromWWWAuthenticate(fmt.Sprintf(`Bearer realm="%s", error="invalid_token", error_description="Invalid access token"`, http.MainRealm))
)

type UserInfoService struct {
	AccessTokenResponseRepository AccessTokenResponseRepository
	ResourceOwnerRepository       ResourceOwnerRepository
}

func (s *UserInfoService) Execute(r http.Request) http.Response {
	authHeader := r.Header("Authorization")
	authHeaderSplit := strings.Split(authHeader, " ")

	if len(authHeaderSplit) < 2 || authHeaderSplit[0] != "Bearer" {
		return errInvalidAuthorizationHeader
	}

	authAccessToken := authHeaderSplit[1]
	var accessTokenResponse core.AccessTokenResponse

	if err := s.AccessTokenResponseRepository.GetByAccessToken(&accessTokenResponse, authAccessToken); err != nil {
		return errInvalidToken
	}

	if accessTokenResponse.AccessToken.ExpiresAt < time.Now().Unix() {
		return errExpiratedToken
	}

	resourceOwnerID := string(accessTokenResponse.AccessToken.Subject)
	var resourceOwner core.ResourceOwner

	if err := s.ResourceOwnerRepository.GetByID(&resourceOwner, resourceOwnerID); err != nil {
		return errInvalidToken
	}

	scopes := core.ProcessScopes(accessTokenResponse.Scope)

	// TODO: Create core function to generate idToken JSON
	idToken := core.IdTokenClaimsFromAccessToken(accessTokenResponse.AccessToken).BindUserData(resourceOwner, scopes)

	return http.ResponseJSON(idToken.UserInfo)
}
