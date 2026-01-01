package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"project-app-inventory/dto"
	"project-app-inventory/model"
	"project-app-inventory/repository"
	"project-app-inventory/utils"
	"time"
)

type AssignmentService interface {
	Create(assignmenet *model.Assignment) error
	GetAllAssignments(page, limit int) (*[]model.Assignment, *dto.Pagination, error)
	SubmitAssignment(studentID, assignmentID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	// GetGradeFormData() ([]model.User, []model.Assignment, error)
	GetAssignmentByID(id int) (*model.Assignment, error)
	Update(id int, data *model.Assignment) error
	Delete(id int) error
}

type assignmentService struct {
	Repo repository.Repository
}

func NewAssignmentService(repo repository.Repository) AssignmentService {
	return &assignmentService{Repo: repo}
}

func (as *assignmentService) Create(assignmenet *model.Assignment) error {
	return as.Repo.AssignmentRepo.Create(assignmenet)
}

func (as *assignmentService) GetAllAssignments(page, limit int) (*[]model.Assignment, *dto.Pagination, error) {
	assignments, total, err := as.Repo.AssignmentRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}
	return &assignments, &pagination, nil
}

func (as *assignmentService) GetAssignmentByID(id int) (*model.Assignment, error) {
	return as.Repo.AssignmentRepo.FindByID(id)
}

func (as *assignmentService) SubmitAssignment(studentID, assignmentID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	assignment, err := as.Repo.AssignmentRepo.FindByID(assignmentID)
	if err != nil {
		return "", err
	}

	count, err := as.Repo.SubmissionRepo.CountByStudentAndAssignment(studentID, assignmentID)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "already submitted", nil
	}

	// save file to disk
	uploadDir := "uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	filename := fmt.Sprintf("%d_%d_%s", assignmentID, studentID, fileHeader.Filename)
	filepath := fmt.Sprintf("%s/%s", uploadDir, filename)
	accessURL := fmt.Sprintf("http://localhost:8080/%s/%s", uploadDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	status := "submitted"
	if time.Now().After(assignment.Deadline) {
		status = "late"
	}

	sub := &model.Submission{
		AssignmentID: assignmentID,
		StudentID:    studentID,
		SubmittedAt:  time.Now(),
		FileURL:      accessURL,
		Status:       status,
	}

	return status, as.Repo.SubmissionRepo.Create(sub)
}

func (as *assignmentService) Update(id int, assignment *model.Assignment) error {
	return as.Repo.AssignmentRepo.Update(id, assignment)
}

func (as *assignmentService) Delete(id int) error {
	return as.Repo.AssignmentRepo.Delete(id)
}

// func (s *assignmentService) GetGradeFormData() ([]model.User, []model.Assignment, error) {
// 	students, err := s.Repo.UserRepo.FindAllStudents()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	assignments, err := s.Repo.AssignmentRepo.FindAll()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return students, assignments, nil
// }
