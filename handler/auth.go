package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory/dto"
	"project-app-inventory/service"
	"project-app-inventory/utils"
)

type AuthHandler struct {
	AuthService service.Service
}

func NewAuthHandler(authService service.Service) AuthHandler {
	return AuthHandler{
		AuthService: authService,
	}
}

// func (h *AuthHandler) LoginView(w http.ResponseWriter, r *http.Request) {
// 	if err := h.Templates.ExecuteTemplate(w, "login", nil); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data", nil)
		return
	}

	result, err := h.AuthService.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusUnauthorized, "email or password failed", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "login success", result)
}

// func (h *AuthHandler) LogoutView(w http.ResponseWriter, r *http.Request) {
// 	if err := h.Templates.ExecuteTemplate(w, "logout", nil); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
// 	// cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "session",
// 		Value:    "",
// 		Path:     "/",
// 		MaxAge:   -1,
// 		HttpOnly: true,
// 	})
// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
// }
