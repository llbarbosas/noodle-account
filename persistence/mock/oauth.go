package mock

import (
	"time"

	"github.com/llbarbosas/noodle-account/core/infra"
	core "github.com/llbarbosas/noodle-account/core/oauth"
	"github.com/llbarbosas/noodle-account/core/util"
	"github.com/llbarbosas/noodle-account/oauth"
)

type ClientRepository struct {
	Repository
}

func (r *ClientRepository) Create(entity *core.Client) error {
	// entity.ID = core.ClientID(util.NewUUIDStr())

	return r.Repository.Create(entity)
}

func (r *ClientRepository) GetByID(out *core.Client, id core.ClientID) error {
	return r.GetOne(out, infra.QueryByStringField("ID", string(id)))
}

func NewClientRepository() ClientRepository {
	return ClientRepository{
		Repository: NewRepository(core.Client{}),
	}
}

type AuthorizationRequestRepository struct {
	Repository
}

func (r *AuthorizationRequestRepository) Create(entity *oauth.AuthorizationRequestModel) error {
	entity.ID = oauth.AuthorizationRequestID(util.NewUUIDStr())
	entity.CreatedAt = time.Now()

	return r.Repository.Create(entity)
}

func (r *AuthorizationRequestRepository) GetByID(out *oauth.AuthorizationRequestModel, id oauth.AuthorizationRequestID) error {
	return r.GetOne(out, infra.QueryByStringField("ID", string(id)))
}

func NewAuthorizationRequestRepository() AuthorizationRequestRepository {
	return AuthorizationRequestRepository{
		Repository: NewRepository(oauth.AuthorizationRequestModel{}),
	}
}

type AuthorizationResponseRepository struct {
	Repository
}

func (r *AuthorizationResponseRepository) GetByCode(out *oauth.AuthorizationResponseModel, code string) error {
	return r.GetOne(out, infra.QueryByStringField("Code", code))
}

func NewMockAuthorizationResponseRepository() AuthorizationResponseRepository {
	return AuthorizationResponseRepository{
		Repository: NewRepository(oauth.AuthorizationResponseModel{}),
	}
}

type AccessTokenResponseRepository struct {
	Repository
}

func (r *AccessTokenResponseRepository) GetByAccessToken(out *core.AccessTokenResponse, code string) error {
	return r.GetOne(out, infra.QueryByStringField("AccessTokenCode", code))
}

func (r *AccessTokenResponseRepository) GetByRefreshToken(out *core.AccessTokenResponse, code string) error {
	return r.GetOne(out, infra.QueryByStringField("RefreshTokenCode", code))
}

func (r *AccessTokenResponseRepository) DeleteByRefreshToken(code string) error {
	return r.Delete(infra.QueryByStringField("RefreshTokenCode", code))
}

func NewAccessTokenResponseRepository() AccessTokenResponseRepository {
	return AccessTokenResponseRepository{
		Repository: NewRepository(core.AccessTokenResponse{}),
	}
}
