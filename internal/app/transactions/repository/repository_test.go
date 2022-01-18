package repository

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"avito-tech-task/internal/app/models"
)

func TestStorage_DoesUserExist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Errorf("Could not mock database connection: %s", err)
	}
	storage := NewStorage(mock)

	dbErr := errors.New("Error in database")
	tests := []struct {
		name        string
		userID      int64
		mock        func()
		expected    bool
		expectedErr bool
		err         error
	}{
		{
			name:   "Successfully found user in database",
			userID: 1,
			mock: func() {
				var userID int64 = 1
				rows := pgxmock.NewRows([]string{"user_id"})
				rows.AddRow(userID)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserID)).WithArgs(userID).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expected: true,
		},
		{
			name:   "User was not found in database",
			userID: 1,
			mock: func() {
				var userID int64 = 1
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserID)).WithArgs(userID).WillReturnError(pgx.ErrNoRows)
				mock.ExpectCommit()
			},
			expected: false,
		},
		{
			name:   "Error in database",
			userID: 1,
			mock: func() {
				var userID int64 = 1
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserID)).WithArgs(userID).WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
	}

	var got bool
	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err = storage.DoesUserExist(test.userID)

			if test.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}

func TestStorage_GetUserTransactions(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Errorf("Could not mock database connection: %s", err)
	}
	storage := NewStorage(mock)

	timeNow := time.Now()
	dbErr := errors.New("Error in database")
	tests := []struct {
		name        string
		userID      int64
		params      *models.TransactionsSelectionParams
		mock        func()
		expected    models.Transactions
		expectedErr bool
		err         error
	}{
		{
			name:   "Successfully get transactions list",
			userID: 1,
			params: &models.TransactionsSelectionParams{
				Limit:         10,
				OperationType: 1,
				Since:         "",
				OrderAmount:   false,
				OrderDate:     false,
			},
			mock: func() {
				var (
					userID        int64 = 1
					limit               = 10
					operationType       = "add"
					receiver      int64
					amount        float64 = 1000
					created               = timeNow
				)
				query := `SELECT operation_type, receiver, amount, created FROM transactions WHERE sender = $1 
				AND operation_type = 'add' LIMIT NULLIF($2, 0)`
				rows := pgxmock.NewRows([]string{"operation_type", "receiver", "amount", "created"})
				rows.AddRow(operationType, receiver, amount, created)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, limit).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expected: models.Transactions{
				&models.Transaction{
					OperationType: "add",
					Amount:        1000,
					Created:       timeNow,
				},
			},
		},
		{
			name:   "Error in database",
			userID: 1,
			params: &models.TransactionsSelectionParams{
				Limit:         10,
				OperationType: 1,
				Since:         "",
				OrderAmount:   false,
				OrderDate:     false,
			},
			mock: func() {
				var (
					userID int64 = 1
					limit        = 10
				)
				query := `SELECT operation_type, receiver, amount, created FROM transactions WHERE sender = $1 
				AND operation_type = 'add' LIMIT NULLIF($2, 0)`
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, limit).WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
	}

	var got models.Transactions
	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err = storage.GetUserTransactions(test.userID, test.params)

			if test.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}
