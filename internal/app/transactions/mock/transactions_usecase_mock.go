// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/app/transactions"
	"sync"
)

// Ensure, that MockService does implement transactions.Service.
// If this is not the case, regenerate this file with moq.
var _ transactions.Service = &MockService{}

// MockService is a mock implementation of transactions.Service.
//
// 	func TestSomethingThatUsesService(t *testing.T) {
//
// 		// make and configure a mocked transactions.Service
// 		mockedService := &MockService{
// 			GetUserTransactionsFunc: func(n int64, transactionsSelectionParams *models.TransactionsSelectionParams) (models.Transactions, error) {
// 				panic("mock out the GetUserTransactions method")
// 			},
// 		}
//
// 		// use mockedService in code that requires transactions.Service
// 		// and then make assertions.
//
// 	}
type MockService struct {
	// GetUserTransactionsFunc mocks the GetUserTransactions method.
	GetUserTransactionsFunc func(n int64, transactionsSelectionParams *models.TransactionsSelectionParams) (models.Transactions, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetUserTransactions holds details about calls to the GetUserTransactions method.
		GetUserTransactions []struct {
			// N is the n argument value.
			N int64
			// TransactionsSelectionParams is the transactionsSelectionParams argument value.
			TransactionsSelectionParams *models.TransactionsSelectionParams
		}
	}
	lockGetUserTransactions sync.RWMutex
}

// GetUserTransactions calls GetUserTransactionsFunc.
func (mock *MockService) GetUserTransactions(n int64, transactionsSelectionParams *models.TransactionsSelectionParams) (models.Transactions, error) {
	if mock.GetUserTransactionsFunc == nil {
		panic("MockService.GetUserTransactionsFunc: method is nil but Service.GetUserTransactions was just called")
	}
	callInfo := struct {
		N                           int64
		TransactionsSelectionParams *models.TransactionsSelectionParams
	}{
		N:                           n,
		TransactionsSelectionParams: transactionsSelectionParams,
	}
	mock.lockGetUserTransactions.Lock()
	mock.calls.GetUserTransactions = append(mock.calls.GetUserTransactions, callInfo)
	mock.lockGetUserTransactions.Unlock()
	return mock.GetUserTransactionsFunc(n, transactionsSelectionParams)
}

// GetUserTransactionsCalls gets all the calls that were made to GetUserTransactions.
// Check the length with:
//     len(mockedService.GetUserTransactionsCalls())
func (mock *MockService) GetUserTransactionsCalls() []struct {
	N                           int64
	TransactionsSelectionParams *models.TransactionsSelectionParams
} {
	var calls []struct {
		N                           int64
		TransactionsSelectionParams *models.TransactionsSelectionParams
	}
	mock.lockGetUserTransactions.RLock()
	calls = mock.calls.GetUserTransactions
	mock.lockGetUserTransactions.RUnlock()
	return calls
}
