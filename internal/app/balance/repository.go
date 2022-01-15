package balance

import "avito-tech-task/internal/app/models"

type Storage interface {
	UpdateBalance(int64, float64) (float64, error)
	GetUserData(int64) (*models.UserData, error)
	CreateAccount(int64) error
	MakeTransfer(int64, int64, float64) error
}
