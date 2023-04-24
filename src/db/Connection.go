package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseService interface {
	createConnection()
	GetConnection() *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mypassword"
	dbname   = "postgres"
)

func CreateConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func DbClosedError() error {
	return fmt.Errorf("Database connection closed")
}
