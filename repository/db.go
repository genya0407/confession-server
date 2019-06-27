package repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/doug-martin/goqu"
	_ "github.com/doug-martin/goqu/adapters/postgres"
	_ "github.com/lib/pq"
)

func GetDB() *goqu.Database {
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		user = "confession"
	}

	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		dbname = "confession-development"
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		password = "confession"
	}

	pgDb, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable ", user, dbname, password))
	if err != nil {
		panic(err.Error())
	}

	return goqu.New("postgres", pgDb)
}
