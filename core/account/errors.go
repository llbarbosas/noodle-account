package account

import "github.com/llbarbosas/noodle-account/core/infra"

var (
	accountErrors             = make(map[string]AccountError)
	ErrVerifyEmailInvalid     = NewAccountError("verify_email_invalid", "Link de confirmação de e-mail inválido", infra.VerifyEmailRoute)
	ErrVerifyEmailExpirated   = NewAccountError("verify_email_expirated", "O link de confirmação de e-mail expirou. Um novo será enviado", infra.VerifyEmailRoute)
	ErrVerifyEmailServerError = NewAccountError("verify_email_error", "Não foi possível confirmar seu e-mail. Tente novamente", infra.VerifyEmailRoute)

	ErrPasswordResetExpirated      = NewAccountError("password_reset_expirated", "A solicitação de recuperação de senha expirou", infra.PasswordResetRoute, infra.RequestPasswordResetRoute)
	ErrPasswordResetInvalid        = NewAccountError("password_reset_invalid", "Solicitação de recuperação de senha inválida", infra.PasswordResetRoute, infra.RequestPasswordResetRoute)
	ErrPasswordResetUnknownAccount = NewAccountError("password_reset_unknown", "Essa conta não existe", infra.PasswordResetRoute, infra.RequestPasswordResetRoute)
	ErrPasswordResetServerError    = NewAccountError("password_reset_error", "Não foi possível gerar a sua solicitação. Tente novamente", infra.PasswordResetRoute, infra.RequestPasswordResetRoute)
)

type AccountError struct {
	Code        string
	Description string
	Routes      []string
}

func (e AccountError) Error() string {
	return e.Description
}

func NewAccountError(code, description string, routes ...string) AccountError {
	err := AccountError{
		Code:        code,
		Description: description,
		Routes:      routes,
	}

	accountErrors[code] = err

	return err
}

func GetErrorDescription(code string) string {
	return accountErrors[code].Description
}

func ErrorMatchRoute(code, route string) bool {
	for _, v := range accountErrors[code].Routes {
		if v == route {
			return true
		}
	}

	return false
}
