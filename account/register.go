package account

import (
	"errors"
	"log"

	core "github.com/llbarbosas/noodle-account/core/account"
	infra "github.com/llbarbosas/noodle-account/core/infra"
	"github.com/llbarbosas/noodle-account/smtp/templates"
)

var (
	errRegisterDuplicatedEmail = errors.New("Este e-mail já foi cadastrado")
	errRegisterServerError     = errors.New("Não foi possível realizar o cadastro. Tente novamente mais tarde")
)

type RegisterService struct {
	UserRepository               UserRepository
	EmailVerificationRepository  EmailVerificationRepository
	UserAuthenticationRepository UserAuthenticationRepository
	UserTokenRepository          UserTokenRepository
	GReCAPTCHAVerifier           infra.GReCAPTCHAVerifier
	MailDriver                   infra.MailDriver
}

type RegisterPayload struct {
	Name            string `form:"name" json:"name"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	GRecaptchaToken string `form:"_grct" json:"_grct"`
}

func (s *RegisterService) Execute(r infra.Request) infra.Response {
	var userData RegisterPayload

	if err := r.Body(&userData); err != nil {
		return infra.ResponseFromErr(err)
	}

	// TODO: Allow reCAPTCHA validation
	if err := s.GReCAPTCHAVerifier(userData.GRecaptchaToken, r.IP()); err != nil {
		log.Println(err)
		// return responseFromAuthenticationErr(r, err, 400)
	}

	err := s.UserRepository.GetByEmail(new(UserModel), userData.Email)

	isErrNotFound := err == infra.ErrNotFound

	if !isErrNotFound {
		if err != nil {
			return responseFromRegisterErr(r, errRegisterServerError, 500)
		}

		return responseFromRegisterErr(r, errRegisterDuplicatedEmail, 400)
	}

	// TODO: Validate user creation
	user := core.NewUser(userData.Name, userData.Email, userData.Password)
	userModel := UserModel{User: user}

	if err := s.UserRepository.Create(&userModel); err != nil {
		return responseFromRegisterErr(r, errRegisterServerError, 500)
	}

	if err := sendConfirmationMail(s.EmailVerificationRepository, s.MailDriver, userModel.User); err != nil {
		return responseFromRegisterErr(r, errRegisterServerError, 500)
	}

	token, err := core.NewUserToken(userModel.User.ID, user.Roles)

	if err != nil {
		return responseFromAuthenticationErr(r, errAuthenticateServerError, 500)
	}

	if err := s.UserTokenRepository.Create(token); err != nil {
		return responseFromAuthenticationErr(r, errAuthenticateServerError, 500)
	}

	r.Cookie(createUserTokenCookie(*token))

	return infra.ResponseRedirect(infra.DashboardRoute)
}

func sendConfirmationMail(repository EmailVerificationRepository, mailDriver infra.MailDriver, user core.User) error {
	confirmation := core.NewEmailVerification(user.ID)

	if err := repository.Create(&confirmation); err != nil {
		return err
	}

	confirmationMail, err := templates.CreateConfirmationMail(user.Email, confirmation.Code)

	if err != nil {
		return err
	}

	sendMailAsync(mailDriver, *confirmationMail)

	return nil
}
