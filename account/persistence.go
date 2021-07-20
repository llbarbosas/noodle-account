package account

import (
	core "github.com/llbarbosas/noodle-account/core/account"
	"github.com/llbarbosas/noodle-account/core/oauth"
)

type UserModel struct {
	core.User
	UpdatedAtModel int64
}

func (u UserModel) ID() oauth.ResourceOwnerID {
	return oauth.ResourceOwnerID(u.User.ID)
}

func (u UserModel) UpdatedAt() int64 {
	return u.UpdatedAtModel
}

func (u UserModel) Name() string {
	return u.User.Name
}

func (u UserModel) GivenName() string {
	return u.User.FirstName()
}

func (u UserModel) LastName() string {
	return u.User.Name
}

func (u UserModel) Picture() string {
	return u.User.Picture
}

func (u UserModel) Email() string {
	return u.User.Email
}

func (u UserModel) EmailVerified() bool {
	return u.User.EmailVerified
}

func (u UserModel) PhoneNumber() string {
	return u.User.PhoneNumber
}

func (u UserModel) PhoneNumberVerified() bool {
	return u.User.PhoneNumberVerified
}

func (u UserModel) Roles() []oauth.ResourceOwnerRole {
	roles := make([]oauth.ResourceOwnerRole, len(u.User.Roles))

	for i, role := range u.User.Roles {
		roles[i] = oauth.ResourceOwnerRole(role)
	}

	return roles
}

type UserRepository interface {
	Create(*UserModel) error
	GetByEmail(*UserModel, string) error
	// GetByID(*UserModel, core.UserID) error
	GetByID(interface{}, string) error
	GetAll(interface{}) error
	VerifyEmail(core.UserID) error
	UpdatePassword(string, core.UserID) error
}

type UserAuthenticationRepository interface {
	Create(*core.AuthenticationRequest) error
	GetByAuthID(*core.AuthenticationRequest, string) error
}

type EmailVerificationRepository interface {
	Create(*core.EmailVerification) error
	GetByCode(*core.EmailVerification, string) error
}

type PasswordResetRepository interface {
	Create(*core.PasswordReset) error
	DeleteByCode(string) error
	GetByCode(*core.PasswordReset, string) error
}

type UserTokenRepository interface {
	Create(*core.UserToken) error
	GetByCode(*core.UserToken, string) error
}
