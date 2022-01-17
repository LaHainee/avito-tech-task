package balance

import "avito-tech-task/internal/app/models"

//go:generate moq -out ./mock/balance_usecase_mock.go -pkg mock . Service:MockService
type Service interface {
	GetBalance(int64, string) (*models.UserData, error)
	MakeTransfer(*models.TransferRequest) (*models.TransferUsersData, error)
	UpdateBalance(*models.RequestUpdateBalance) (*models.UserData, error)
}
