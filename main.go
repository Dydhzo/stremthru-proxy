// Package main provides the entry point for StremThru proxy server
package main

import (
	"log"
	"net/http"

	"github.com/Dydhzo/stremthru-proxy/internal/config"
	"github.com/Dydhzo/stremthru-proxy/internal/endpoint"
	"github.com/Dydhzo/stremthru-proxy/internal/shared"
)

func main() {
	// SECURITY: Proxy authentication is MANDATORY
	if len(config.ProxyAuth) == 0 {
		log.Fatalf("❌ FATAL: STREMTHRU_PROXY_AUTH is required but not configured!\n" +
			"   A proxy server MUST have authentication for security.\n" +
			"   Please set STREMTHRU_PROXY_AUTH=username:password")
	}

	config.PrintConfig(&config.AppState{})

	mux := http.NewServeMux()

	// Only keep essential endpoints
	endpoint.AddRootEndpoint(mux)
	endpoint.AddHealthEndpoints(mux)
	endpoint.AddProxyEndpoints(mux)
	endpoint.AddStatsEndpoint(mux)

	handler := shared.RootServerContext(mux)

	addr := ":" + config.Port
	if config.Environment == config.EnvDev {
		addr = "localhost" + addr
	}
	server := &http.Server{Addr: addr, Handler: handler}

	// Authentication is guaranteed to exist (checked at startup)

	log.Println("StremThru Proxy listening on " + addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start proxy: %v", err)
	}
}
