package oauth

import (
	http "github.com/llbarbosas/noodle-account/core/infra"
)

type ResourceOwnerID string

type ResourceOwnerRole string
type ResourceOwner interface {
	ID() ResourceOwnerID
	UpdatedAt() int64
	Name() string
	GivenName() string
	LastName() string
	Picture() string
	Email() string
	EmailVerified() bool
	PhoneNumber() string
	PhoneNumberVerified() bool
	Roles() []ResourceOwnerRole
}

type AuthenticatedResourceOwnerGetter func(http.Request) (ResourceOwner, error)
