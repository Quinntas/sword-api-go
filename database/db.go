package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/quinntas/go-fiber-template/database/repository"
	"log"
	"os"
)

var (
	Repo *repository.Queries
)

func SetupDatabase() {
	databaseUrl := os.Getenv("DATABASE_URL")
	conn, err := sql.Open("mysql", databaseUrl+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	Repo = repository.New(conn)
}
