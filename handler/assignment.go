package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/service"
	"project-app-inventory/utils"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type AssignmentHandler struct {
	AssignmentService service.AssignmentService
	Config            utils.Configuration
}

func NewAssignmentHandler(assignmentService service.AssignmentService, config utils.Configuration) AssignmentHandler {
	return AssignmentHandler{
		AssignmentService: assignmentService,
		Config:            config,
	}
}

func (assignmentHandler *AssignmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.AssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data", nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing time
	deadLine, err := time.ParseInLocation("2006-01-02 15:04:05", req.Deadline, time.Local)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error parsing deadline time :"+err.Error(), nil)
		return
	}

	// parsing to model assignment
	assignment := model.Assignment{
		CourseID:   req.CourseID,
		LecturerID: req.LecturerID,
		Title:      req.Title,
		Deadline:   deadLine,
	}

	// create assignment service
	err = assignmentHandler.AssignmentService.Create(&assignment)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success created assignment", nil)
}

func (assignmentHandler *AssignmentHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	// config limit pagination
	limit := assignmentHandler.Config.Limit

	// Get data assignment form service all assignment
	assignments, pagination, err := assignmentHandler.AssignmentService.GetAllAssignments(page, limit)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch assignments: "+err.Error(), nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", assignments, *pagination)

}

func (assignmentHandler *AssignmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	assignmentIDstr := chi.URLParam(r, "assignment_id")

	assignmentID, err := strconv.Atoi(assignmentIDstr)
	if err != nil {
		return
	}

	response, err := assignmentHandler.AssignmentService.GetAssignmentByID(assignmentID)
	if err != nil {
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data assignment by id", response)
}

func (assignmentHandler *AssignmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	assignmentIDstr := chi.URLParam(r, "assignment_id")

	assignmentID, err := strconv.Atoi(assignmentIDstr)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error param assignment id :"+err.Error(), nil)
		return
	}

	var req dto.AssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error data :"+err.Error(), nil)
		return
	}

	// validation
	messages, err := utils.ValidateErrors(req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, err.Error(), messages)
		return
	}

	// parsing time
	deadLine, err := time.ParseInLocation("2006-01-02 15:04:05", req.Deadline, time.Local)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "error parsing deadline time :"+err.Error(), nil)
		return
	}

	// parsing to model assignment
	assignment := model.Assignment{
		CourseID:   req.CourseID,
		LecturerID: req.LecturerID,
		Title:      req.Title,
		Deadline:   deadLine,
	}

	err = assignmentHandler.AssignmentService.Update(assignmentID, &assignment)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Error update :"+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Updated Success", nil)
}

func (assignmentHandler *AssignmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	assignmentIDstr := chi.URLParam(r, "assignment_id")

	assignmentID, err := strconv.Atoi(assignmentIDstr)
	if err != nil {
		return
	}

	err = assignmentHandler.AssignmentService.Delete(assignmentID)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadRequest, "Error delete :"+err.Error(), nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Deleted Success", nil)
}

func (AssignmentHandler *AssignmentHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "error file size", http.StatusBadRequest)
			return
		}
	}

	// get assignment id
	assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	// get student id
	c, _ := r.Cookie("session")
	idStr := strings.TrimPrefix(c.Value, "lumos-")
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// get file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "error file", http.StatusBadRequest)
		return
	}

	status, err := AssignmentHandler.AssignmentService.SubmitAssignment(studentID, assignmentID, file, fileHeader)
	if err != nil {
		http.Error(w, "error submit", http.StatusBadRequest)
		return
	}

	fmt.Println(status)
	http.Redirect(w, r, "/user/success-submit", http.StatusSeeOther)
}
