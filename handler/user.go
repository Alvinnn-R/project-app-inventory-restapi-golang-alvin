package handler

import (
	"encoding/json"
	"net/http"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/service"
	"project-app-inventory/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService service.UserService
	Config      utils.Configuration
}

func NewUserHandler(userService service.UserService, config utils.Configuration) UserHandler {
	return UserHandler{
		UserService: userService,
		Config:      config,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request data", nil)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// Hash password
	hashedPassword := utils.HashPassword(req.Password)

	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RoleID:       req.RoleID,
		IsActive:     req.IsActive,
	}

	err = h.UserService.Create(&user)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	utils.ResponseSuccess(w, http.StatusCreated, "user created successfully", user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := h.Config.Limit

	users, pagination, err := h.UserService.GetAllUsers(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "failed to fetch users: "+err.Error(), nil)
		return
	}

	// Remove password hashes from response
	safeUsers := *users
	for i := range safeUsers {
		safeUsers[i].PasswordHash = ""
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", &safeUsers, *pagination)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "user_id")

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	user, err := h.UserService.GetUserByIDDetailed(userID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	utils.ResponseSuccess(w, http.StatusOK, "success get user by id", user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "user_id")

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	var req dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid request data", nil)
		return
	}

	// Validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	user := model.User{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: req.RoleID,
	}

	// Hash new password if provided
	if req.Password != "" {
		hashedPassword := utils.HashPassword(req.Password)
		user.PasswordHash = hashedPassword
	}

	// Handle IsActive field (pointer to bool)
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	err = h.UserService.Update(userID, &user)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "user updated successfully", nil)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "user_id")

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	err = h.UserService.Delete(userID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "user deleted successfully", nil)
}
