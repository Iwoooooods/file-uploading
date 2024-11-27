package database

import (
	"database/sql"
	"log"

	"github.com/Iwoooooods/fs-upload-go/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DSN)
	if err != nil {
		log.Printf("failed to open database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("failed to ping database: %v", err)
		return nil, err
	}
	log.Printf("successfully connected to database: %v", cfg.DSN)
	return db, nil
}
