package balance

import "avito-tech-task/internal/app/models"

type Service interface {
	GetBalance(int64, string) (*models.UserData, error)
	MakeTransfer(*models.TransferRequest) (*models.TransferResponse, error)
	UpdateBalance(*models.RequestUpdateBalance) (*models.UserData, error)
}
