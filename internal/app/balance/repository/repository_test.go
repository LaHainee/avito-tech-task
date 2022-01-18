package repository

import (
	createdErrors "avito-tech-task/internal/pkg/errors"
	"errors"
	"regexp"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"avito-tech-task/internal/app/models"
)

func TestStorage_GetUserData(t *testing.T) {
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
		expected    *models.UserData
		expectedErr bool
		err         error
	}{
		{
			name:   "Successfully get user data from database",
			userID: 1,
			mock: func() {
				var (
					balance float64 = 1000
					userID  int64   = 1
				)
				rows := pgxmock.NewRows([]string{"balance"})
				rows.AddRow(balance)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBalance)).WithArgs(userID).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expected: &models.UserData{
				UserID:  1,
				Balance: 1000,
			},
		},
		{
			name:   "User not found in database",
			userID: 1,
			mock: func() {
				var userID int64 = 1
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBalance)).WithArgs(userID).WillReturnError(pgx.ErrNoRows)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         createdErrors.ErrUserDoesNotExist,
		},
		{
			name:   "Error in database",
			userID: 1,
			mock: func() {
				var userID int64 = 1
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetBalance)).WithArgs(userID).WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
	}

	var got *models.UserData
	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err = storage.GetUserData(test.userID)

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

func TestStorage_CreateAccount(t *testing.T) {
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
		expectedErr bool
		err         error
	}{
		{
			name:   "Successfully created new account",
			userID: 1,
			mock: func() {
				var userID int64 = 1
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryInsertBalance)).WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			},
		},
		{
			name:   "Error in database",
			userID: 1,
			mock: func() {
				var userID int64 = 1
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryInsertBalance)).WithArgs(userID).WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err = storage.CreateAccount(test.userID)

			if test.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStorage_UpdateBalance(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Errorf("Could not mock database connection: %s", err)
	}
	storage := NewStorage(mock)
	dbErr := errors.New("Error in database")

	tests := []struct {
		name        string
		userID      int64
		amount      float64
		mock        func()
		expected    float64
		expectedErr bool
		err         error
	}{
		{
			name:   "Successfully updated balance",
			userID: 1,
			amount: 1000,
			mock: func() {
				var (
					userID         int64   = 1
					amount         float64 = 1000
					updatedBalance float64 = 2000
					operationType          = "add"
				)
				rows := pgxmock.NewRows([]string{"balance"})
				rows.AddRow(updatedBalance)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount, userID).
					WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta(querySaveTransaction)).WithArgs(operationType, userID, 0, amount).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			},
			expected: 2000,
		},
		{
			name:   "Error in database during updating balance",
			userID: 1,
			amount: 1000,
			mock: func() {
				var (
					userID int64   = 1
					amount float64 = 1000
				)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount, userID).WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
		{
			name:   "Error in database during saving transaction info",
			userID: 1,
			amount: -1000,
			mock: func() {
				var (
					userID         int64   = 1
					amount         float64 = -1000
					updatedBalance float64
					operationType  = "write_off"
				)
				rows := pgxmock.NewRows([]string{"balance"})
				rows.AddRow(updatedBalance)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount, userID).
					WillReturnRows(rows)
				mock.ExpectExec(regexp.QuoteMeta(querySaveTransaction)).WithArgs(operationType, userID, 0, amount*-1).
					WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
	}

	var got float64
	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err = storage.UpdateBalance(test.userID, test.amount)

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

func TestStorage_MakeTransfer(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Errorf("Could not mock database connection: %s", err)
	}
	storage := NewStorage(mock)
	dbErr := errors.New("Error in database")

	tests := []struct {
		name        string
		senderID    int64
		receiverID  int64
		amount      float64
		mock        func()
		expectedErr bool
		err         error
	}{
		{
			name:       "Successfully transferred money",
			senderID:   1,
			receiverID: 2,
			amount:     1000,
			mock: func() {
				var (
					senderID      int64   = 1
					receiverID    int64   = 2
					amount        float64 = 1000
					operationType         = "transfer"
				)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount*-1, senderID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount, receiverID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(querySaveTransaction)).
					WithArgs(operationType, senderID, receiverID, amount).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			},
		},
		{
			name:       "Error in database during writing off money from sender",
			senderID:   1,
			receiverID: 2,
			amount:     1000,
			mock: func() {
				var (
					senderID int64   = 1
					amount   float64 = 1000
				)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount*-1, senderID).
					WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
		{
			name:       "Error in database during adding money to receiver",
			senderID:   1,
			receiverID: 2,
			amount:     1000,
			mock: func() {
				var (
					senderID   int64   = 1
					receiverID int64   = 2
					amount     float64 = 1000
				)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount*-1, senderID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount, receiverID).
					WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
		{
			name:       "Error in database during saving sender transaction",
			senderID:   1,
			receiverID: 2,
			amount:     1000,
			mock: func() {
				var (
					senderID      int64   = 1
					receiverID    int64   = 2
					amount        float64 = 1000
					operationType         = "transfer"
				)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount*-1, senderID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(queryUpdateBalance)).WithArgs(amount, receiverID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(regexp.QuoteMeta(querySaveTransaction)).
					WithArgs(operationType, senderID, receiverID, amount).
					WillReturnError(dbErr)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbErr,
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err = storage.MakeTransfer(test.senderID, test.receiverID, test.amount)

			if test.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStorage_GetTransferUsersData(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Errorf("Could not mock database connection: %s", err)
	}
	storage := NewStorage(mock)
	dbError := errors.New("Error in database")

	tests := []struct {
		name        string
		senderID    int64
		receiverID  int64
		mock        func()
		expected    *models.TransferUsersData
		expectedErr bool
		err         error
	}{
		{
			name:       "Successfully get data about users in transfer",
			senderID:   1,
			receiverID: 2,
			mock: func() {
				var (
					senderID        int64   = 1
					senderBalance   float64 = 1000
					receiverID      int64   = 2
					receiverBalance float64 = 1000
				)
				mock.ExpectBegin()
				rows := pgxmock.NewRows([]string{"user_id", "balance"})
				rows.AddRow(senderID, senderBalance)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(senderID).WillReturnRows(rows)
				rows = pgxmock.NewRows([]string{"user_id", "balance"})
				rows.AddRow(receiverID, receiverBalance)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(receiverID).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expected: &models.TransferUsersData{
				Sender: &models.UserData{
					UserID:  1,
					Balance: 1000,
				},
				Receiver: &models.UserData{
					UserID:  2,
					Balance: 1000,
				},
			},
		},
		{
			name:       "Sender not found",
			senderID:   1,
			receiverID: 2,
			mock: func() {
				var (
					senderID        int64   = 1
					receiverID      int64   = 2
					receiverBalance float64 = 1000
				)
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(senderID).WillReturnError(pgx.ErrNoRows)
				rows := pgxmock.NewRows([]string{"user_id", "balance"})
				rows.AddRow(receiverID, receiverBalance)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(receiverID).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expected: &models.TransferUsersData{
				Sender: nil,
				Receiver: &models.UserData{
					UserID:  2,
					Balance: 1000,
				},
			},
		},
		{
			name:       "Receiver not found",
			senderID:   1,
			receiverID: 2,
			mock: func() {
				var (
					senderID      int64   = 1
					senderBalance float64 = 1000
					receiverID    int64   = 2
				)
				mock.ExpectBegin()
				rows := pgxmock.NewRows([]string{"user_id", "balance"})
				rows.AddRow(senderID, senderBalance)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(senderID).WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(receiverID).WillReturnError(pgx.ErrNoRows)
				mock.ExpectCommit()
			},
			expected: &models.TransferUsersData{
				Sender: &models.UserData{
					UserID:  1,
					Balance: 1000,
				},
				Receiver: nil,
			},
		},
		{
			name:       "Error in database occurred during getting data about sender",
			senderID:   1,
			receiverID: 2,
			mock: func() {
				var senderID int64 = 1
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(senderID).WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbError,
		},
		{
			name:       "Error in database occurred during getting data about receiver",
			senderID:   1,
			receiverID: 2,
			mock: func() {
				var (
					senderID      int64   = 1
					senderBalance float64 = 1000
					receiverID    int64   = 2
				)
				mock.ExpectBegin()
				rows := pgxmock.NewRows([]string{"user_id", "balance"})
				rows.AddRow(senderID, senderBalance)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(senderID).WillReturnRows(rows)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUser)).WithArgs(receiverID).WillReturnError(dbError)
				mock.ExpectRollback()
			},
			expectedErr: true,
			err:         dbError,
		},
	}

	var got *models.TransferUsersData
	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err = storage.GetTransferUsersData(test.senderID, test.receiverID)

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
