package usecase

import (
	"avito-tech-task/internal/app/balance"
	"avito-tech-task/internal/app/models"
	"avito-tech-task/internal/pkg/currency"
	createdErrors "avito-tech-task/internal/pkg/errors"
	"avito-tech-task/internal/pkg/utils"
	"log"
)

type Service struct {
	storage   balance.Storage
	validator *utils.Validation
	converter currency.Service
}

func NewService(storage balance.Storage, validator *utils.Validation, converter *currency.Converter) *Service {
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

func (s *Service) MakeTransfer(data *models.TransferRequest) (*models.TransferResponse, error) {
	errors := s.validator.Validate(data) // validation
	for _, err := range errors {
		log.Println(err)
		switch err.Field() {
		case "SenderID":
			return nil, createdErrors.ErrSenderIDisRequired
		case "ReceiverID":
			return nil, createdErrors.ErrReceiverIDisRequired
		case "Amount":
			return nil, createdErrors.ErrAmountFiledIsRequired
		}
	}

	senderData, err := s.storage.GetUserData(data.SenderID) // get sender data
	if err != nil {
		return nil, err
	}
	if senderData == nil {
		return nil, createdErrors.ErrSenderDoesNotExist
	}

	receiverData, err := s.storage.GetUserData(data.ReceiverID) // get receiver data
	if err != nil {
		return nil, err
	}
	if receiverData == nil {
		return nil, createdErrors.ErrReceiverDoesNotExist
	}

	if senderData.Balance < data.Amount {
		return nil, createdErrors.ErrNotEnoughMoney
	}

	if err = s.storage.MakeTransfer(data.SenderID, data.ReceiverID, data.Amount); err != nil {
		return nil, err
	}

	response := &models.TransferResponse{
		Sender: models.UserData{
			UserID:  data.SenderID,
			Balance: senderData.Balance - data.Amount,
		},
		Receiver: models.UserData{
			UserID:  data.ReceiverID,
			Balance: receiverData.Balance + data.Amount,
		},
	}

	return response, nil
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
