package db

import (
	"News_site/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DB struct {
	db *sql.DB
}

func New(config *config.Config) (*DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	mysqlDB, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL: %v", err)
	}

	if err := mysqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging MySQL: %v", err)
	}

	log.Println("Connected to MySQL")
	return &DB{db: mysqlDB}, nil
}

func (mysql *DB) Close() error {
	return mysql.db.Close()
}

func (mysql *DB) GetDB() *sql.DB {
	return mysql.db
}
