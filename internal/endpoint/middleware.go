package endpoint

import (
	"net/http"
	"strings"

	"github.com/Dydhzo/stremthru-proxy/core"
	"github.com/Dydhzo/stremthru-proxy/internal/config"
	"github.com/Dydhzo/stremthru-proxy/internal/context"
	"github.com/Dydhzo/stremthru-proxy/internal/server"
)

func extractProxyAuthToken(r *http.Request, readQuery bool) (token string, hasToken bool) {
	token = r.Header.Get(server.HEADER_STREMTHRU_AUTHORIZATION)
	if token == "" {
		token = r.Header.Get(server.HEADER_PROXY_AUTHORIZATION)
		if token != "" {
			r.Header.Del(server.HEADER_PROXY_AUTHORIZATION)
		}
	}
	if token == "" && readQuery {
		token = r.URL.Query().Get("token")
	}
	token = strings.TrimPrefix(token, "Basic ")
	return token, token != ""
}

func getProxyAuthorization(r *http.Request, readQuery bool) (isAuthorized bool, user, pass string) {
	token, hasToken := extractProxyAuthToken(r, readQuery)
	auth, err := core.ParseBasicAuth(token)
	isAuthorized = hasToken && err == nil && config.ProxyAuth.IsAuthorized(auth.Username, auth.Password)
	user = auth.Username
	pass = auth.Password
	return isAuthorized, user, pass
}

func ProxyAuthContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.GetProxyContext(r)
		ctx.IsProxyAuthorized, ctx.ProxyAuthUser, ctx.ProxyAuthPassword = getProxyAuthorization(r, false)
		next.ServeHTTP(w, r)
	})
}
