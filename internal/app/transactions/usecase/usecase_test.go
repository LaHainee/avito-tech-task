package usecase

import (
	"avito-tech-task/internal/app/models"
	storageMock "avito-tech-task/internal/app/transactions/mock"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_GetUserTransactions(t *testing.T) {
	storageError := errors.New("Error in storage")

	timeNow := time.Now()
	tests := []struct {
		name        string
		userID      int64
		params      *models.TransactionsSelectionParams
		storageMock *storageMock.MockStorage
		expected    models.Transactions
		expectedErr bool
		err         error
	}{
		{
			name:   "Successfully get transactions list",
			userID: 1,
			params: &models.TransactionsSelectionParams{
				Limit:       0,
				Since:       "some datetime",
				OrderAmount: false,
				OrderDate:   false,
			},
			storageMock: &storageMock.MockStorage{
				DoesUserExistFunc: func(n int64) (bool, error) {
					return true, nil
				},
				GetUserTransactionsFunc: func(n int64, transactionsSelectionParams *models.TransactionsSelectionParams) (models.Transactions, error) {
					return models.Transactions{
						&models.Transaction{
							Description: "description",
							Amount:      1000,
							Created:     timeNow,
						},
					}, nil
				},
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
			name:   "User does not exist",
			userID: 1,
			params: &models.TransactionsSelectionParams{
				Limit:       0,
				Since:       "",
				OrderAmount: false,
				OrderDate:   false,
			},
			storageMock: &storageMock.MockStorage{
				DoesUserExistFunc: func(n int64) (bool, error) {
					return false, nil
				},
			},
			expectedErr: true,
			err:         createdErrors.ErrUserDoesNotExist,
		},
		{
			name:   "Error occurred in storage during checking user existence",
			userID: 1,
			params: &models.TransactionsSelectionParams{
				Limit:       0,
				Since:       "",
				OrderAmount: false,
				OrderDate:   false,
			},
			storageMock: &storageMock.MockStorage{
				DoesUserExistFunc: func(n int64) (bool, error) {
					return false, storageError
				},
			},
			expectedErr: true,
			err:         storageError,
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			service := NewService(test.storageMock)

			got, err := service.GetUserTransactions(test.userID, test.params)

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
