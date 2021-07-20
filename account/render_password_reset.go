package account

import (
	"os"
	"time"

	core "github.com/llbarbosas/noodle-account/core/account"
	infra "github.com/llbarbosas/noodle-account/core/infra"
)

type RenderPasswordResetService struct {
	PasswordResetRepository PasswordResetRepository
}

func (s *RenderPasswordResetService) Execute(r infra.Request) infra.Response {
	code := r.QueryOne("code")

	if code == "" {
		return infra.ResponseRedirect(infra.RequestPasswordResetRoute)
	}

	var passwordReset core.PasswordReset

	if err := s.PasswordResetRepository.GetByCode(&passwordReset, code); err != nil {
		if err == infra.ErrNotFound {
			return responseFromRequestPasswordErr(core.ErrPasswordResetInvalid, 400)
		}

		return responseFromRequestPasswordErr(core.ErrPasswordResetServerError, 500)
	}

	if passwordReset.ExpiresIn < time.Now().Unix() {
		s.PasswordResetRepository.DeleteByCode(code)
		return responseFromRequestPasswordErr(core.ErrPasswordResetExpirated, 400)
	}

	bindData := map[string]string{
		"CodeVerifier":      passwordReset.CodeVerifier,
		"GReCAPTCHASiteKey": os.Getenv("RECAPTCHA_SITEKEY"),
	}

	return infra.ResponseFile("password_reset", bindData)
}
