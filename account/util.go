package account

import (
	"fmt"
	"log"
	"os"

	"github.com/llbarbosas/noodle-account/core/account"
	infra "github.com/llbarbosas/noodle-account/core/infra"
)

func responseFromARErr(file string) func(r infra.Request, err error, code ...int) infra.Response {
	gReCAPTCHASiteKey := os.Getenv("RECAPTCHA_SITEKEY")

	return func(r infra.Request, err error, code ...int) infra.Response {
		return infra.ResponseFile(file, map[string]interface{}{
			"Error":             err.Error(),
			"ReturnTo":          r.QueryOne("return_to"),
			"GReCAPTCHASiteKey": gReCAPTCHASiteKey,
		}, code...)
	}
}

var (
	responseFromAuthenticationErr = responseFromARErr("signin")
	responseFromRegisterErr       = responseFromARErr("signup")
)

func sendMailAsync(driver infra.MailDriver, email infra.Email) {
	go func(d infra.MailDriver, e infra.Email) {
		if err := d.SendMail(e); err != nil {
			log.Println(err)
		}
	}(driver, email)
}

func responseFromAccountErr(path string) func(err account.AccountError, code ...int) infra.Response {
	return func(err account.AccountError, code ...int) infra.Response {
		url := fmt.Sprintf("%s?action=%s", path, err.Code)
		return infra.ResponseRedirect(url, code...)
	}
}

var (
	responseFromVerifyEmailErr     = responseFromAccountErr(infra.VerifyEmailRoute)
	responseFromRequestPasswordErr = responseFromAccountErr(infra.RequestPasswordResetRoute)
	responseFromPasswordResetErr   = responseFromAccountErr(infra.PasswordResetRoute)
)
