package middleware

import (
	"net/http"

	"github.com/rhinosc/web-market/code/internal/auth"
	"github.com/rhinosc/web-market/code/platform/web/response"
)

type Authenticator struct {
	// au is the authenticator service.
	au auth.AuthToken
}

func NewAuthenticator(au auth.AuthToken) *Authenticator {
	return &Authenticator{
		au: au,
	}
}

func (a *Authenticator) Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		// get token
		token := r.Header.Get("Authorization")

		// validate token
		if err := a.au.Auth(token); err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// call
		handler.ServeHTTP(w, r)

		// after
		// ...
	})
}
