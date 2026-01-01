package dto

type AssignmentRequest struct {
	CourseID    int    `json:"course_id" validate:"required,gt=0"`
	LecturerID  int    `json:"lecturer_id" validate:"required,gt=0"`
	Title       string `json:"title" validate:"required,min=3,max=150"`
	Description string `json:"description" validate:"required,min=10"`
	Deadline    string `json:"deadline" validate:"required"`
}

type AssignmentResponse struct {
	CourseID    int    `json:"course_id"`
	LecturerID  int    `json:"lecturer_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}
