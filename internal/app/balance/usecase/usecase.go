package usecase

import (
	"avito-tech-task/internal/app/balance"
	"avito-tech-task/internal/app/models"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"avito-tech-task/internal/pkg/utils"
	"log"
)

type Service struct {
	storage   balance.Storage
	validator *utils.Validation
}

func NewService(storage balance.Storage, validator *utils.Validation) *Service {
	return &Service{
		storage:   storage,
		validator: validator,
	}
}

func (s *Service) GetUserData(id int64) (*models.UserData, error) {
	userData, err := s.storage.GetUserData(id)
	if err != nil {
		return nil, err
	}
	if userData == nil {
		return nil, createdErrors.ErrUserDoesNotExist
	}
	
	return userData, nil
}

func (s *Service) UpdateBalance(data *models.RequestUpdateBalance) (*models.UserData, error) {
	errors := s.validator.Validate(data) // validation
	for _, err := range errors {
		log.Println(err)
		switch err.Field() {
		case "UserID":
			return nil, createdErrors.ErrNegativeUserID
		case "OperationType":
			return nil, createdErrors.ErrNotSupportedOperationType
		case "Amount":
			return nil, createdErrors.ErrAmountFiledIsRequired
		}
	}

	userData, err := s.storage.GetUserData(data.UserID) // check if user exists
	if err != nil {
		return nil, err
	}

	if userData == nil { // user does not exist
		if err = s.storage.CreateAccount(data.UserID); err != nil { // create account if user does not exist
			return nil, err
		}

		if data.OperationType == models.REDUCE { // trying to write off money from new account (balance = 0)
			return nil, createdErrors.ErrNotEnoughMoney
		}
	}

	if data.OperationType == models.REDUCE {
		if userData.Balance < data.Amount {
			return nil, createdErrors.ErrNotEnoughMoney
		}
		data.Amount *= -1
	}

	updatedUserData, err := s.storage.UpdateBalance(data.UserID, data.Amount)
	if err != nil {
		return nil, err
	}

	return updatedUserData, nil
}
