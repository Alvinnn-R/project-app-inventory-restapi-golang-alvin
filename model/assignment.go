package model

import "time"

type Assignment struct {
	Model
	CourseID    int       `json:"course_id"`
	LecturerID  int       `json:"lecturer_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}
