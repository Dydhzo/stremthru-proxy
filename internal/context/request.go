package context

import (
	"context"
	"net/http"
)

type proxyContextKey struct{}

type ProxyContext struct {
	IsProxyAuthorized bool
	ProxyAuthUser     string
	ProxyAuthPassword string
}

func SetProxyContext(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), proxyContextKey{}, &ProxyContext{})
	return r.WithContext(ctx)
}

func GetProxyContext(r *http.Request) *ProxyContext {
	return r.Context().Value(proxyContextKey{}).(*ProxyContext)
}
