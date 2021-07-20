package smtp

import (
	"net/smtp"
	"os"

	"github.com/llbarbosas/noodle-account/core/infra"
)

type SimpleMailDriver struct {
	addr string
	auth smtp.Auth
}

// TODO: Add mail to queue
func (d SimpleMailDriver) SendMail(e infra.Email) error {
	from := e.From

	if from == "" {
		from = os.Getenv("SMTP_USER")
	}

	return smtp.SendMail(d.addr, d.auth, from, e.To, e.Message)
}

func NewSimpleMailDriver() SimpleMailDriver {
	host := os.Getenv("SMTP_HOST")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.CRAMMD5Auth(user, pass)

	return SimpleMailDriver{
		addr: host + ":" + port,
		auth: auth,
	}
}
