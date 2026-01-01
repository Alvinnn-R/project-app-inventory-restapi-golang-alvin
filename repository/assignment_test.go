package repository

import (
	"errors"
	"project-app-inventory/model"
	"testing"
	"time"

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
