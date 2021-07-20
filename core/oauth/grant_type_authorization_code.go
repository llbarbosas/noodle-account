package oauth

func GrantTypeAuthorizationCode(tokenRequest *AccessTokenRequest, authRequest *AuthorizationRequest, subject ResourceOwner) (*AccessTokenResponse, error) {
	switch authRequest.CodeChallengeMethod {
	case "plain":
		if authRequest.CodeChallenge != tokenRequest.CodeVerifier {
			return nil, NewInvalidRequestError("PKCE verification failed")
		}
	case "S256":
		challengeFromVerfier := GenerateCodeChallenge(tokenRequest.CodeVerifier)

		if authRequest.CodeChallenge != challengeFromVerfier {
			return nil, NewInvalidRequestError("PKCE verification failed")
		}
	default:
		return nil, NewServerError("Internal error")
	}

	accessTokenResponse, err := newAccessTokenResponse(subject, authRequest.ClientID, authRequest.Scope)

	if err != nil {
		return nil, NewServerError("Internal error")
	}

	return accessTokenResponse, nil
}
