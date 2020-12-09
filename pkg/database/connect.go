package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(filePath string) (*sql.DB, error) {
	// open or create db file...
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil, fmt.Errorf("create file (%s): %w", filePath, err)
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close file (%s): %w", filePath, err)
	}

	// open...
	db, err := sql.Open("sqlite3", filePath)

	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	return db, nil
}
