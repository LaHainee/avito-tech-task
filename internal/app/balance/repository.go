package balance

import "avito-tech-task/internal/app/models"

type Storage interface {
	UpdateBalance(int64, float64) (*models.UserData, error)
	GetUserData(int64) (*models.UserData, error)
	CreateAccount(int64) error
}
