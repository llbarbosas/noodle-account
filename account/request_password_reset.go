package account

import (
	"log"

	core "github.com/llbarbosas/noodle-account/core/account"
	infra "github.com/llbarbosas/noodle-account/core/infra"
	"github.com/llbarbosas/noodle-account/smtp/templates"
)

type RequestPasswordResetService struct {
	UserRepository          UserRepository
	PasswordResetRepository PasswordResetRepository
	MailDriver              infra.MailDriver
	GReCAPTCHAVerifier      infra.GReCAPTCHAVerifier
}

type requestPasswordResetPayload struct {
	Email           string `json:"email" form:"email"`
	GRecaptchaToken string `json:"_grct" form:"_grct"`
}

func (s *RequestPasswordResetService) Execute(r infra.Request) infra.Response {
	var payload requestPasswordResetPayload

	if err := r.Body(&payload); err != nil {
		return responseFromRequestPasswordErr(core.ErrPasswordResetInvalid, 400)
	}

	// TODO: Allow reCAPTCHA validation
	if err := s.GReCAPTCHAVerifier(payload.GRecaptchaToken, r.IP()); err != nil {
		log.Println(err)
		// return infra.ResponseRedirect(infra.RequestPasswordResetRoute+"?action=password_reset_invalid", 400)
	}

	var user UserModel

	if err := s.UserRepository.GetByEmail(&user, payload.Email); err != nil {
		return responseFromRequestPasswordErr(core.ErrPasswordResetUnknownAccount, 400)
	}

	passwordReset := core.NewPasswordReset(user.User.ID)

	if err := s.PasswordResetRepository.Create(&passwordReset); err != nil {
		return responseFromRequestPasswordErr(core.ErrPasswordResetServerError, 500)
	}

	passwordResetMail, _ := templates.CreatePasswordResetMail(user.User.Email, passwordReset.Code)

	sendMailAsync(s.MailDriver, *passwordResetMail)

	return infra.ResponseRedirect(infra.RequestPasswordResetRoute + "?action=password_reset_sent")
}
