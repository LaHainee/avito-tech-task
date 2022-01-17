package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	storageMock "avito-tech-task/internal/app/balance/mock"
	"avito-tech-task/internal/app/models"
	converterMock "avito-tech-task/internal/pkg/currency/mock"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"avito-tech-task/internal/pkg/utils"
)

func TestService_GetBalance(t *testing.T) {
	storageError := errors.New("Storage error")
	converterError := errors.New("Unsupported currency")

	tests := []struct {
		name          string
		userID        int64
		currency      string
		storageMock   *storageMock.MockStorage
		converterMock *converterMock.MockConverterIface
		expected      *models.UserData
		expectedErr   bool
		err           error
	}{
		{
			name:     "Successfully got user balance",
			userID:   1,
			currency: "USD",
			storageMock: &storageMock.MockStorage{GetUserDataFunc: func(n int64) (*models.UserData, error) {
				return &models.UserData{
					UserID:  1,
					Balance: 1000,
				}, nil
			}},
			converterMock: &converterMock.MockConverterIface{GetFunc: func(s string) (float64, error) {
				return 0.5, nil
			}},
			expected: &models.UserData{
				UserID:  1,
				Balance: 500,
			},
		},
		{
			name:     "Error occurred in storage",
			userID:   1,
			currency: "USD",
			storageMock: &storageMock.MockStorage{GetUserDataFunc: func(n int64) (*models.UserData, error) {
				return nil, storageError
			}},
			expectedErr: true,
			err:         storageError,
		},
		{
			name:     "No user data returned from storage",
			userID:   1,
			currency: "USD",
			storageMock: &storageMock.MockStorage{GetUserDataFunc: func(n int64) (*models.UserData, error) {
				return nil, createdErrors.ErrUserDoesNotExist
			}},
			expectedErr: true,
			err:         createdErrors.ErrUserDoesNotExist,
		},
		{
			name:   "Unsupported currency",
			userID: 1,
			storageMock: &storageMock.MockStorage{GetUserDataFunc: func(n int64) (*models.UserData, error) {
				return &models.UserData{
					UserID:  1,
					Balance: 1000,
				}, nil
			}},
			converterMock: &converterMock.MockConverterIface{GetFunc: func(s string) (float64, error) {
				return 0, converterError
			}},
			expectedErr: true,
			err:         converterError,
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			validator := utils.NewValidator()
			service := NewService(test.storageMock, validator, test.converterMock)

			got, err := service.GetBalance(test.userID, test.currency)

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

func TestService_UpdateBalance(t *testing.T) {
	storageError := errors.New("Error in storage")

	tests := []struct {
		name        string
		data        *models.RequestUpdateBalance
		storageMock *storageMock.MockStorage
		expected    *models.UserData
		expectedErr bool
		err         error
	}{
		{
			name: "Successfully updated balance",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 0,
				Amount:        1000,
			},
			storageMock: &storageMock.MockStorage{
				GetUserDataFunc: func(n int64) (*models.UserData, error) {
					return &models.UserData{
						UserID:  1,
						Balance: 1000,
					}, nil
				},
				UpdateBalanceFunc: func(n int64, f float64) (float64, error) {
					return 2000, nil
				},
			},
			expected: &models.UserData{
				UserID:  1,
				Balance: 2000,
			},
		},
		{
			name: "Error in storage, GetUserData",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 0,
				Amount:        1000,
			},
			storageMock: &storageMock.MockStorage{
				GetUserDataFunc: func(n int64) (*models.UserData, error) {
					return nil, storageError
				},
			},
			expectedErr: true,
			err:         storageError,
		},
		{
			name: "Error in storage, CreateAccount",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 0,
				Amount:        1000,
			},
			storageMock: &storageMock.MockStorage{
				GetUserDataFunc: func(n int64) (*models.UserData, error) {
					return nil, createdErrors.ErrUserDoesNotExist
				},
				CreateAccountFunc: func(n int64) error {
					return storageError
				},
			},
			expectedErr: true,
			err:         storageError,
		},
		{
			name: "Error in storage, UpdateBalance",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 1,
				Amount:        1000,
			},
			storageMock: &storageMock.MockStorage{
				GetUserDataFunc: func(n int64) (*models.UserData, error) {
					return &models.UserData{
						UserID:  1,
						Balance: 1500,
					}, nil
				},
				UpdateBalanceFunc: func(n int64, f float64) (float64, error) {
					return 0, storageError
				},
			},
			expectedErr: true,
			err:         storageError,
		},
		{
			name: "Write off money from created account with zero balance",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 1,
				Amount:        1000,
			},
			storageMock: &storageMock.MockStorage{
				GetUserDataFunc: func(n int64) (*models.UserData, error) {
					return nil, createdErrors.ErrUserDoesNotExist
				},
				CreateAccountFunc: func(n int64) error {
					return nil
				},
			},
			expectedErr: true,
			err:         createdErrors.ErrNotEnoughMoney,
		},
		{
			name: "Not enough money to write off",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 1,
				Amount:        1000,
			},
			storageMock: &storageMock.MockStorage{
				GetUserDataFunc: func(n int64) (*models.UserData, error) {
					return &models.UserData{
						UserID:  1,
						Balance: 500,
					}, nil
				},
			},
			expectedErr: true,
			err:         createdErrors.ErrNotEnoughMoney,
		},
		{
			name: "Negative user ID",
			data: &models.RequestUpdateBalance{
				UserID:        -1,
				OperationType: 0,
				Amount:        1000,
			},
			expectedErr: true,
			err:         createdErrors.ErrNegativeUserID,
		},
		{
			name: "Not supported operation type",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 10,
				Amount:        1000,
			},
			expectedErr: true,
			err:         createdErrors.ErrNotSupportedOperationType,
		},
		{
			name: "Field amount is empty",
			data: &models.RequestUpdateBalance{
				UserID:        1,
				OperationType: 1,
			},
			expectedErr: true,
			err:         createdErrors.ErrAmountFiledIsRequired,
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			validator := utils.NewValidator()
			service := NewService(test.storageMock, validator, nil)

			got, err := service.UpdateBalance(test.data)

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

func TestService_MakeTransfer(t *testing.T) {
	storageError := errors.New("Error in storage")

	tests := []struct {
		name        string
		data        *models.TransferRequest
		storageMock *storageMock.MockStorage
		expected    *models.TransferUsersData
		expectedErr bool
		err         error
	}{
		{
			name: "Successfully transferred money",
			data: &models.TransferRequest{
				SenderID:   1,
				ReceiverID: 2,
				Amount:     500,
			},
			storageMock: &storageMock.MockStorage{
				GetTransferUsersDataFunc: func(n1 int64, n2 int64) (*models.TransferUsersData, error) {
					return &models.TransferUsersData{
						Sender: &models.UserData{
							UserID:  1,
							Balance: 1000,
						},
						Receiver: &models.UserData{
							UserID:  2,
							Balance: 1000,
						},
					}, nil
				},
				MakeTransferFunc: func(n1 int64, n2 int64, f float64) error {
					return nil
				},
			},
			expected: &models.TransferUsersData{
				Sender: &models.UserData{
					UserID:  1,
					Balance: 500,
				},
				Receiver: &models.UserData{
					UserID:  2,
					Balance: 1500,
				},
			},
		},
		{
			name: "Error in storage, GetTransferUsersData",
			data: &models.TransferRequest{
				SenderID:   1,
				ReceiverID: 2,
				Amount:     500,
			},
			storageMock: &storageMock.MockStorage{
				GetTransferUsersDataFunc: func(n1 int64, n2 int64) (*models.TransferUsersData, error) {
					return nil, storageError
				},
			},
			expectedErr: true,
			err:         storageError,
		},
		{
			name: "Error in storage, MakeTransfer",
			data: &models.TransferRequest{
				SenderID:   1,
				ReceiverID: 2,
				Amount:     500,
			},
			storageMock: &storageMock.MockStorage{
				GetTransferUsersDataFunc: func(n1 int64, n2 int64) (*models.TransferUsersData, error) {
					return &models.TransferUsersData{
						Sender: &models.UserData{
							UserID:  1,
							Balance: 1000,
						},
						Receiver: &models.UserData{
							UserID:  2,
							Balance: 1000,
						},
					}, nil
				},
				MakeTransferFunc: func(n1 int64, n2 int64, f float64) error {
					return storageError
				},
			},
			expectedErr: true,
			err:         storageError,
		},
		{
			name: "Sender does not exist",
			data: &models.TransferRequest{
				SenderID:   1,
				ReceiverID: 2,
				Amount:     500,
			},
			storageMock: &storageMock.MockStorage{
				GetTransferUsersDataFunc: func(n1 int64, n2 int64) (*models.TransferUsersData, error) {
					return &models.TransferUsersData{
						Sender: nil,
						Receiver: &models.UserData{
							UserID:  2,
							Balance: 1000,
						},
					}, nil
				},
			},
			expectedErr: true,
			err:         createdErrors.ErrSenderDoesNotExist,
		},
		{
			name: "Receiver does not exist",
			data: &models.TransferRequest{
				SenderID:   1,
				ReceiverID: 2,
				Amount:     500,
			},
			storageMock: &storageMock.MockStorage{
				GetTransferUsersDataFunc: func(n1 int64, n2 int64) (*models.TransferUsersData, error) {
					return &models.TransferUsersData{
						Sender: &models.UserData{
							UserID:  1,
							Balance: 1000,
						},
						Receiver: nil,
					}, nil
				},
			},
			expectedErr: true,
			err:         createdErrors.ErrReceiverDoesNotExist,
		},
		{
			name: "Not enough money",
			data: &models.TransferRequest{
				SenderID:   1,
				ReceiverID: 2,
				Amount:     2000,
			},
			storageMock: &storageMock.MockStorage{
				GetTransferUsersDataFunc: func(n1 int64, n2 int64) (*models.TransferUsersData, error) {
					return &models.TransferUsersData{
						Sender: &models.UserData{
							UserID:  1,
							Balance: 1000,
						},
						Receiver: &models.UserData{
							UserID:  2,
							Balance: 1000,
						},
					}, nil
				},
			},
			expectedErr: true,
			err:         createdErrors.ErrNotEnoughMoney,
		},
		{
			name: "Sender ID field is required",
			data: &models.TransferRequest{
				ReceiverID: 2,
				Amount:     2000,
			},
			storageMock: &storageMock.MockStorage{},
			expectedErr: true,
			err:         createdErrors.ErrSenderIDisRequired,
		},
		{
			name: "Receiver ID field is required",
			data: &models.TransferRequest{
				SenderID: 1,
				Amount:   2000,
			},
			storageMock: &storageMock.MockStorage{},
			expectedErr: true,
			err:         createdErrors.ErrReceiverIDisRequired,
		},
		{
			name: "Amount field is required",
			data: &models.TransferRequest{
				SenderID:   1,
				ReceiverID: 2,
			},
			storageMock: &storageMock.MockStorage{},
			expectedErr: true,
			err:         createdErrors.ErrAmountFiledIsRequired,
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			validator := utils.NewValidator()
			service := NewService(test.storageMock, validator, nil)

			got, err := service.MakeTransfer(test.data)

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
