package oauth

import (
	http "github.com/llbarbosas/noodle-account/core/infra"
	core "github.com/llbarbosas/noodle-account/core/oauth"
)

type TokenService struct {
	ClientRepository                ClientRepository
	AuthorizationRequestRepository  AuthorizationRequestRepository
	AuthorizationResponseRepository AuthorizationResponseRepository
	AccessTokenResponseRepository   AccessTokenResponseRepository
	ResourceOwnerRepository         ResourceOwnerRepository
}

func (s *TokenService) Execute(r http.Request) http.Response {
	tokenRequest := new(core.AccessTokenRequest)

	if err := r.Body(tokenRequest); err != nil {
		if err := r.Query(tokenRequest); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError(err.Error()))
		}
	}

	var (
		tokenResponse *core.AccessTokenResponse
		err           error
	)

	switch tokenRequest.GrantType {
	case "authorization_code":
		authorizationResponse := new(AuthorizationResponseModel)

		if err := s.AuthorizationResponseRepository.GetByCode(authorizationResponse, tokenRequest.Code); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError("Invalid authorization code"))
		}

		authorizationRequest := new(AuthorizationRequestModel)

		if err := s.AuthorizationRequestRepository.GetByID(authorizationRequest, authorizationResponse.AuthorizationRequestID); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError("Invalid authorization code"))
		}

		resourceOwner := new(core.ResourceOwner)

		if err := s.ResourceOwnerRepository.GetByID(resourceOwner, string(authorizationResponse.ResourceOwnerID)); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError("Invalid authorization"))
		}

		tokenResponse, err = core.GrantTypeAuthorizationCode(tokenRequest, &authorizationRequest.AuthorizationRequest, *resourceOwner)
	case "refresh_token":
		lastAccessTokenResponse := new(core.AccessTokenResponse)

		if err := s.AccessTokenResponseRepository.GetByRefreshToken(lastAccessTokenResponse, tokenRequest.RefreshToken); err != nil {
			return responseFromErrorResponse(core.NewServerError("Cannot connect to database"))
		}

		if err := s.AccessTokenResponseRepository.DeleteByRefreshToken(tokenRequest.RefreshToken); err != nil {
			return responseFromErrorResponse(core.NewServerError("Cannot connect to database"))
		}

		authorizationResponse := new(AuthorizationResponseModel)

		if err := s.AuthorizationResponseRepository.GetByCode(authorizationResponse, tokenRequest.Code); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError("Invalid authorization code"))
		}

		resourceOwner := new(core.ResourceOwner)

		if err := s.ResourceOwnerRepository.GetByID(resourceOwner, string(authorizationResponse.ResourceOwnerID)); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError("Invalid authorization"))
		}

		tokenResponse, err = core.GrantTypeRefreshToken(lastAccessTokenResponse, tokenRequest, *resourceOwner)
	case "client_credentials":
		tokenResponse, err = nil, core.NewInvalidRequestError("Soon")
	default:
		tokenResponse, err = nil, core.NewInvalidRequestError("Invalid grant_type")
	}

	if err != nil {
		return responseFromErrorResponse(err.(core.ErrorResponse))
	}

	if err := s.AccessTokenResponseRepository.Create(tokenResponse); err != nil {
		return responseFromErrorResponse(core.NewServerError("Cannot connect to database"))
	}

	// TODO: Reimplement AccessTokenResponse json handling
	tokenResponseJSON := AccessTokenResponseJSON{
		AccessToken:           tokenResponse.AccessTokenCode,
		TokenType:             tokenResponse.TokenType,
		ExpiresIn:             tokenResponse.ExpiresIn,
		RefreshToken:          tokenResponse.RefreshTokenCode,
		RefreshTokenExpiresIn: tokenResponse.RefreshTokenExpiresIn,
		Scope:                 tokenResponse.Scope,
		IDToken:               tokenResponse.IDToken,
	}

	return http.ResponseJSON(tokenResponseJSON)
}

type AccessTokenResponseJSON struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
	IDToken               string `json:"id_token,omitempty"`
}
