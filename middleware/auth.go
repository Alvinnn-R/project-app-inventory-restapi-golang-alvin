package middleware

import (
	"context"
	"net/http"
	"project-app-inventory/utils"
	"strings"
)

// AuthMiddleware validates Bearer token from Authorization header
func (mw *MiddlewareCostume) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ResponseBadRequest(w, http.StatusUnauthorized, "missing authorization header", nil)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ResponseBadRequest(w, http.StatusUnauthorized, "invalid authorization header format", nil)
			return
		}

		token := parts[1]

		// Validate token
		user, err := mw.Service.AuthService.ValidateToken(token)
		if err != nil {
			utils.ResponseBadRequest(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		// Store user in context for use in handlers
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware checks if user has required role
func (mw *MiddlewareCostume) RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (set by AuthMiddleware)
			user := r.Context().Value("user")
			if user == nil {
				utils.ResponseBadRequest(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}

			// Check if user role is in allowed roles
			userRole := user.(*struct {
				ID           int
				Name         string
				Email        string
				PasswordHash string
				RoleID       int
				RoleName     string
				IsActive     bool
			}).RoleName

			allowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				utils.ResponseBadRequest(w, http.StatusForbidden, "access denied", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
