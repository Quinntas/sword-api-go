package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/quinntas/go-fiber-template/database/repository"
	"log"
	"os"
)

var (
	Repo *repository.Queries
)

func SetupDatabase() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	Repo = repository.New(conn)
}
