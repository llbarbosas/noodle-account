package http

import (
	"github.com/llbarbosas/noodle-account/account"
	"github.com/llbarbosas/noodle-account/oauth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type App struct {
	fiberApp *fiber.App

	OAuthContainer   *oauth.Container
	AccountContainer *account.Container
}
type AppConfig struct {
	OAuthDependencies   oauth.ContainerDependencies
	AccountDependencies account.ContainerDependencies
}

func (a *App) Listen(addr string) error {
	return a.fiberApp.Listen(addr)
}

func (a *App) ListenHTTPS(addr string) error {
	ln, err := NewTLSListener("./http/certs")

	if err != nil {
		return err
	}

	return a.fiberApp.Listener(ln)
}

// TODO: Validate config
func NewApp(config AppConfig) *App {
	engine := html.New("./http/views", ".html")

	fiberApp := fiber.New(fiber.Config{
		Views: engine,
	})

	fiberApp.Static("/", "./http/public")

	oauthContainer := oauth.CreateContainer(oauth.ContainerDependencies{
		ClientRepository:                 config.OAuthDependencies.ClientRepository,
		AuthorizationRequestRepository:   config.OAuthDependencies.AuthorizationRequestRepository,
		AuthorizationResponseRepository:  config.OAuthDependencies.AuthorizationResponseRepository,
		AccessTokenResponseRepository:    config.OAuthDependencies.AccessTokenResponseRepository,
		ResourceOwnerRepository:          config.OAuthDependencies.ResourceOwnerRepository,
		AuthenticatedResourceOwnerGetter: config.OAuthDependencies.AuthenticatedResourceOwnerGetter,
	})

	accountContainer := account.CreateContainer(account.ContainerDependencies{
		UserRepository:               config.AccountDependencies.UserRepository,
		UserAuthenticationRepository: config.AccountDependencies.UserAuthenticationRepository,
		UserTokenRepository:          config.AccountDependencies.UserTokenRepository,
		EmailVerificationRepository:  config.AccountDependencies.EmailVerificationRepository,
		GReCAPTCHAVerifier:           config.AccountDependencies.GReCAPTCHAVerifier,
		MailDriver:                   config.AccountDependencies.MailDriver,
		PasswordResetRepository:      config.AccountDependencies.PasswordResetRepository,
	})

	a := &App{
		fiberApp: fiberApp,

		AccountContainer: accountContainer,
		OAuthContainer:   oauthContainer,
	}

	a.registerRoutes()

	return a
}
