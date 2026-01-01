package middleware

import (
	"net/http"
	"strconv"
)

func (middlewareCostume *MiddlewareCostume) RequirePermission(code string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDstr, _ := r.Cookie("session")
			userID, _ := strconv.Atoi(userIDstr.Value)

			allowed, err := middlewareCostume.Service.PermissionService.Allowed(userID, code)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
