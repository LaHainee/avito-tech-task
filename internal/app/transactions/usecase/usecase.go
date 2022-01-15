package usecase

import (
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/app/transactions"
	createdErrors "avito-tech-task/internal/pkg/errors"
)

type Service struct {
	storage transactions.Storage
}

func NewService(storage transactions.Storage) *Service {
	return &Service{storage}
}

func (s *Service) GetUserTransactions(id int64, params *models.TransactionsSelectionParams) (models.Transactions, error) {
	doesUserExist, err := s.storage.DoesUserExist(id)
	if err != nil {
		return nil, err
	}
	if !doesUserExist {
		return nil, createdErrors.ErrUserDoesNotExist
	}

	return s.storage.GetUserTransactions(id, params)
}
