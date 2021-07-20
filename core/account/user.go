package account

import (
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserID string

type UserRole string

const (
	BasicRole         UserRole = "user"
	AdministratorRole          = "admin"
)

var (
	EmailConfirmationExpiration = time.Hour * 2
)

type User struct {
	ID                  UserID     `json:"id"`
	Name                string     `form:"name" json:"name"`
	Email               string     `form:"email" json:"email"`
	EmailVerified       bool       `json:"email_verified"`
	Password            string     `form:"password" json:"password"`
	Picture             string     `form:"picture" json:"picture"`
	PhoneNumber         string     `form:"phone_number" json:"phone_number"`
	PhoneNumberVerified bool       `form:"phone_number_verified" json:"phone_number_verified"`
	Roles               []UserRole `json:"roles"`
}

func NewUser(name, email, password string) User {
	email = strings.TrimSpace(email)
	roles := []UserRole{BasicRole}
	matched, _ := regexp.MatchString(`@noodle\.com$`, email)

	if matched {
		roles = append(roles, AdministratorRole)
	}

	return User{
		Name:                name,
		Email:               email,
		Password:            password,
		EmailVerified:       false,
		Picture:             "",
		PhoneNumber:         "",
		PhoneNumberVerified: false,
		Roles:               roles,
	}
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		return false
	}

	return true
}

func (u *User) FirstName() string {
	return strings.Split(u.Name, " ")[0]
}

type AuthenticationRequest struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	AuthorizationID string `form:"_aid" query:"_aid"`
	GRecaptchaToken string `form:"_grct" query:"_grct"`
}

type EmailVerification struct {
	Code      string
	UserID    UserID
	ExpiresIn int64
}

func NewEmailVerification(userID UserID) EmailVerification {
	return EmailVerification{
		UserID:    userID,
		ExpiresIn: time.Now().Add(EmailConfirmationExpiration).Unix(),
	}
}

type PasswordReset struct {
	Code         string
	CodeVerifier string
	UserID       UserID
	ExpiresIn    int64
}

func NewPasswordReset(userID UserID) PasswordReset {
	return PasswordReset{
		UserID:    userID,
		ExpiresIn: time.Now().Add(EmailConfirmationExpiration).Unix(),
	}
}
