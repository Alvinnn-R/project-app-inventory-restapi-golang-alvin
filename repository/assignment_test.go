package repository

import (
	"errors"
	"project-app-inventory/model"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAssignmentRepository_Create_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	deadline := time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC)

	assignment := &model.Assignment{
		CourseID:    1,
		LecturerID:  2,
		Title:       "Golang advance",
		Description: "Golang course advance",
		Deadline:    deadline,
	}

	mockDB.
		ExpectQuery(`INSERT INTO assignments`).
		WithArgs(assignment.CourseID, assignment.LecturerID, assignment.Title, assignment.Description, assignment.Deadline).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(99))

	err = repo.Create(assignment)
	require.NoError(t, err)
	require.Equal(t, 99, assignment.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestAssignmentRepository_Create_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	deadline := time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC)

	assignment := &model.Assignment{
		CourseID:    1,
		LecturerID:  2,
		Title:       "Golang advance",
		Description: "Golang course advance",
		Deadline:    deadline,
	}

	mockDB.
		ExpectQuery(`INSERT INTO assignments`).
		WithArgs(assignment.CourseID, assignment.LecturerID, assignment.Title, assignment.Description, assignment.Deadline).
		WillReturnError(errors.New("database error"))

	err = repo.Create(assignment)
	require.Error(t, err)
	require.Equal(t, 0, assignment.ID)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_FindByID_Success tests finding assignment by ID
func TestAssignmentRepository_FindByID_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	deadline := time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC)
	mockDB.
		ExpectQuery(`SELECT (.+) FROM assignments WHERE id`).
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "course_id", "lecturer_id", "title", "description", "deadline",
		}).AddRow(1, time.Now(), time.Now(), nil, 1, 2, "Test Assignment", "Description", deadline))

	assignment, err := repo.FindByID(1)
	require.NoError(t, err)
	require.NotNil(t, assignment)
	require.Equal(t, 1, assignment.ID)
	require.Equal(t, "Test Assignment", assignment.Title)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_FindByID_NotFound tests assignment not found
func TestAssignmentRepository_FindByID_NotFound(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT (.+) FROM assignments WHERE id`).
		WithArgs(999).
		WillReturnError(pgx.ErrNoRows)

	assignment, err := repo.FindByID(999)
	require.NoError(t, err)
	require.Nil(t, assignment)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_FindAll_Success tests finding all assignments with pagination
func TestAssignmentRepository_FindAll_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	deadline := time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC)
	rows := pgxmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "course_id", "lecturer_id", "title", "description", "deadline",
	}).
		AddRow(1, time.Now(), time.Now(), nil, 1, 2, "Assignment 1", "Desc 1", deadline).
		AddRow(2, time.Now(), time.Now(), nil, 1, 2, "Assignment 2", "Desc 2", deadline)

	mockDB.
		ExpectQuery(`SELECT COUNT`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(2))

	mockDB.
		ExpectQuery(`SELECT (.+) FROM assignments WHERE deleted_at IS NULL`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	assignments, total, err := repo.FindAll(1, 10)
	require.NoError(t, err)
	require.Equal(t, 2, total)
	require.Len(t, assignments, 2)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_FindAll_Error tests findall error
func TestAssignmentRepository_FindAll_Error(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectQuery(`SELECT COUNT`).
		WillReturnError(errors.New("db error"))

	assignments, total, err := repo.FindAll(1, 10)
	require.Error(t, err)
	require.Equal(t, 0, total)
	require.Nil(t, assignments)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_Update_Success tests successful update
func TestAssignmentRepository_Update_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	deadline := time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC)
	assignment := &model.Assignment{
		CourseID:    1,
		LecturerID:  2,
		Title:       "Updated",
		Description: "Updated desc",
		Deadline:    deadline,
	}

	mockDB.
		ExpectExec(`UPDATE assignments`).
		WithArgs(assignment.CourseID, assignment.LecturerID, assignment.Title, assignment.Description, assignment.Deadline, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(1, assignment)
	require.NoError(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_Update_NotFound tests update not found
func TestAssignmentRepository_Update_NotFound(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	deadline := time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC)
	assignment := &model.Assignment{
		CourseID:    1,
		LecturerID:  2,
		Title:       "Updated",
		Description: "Updated desc",
		Deadline:    deadline,
	}

	mockDB.
		ExpectExec(`UPDATE assignments`).
		WithArgs(assignment.CourseID, assignment.LecturerID, assignment.Title, assignment.Description, assignment.Deadline, 999).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.Update(999, assignment)
	require.Error(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_Delete_Success tests successful delete
func TestAssignmentRepository_Delete_Success(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`UPDATE assignments SET deleted_at`).
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Delete(1)
	require.NoError(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

// TestAssignmentRepository_Delete_NotFound tests delete not found
func TestAssignmentRepository_Delete_NotFound(t *testing.T) {
	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := NewAssignmentRepository(mockDB, zap.NewNop())

	mockDB.
		ExpectExec(`UPDATE assignments SET deleted_at`).
		WithArgs(999).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.Delete(999)
	require.Error(t, err)
	require.NoError(t, mockDB.ExpectationsWereMet())
}
