package repository

import (
	"errors"
	"project-app-inventory/model"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

// TestSubmissionRepository_CountByStudentAndAssignment_Success tests counting submissions
func TestSubmissionRepository_CountByStudentAndAssignment_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM submissions WHERE student_id`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(3))

	count, err := repo.CountByStudentAndAssignment(1, 1)

	require.NoError(t, err)
	require.Equal(t, int64(3), count)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_CountByStudentAndAssignment_Error tests count error
func TestSubmissionRepository_CountByStudentAndAssignment_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)

	mockDB.
		ExpectQuery(`SELECT COUNT\(\*\) FROM submissions WHERE student_id`).
		WithArgs(1, 1).
		WillReturnError(errors.New("db error"))

	count, err := repo.CountByStudentAndAssignment(1, 1)

	require.Error(t, err)
	require.Equal(t, int64(0), count)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_Create_Success tests submission creation
func TestSubmissionRepository_Create_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)
	submittedAt := time.Now()

	submission := &model.Submission{
		AssignmentID: 1,
		StudentID:    1,
		SubmittedAt:  submittedAt,
		FileURL:      "http://example.com/file.pdf",
		Status:       "submitted",
	}

	mockDB.
		ExpectExec(`INSERT INTO submissions`).
		WithArgs(submission.AssignmentID, submission.StudentID, submission.SubmittedAt, submission.FileURL, submission.Status).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(submission)

	require.NoError(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_Create_Error tests creation error
func TestSubmissionRepository_Create_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)
	submittedAt := time.Now()
	submission := &model.Submission{
		AssignmentID: 1,
		StudentID:    1,
		SubmittedAt:  submittedAt,
		FileURL:      "file.pdf",
		Status:       "submitted",
	}

	mockDB.
		ExpectExec(`INSERT INTO submissions`).
		WithArgs(submission.AssignmentID, submission.StudentID, submission.SubmittedAt, submission.FileURL, submission.Status).
		WillReturnError(errors.New("db error"))

	err = repo.Create(submission)

	require.Error(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_GetAllWithStudentAndAssignment_Success tests getting all submissions
func TestSubmissionRepository_GetAllWithStudentAndAssignment_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)
	grade1 := 85.0
	grade2 := 90.0

	rows := pgxmock.NewRows([]string{
		"id", "assignment_id", "student_id", "student_name", "assignment_title", "file_url", "status", "grade",
	}).
		AddRow(1, 1, 1, "John Doe", "Assignment 1", "file.pdf", "submitted", &grade1).
		AddRow(2, 1, 2, "Jane Doe", "Assignment 1", "file2.pdf", "graded", &grade2)

	mockDB.
		ExpectQuery(`SELECT (.+) FROM submissions s JOIN users u`).
		WillReturnRows(rows)

	submissions, err := repo.GetAllWithStudentAndAssignment()

	require.NoError(t, err)
	require.NotNil(t, submissions)
	require.Equal(t, 2, len(submissions))
	require.Equal(t, "John Doe", submissions[0].StudentName)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_GetAllWithStudentAndAssignment_Error tests query error
func TestSubmissionRepository_GetAllWithStudentAndAssignment_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)

	mockDB.
		ExpectQuery(`SELECT (.+) FROM submissions s JOIN users u`).
		WillReturnError(errors.New("db error"))

	submissions, err := repo.GetAllWithStudentAndAssignment()

	require.Error(t, err)
	require.Nil(t, submissions)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_FindByStudentAndAssignment_Success tests finding submission
func TestSubmissionRepository_FindByStudentAndAssignment_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)
	submittedAt := time.Now()
	grade := 85.0

	mockDB.
		ExpectQuery(`SELECT (.+) FROM submissions WHERE student_id`).
		WithArgs(1, 1).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "assignment_id", "student_id", "submitted_at", "file_url", "status", "grade",
		}).AddRow(1, 1, 1, submittedAt, "file.pdf", "submitted", &grade))

	submission, err := repo.FindByStudentAndAssignment(1, 1)

	require.NoError(t, err)
	require.NotNil(t, submission)
	require.Equal(t, 1, submission.StudentID)
	require.NotNil(t, submission.Grade)
	require.Equal(t, 85.0, *submission.Grade)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_FindByStudentAndAssignment_NotFound tests not found
func TestSubmissionRepository_FindByStudentAndAssignment_NotFound(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)

	mockDB.
		ExpectQuery(`SELECT (.+) FROM submissions WHERE student_id`).
		WithArgs(999, 999).
		WillReturnError(errors.New("no rows"))

	submission, err := repo.FindByStudentAndAssignment(999, 999)

	require.Error(t, err)
	require.Nil(t, submission)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_UpdateGrade_Success tests updating grade
func TestSubmissionRepository_UpdateGrade_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)
	grade := 95.0
	submission := &model.Submission{
		StudentID:    1,
		AssignmentID: 1,
		Grade:        &grade,
	}

	mockDB.
		ExpectExec(`UPDATE submissions SET grade`).
		WithArgs(submission.Grade, submission.StudentID, submission.AssignmentID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateGrade(submission)

	require.NoError(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestSubmissionRepository_UpdateGrade_Error tests update error
func TestSubmissionRepository_UpdateGrade_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewSubmissionRepo(mockDB)
	grade := 95.0
	submission := &model.Submission{StudentID: 1, AssignmentID: 1, Grade: &grade}

	mockDB.
		ExpectExec(`UPDATE submissions SET grade`).
		WithArgs(submission.Grade, submission.StudentID, submission.AssignmentID).
		WillReturnError(errors.New("db error"))

	err = repo.UpdateGrade(submission)

	require.Error(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}
