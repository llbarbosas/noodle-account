package account

import (
	"errors"
	"log"
	"time"

	"github.com/llbarbosas/noodle-account/core/account"
	core "github.com/llbarbosas/noodle-account/core/account"
	infra "github.com/llbarbosas/noodle-account/core/infra"
)

var (
	errAuthenticateWrongCredentials = errors.New("E-mail ou senha incorretos")
	errAuthenticateServerError      = errors.New("Não foi possível realizar o login. Tente novamente mais tarde")
)

type AuthenticateService struct {
	UserRepository               UserRepository
	UserAuthenticationRepository UserAuthenticationRepository
	UserTokenRepository          UserTokenRepository
	GReCAPTCHAVerifier           infra.GReCAPTCHAVerifier
}

func (s *AuthenticateService) Execute(r infra.Request) infra.Response {
	// TODO: account.AuthenticationRequest -> AuthenticatePayload
	var authRequest account.AuthenticationRequest

	if err := r.Body(&authRequest); err != nil {
		return responseFromAuthenticationErr(r, errAuthenticateServerError)
	}

	// TODO: Allow reCAPTCHA validation
	if err := s.GReCAPTCHAVerifier(authRequest.GRecaptchaToken, r.IP()); err != nil {
		log.Println(err)
		// return responseFromAuthenticationErr(r, err, 400)
	}

	var modelUser UserModel

	if err := s.UserRepository.GetByEmail(&modelUser, authRequest.Email); err != nil {
		if err == infra.ErrNotFound {
			return responseFromAuthenticationErr(r, errAuthenticateWrongCredentials)
		}

		return responseFromAuthenticationErr(r, errAuthenticateServerError, 500)
	}

	if !modelUser.ComparePassword(authRequest.Password) {
		return responseFromAuthenticationErr(r, errAuthenticateWrongCredentials, 400)
	}

	if err := s.UserAuthenticationRepository.Create(&authRequest); err != nil {
		return responseFromAuthenticationErr(r, errAuthenticateServerError, 500)
	}

	token, err := core.NewUserToken(modelUser.User.ID, modelUser.User.Roles)

	if err != nil {
		return responseFromAuthenticationErr(r, errAuthenticateServerError, 500)
	}

	if err := s.UserTokenRepository.Create(token); err != nil {
		return responseFromAuthenticationErr(r, errAuthenticateServerError, 500)
	}

	r.Cookie(createUserTokenCookie(*token))

	if returnTo := r.QueryOne("return_to"); returnTo != "" {
		return infra.ResponseSafeRedirect(returnTo)
	}

	return infra.ResponseRedirect(infra.DashboardRoute)
}

func createUserTokenCookie(token core.UserToken) *infra.Cookie {
	utCookie := new(infra.Cookie)
	utCookie.Name = "_ut"
	utCookie.Value = token.Code
	utCookie.Secure = true
	utCookie.Expires = time.Now().Add(core.DefaultUserTokenExpiration)

	return utCookie
}
