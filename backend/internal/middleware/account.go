package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type AccountAuthorizationMiddlewareClient struct {
	logger *zap.Logger
}

func NewAccountAuthorizationMiddleware(logger *zap.Logger) *AccountAuthorizationMiddlewareClient {
	return &AccountAuthorizationMiddlewareClient{
		logger: logger,
	}
}

type MiddlewareExemption struct {
	Path   string
	Method string
}

var (
	exemptPaths = []MiddlewareExemption{
		{
			Path:   "/api/v1/accounts/2fa/2fa",
			Method: "GET",
		},
		{
			Path:   "/api/v1/accounts/2fa/verify",
			Method: "GET",
		},
	}
)

func (aa *AccountAuthorizationMiddlewareClient) IsAccountAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range exemptPaths {
			if r.URL.Path == path.Path && r.Method == path.Method {
				aa.logger.Debug("exempting request from account authorization middleware")
				next.ServeHTTP(w, r)
				return
			}
		}

		email := r.URL.Query().Get("email")
		provider := r.URL.Query().Get("provider")

		jwt, err := GetJWTFromRequestContext(r)
		if err != nil {
			aa.logger.Info("error getting jwt from request context", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if jwt.Email != email || jwt.Provider != provider {
			aa.logger.Info("error getting 2fa: email and provider do not match jwt",
				zap.String("2fa_email", email),
				zap.String("2fa_provider", provider),
				zap.String("jwt_email", jwt.Email),
				zap.String("jwt_provider", jwt.Provider))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		aa.logger.Info("account authorized")

		next.ServeHTTP(w, r)
	})
}
