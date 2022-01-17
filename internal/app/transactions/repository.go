package transactions

import "avito-tech-task/internal/app/models"

//go:generate moq -out ./mock/transactions_repo_mock.go -pkg mock . Storage:MockStorage
type Storage interface {
	DoesUserExist(int64) (bool, error)
	GetUserTransactions(int64, *models.TransactionsSelectionParams) (models.Transactions, error)
}
