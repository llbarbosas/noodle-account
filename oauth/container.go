package oauth

import (
	core "github.com/llbarbosas/noodle-account/core/oauth"
)

type ContainerDependencies struct {
	ClientRepository                 ClientRepository
	AuthorizationRequestRepository   AuthorizationRequestRepository
	AuthorizationResponseRepository  AuthorizationResponseRepository
	AccessTokenResponseRepository    AccessTokenResponseRepository
	ResourceOwnerRepository          ResourceOwnerRepository
	AuthenticatedResourceOwnerGetter core.AuthenticatedResourceOwnerGetter
}

type Container struct {
	ContainerDependencies
	TokenService                 *TokenService
	AuthorizationService         *AuthorizationService
	AuthorizationResponseService *AuthorizationResponseService
	UserInfoService              *UserInfoService
}

func CreateContainer(deps ContainerDependencies) *Container {
	authorizationResponseService := &AuthorizationResponseService{
		AuthorizationRequestRepository:   deps.AuthorizationRequestRepository,
		AuthorizationResponseRepository:  deps.AuthorizationResponseRepository,
		AuthenticatedResourceOwnerGetter: deps.AuthenticatedResourceOwnerGetter,
	}

	return &Container{
		ContainerDependencies: deps,
		TokenService: &TokenService{
			ClientRepository:                deps.ClientRepository,
			AuthorizationRequestRepository:  deps.AuthorizationRequestRepository,
			AuthorizationResponseRepository: deps.AuthorizationResponseRepository,
			AccessTokenResponseRepository:   deps.AccessTokenResponseRepository,
			ResourceOwnerRepository:         deps.ResourceOwnerRepository,
		},
		AuthorizationResponseService: authorizationResponseService,
		AuthorizationService: &AuthorizationService{
			ClientRepository:                 deps.ClientRepository,
			AuthorizationRequestRepository:   deps.AuthorizationRequestRepository,
			AuthorizationResponseService:     authorizationResponseService,
			AuthenticatedResourceOwnerGetter: deps.AuthenticatedResourceOwnerGetter,
		},
		UserInfoService: &UserInfoService{
			AccessTokenResponseRepository: deps.AccessTokenResponseRepository,
			ResourceOwnerRepository:       deps.ResourceOwnerRepository,
		},
	}
}
