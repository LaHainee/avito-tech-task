package usecase

import (
	"avito-tech-task/internal/pkg/utils"
	"strings"

	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/app/transactions"
	createdErrors "avito-tech-task/internal/pkg/errors"
)

type Service struct {
	storage   transactions.Storage
	validator *utils.Validation
}

func NewService(storage transactions.Storage, validator *utils.Validation) *Service {
	return &Service{
		storage:   storage,
		validator: validator,
	}
}

func (s *Service) GetUserTransactions(userID int64, params *models.TransactionsSelectionParams) (models.Transactions, error) {
	errs := s.validator.Validate(params) // validation
	for _, err := range errs {
		if err.Field() == "Limit" {
			return nil, createdErrors.ErrNegativeLimit
		}
	}

	doesUserExist, err := s.storage.DoesUserExist(userID)
	if err != nil {
		return nil, err
	}
	if !doesUserExist {
		return nil, createdErrors.ErrUserDoesNotExist
	}

	if params.Since != "" {
		// without this parsing we get timestamp like that: 2022-01-15T21:37:23.822151 03:00
		// but database will process only this format: 2022-01-15T21:37:23.822151 +03:00
		// so we need to add "+"
		params.Since = strings.Join(strings.Split(params.Since, " "), " +")
	}

	return s.storage.GetUserTransactions(userID, params)
}
