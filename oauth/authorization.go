package oauth

import (
	http "github.com/llbarbosas/noodle-account/core/infra"
	infra "github.com/llbarbosas/noodle-account/core/infra"
	core "github.com/llbarbosas/noodle-account/core/oauth"
)

type AuthorizationService struct {
	ClientRepository                 ClientRepository
	AuthorizationRequestRepository   AuthorizationRequestRepository
	AuthenticatedResourceOwnerGetter core.AuthenticatedResourceOwnerGetter
	AuthorizationResponseService     *AuthorizationResponseService
}

func (s *AuthorizationService) Execute(r infra.Request) infra.Response {
	authRequest := new(core.AuthorizationRequest)

	if err := r.Body(authRequest); err != nil {
		if err := r.Query(authRequest); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError("Malformed request"))
		}
	}

	resourceOwner, err := s.AuthenticatedResourceOwnerGetter(r)

	if err != nil {
		if err == infra.ErrNotFound {
			return redirectSignin(r.RawQuery())
		}

		return bindErrorResponseAuthRequest(clientRedirectURIGetter(s.ClientRepository), authRequest, core.NewServerError("Cannot get request data"))
	}

	client := new(core.Client)

	if err := s.ClientRepository.GetByID(client, authRequest.ClientID); err != nil {
		return bindErrorResponseAuthRequest(authorizationRedirectURIGetter(authRequest), authRequest, core.NewServerError("Unable to get client data"))
	}

	authRenderData, err := core.AuthorizationCodeGrant(authRequest, client)

	if err != nil {
		errResponse := err.(core.ErrorResponse)

		return responseFromErrorResponse(errResponse.Bind(authRequest.State, client.RedirectURI))
	}

	modelAuthRequest := &AuthorizationRequestModel{
		AuthorizationRequest:    *authRequest,
		AuthorizationRenderData: *authRenderData,
	}

	if err := s.AuthorizationRequestRepository.Create(modelAuthRequest); err != nil {
		return responseFromErrorResponse(core.NewServerError(err.Error()).
			Bind(authRequest.State, client.RedirectURI))
	}

	if client.IsInternal() {
		request := r.WithContext(NewContextWithAuthorizationRequest(*modelAuthRequest))

		return s.AuthorizationResponseService.Execute(request)
	}

	bindData := map[string]interface{}{
		"UserName":   resourceOwner.GivenName(),
		"ClientName": modelAuthRequest.AuthorizationRenderData.ClientName,
		"Scopes":     modelAuthRequest.AuthorizationRenderData.Scopes,
		"AID":        modelAuthRequest.ID,
	}

	return infra.ResponseFile("authorize", bindData)
}

type redirectURIGetter func(clientID core.ClientID) (string, error)

func clientRedirectURIGetter(cr ClientRepository) redirectURIGetter {
	return func(clientID core.ClientID) (string, error) {
		client := new(core.Client)

		if err := cr.GetByID(client, clientID); err != nil {
			return "", err
		}

		return client.RedirectURI, nil
	}
}

func authorizationRedirectURIGetter(ar *core.AuthorizationRequest) redirectURIGetter {
	return func(clientID core.ClientID) (string, error) {
		return ar.RedirectURI, nil
	}
}

func bindErrorResponseAuthRequest(getter redirectURIGetter, authRequest *core.AuthorizationRequest, errResponse core.ErrorResponse) http.Response {
	redirectURI, err := getter(authRequest.ClientID)

	if err != nil {
		return responseFromErrorResponse(errResponse.
			BindAR(*authRequest))
	}

	return responseFromErrorResponse(errResponse.Bind(authRequest.State, redirectURI))
}
