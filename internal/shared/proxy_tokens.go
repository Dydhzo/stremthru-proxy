package shared

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Dydhzo/stremthru-proxy/core"
	"github.com/Dydhzo/stremthru-proxy/internal/cache"
	"github.com/Dydhzo/stremthru-proxy/internal/config"
)

type proxyLinkTokenData struct {
	EncLink    string            `json:"enc_link"`
	EncFormat  string            `json:"enc_format"`
	TunnelType config.TunnelType `json:"tunt,omitempty"`
}

type proxyLinkData struct {
	User    string            `json:"u"`
	Value   string            `json:"v"`
	Headers map[string]string `json:"reqh,omitempty"`
	TunT    config.TunnelType `json:"tunt,omitempty"`
}

var proxyLinkTokenCache = func() cache.Cache[proxyLinkData] {
	return cache.NewCache[proxyLinkData](&cache.CacheConfig{
		Name:     "store:proxyLinkToken",
		Lifetime: 30 * time.Minute,
	})
}()

func CreateProxyLink(r *http.Request, link string, headers map[string]string, tunnelType config.TunnelType, expiresIn time.Duration, user, password string, shouldEncrypt bool, filename string) (string, error) {
	var encodedToken string

	if !shouldEncrypt && expiresIn == 0 {
		blob, err := json.Marshal(proxyLinkData{
			User:    user + ":" + password,
			Value:   link,
			Headers: headers,
			TunT:    tunnelType,
		})
		if err != nil {
			return "", err
		}
		encodedToken = "base64." + core.Base64EncodeByte(blob)
	} else {
		linkBlob := link
		if headers != nil {
			for k, v := range headers {
				linkBlob += "\n" + k + ": " + v
			}
		}

		var encLink string
		var encFormat string

		if shouldEncrypt {
			encryptedLink, err := core.Encrypt(password, linkBlob)
			if err != nil {
				return "", err
			}
			encLink = encryptedLink
			encFormat = core.EncryptionFormat
		} else {
			encLink = core.Base64Encode(linkBlob)
			encFormat = "base64"
		}

		claims := core.JWTClaims[proxyLinkTokenData]{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:  "stremthru",
				Subject: user,
			},
			Data: &proxyLinkTokenData{
				EncLink:    encLink,
				EncFormat:  encFormat,
				TunnelType: tunnelType,
			},
		}
		if expiresIn != 0 {
			claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiresIn))
		}

		token, err := core.GenerateJWT(claims)
		if err != nil {
			return "", err
		}
		encodedToken = token
	}

	baseURL := ExtractRequestBaseURL(r)
	proxyURL := baseURL.String() + "/v0/proxy/" + encodedToken
	if filename != "" {
		proxyURL += "/" + filename
	}

	return proxyURL, nil
}

func UnwrapProxyLinkToken(encodedToken string) (user string, link string, headers map[string]string, tunnelType config.TunnelType, err error) {
	if cached, ok := proxyLinkTokenCache.Get(encodedToken); ok {
		return cached.User, cached.Value, cached.Headers, cached.TunT, nil
	}

	proxyLink := &proxyLinkData{}

	if strings.HasPrefix(encodedToken, "base64.") {
		blob, err := core.Base64DecodeByte(strings.TrimPrefix(encodedToken, "base64."))
		if err != nil {
			return "", "", nil, "", err
		}
		if err := json.Unmarshal(blob, proxyLink); err != nil {
			return "", "", nil, "", err
		}
		user, pass, _ := strings.Cut(proxyLink.User, ":")
		if pass != config.ProxyAuth[user] {
			err := core.NewAPIError("unauthorized")
			err.StatusCode = http.StatusUnauthorized
			return "", "", nil, "", err
		}
		proxyLink.User = user
	} else {
		// JWT token - parse with our existing function
		claims, err := core.ParseJWT[proxyLinkTokenData](encodedToken)
		if err != nil {
			rerr := core.NewAPIError("unauthorized")
			rerr.StatusCode = http.StatusUnauthorized
			rerr.Cause = err
			return "", "", nil, "", rerr
		}

		user = claims.Subject

		// For JWT tokens, we need the user's password for decryption
		password := config.ProxyAuth[user]

		var linkBlob string
		if claims.Data.EncFormat == "base64" {
			blob, err := core.Base64Decode(claims.Data.EncLink)
			if err != nil {
				return "", "", nil, "", err
			}
			linkBlob = blob
		} else {
			blob, err := core.Decrypt(password, claims.Data.EncLink)
			if err != nil {
				return "", "", nil, "", err
			}
			linkBlob = blob
		}

		link, headersBlob, hasHeaders := strings.Cut(linkBlob, "\n")

		proxyLink.User = user
		proxyLink.TunT = claims.Data.TunnelType
		proxyLink.Value = link

		if hasHeaders {
			proxyLink.Headers = map[string]string{}
			for _, header := range strings.Split(headersBlob, "\n") {
				if k, v, ok := strings.Cut(header, ": "); ok {
					proxyLink.Headers[k] = v
				}
			}
		}
	}

	proxyLinkTokenCache.Set(encodedToken, *proxyLink)

	return proxyLink.User, proxyLink.Value, proxyLink.Headers, proxyLink.TunT, nil
}