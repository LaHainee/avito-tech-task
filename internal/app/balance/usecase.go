package balance

import "avito-tech-task/internal/app/models"

type Service interface {
	UpdateBalance(*models.RequestUpdateBalance) (*models.UserData, error)
	GetUserData(int64) (*models.UserData, error)
}
