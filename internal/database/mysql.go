package database;

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka database MySQL: %w", err)
	}
	return db, nil
}


