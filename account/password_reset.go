package account

import (
	"log"
	"time"

	core "github.com/llbarbosas/noodle-account/core/account"
	infra "github.com/llbarbosas/noodle-account/core/infra"
	"github.com/llbarbosas/noodle-account/core/util"
	"github.com/llbarbosas/noodle-account/smtp/templates"
)

type PasswordResetService struct {
	UserRepository          UserRepository
	PasswordResetRepository PasswordResetRepository
	MailDriver              infra.MailDriver
	GReCAPTCHAVerifier      infra.GReCAPTCHAVerifier
}

type passwordResetPayload struct {
	CodeVerifier    string `json:"verifier" form:"verifier"`
	Password        string `json:"password" form:"password"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm"`
	GRecaptchaToken string `json:"_grct" form:"_grct"`
}

func (s *PasswordResetService) Execute(r infra.Request) infra.Response {
	var payload passwordResetPayload

	if err := r.Body(&payload); err != nil {
		return responseFromPasswordResetErr(core.ErrPasswordResetInvalid)
	}

	// TODO: Allow reCAPTCHA validation
	if err := s.GReCAPTCHAVerifier(payload.GRecaptchaToken, r.IP()); err != nil {
		log.Println(err)
		// return responseFromAuthenticationErr(r, err, 400)
	}

	code := util.Sha256Base64(payload.CodeVerifier)

	var passwordReset core.PasswordReset

	if err := s.PasswordResetRepository.GetByCode(&passwordReset, code); err != nil {
		return responseFromPasswordResetErr(core.ErrPasswordResetServerError, 400)
	}

	if passwordReset.ExpiresIn < time.Now().Unix() {
		s.PasswordResetRepository.DeleteByCode(code)
		return responseFromPasswordResetErr(core.ErrPasswordResetExpirated, 400)
	}

	if err := s.UserRepository.UpdatePassword(payload.Password, passwordReset.UserID); err != nil {
		return responseFromPasswordResetErr(core.ErrPasswordResetServerError, 400)
	}

	s.PasswordResetRepository.DeleteByCode(code)

	var user UserModel

	if err := s.UserRepository.GetByID(&user, string(passwordReset.UserID)); err != nil {
		return responseFromPasswordResetErr(core.ErrPasswordResetServerError, 400)
	}

	successfulReset, _ := templates.CreateSuccessPasswordResetMail(user.User.Email)

	sendMailAsync(s.MailDriver, *successfulReset)

	return infra.ResponseRedirect(infra.AuthenticateRoute + "?action=password_reset")
}
