package main

import (
	deliveryBalance "avito-tech-task/internal/app/balance/delivery"
	repositoryBalance "avito-tech-task/internal/app/balance/repository"
	usecaseBalance "avito-tech-task/internal/app/balance/usecase"
	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"avito-tech-task/config"
	"avito-tech-task/internal/pkg/constants"
	"avito-tech-task/internal/pkg/utils"
)

type Handlers struct {
	BalanceHandlers deliveryBalance.Handlers
}

func NewHandlers(pool *pgx.ConnPool, logger *logrus.Logger, validator *utils.Validation) *Handlers {
	balanceStorage := repositoryBalance.NewStorage(pool)
	balanceService := usecaseBalance.NewService(balanceStorage, validator)
	balanceHandlers := deliveryBalance.NewHandlers(balanceService, logger)

	return &Handlers{
		BalanceHandlers: *balanceHandlers,
	}
}

func main() {
	server := echo.New()

	config := config.NewConfig()
	if _, err := toml.DecodeFile(constants.ConfigPath, &config); err != nil {
		logrus.Fatalf("Could not decode config: %s", err)
	}

	pool := utils.NewPostgresConnection(config)

	logger, fileClose := utils.NewLogger(config)
	defer func(close func() error) {
		if err := close(); err != nil {
			logrus.Fatalf("Could not close file: %s", err)
		}
	}(fileClose)

	validator := utils.NewValidator()

	api := NewHandlers(pool, logger, validator)
	api.BalanceHandlers.InitHandlers(server)

	server.Logger.Fatal(server.Start("0.0.0.0:5000"))
}
