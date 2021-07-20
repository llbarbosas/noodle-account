package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/llbarbosas/noodle-account/account"
	accountcore "github.com/llbarbosas/noodle-account/core/account"
	"github.com/llbarbosas/noodle-account/core/infra"
	oauthcore "github.com/llbarbosas/noodle-account/core/oauth"
	"github.com/llbarbosas/noodle-account/http"
	"github.com/llbarbosas/noodle-account/oauth"
	mockrepository "github.com/llbarbosas/noodle-account/persistence/mock"
	"github.com/llbarbosas/noodle-account/smtp"
)

// https://github.com/joho/godotenv/blob/ddf83eb33bbb136f62617a409142b74b91dbcff3/godotenv.go#L79
func loadEnv() {
	envPath := ".env.development"
	env := os.Getenv("ENV")

	if env == "prod" {
		envPath = ".env"
	}

	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	loadEnv()

	log.SetFlags(log.Lshortfile)

	clientRepository := mockrepository.NewClientRepository()
	clientRepository.Create(&oauthcore.Client{
		ID:                "123",
		Label:             "RevendaTech",
		Type:              "SPA",
		RedirectURI:       "http://revendatech.com/cb",
		LogoutRedirectURI: "http://revendatech.com",
	})
	clientRepository.Create(&oauthcore.Client{
		ID:                "456",
		Label:             "NoodleCloud",
		Type:              "internal:SPA",
		RedirectURI:       "http://cloud.noodle.com/cb",
		LogoutRedirectURI: "http://cloud.noodle.com",
	})

	authRequestRepository := mockrepository.NewAuthorizationRequestRepository()
	authorizationResponseRepository := mockrepository.NewMockAuthorizationResponseRepository()
	accessTokenResponseRepository := mockrepository.NewAccessTokenResponseRepository()

	userRepository := mockrepository.NewUserRepository()

	userRepository.Create(&account.UserModel{
		User: accountcore.User{
			Name:     "User",
			Email:    "user@example.com",
			Password: "1234567",
			Roles:    []accountcore.UserRole{accountcore.BasicRole},
		},
	})

	userRepository.Create(&account.UserModel{
		User: accountcore.User{
			Name:          "Admin User",
			Email:         "admin@noodle.com",
			EmailVerified: true,
			Password:      "1234567",
			Roles:         []accountcore.UserRole{accountcore.BasicRole, accountcore.AdministratorRole},
		},
	})

	userAuthenticationRepository := mockrepository.NewUserAuthenticationRepository()
	userTokenRepository := mockrepository.NewUserTokenRepository()
	emailVerificationRepository := mockrepository.NewEmailVerificationRepository()
	passwordResetRepository := mockrepository.NewPasswordResetRepository()

	authenticatedResourceOwnerGetter := func(r infra.Request) (oauthcore.ResourceOwner, error) {
		userTokenCode := r.Cookies("_ut")

		var userToken accountcore.UserToken

		if err := userTokenRepository.GetByCode(&userToken, userTokenCode); err != nil {
			if err == infra.ErrNotFound {
				r.ClearCookie("_ut")
			}

			return nil, err
		}

		var user account.UserModel

		if err := userRepository.GetByID(&user, string(userToken.Subject)); err != nil {
			return nil, err
		}

		return &user, nil
	}

	mailDriver := smtp.NewSimpleMailDriver()

	config := http.AppConfig{
		OAuthDependencies: oauth.ContainerDependencies{
			ClientRepository:                 &clientRepository,
			AuthorizationRequestRepository:   &authRequestRepository,
			AuthorizationResponseRepository:  &authorizationResponseRepository,
			AccessTokenResponseRepository:    &accessTokenResponseRepository,
			ResourceOwnerRepository:          &userRepository,
			AuthenticatedResourceOwnerGetter: authenticatedResourceOwnerGetter,
		},

		AccountDependencies: account.ContainerDependencies{
			UserRepository:               &userRepository,
			UserAuthenticationRepository: &userAuthenticationRepository,
			UserTokenRepository:          &userTokenRepository,
			EmailVerificationRepository:  &emailVerificationRepository,
			PasswordResetRepository:      &passwordResetRepository,
			MailDriver:                   mailDriver,
			GReCAPTCHAVerifier:           http.GReCAPTCHAVerifier,
		},
	}

	app := http.NewApp(config)

	log.Fatal(app.Listen(":3001"))
}
