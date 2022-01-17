package repository

import (
	"avito-tech-task/internal/app/models"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestStorage_DoesUserExist(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Errorf("Could not mock database connection: %s", err)
	}
	query := `SELECT user_id FROM balance WHERE user_id = $1`
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
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID).WillReturnRows(rows)
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
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID).WillReturnError(pgx.ErrNoRows)
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
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID).WillReturnError(dbErr)
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
	dbErr := errors.New("Error in database")

	timeNow := time.Now()
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
				Limit:       10,
				Since:       "",
				OrderAmount: false,
				OrderDate:   false,
			},
			mock: func() {
				var (
					userID      int64   = 1
					limit       int     = 10
					description         = "description"
					amount      float64 = 1000
					time                = timeNow
				)
				query := `SELECT description, amount, created FROM transactions WHERE user_id = $1 LIMIT NULLIF($2, 0)`
				rows := pgxmock.NewRows([]string{"description", "amount", "created"})
				rows.AddRow(description, amount, time)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, limit).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expected: models.Transactions{
				&models.Transaction{
					Description: "description",
					Amount:      1000,
					Created:     timeNow,
				},
			},
		},
		{
			name:   "Error in database",
			userID: 1,
			params: &models.TransactionsSelectionParams{
				Limit:       10,
				Since:       "",
				OrderAmount: false,
				OrderDate:   false,
			},
			mock: func() {
				var (
					userID int64 = 1
					limit  int   = 10
				)
				query := `SELECT description, amount, created FROM transactions WHERE user_id = $1 LIMIT NULLIF($2, 0)`
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
