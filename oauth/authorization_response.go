package oauth

import (
	"log"

	http "github.com/llbarbosas/noodle-account/core/infra"
	infra "github.com/llbarbosas/noodle-account/core/infra"
	core "github.com/llbarbosas/noodle-account/core/oauth"
	"github.com/llbarbosas/noodle-account/core/util"
)

type AuthorizationResponseService struct {
	AuthorizationRequestRepository   AuthorizationRequestRepository
	AuthorizationResponseRepository  AuthorizationResponseRepository
	AuthenticatedResourceOwnerGetter core.AuthenticatedResourceOwnerGetter
}

type GrantAuthorizationRequest struct {
	AuthorizationID AuthorizationRequestID `form:"_aid" query:"_aid"`
	Authorize       string                 `form:"authorize" query:"authorize"`
}

func (s *AuthorizationResponseService) Execute(r http.Request) http.Response {
	var authorizationRequest AuthorizationRequestModel

	if authRequest, ok := AuthorizationRequestFromContext(r.Context()); ok {
		authorizationRequest = authRequest
	} else {
		var grantRequest GrantAuthorizationRequest

		if err := r.Body(&grantRequest); err != nil {
			return responseFromErrorResponse(core.NewInvalidRequestError(err.Error()))
		}

		if grantRequest.Authorize == "0" {
			return bindErrorResponseGrantRequest(authorizationRequestGetter(s.AuthorizationRequestRepository), &grantRequest, core.NewAccessDeniedError("User denied access"))
		}

		if err := s.AuthorizationRequestRepository.GetByID(&authorizationRequest, grantRequest.AuthorizationID); err != nil {
			return responseFromErrorResponse(core.NewServerError("Cannot connect to database"))
		}
	}

	code, err := util.RandomString(60)

	if err != nil {
		return responseFromErrorResponse(core.NewServerError(err.Error()).BindAR(authorizationRequest.AuthorizationRequest))
	}

	resourceOwner, err := s.AuthenticatedResourceOwnerGetter(r)

	if err != nil {
		// TODO: Debug
		if err == infra.ErrNotFound {
			return redirectSignin(r.RawQuery())
		}

		return responseFromErrorResponse(core.NewServerError(err.Error()).BindAR(authorizationRequest.AuthorizationRequest))
	}

	authorizationResponse := &AuthorizationResponseModel{
		AuthorizationResponse:  core.NewAuthorizationResponse(code, authorizationRequest.AuthorizationRequest.State),
		AuthorizationRequestID: authorizationRequest.ID,
		ResourceOwnerID:        resourceOwner.ID(),
	}

	if err := s.AuthorizationResponseRepository.Create(authorizationResponse); err != nil {
		return responseFromErrorResponse(core.NewServerError(err.Error()).BindAR(authorizationRequest.AuthorizationRequest))
	}

	return http.ResponseRedirect(authorizationRequest.AuthorizationRequest.RedirectURI + "?" + authorizationResponse.URLQuery())
}

type authorizationGetter func(authID AuthorizationRequestID) (*AuthorizationRequestModel, error)

func authorizationRequestGetter(ar AuthorizationRequestRepository) authorizationGetter {
	return func(authID AuthorizationRequestID) (*AuthorizationRequestModel, error) {
		authorizationRequest := new(AuthorizationRequestModel)

		if err := ar.GetByID(authorizationRequest, authID); err != nil {
			return nil, err
		}

		return authorizationRequest, nil
	}
}

func bindErrorResponseGrantRequest(getter authorizationGetter, grantRequest *GrantAuthorizationRequest, errResponse core.ErrorResponse) http.Response {
	authRequest, err := getter(grantRequest.AuthorizationID)

	if err != nil {
		return responseFromErrorResponse(errResponse)
	}

	log.Println(errResponse.BindAR(authRequest.AuthorizationRequest))

	return responseFromErrorResponse(errResponse.BindAR(authRequest.AuthorizationRequest))
}
