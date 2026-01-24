package resource

import (
	"database/sql"
	"fmt"
	"os"
)

func NewPostgreSqlDb() (*sql.DB, error) {
	host := os.Getenv("LOCAL_DB_HOST")
	port := os.Getenv("LOCAL_DB_PORT")
	user := os.Getenv("LOCAL_DB_USER")
	password := os.Getenv("LOCAL_DB_PASSWORD")
	dbname := os.Getenv("LOCAL_DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	return sql.Open("postgres", connStr)
}
