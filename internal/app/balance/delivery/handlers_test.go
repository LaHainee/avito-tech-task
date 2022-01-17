package delivery

import (
	"avito-tech-task/config"
	"avito-tech-task/internal/app/balance/mock"
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/constants"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"avito-tech-task/internal/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandlers_GetBalance(t *testing.T) {
	const removeLogs = true // set false to deny deleting logs after test

	config := &config.Config{
		LoggingLevel:    "debug",
		LoggingFilePath: "./logs/",
		CurrencyApiURL:  "",
		Server:          config.ServerConfig{},
	}
	logger, fileClose := utils.NewLogger(config)
	defer func(close func() error) {
		if err := close(); err != nil {
			t.Errorf("Could not close file: %s", err)
		}
	}(fileClose)
	internalServerErr := errors.New("Internal server error")

	if removeLogs {
		defer func() {
			if err := os.RemoveAll("./logs/"); err != nil {
				t.Errorf("Could not remove temporary logs directory: %s", err)
			}
		}()
	}

	tests := []struct {
		name           string
		serviceMock    *mock.MockService
		userIDParam    string
		expectedStatus int
		expected       interface{}
	}{
		{
			name: "Successfully get user balance",
			serviceMock: &mock.MockService{
				GetBalanceFunc: func(n int64, s string) (*models.UserData, error) {
					return &models.UserData{
						UserID:  1,
						Balance: 1000,
					}, nil
				},
			},
			userIDParam:    "1",
			expectedStatus: http.StatusOK,
			expected: &models.UserData{
				UserID:  1,
				Balance: 1000,
			},
		},
		{
			name:           "Invalid user id in param",
			userIDParam:    "string???",
			expectedStatus: http.StatusBadRequest,
			expected:       &models.ResponseMessage{Message: constants.InvalidUserIDMessage},
		},
		{
			name: "Not supported currency",
			serviceMock: &mock.MockService{
				GetBalanceFunc: func(n int64, s string) (*models.UserData, error) {
					return nil, createdErrors.ErrNotSupportedCurrency
				},
			},
			userIDParam:    "1",
			expectedStatus: http.StatusUnprocessableEntity,
			expected:       &models.ResponseMessage{Message: createdErrors.ErrNotSupportedCurrency.Error()},
		},
		{
			name: "User does not exist",
			serviceMock: &mock.MockService{
				GetBalanceFunc: func(n int64, s string) (*models.UserData, error) {
					return nil, createdErrors.ErrUserDoesNotExist
				},
			},
			userIDParam:    "1",
			expectedStatus: http.StatusNotFound,
			expected:       &models.ResponseMessage{Message: createdErrors.ErrUserDoesNotExist.Error()},
		},
		{
			name: "Internal server error",
			serviceMock: &mock.MockService{
				GetBalanceFunc: func(n int64, s string) (*models.UserData, error) {
					return nil, internalServerErr
				},
			},
			userIDParam:    "1",
			expectedStatus: http.StatusInternalServerError,
			expected:       &models.ResponseMessage{Message: internalServerErr.Error()},
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/balance/:user_id")
			ctx.SetParamNames("user_id")
			ctx.SetParamValues(test.userIDParam)

			handlers := NewHandlers(test.serviceMock, logger)
			if assert.NoError(t, handlers.GetBalance(ctx)) {
				assert.Equal(t, test.expectedStatus, rec.Code)

				expectedString, _ := json.Marshal(test.expected)
				assert.Equal(t, string(expectedString)+"\n", rec.Body.String())
			}
		})
	}
}

func TestHandlers_UpdateBalance(t *testing.T) {
	const removeLogs = true // set false to deny deleting logs after test

	config := &config.Config{
		LoggingLevel:    "debug",
		LoggingFilePath: "./logs/",
		CurrencyApiURL:  "",
		Server:          config.ServerConfig{},
	}
	logger, fileClose := utils.NewLogger(config)
	defer func(close func() error) {
		if err := close(); err != nil {
			t.Errorf("Could not close file: %s", err)
		}
	}(fileClose)
	internalServerErr := errors.New("Internal server error")

	if removeLogs {
		defer func() {
			if err := os.RemoveAll("./logs/"); err != nil {
				t.Errorf("Could not remove temporary logs directory: %s", err)
			}
		}()
	}

	tests := []struct {
		name           string
		serviceMock    *mock.MockService
		userIDParam    string
		body           string
		expectedStatus int
		expected       interface{}
	}{
		{
			name: "Successfully updated user balance",
			serviceMock: &mock.MockService{
				UpdateBalanceFunc: func(requestUpdateBalance *models.RequestUpdateBalance) (*models.UserData, error) {
					return &models.UserData{
						UserID:  1,
						Balance: 1000,
					}, nil
				},
			},
			userIDParam:    "1",
			body:           `{"operation_type": 0, "amount": 1000}`,
			expectedStatus: http.StatusOK,
			expected: &models.UserData{
				UserID:  1,
				Balance: 1000,
			},
		},
		{
			name:           "Invalid body",
			userIDParam:    "1",
			body:           `{"operation_type": 0, "amount": "hello"}`,
			expectedStatus: http.StatusBadRequest,
			expected:       &models.ResponseMessage{Message: constants.InvalidBodyMessage},
		},
		{
			name: "Not enough money | Not supported operation type | Amount field was not set | Negative user ID",
			serviceMock: &mock.MockService{
				UpdateBalanceFunc: func(requestUpdateBalance *models.RequestUpdateBalance) (*models.UserData, error) {
					return nil, createdErrors.ErrNotEnoughMoney
				},
			},
			userIDParam:    "1",
			body:           `{"operation_type": 1, "amount": 1000}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expected:       &models.ResponseMessage{Message: createdErrors.ErrNotEnoughMoney.Error()},
		},
		{
			name: "Internal server error",
			serviceMock: &mock.MockService{
				UpdateBalanceFunc: func(requestUpdateBalance *models.RequestUpdateBalance) (*models.UserData, error) {
					return nil, internalServerErr
				},
			},
			userIDParam:    "1",
			body:           `{"operation_type": 1, "amount": 1000}`,
			expectedStatus: http.StatusInternalServerError,
			expected:       &models.ResponseMessage{Message: internalServerErr.Error()},
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			server := echo.New()

			req := httptest.NewRequest(echo.POST, "/", strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/balance/:user_id")
			ctx.SetParamNames("user_id")
			ctx.SetParamValues(test.userIDParam)

			handlers := NewHandlers(test.serviceMock, logger)
			if assert.NoError(t, handlers.UpdateBalance(ctx)) {
				assert.Equal(t, test.expectedStatus, rec.Code)

				expectedString, _ := json.Marshal(test.expected)
				assert.Equal(t, string(expectedString)+"\n", rec.Body.String())
			}
		})
	}
}

func TestHandlers_Transfer(t *testing.T) {
	const removeLogs = true // set false to deny deleting logs after test

	config := &config.Config{
		LoggingLevel:    "debug",
		LoggingFilePath: "./logs/",
		CurrencyApiURL:  "",
		Server:          config.ServerConfig{},
	}
	logger, fileClose := utils.NewLogger(config)
	defer func(close func() error) {
		if err := close(); err != nil {
			t.Errorf("Could not close file: %s", err)
		}
	}(fileClose)
	internalServerErr := errors.New("Internal server error")

	if removeLogs {
		defer func() {
			if err := os.RemoveAll("./logs/"); err != nil {
				t.Errorf("Could not remove temporary logs directory: %s", err)
			}
		}()
	}

	tests := []struct {
		name           string
		serviceMock    *mock.MockService
		body           string
		expectedStatus int
		expected       interface{}
	}{
		{
			name: "Successfully transferred money",
			serviceMock: &mock.MockService{
				MakeTransferFunc: func(transferRequest *models.TransferRequest) (*models.TransferUsersData, error) {
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
			body:           `{"sender_id": 1, "receiver_id": 2, "amount": 1000}`,
			expectedStatus: http.StatusOK,
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
			name: "Not enough money to make transfer",
			serviceMock: &mock.MockService{
				MakeTransferFunc: func(transferRequest *models.TransferRequest) (*models.TransferUsersData, error) {
					return nil, createdErrors.ErrNotEnoughMoney
				},
			},
			body:           `{"sender_id": 1, "receiver_id": 2, "amount": 1000}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expected:       &models.ResponseMessage{Message: createdErrors.ErrNotEnoughMoney.Error()},
		},
		{
			name: "Sender not found | Receiver not found",
			serviceMock: &mock.MockService{
				MakeTransferFunc: func(transferRequest *models.TransferRequest) (*models.TransferUsersData, error) {
					return nil, createdErrors.ErrSenderDoesNotExist
				},
			},
			body:           `{"sender_id": 1, "receiver_id": 2, "amount": 1000}`,
			expectedStatus: http.StatusNotFound,
			expected:       &models.ResponseMessage{Message: createdErrors.ErrSenderDoesNotExist.Error()},
		},
		{
			name:           "Invalid body",
			body:           `{"sender_id": 1, "receiver_id": 2, "amount": "hello"}`,
			expectedStatus: http.StatusBadRequest,
			expected:       &models.ResponseMessage{Message: constants.InvalidBodyMessage},
		},
		{
			name: "Internal server error",
			serviceMock: &mock.MockService{
				MakeTransferFunc: func(transferRequest *models.TransferRequest) (*models.TransferUsersData, error) {
					return nil, internalServerErr
				},
			},
			body:           `{"sender_id": 1, "receiver_id": 2, "amount": 1000}`,
			expectedStatus: http.StatusInternalServerError,
			expected:       &models.ResponseMessage{Message: internalServerErr.Error()},
		},
	}

	for _, current := range tests {
		test := current
		t.Run(test.name, func(t *testing.T) {
			server := echo.New()

			req := httptest.NewRequest(echo.POST, "/", strings.NewReader(test.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/transfer")

			handlers := NewHandlers(test.serviceMock, logger)
			if assert.NoError(t, handlers.Transfer(ctx)) {
				assert.Equal(t, test.expectedStatus, rec.Code)

				expectedString, _ := json.Marshal(test.expected)
				assert.Equal(t, string(expectedString)+"\n", rec.Body.String())
			}
		})
	}
}
