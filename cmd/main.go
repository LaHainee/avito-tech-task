package main

import (
	"avito-tech-task/internal/pkg/currency"
	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"

	"avito-tech-task/config"
	deliveryBalance "avito-tech-task/internal/app/balance/delivery"
	repositoryBalance "avito-tech-task/internal/app/balance/repository"
	usecaseBalance "avito-tech-task/internal/app/balance/usecase"
	deliveryTransactions "avito-tech-task/internal/app/transactions/delivery"
	repositoryTransactions "avito-tech-task/internal/app/transactions/repository"
	usecaseTransactions "avito-tech-task/internal/app/transactions/usecase"
	"avito-tech-task/internal/pkg/constants"
	"avito-tech-task/internal/pkg/utils"
)

type Handlers struct {
	BalanceHandlers      deliveryBalance.Handlers
	TransactionsHandlers deliveryTransactions.Handlers
}

func NewHandlers(pool *pgx.ConnPool, logger *logrus.Logger, validator *utils.Validation, converter *currency.Converter) *Handlers {
	balanceStorage := repositoryBalance.NewStorage(pool)
	balanceService := usecaseBalance.NewService(balanceStorage, validator, converter)
	balanceHandlers := deliveryBalance.NewHandlers(balanceService, logger)

	transactionsStorage := repositoryTransactions.NewStorage(pool)
	transactionsService := usecaseTransactions.NewService(transactionsStorage)
	transactionsHandlers := deliveryTransactions.NewHandlers(transactionsService, logger)

	return &Handlers{
		BalanceHandlers:      *balanceHandlers,
		TransactionsHandlers: *transactionsHandlers,
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

	converter := currency.NewConverter(config, logger)

	api := NewHandlers(pool, logger, validator, converter)
	api.BalanceHandlers.InitHandlers(server)
	api.TransactionsHandlers.InitHandlers(server)

	go func() {
		server.Logger.Fatal(server.Start("0.0.0.0:5000"))
	}()

	cancel := make(chan struct{})
	go currency.UpdateCurrency(converter, cancel)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(done, os.Kill)
	log.Println("Graceful shutdown")
	<-done
	cancel <- struct{}{}
}
