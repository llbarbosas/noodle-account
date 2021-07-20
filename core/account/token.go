package account

import (
	"time"

	"github.com/llbarbosas/noodle-account/core/util"
)

var (
	DefaultUserTokenExpiration = time.Hour * 24
)

type UserToken struct {
	Code      string
	Subject   UserID
	ExpiresIn int64
	IssuedAt  int64
	Roles     []UserRole
}

func NewUserToken(userID UserID, roles []UserRole) (*UserToken, error) {
	code, err := util.RandomString(100)

	if err != nil {
		return nil, err
	}

	token := UserToken{
		Code:      code,
		Subject:   userID,
		ExpiresIn: DefaultUserTokenExpiration.Milliseconds(),
		IssuedAt:  time.Now().Unix(),
		Roles:     roles,
	}

	return &token, nil
}
