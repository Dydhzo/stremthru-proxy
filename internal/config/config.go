package config

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Dydhzo/stremthru-proxy/core"
)

const (
	EnvDev  string = "dev"
	EnvProd string = "prod"
	EnvTest string = "test"
)

var Environment = func() string {
	if testing.Testing() {
		return EnvTest
	}

	value, _ := os.LookupEnv("STREMTHRU_ENV")
	switch value {
	case "dev", "development":
		return EnvDev
	case "prod", "production":
		return EnvProd
	case "test":
		return EnvTest
	default:
		return ""
	}
}()

var defaultValueByEnv = map[string]map[string]string{
	EnvDev: {
		"STREMTHRU_LOG_FORMAT": "text",
		"STREMTHRU_LOG_LEVEL":  "DEBUG",
	},
	EnvProd: {},
	EnvTest: {
		"STREMTHRU_LOG_FORMAT": "text",
		"STREMTHRU_LOG_LEVEL":  "DEBUG",
	},
	"": {
		"STREMTHRU_BASE_URL":   "http://localhost:8080",
		"STREMTHRU_LOG_FORMAT": "json",
		"STREMTHRU_LOG_LEVEL":  "INFO",
		"STREMTHRU_PORT":       "8080",
		"STREMTHRU_LANDING_PAGE": "{}",
		"STREMTHRU_IP_CHECKER": "akamai",
	},
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists && len(value) > 0 {
		return value
	}
	if val, found := defaultValueByEnv[Environment][key]; found && len(val) > 0 {
		return val
	}
	if Environment != "" {
		if val, found := defaultValueByEnv[""][key]; found && len(val) > 0 {
			return val
		}
	}
	return ""
}

type AppState struct{}

type ProxyAuthMap map[string]string

func (pam ProxyAuthMap) IsAuthorized(user, password string) bool {
	if storedPassword, exists := pam[user]; exists {
		return storedPassword == password
	}
	return false
}

var config = func() struct {
	BaseURL        string
	Port           string
	LogLevel       string
	LogFormat      string
	ProxyAuth      ProxyAuthMap
	IsPublicInstance bool
} {
	// Parse proxy auth
	proxyAuthMap := make(ProxyAuthMap)
	proxyAuthCredList := strings.FieldsFunc(getEnv("STREMTHRU_PROXY_AUTH"), func(c rune) bool {
		return c == ','
	})

	for _, cred := range proxyAuthCredList {
		if basicAuth, err := core.ParseBasicAuth(cred); err == nil {
			proxyAuthMap[basicAuth.Username] = basicAuth.Password
		}
	}

	return struct {
		BaseURL        string
		Port           string
		LogLevel       string
		LogFormat      string
		ProxyAuth      ProxyAuthMap
		IsPublicInstance bool
	}{
		BaseURL:          getEnv("STREMTHRU_BASE_URL"),
		Port:             getEnv("STREMTHRU_PORT"),
		LogLevel:         getEnv("STREMTHRU_LOG_LEVEL"),
		LogFormat:        getEnv("STREMTHRU_LOG_FORMAT"),
		ProxyAuth:        proxyAuthMap,
		IsPublicInstance: len(proxyAuthMap) == 0,
	}
}()

// Exported variables
var BaseURL = config.BaseURL
var Port = config.Port
var LogLevel = config.LogLevel
var LogFormat = config.LogFormat
var ProxyAuth = config.ProxyAuth
var ProxyAuthPassword = func() []string {
	var passwords []string
	for _, pass := range ProxyAuth {
		passwords = append(passwords, pass)
	}
	return passwords
}()
var IsPublicInstance = config.IsPublicInstance
var LandingPage = getEnv("STREMTHRU_LANDING_PAGE")
var Version = "v1.0.0"

func PrintConfig(state *AppState) {
	l := log.New(os.Stderr, "=   ", 0)

	l.Println("=== StremThru Proxy ===")
	l.Println()
	l.Println(" Proxy:")
	l.Println("   base_url: " + BaseURL)
	l.Println("      port: " + Port)
	l.Println("  log_level: " + LogLevel)
	l.Println(" log_format: " + LogFormat)

	if len(ProxyAuth) > 0 {
		l.Println("      users:", len(ProxyAuth))
	} else {
		l.Println("  auth: disabled (public)")
	}

	l.Println()
	l.Print("=======================\n\n")
}