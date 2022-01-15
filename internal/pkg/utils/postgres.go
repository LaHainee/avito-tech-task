package utils

import (
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"

	"avito-tech-task/config"
)

func NewPostgresConnection(config *config.Config) *pgx.ConnPool {
	pgxConnectionConfig, err := pgx.ParseConnectionString(config.Server.DatabaseConnString)
	if err != nil {
		logrus.Fatalf("Could not parse connection string %s: %s", config.Server.DatabaseConnString, err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxConnectionConfig,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		logrus.Fatalf("Could not establish connection to database: %s", err)
	}

	return pool
}
