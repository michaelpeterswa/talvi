package handlers

import (
	"fmt"
	"net/http"

	"github.com/michaelpeterswa/talvi/backend/internal/middleware"
)

func getJWTFromRequestContext(r *http.Request) (*middleware.JWT, error) {
	jwtVal := r.Context().Value(middleware.JWTContextKey)
	if jwtVal == nil {
		return nil, fmt.Errorf("jwt not found in context")
	}

	jwt, ok := jwtVal.(middleware.JWT)
	if !ok {
		return nil, fmt.Errorf("error casting jwt from context")
	}

	return &jwt, nil
}
