package mock

import (
	"reflect"
	"time"

	"github.com/llbarbosas/noodle-account/account"
	core "github.com/llbarbosas/noodle-account/core/account"
	"github.com/llbarbosas/noodle-account/core/infra"
	"github.com/llbarbosas/noodle-account/core/util"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticationRepository struct {
	Repository
}

func (r *UserAuthenticationRepository) Create(userAuthentication *core.AuthenticationRequest) error {
	userAuthentication.Password = ""

	return r.Repository.Create(userAuthentication)
}

func (r *UserAuthenticationRepository) GetByAuthID(out *core.AuthenticationRequest, authorizationID string) error {
	return r.GetOne(out, infra.QueryByStringField("AuthorizationID", authorizationID))
}

func NewUserAuthenticationRepository() UserAuthenticationRepository {
	return UserAuthenticationRepository{
		Repository: NewRepository(core.AuthenticationRequest{}),
	}
}

type UserRepository struct {
	Repository
}

func (r *UserRepository) Create(user *account.UserModel) error {
	password := []byte(user.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.User.ID = core.UserID(util.NewUUIDStr())
	user.UpdatedAtModel = time.Now().Unix()

	return r.Repository.Create(user)
}

func (r *UserRepository) GetByEmail(out *account.UserModel, email string) error {
	return r.GetOne(out, infra.QueryByStringField("Email", email))
}

// func (r *UserRepository) GetByID(out *account.UserModel, id core.UserID) error {
func (r *UserRepository) GetByID(out interface{}, id string) error {
	return r.GetOne(out, infra.QueryByStringField("ID", id))
}

func (r *UserRepository) VerifyEmail(id core.UserID) error {
	updateFunc := func(v reflect.Value) reflect.Value {
		v.FieldByName("EmailVerified").Set(reflect.ValueOf(true))
		v.FieldByName("UpdatedAtModel").Set(reflect.ValueOf(time.Now().Unix()))
		return v
	}

	return r.UpdateOne(infra.QueryByStringField("ID", string(id)), updateFunc)
}

func (r *UserRepository) UpdatePassword(newPassword string, id core.UserID) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	password := string(hashedPassword)

	updateFunc := func(v reflect.Value) reflect.Value {
		v.FieldByName("Password").Set(reflect.ValueOf(password))
		v.FieldByName("UpdatedAtModel").Set(reflect.ValueOf(time.Now().Unix()))
		return v
	}

	return r.UpdateOne(infra.QueryByStringField("ID", string(id)), updateFunc)
}

func NewUserRepository() UserRepository {
	return UserRepository{
		Repository: NewRepository(account.UserModel{}),
	}
}

type EmailVerificationRepository struct {
	Repository
}

func (r *EmailVerificationRepository) Create(confirmation *core.EmailVerification) error {
	code, err := util.RandomString(100)

	if err != nil {
		return err
	}

	confirmation.Code = code

	return r.Repository.Create(confirmation)
}

func (r *EmailVerificationRepository) GetByCode(out *core.EmailVerification, code string) error {
	return r.GetOne(out, infra.QueryByStringField("Code", code))
}

func NewEmailVerificationRepository() EmailVerificationRepository {
	return EmailVerificationRepository{
		Repository: NewRepository(core.EmailVerification{}),
	}
}

type PasswordResetRepository struct {
	Repository
}

func (r *PasswordResetRepository) Create(passwordReset *core.PasswordReset) error {
	codeVerifier, err := util.RandomString(100)

	if err != nil {
		return err
	}

	passwordReset.CodeVerifier = codeVerifier
	passwordReset.Code = util.Sha256Base64(codeVerifier)

	return r.Repository.Create(passwordReset)
}

func (r *PasswordResetRepository) GetByCode(out *core.PasswordReset, code string) error {
	return r.GetOne(out, infra.QueryByStringField("Code", code))
}

func (r *PasswordResetRepository) DeleteByCode(code string) error {
	return r.DeleteOne(infra.QueryByStringField("Code", code))
}

func NewPasswordResetRepository() PasswordResetRepository {
	return PasswordResetRepository{
		Repository: NewRepository(core.PasswordReset{}),
	}
}

type UserTokenRepository struct {
	Repository
}

func (r *UserTokenRepository) Create(token *core.UserToken) error {
	return r.Repository.Create(token)
}

func (r *UserTokenRepository) GetByCode(out *core.UserToken, code string) error {
	return r.GetOne(out, infra.QueryByStringField("Code", code))
}

func NewUserTokenRepository() UserTokenRepository {
	return UserTokenRepository{
		Repository: NewRepository(core.UserToken{}),
	}
}
