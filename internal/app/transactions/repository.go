package transactions

import "avito-tech-task/internal/app/models"

type Storage interface {
	DoesUserExist(int64) (bool, error)
	GetUserTransactions(int64, *models.TransactionsSelectionParams) (models.Transactions, error)
}