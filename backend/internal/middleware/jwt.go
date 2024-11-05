package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-jose/go-jose/v3"
	"go.uber.org/zap"
)

var (
	bearerTokenRegex = regexp.MustCompile(`^Bearer\s(.*)$`)
	JWTContextKey    = JWTContextKeyType("jwt")
)

type JWTContextKeyType string

type JWT struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	Role     string `json:"role"`
	Sub      string `json:"sub"`
	ID       string `json:"id"`
	Provider string `json:"provider"`
	Iat      int    `json:"iat"`
	Exp      int    `json:"exp"`
	Jti      string `json:"jti"`
}

type JWTMiddlewareClient struct {
	secret []byte
	logger *zap.Logger
}

func NewJWTMiddleware(logger *zap.Logger, secret string) (*JWTMiddlewareClient, error) {
	decSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("error decoding secret: %w", err)
	}

	return &JWTMiddlewareClient{
		logger: logger,
		secret: decSecret,
	}, nil
}

func (jwtmc *JWTMiddlewareClient) JWTMiddleware(next http.Handler) http.Handler {
	jwtmc.logger.Debug("request received")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jwtmc.logger.Info("no authorization header found")
			next.ServeHTTP(w, r)
			return
		}
		jwe := bearerTokenRegex.FindStringSubmatch(authHeader)[1]

		object, err := jose.ParseEncrypted(jwe)
		if err != nil {
			jwtmc.logger.Info("error parsing encrypted jwt", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		decrypted, err := object.Decrypt(jwtmc.secret)
		if err != nil {
			jwtmc.logger.Info("error decrypting jwt", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var jwt JWT
		err = json.Unmarshal(decrypted, &jwt)
		if err != nil {
			jwtmc.logger.Info("error unmarshalling jwt", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtContext := context.WithValue(r.Context(), JWTContextKey, jwt)
		next.ServeHTTP(w, r.WithContext(jwtContext))
	})
}

func GetJWTFromRequestContext(r *http.Request) (*JWT, error) {
	jwtVal := r.Context().Value(JWTContextKey)
	if jwtVal == nil {
		return nil, fmt.Errorf("jwt not found in context")
	}

	jwt, ok := jwtVal.(JWT)
	if !ok {
		return nil, fmt.Errorf("error getting jwt from context")
	}

	return &jwt, nil
}
