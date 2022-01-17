package transactions

import "avito-tech-task/internal/app/models"

//go:generate moq -out ./mock/transactions_usecase_mock.go -pkg mock . Service:MockService
type Service interface {
	GetUserTransactions(int64, *models.TransactionsSelectionParams) (models.Transactions, error)
}
