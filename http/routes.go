package http

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/llbarbosas/noodle-account/account"
	accountcore "github.com/llbarbosas/noodle-account/core/account"
	"github.com/llbarbosas/noodle-account/core/infra"
)

func (a *App) registerRoutes() {
	a.fiberApp.Get(infra.AuthorizationRoute, FiberHandlerAdapter(a.OAuthContainer.AuthorizationService.Execute))
	a.fiberApp.Post(infra.AuthorizationRoute, FiberHandlerAdapter(a.OAuthContainer.AuthorizationResponseService.Execute))
	a.fiberApp.Post(infra.TokenRoute, FiberHandlerAdapter(a.OAuthContainer.TokenService.Execute))
	a.fiberApp.Get(infra.UserInfoRoute, FiberHandlerAdapter(a.OAuthContainer.UserInfoService.Execute))

	a.fiberApp.Post(infra.AuthenticateRoute, FiberHandlerAdapter(a.AccountContainer.AuthenticateService.Execute))
	a.fiberApp.Get(infra.AuthenticateRoute, handleAccountView("signin", a.AccountContainer.UserTokenRepository))
	a.fiberApp.Get(infra.RequestPasswordResetRoute, handleAccountView("request_password_reset", a.AccountContainer.UserTokenRepository))
	a.fiberApp.Post(infra.RequestPasswordResetRoute, FiberHandlerAdapter(a.AccountContainer.RequestPasswordResetService.Execute))
	a.fiberApp.Get(infra.PasswordResetRoute, FiberHandlerAdapter(a.AccountContainer.RenderPasswordResetService.Execute))
	a.fiberApp.Post(infra.PasswordResetRoute, FiberHandlerAdapter(a.AccountContainer.PasswordResetService.Execute))
	a.fiberApp.Post(infra.RegisterRoute, FiberHandlerAdapter(a.AccountContainer.RegisterService.Execute))
	a.fiberApp.Get(infra.RegisterRoute, handleAccountView("signup", a.AccountContainer.UserTokenRepository))
	a.fiberApp.Get(infra.VerifyEmailRoute, FiberHandlerAdapter(a.AccountContainer.VerifyEmailService.Execute))

	a.fiberApp.Get(infra.DashboardRoute, handleDashboard(a.AccountContainer.UserRepository, a.AccountContainer.UserTokenRepository))

	a.fiberApp.Get("/docs", FiberRenderFile("docs"))
}

func handleAccountView(view string, tokenRepo account.UserTokenRepository) func(c *fiber.Ctx) error {
	gReCAPTCHASiteKey := os.Getenv("RECAPTCHA_SITEKEY")

	return func(c *fiber.Ctx) error {
		var (
			successMessage string
			errorMessage   string
		)

		action := c.Query("action", "")
		returnTo := c.Query("return_to", "")

		switch {
		case action == "password_reset" && view == "signin":
			successMessage = "Sua senha foi alterada com sucesso"

		case action == "password_reset_sent" && view == "request_password_reset":
			successMessage = "E-mail de restauração enviado"

		default:
			if accountcore.ErrorMatchRoute(action, c.Path()) {
				errorMessage = accountcore.GetErrorDescription(action)
			}
		}

		if userTokenCode := c.Cookies("_ut", ""); userTokenCode != "" {
			var userToken accountcore.UserToken

			if err := tokenRepo.GetByCode(&userToken, userTokenCode); err == nil {
				if returnTo != "" {
					return c.Redirect(SafeRedirect(returnTo, infra.DashboardRoute, c.Hostname()))
				}

				return c.Redirect(infra.DashboardRoute)
			}

			c.ClearCookie("_ut")
		}

		return c.Render(view, map[string]string{
			"ReturnTo":          returnTo,
			"GReCAPTCHASiteKey": gReCAPTCHASiteKey,
			"Success":           successMessage,
			"Error":             errorMessage,
		})
	}
}

func handleDashboard(userRepo account.UserRepository, tokenRepo account.UserTokenRepository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userTokenCode := c.Cookies("_ut", "")

		if userTokenCode == "" {
			return c.Redirect(infra.AuthenticateRoute)
		}

		var userToken accountcore.UserToken

		if err := tokenRepo.GetByCode(&userToken, userTokenCode); err != nil {
			c.ClearCookie("_ut")
			return c.Redirect(infra.AuthenticateRoute)
		}

		var user account.UserModel

		if err := userRepo.GetByID(&user, string(userToken.Subject)); err != nil {
			c.ClearCookie("_ut")
			return c.Redirect(infra.AuthenticateRoute)
		}

		return c.Render("dashboard", user)
	}
}
