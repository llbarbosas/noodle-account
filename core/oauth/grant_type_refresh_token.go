package oauth

func GrantTypeRefreshToken(lastAccessTokenResponse *AccessTokenResponse, tokenRequest *AccessTokenRequest, resourceOwner ResourceOwner) (*AccessTokenResponse, error) {
	if tokenRequest.RefreshToken == "" {
		return nil, NewInvalidRequestError("refresh_token required")
	}

	accessTokenResponse, err := newAccessTokenResponseFromToken(&lastAccessTokenResponse.AccessToken, lastAccessTokenResponse.Scope, resourceOwner)

	if err != nil {
		return nil, NewServerError("Internal error")
	}

	return accessTokenResponse, nil
}
