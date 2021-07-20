package account

import (
	infra "github.com/llbarbosas/noodle-account/core/infra"
)

type ContainerDependencies struct {
	UserRepository               UserRepository
	UserAuthenticationRepository UserAuthenticationRepository
	UserTokenRepository          UserTokenRepository
	EmailVerificationRepository  EmailVerificationRepository
	PasswordResetRepository      PasswordResetRepository
	GReCAPTCHAVerifier           infra.GReCAPTCHAVerifier
	MailDriver                   infra.MailDriver
}

type Container struct {
	ContainerDependencies
	RegisterService             *RegisterService
	AuthenticateService         *AuthenticateService
	VerifyEmailService          *VerifyEmailService
	PasswordResetService        *PasswordResetService
	RenderPasswordResetService  *RenderPasswordResetService
	RequestPasswordResetService *RequestPasswordResetService
}

func CreateContainer(deps ContainerDependencies) *Container {
	return &Container{
		ContainerDependencies: deps,
		RegisterService: &RegisterService{
			UserRepository:              deps.UserRepository,
			UserTokenRepository:         deps.UserTokenRepository,
			EmailVerificationRepository: deps.EmailVerificationRepository,
			GReCAPTCHAVerifier:          deps.GReCAPTCHAVerifier,
			MailDriver:                  deps.MailDriver,
		},
		AuthenticateService: &AuthenticateService{
			UserRepository:               deps.UserRepository,
			UserAuthenticationRepository: deps.UserAuthenticationRepository,
			UserTokenRepository:          deps.UserTokenRepository,
			GReCAPTCHAVerifier:           deps.GReCAPTCHAVerifier,
		},
		VerifyEmailService: &VerifyEmailService{
			UserRepository:              deps.UserRepository,
			EmailVerificationRepository: deps.EmailVerificationRepository,
			MailDriver:                  deps.MailDriver,
		},
		PasswordResetService: &PasswordResetService{
			UserRepository:          deps.UserRepository,
			PasswordResetRepository: deps.PasswordResetRepository,
			MailDriver:              deps.MailDriver,
			GReCAPTCHAVerifier:      deps.GReCAPTCHAVerifier,
		},
		RenderPasswordResetService: &RenderPasswordResetService{
			PasswordResetRepository: deps.PasswordResetRepository,
		},
		RequestPasswordResetService: &RequestPasswordResetService{
			UserRepository:          deps.UserRepository,
			PasswordResetRepository: deps.PasswordResetRepository,
			MailDriver:              deps.MailDriver,
			GReCAPTCHAVerifier:      deps.GReCAPTCHAVerifier,
		},
	}
}
