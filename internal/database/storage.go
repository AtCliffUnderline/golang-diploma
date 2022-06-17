package database

import (
	"context"
	"github.com/golang-migrate/migrate/v4"

	"github.com/AtCliffUnderline/golang-diploma/internal/config"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

type Storage struct {
	DBConn *pgx.Conn
}

func InitStorage(c config.Config) *Storage {
	makeMigration(c.DatabaseURI)
	conn, err := pgx.Connect(context.Background(), c.DatabaseURI)
	if err != nil {
		panic(err)
	}

	return &Storage{
		DBConn: conn,
	}
}

func makeMigration(uri string) {
	m, err := migrate.New(
		"file://internal/database/migrations",
		uri)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
