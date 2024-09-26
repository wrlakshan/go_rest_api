package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLConnection creates a new MySQL connection.
func NewMySQLConnection() (*sql.DB, error) {
	dbUser := "wrlakshan"
	dbPass := "wrlakshan"
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "go_rest_api"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

