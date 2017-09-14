package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type (
	authPair struct {
		Credential string
		User       string
	}
	authPairs []authPair
)

func (a authPairs) SearchCredential(credential string) (string, bool) {
	for _, pair := range a {
		if pair.Credential == credential {
			return pair.User, true
		}
	}
	return "", false
}

func buildAuthPairs(auths map[string]string) authPairs {
	pairs := make(authPairs, 0, len(auths))
	for user, pass := range auths {
		base := fmt.Sprintf("%s:%s", user, pass)
		credential := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(base)))
		pairs = append(pairs, authPair{
			User:       user,
			Credential: credential,
		})
	}
	return pairs
}

type basicAuthHandler struct {
	http.Handler
	realm     string
	authPairs authPairs
}

// ServeHTTP implements ServeHTTP method of http.Handler.
func (b basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if len(authorization) <= 0 {
		// No authorization provided, ask for it.
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%s", b.realm))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	_, found := b.authPairs.SearchCredential(r.Header.Get("Authorization"))
	if !found {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	b.Handler.ServeHTTP(w, r)
}

func NewBasicAuthHandler(realm string, auths map[string]string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return basicAuthHandler{
			Handler:   h,
			realm:     realm,
			authPairs: buildAuthPairs(auths),
		}
	}
}
