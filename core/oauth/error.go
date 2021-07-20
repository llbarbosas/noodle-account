package oauth

import (
	"fmt"
	"net/url"
)

// ErrorResponse represents errors that can occour on
// authorization flow described in
// [section 4.1.2.1](https://tools.ietf.org/html/rfc6749#section-4.1.2.1)
// extended by [PKFC](https://tools.ietf.org/html/rfc7636#section-4.4.1)
type ErrorResponse struct {
	Code        string
	Description string
	URI         string
	State       string
	CallbackURL string
}

func NewErrorResponse(code, description string) ErrorResponse {
	return ErrorResponse{
		Code:        code,
		Description: description,
	}
}

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", er.Code, er.Description)
}

func (er ErrorResponse) URLQuery() string {
	q := url.Values{}
	q.Set("error", er.Code)
	q.Set("error_description", er.Description)

	if er.URI != "" {
		q.Set("error_uri", er.URI)
	}

	if er.State != "" {
		q.Set("state", er.State)
	}

	return q.Encode()
}

func (er ErrorResponse) Bind(state, callbackURL string) ErrorResponse {
	er.State = state
	er.CallbackURL = callbackURL

	return er
}

func (er ErrorResponse) BindAR(ar AuthorizationRequest) ErrorResponse {
	er.State = ar.State
	er.CallbackURL = ar.RedirectURI

	return er
}

// NewInvalidRequestError occurs when "The request is missing
// a required parameter, includes an invalid parameter value,
// includes a parameter more than once, or is otherwise
// malformed" ([section 4.1.2.1](https://tools.ietf.org/html/rfc6749#section-4.1.2.1))
func NewInvalidRequestError(description string) ErrorResponse {
	return NewErrorResponse("invalid_request", description)
}

// NewUnauthorizedClientError occurs when "The client is not
// authorized to request an authorization code using this method."
// ([section 4.1.2.1](https://tools.ietf.org/html/rfc6749#section-4.1.2.1))
func NewUnauthorizedClientError(description string) ErrorResponse {
	return NewErrorResponse("unauthorized_client", description)
}

// NewAccessDeniedError occurs when "The resource owner
// or authorization server denied the request."
func NewAccessDeniedError(description string) ErrorResponse {
	return NewErrorResponse("access_denied", description)
}

// NewUnsupportedResponseTypeError occurs when "The
// authorization server does not support obtaining
// an authorization code using this method."
func NewUnsupportedResponseTypeError(description string) ErrorResponse {
	return NewErrorResponse("unsupported_response_type", description)
}

// NewInvalidScopeError occurs when "The requested scope is
// invalid, unknown, or malformed."
func NewInvalidScopeError(description string) ErrorResponse {
	return NewErrorResponse("invalid_scope", description)
}

// NewServerError occurs when "The authorization server
// encountered an unexpected condition that prevented it from
// fulfilling the request."
func NewServerError(description string) ErrorResponse {
	return NewErrorResponse("server_error", description)
}

// NewTemporarilyUnavailableError occurs when "The authorization
// server is currently unable to handle the request due to a temporary
// overloading or maintenance of the server."
func NewTemporarilyUnavailableError(description string) ErrorResponse {
	return NewErrorResponse("temporarily_unavailable", description)
}
