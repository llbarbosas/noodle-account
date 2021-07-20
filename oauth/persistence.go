package oauth

import (
	"context"
	"time"

	core "github.com/llbarbosas/noodle-account/core/oauth"
)

type AuthorizationResponseModel struct {
	core.AuthorizationResponse
	AuthorizationRequestID AuthorizationRequestID
	ResourceOwnerID        core.ResourceOwnerID
}

func NewAuthorizationResponseModel(code, state string, requestID AuthorizationRequestID) AuthorizationResponseModel {
	return AuthorizationResponseModel{
		AuthorizationResponse:  core.NewAuthorizationResponse(code, state),
		AuthorizationRequestID: requestID,
	}
}

type AuthorizationRequestID string

type AuthorizationRequestModel struct {
	ID                      AuthorizationRequestID
	CreatedAt               time.Time
	AuthorizationRequest    core.AuthorizationRequest
	AuthorizationRenderData core.AuthorizationRenderData
}

type ClientRepository interface {
	Create(*core.Client) error
	GetByID(*core.Client, core.ClientID) error
}

type AuthorizationRequestRepository interface {
	Create(*AuthorizationRequestModel) error
	GetByID(*AuthorizationRequestModel, AuthorizationRequestID) error
	GetAll(interface{}) error
}

type AuthorizationResponseRepository interface {
	Create(interface{}) error
	GetByCode(*AuthorizationResponseModel, string) error
}

type AccessTokenResponseRepository interface {
	Create(interface{}) error
	GetByAccessToken(*core.AccessTokenResponse, string) error
	GetByRefreshToken(*core.AccessTokenResponse, string) error
	DeleteByRefreshToken(string) error
}

type ResourceOwnerRepository interface {
	// GetByID(*core.ResourceOwner, core.ResourceOwnerID) error
	GetByID(interface{}, string) error
}

type key int

const authorizationRequestKey key = 0

func NewContextWithAuthorizationRequest(authRequest AuthorizationRequestModel) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, authorizationRequestKey, authRequest)
	}
}

func AuthorizationRequestFromContext(ctx context.Context) (AuthorizationRequestModel, bool) {
	authRequest, ok := ctx.Value(authorizationRequestKey).(AuthorizationRequestModel)
	return authRequest, ok
}
