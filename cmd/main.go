package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"avito-tech-task/config"
	deliveryBalance "avito-tech-task/internal/app/balance/delivery"
	repositoryBalance "avito-tech-task/internal/app/balance/repository"
	usecaseBalance "avito-tech-task/internal/app/balance/usecase"
	deliveryTransactions "avito-tech-task/internal/app/transactions/delivery"
	repositoryTransactions "avito-tech-task/internal/app/transactions/repository"
	usecaseTransactions "avito-tech-task/internal/app/transactions/usecase"
	"avito-tech-task/internal/pkg/constants"
	"avito-tech-task/internal/pkg/currency"
	"avito-tech-task/internal/pkg/utils"
)

type Handlers struct {
	BalanceHandlers      deliveryBalance.Handlers
	TransactionsHandlers deliveryTransactions.Handlers
}

func NewHandlers(conn utils.PgxIface, logger *logrus.Logger, validator *utils.Validation, converter *currency.Converter) *Handlers {
	balanceStorage := repositoryBalance.NewStorage(conn)
	balanceService := usecaseBalance.NewService(balanceStorage, validator, converter)
	balanceHandlers := deliveryBalance.NewHandlers(balanceService, logger)

	transactionsStorage := repositoryTransactions.NewStorage(conn)
	transactionsService := usecaseTransactions.NewService(transactionsStorage)
	transactionsHandlers := deliveryTransactions.NewHandlers(transactionsService, logger)

	return &Handlers{
		BalanceHandlers:      *balanceHandlers,
		TransactionsHandlers: *transactionsHandlers,
	}
}

// @title        BalanceApplication
// @version      1.0
// @description  API for BalanceApplication

// @license.name  ""

// @BasePath  /api/v1

// @x-extension-openapi  {"example": "value on a json format"}

func main() {
	server := echo.New()

	config := config.NewConfig()
	if _, err := toml.DecodeFile(constants.ConfigPath, &config); err != nil {
		logrus.Fatalf("Could not decode config: %s", err)
	}

	conn := utils.NewPostgresConnection(config)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			logrus.Fatalf("Could not close database connection: %s", err)
		}
	}(conn, context.Background())

	logger, closeF := utils.NewLogger(config)
	defer func(closeF func() error) {
		if err := closeF(); err != nil {
			logrus.Fatalf("Could not close file: %s", err)
		}
	}(closeF)

	validator := utils.NewValidator()

	converter := currency.NewConverter(config, logger)

	api := NewHandlers(conn, logger, validator, converter)
	api.BalanceHandlers.InitHandlers(server)
	api.TransactionsHandlers.InitHandlers(server)

	go func() {
		server.Logger.Fatal(server.Start("0.0.0.0:5000"))
	}()

	cancel := make(chan struct{})
	go currency.UpdateCurrency(converter, cancel)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	cancel <- struct{}{}
}
