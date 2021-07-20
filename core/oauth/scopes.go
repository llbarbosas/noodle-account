package oauth

import (
	"strings"

	"github.com/llbarbosas/noodle-account/core/util"
)

var (
	ScopeDescriptions = map[string]string{
		"openid":  "Confirmar sua identidade",
		"email":   "Ver seu endereço de e-mail",
		"profile": "Ver dados pessoais básicos",
		"phone":   "Ver seu telefone de contato",
	}
	AllowedScopes = util.GetKeys(ScopeDescriptions)
)

func ProcessScopes(scopesStr string) []string {
	scopes := strings.Split(scopesStr, " ")
	actualScopes := make([]string, 0, len(scopes))

	for _, scope := range scopes {
		if ScopeDescriptions[scope] != "" {
			actualScopes = append(actualScopes, scope)
		}
	}

	return actualScopes
}
