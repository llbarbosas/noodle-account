package templates

import (
	"bytes"
	"html/template"

	"github.com/llbarbosas/noodle-account/core/infra"
)

func createEmail(templPath string, templData map[string]string, to []string) (*infra.Email, error) {
	templ, err := template.ParseFiles(templPath)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := templ.Execute(&buf, templData); err != nil {
		return nil, err
	}

	mail := infra.Email{
		To:      to,
		Message: buf.Bytes(),
	}

	return &mail, nil
}

func CreateConfirmationMail(toMail string, confirmationCode string) (*infra.Email, error) {
	templateData := map[string]string{
		"To":               toMail,
		"ConfirmationCode": confirmationCode,
	}

	return createEmail("./smtp/templates/signup_confirmation.html", templateData, []string{toMail})
}

func CreatePasswordResetMail(toMail string, code string) (*infra.Email, error) {
	templateData := map[string]string{
		"To":   toMail,
		"Code": code,
	}

	return createEmail("./smtp/templates/password_reset.html", templateData, []string{toMail})
}

func CreateSuccessPasswordResetMail(toMail string) (*infra.Email, error) {
	templateData := map[string]string{
		"To": toMail,
	}

	return createEmail("./smtp/templates/successful_password_reset.html", templateData, []string{toMail})
}
