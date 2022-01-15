package transactions

import "avito-tech-task/internal/app/models"

type Service interface {
	GetUserTransactions(int64, *models.TransactionsSelectionParams) (models.Transactions, error)
}
