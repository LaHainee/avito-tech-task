package utils

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"

	"avito-tech-task/config"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close(context.Context) error
}

func NewPostgresConnection(config *config.Config) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config.Server.DatabaseConnString)
	if err != nil {
		logrus.Fatalf("Could not establish connection to database: %s", err)
	}

	return conn
}
