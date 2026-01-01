package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory/dto"
	"project-app-inventory/service"
	"project-app-inventory/utils"
	"strings"
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
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request data", nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Login
	result, err := h.AuthService.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "login success", result)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
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

	// Logout (revoke token)
	err := h.AuthService.AuthService.Logout(token)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "logout success", nil)
}
