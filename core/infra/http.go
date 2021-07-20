package infra

import (
	"context"
	"time"
)

const (
	AuthenticateRoute         string = "/signin"
	RequestPasswordResetRoute        = "/signin/password_reset/request"
	PasswordResetRoute               = "/signin/password_reset"
	RegisterRoute                    = "/signup"
	VerifyEmailRoute                 = "/signup/verify"
	DashboardRoute                   = "/"

	AuthorizationRoute = "/oauth/authorize"
	TokenRoute         = "/oauth/token"
	UserInfoRoute      = "/oauth/userinfo"

	MainRealm = "main@account.noodle.com"
)

type Handler func(Request) Response

type Request interface {
	Host() string
	Body(out interface{}) error
	Query(out interface{}) error
	Header(key string) string
	Params(key string) string
	Cookie(cookie *Cookie)
	Cookies(key string) string
	ClearCookie(key ...string)
	QueryOne(key string) string
	RawQuery() string
	Context() context.Context
	WithContext(ctxFns ...PartialContextFn) Request
	IP() string
}

type PartialContextFn func(context.Context) context.Context

type Response struct {
	StatusCode int
	Type       string
	Headers    map[string]string

	// Type == "json"
	Body interface{}

	// Type == "redirect"
	RedirectURL  string
	RedirectSafe bool

	// Type == "file"
	File         string
	FileBindData interface{}
}

func IsOK(r Response) bool {
	return r.StatusCode >= 200 && r.StatusCode <= 299
}

func getCode(code []int, defaultValue int) int {
	if len(code) < 1 {
		return defaultValue
	}

	return code[0]
}

func ResponseFromWWWAuthenticate(value string, code ...int) Response {
	return Response{
		StatusCode: getCode(code, 401),
		Headers: map[string]string{
			"WWW-Authenticate": value,
		},
	}
}

func ResponseFromErr(err error, code ...int) Response {
	return Response{
		StatusCode: getCode(code, 400),
		Body: map[string]interface{}{
			"error": err.Error(),
		},
	}
}

func ResponseRedirect(url string, code ...int) Response {
	return Response{
		Type:        "redirect",
		StatusCode:  getCode(code, 302),
		RedirectURL: url,
	}
}

func ResponseSafeRedirect(url string, code ...int) Response {
	return Response{
		Type:         "redirect",
		StatusCode:   getCode(code, 302),
		RedirectURL:  url,
		RedirectSafe: true,
	}
}

func ResponseJSON(body interface{}, code ...int) Response {
	return Response{
		Type:       "json",
		StatusCode: getCode(code, 200),
		Body:       body,
	}
}

func ResponseFile(filePath string, bind interface{}, code ...int) Response {
	return Response{
		Type:         "file",
		StatusCode:   getCode(code, 200),
		File:         filePath,
		FileBindData: bind,
	}
}

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Expires  time.Time
	Secure   bool
	HTTPOnly bool
	SameSite string
}

type GReCAPTCHAVerifier func(token string, userip ...string) error
