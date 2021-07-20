package oauth

import (
	"net/url"

	http "github.com/llbarbosas/noodle-account/core/infra"
	core "github.com/llbarbosas/noodle-account/core/oauth"
)

func responseFromErrorResponse(err core.ErrorResponse) http.Response {
	if err.CallbackURL != "" {
		return http.ResponseRedirect(err.CallbackURL + "?" + err.URLQuery())
	}

	return http.ResponseFile("oauth_error", map[string]interface{}{
		"Error":            err.Code,
		"ErrorDescription": err.Description,
	})
}

func redirectSignin(rawQuery string) http.Response {
	q := url.Values{}
	q.Add("return_to", "/oauth/authorize?"+rawQuery)
	return http.ResponseRedirect("/signin?" + q.Encode())
}
