package oauth

type AuthorizationRenderData struct {
	ClientName string
	Scopes     map[string]string
}

func AuthorizationCodeGrant(authorizationRequest *AuthorizationRequest, client *Client) (*AuthorizationRenderData, error) {
	redirectURI := authorizationRequest.RedirectURI

	if err := authorizationRequest.Valid(); err != nil {
		return nil, NewInvalidRequestError("Invalid authorization request: "+err.Error()).
			Bind(authorizationRequest.State, "")
	}

	if redirectURI == "" {
		authorizationRequest.RedirectURI = client.RedirectURI
	} else {
		if client.RedirectURI != authorizationRequest.RedirectURI {
			return nil, NewInvalidRequestError("The redirect_uri isn't valid for this client").
				Bind(authorizationRequest.State, "")
		}
	}

	authorizationRenderData := AuthorizationRenderData{
		ClientName: client.Label,
		Scopes:     getScopeDescriptions(authorizationRequest.Scope),
	}

	return &authorizationRenderData, nil
}

func getScopeDescriptions(scopes string) map[string]string {
	desc := make(map[string]string)

	for _, value := range ProcessScopes(scopes) {
		desc[value] = ScopeDescriptions[value]
	}

	return desc
}
