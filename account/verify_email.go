package account

import (
	"errors"
	"time"

	core "github.com/llbarbosas/noodle-account/core/account"
	infra "github.com/llbarbosas/noodle-account/core/infra"
)

type VerifyEmailService struct {
	UserRepository              UserRepository
	EmailVerificationRepository EmailVerificationRepository
	MailDriver                  infra.MailDriver
}

func (s *VerifyEmailService) Execute(r infra.Request) infra.Response {
	code := r.QueryOne("code")
	action := r.QueryOne("action")

	if action != "" {
		return infra.ResponseFromErr(errors.New("Not implemented"), 500)
	}

	var confirmation core.EmailVerification

	if err := s.EmailVerificationRepository.GetByCode(&confirmation, code); err != nil {
		if err == infra.ErrNotFound {
			return responseFromVerifyEmailErr(core.ErrVerifyEmailInvalid, 400)
		}

		return responseFromVerifyEmailErr(core.ErrVerifyEmailServerError, 400)
	}

	if confirmation.ExpiresIn < time.Now().Unix() {
		var user UserModel

		if err := s.UserRepository.GetByID(&user, string(confirmation.UserID)); err != nil {
			return responseFromVerifyEmailErr(core.ErrVerifyEmailServerError, 400)
		}

		if err := sendConfirmationMail(s.EmailVerificationRepository, s.MailDriver, user.User); err != nil {
			return responseFromVerifyEmailErr(core.ErrVerifyEmailServerError, 400)
		}

		return responseFromVerifyEmailErr(core.ErrVerifyEmailExpirated, 400)
	}

	if err := s.UserRepository.VerifyEmail(confirmation.UserID); err != nil {
		return responseFromVerifyEmailErr(core.ErrVerifyEmailServerError, 400)
	}

	return infra.ResponseRedirect(infra.DashboardRoute)
}
