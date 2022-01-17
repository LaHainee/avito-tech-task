package usecase

import (
	"avito-tech-task/internal/app/balance"
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/currency"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"avito-tech-task/internal/pkg/utils"
)

type Service struct {
	validator *utils.Validation
	storage   balance.Storage
	converter currency.ConverterIface
}

func NewService(storage balance.Storage, validator *utils.Validation, converter currency.ConverterIface) *Service {
	return &Service{
		storage:   storage,
		validator: validator,
		converter: converter,
	}
}

func (s *Service) GetBalance(id int64, currency string) (*models.UserData, error) {
	if len(currency) == 0 {
		currency = "RUB"
	}

	userData, err := s.storage.GetUserData(id)
	if err != nil {
		return nil, err
	}
	if userData == nil {
		return nil, createdErrors.ErrUserDoesNotExist
	}

	convertingCoeff, err := s.converter.Get(currency)
	if err != nil {
		return nil, err
	}
	userData.Balance *= convertingCoeff

	return userData, nil
}

func (s *Service) MakeTransfer(data *models.TransferRequest) (*models.TransferUsersData, error) {
	errors := s.validator.Validate(data) // validation
	for _, err := range errors {
		switch err.Field() {
		case "SenderID":
			return nil, createdErrors.ErrSenderIDisRequired
		case "ReceiverID":
			return nil, createdErrors.ErrReceiverIDisRequired
		case "Amount":
			return nil, createdErrors.ErrAmountFiledIsRequired
		}
	}

	transferUsersData, err := s.storage.GetTransferUsersData(data.SenderID, data.ReceiverID)
	if err != nil {
		return nil, err
	}
	if transferUsersData.Sender == nil { // check if sender exists
		return nil, createdErrors.ErrSenderDoesNotExist
	}
	if transferUsersData.Receiver == nil { // check if receiver exists
		return nil, createdErrors.ErrReceiverDoesNotExist
	}

	if transferUsersData.Sender.Balance < data.Amount {
		return nil, createdErrors.ErrNotEnoughMoney
	}

	if err = s.storage.MakeTransfer(data.SenderID, data.ReceiverID, data.Amount); err != nil {
		return nil, err
	}

	transferUsersData.Sender.Balance -= data.Amount
	transferUsersData.Receiver.Balance += data.Amount

	return transferUsersData, nil
}

func (s *Service) UpdateBalance(data *models.RequestUpdateBalance) (*models.UserData, error) {
	errors := s.validator.Validate(data) // validation
	for _, err := range errors {
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

		if data.OperationType == models.REDUCE { // trying to write off money from created account (balance = 0)
			return nil, createdErrors.ErrNotEnoughMoney
		}
	}

	if data.OperationType == models.REDUCE {
		if userData.Balance < data.Amount {
			return nil, createdErrors.ErrNotEnoughMoney
		}
		data.Amount *= -1
	}

	newBalance, err := s.storage.UpdateBalance(data.UserID, data.Amount)
	if err != nil {
		return nil, err
	}

	return &models.UserData{UserID: data.UserID, Balance: newBalance}, nil
}
