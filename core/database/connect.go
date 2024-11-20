package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

const MigrationDir = "./migrations/sqlite"

type Sqlite struct {
	dbName string
	dsn    string
	conn   *sqlx.DB

	mu *sync.RWMutex
}

func ConnectSqlite(dbName string) *Sqlite {
	return &Sqlite{
		dbName: dbName,
		dsn:    "sqlite3",
		mu:     &sync.RWMutex{},
	}
}

func (s *Sqlite) Connect(ctx context.Context) error {
	db, err := sqlx.Connect(s.dsn, s.dbName)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to sqlite")
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.conn = db
	if err := s.conn.PingContext(ctx); err != nil {
		log.Error().Err(err).Msg("failed to ping sqlite")
		return err
	}

	return nil
}

func (s *Sqlite) Setup(ctx context.Context) error {
	var err error

	schemaFile := fmt.Sprintf("%s/schema.sql", MigrationDir)

	data, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to read schema file")
		return err
	}

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}

	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("failed to setup sqlite")
			tx.Rollback()
		}
	}()

	stmts := strings.Split(string(data), "---")

	for _, stmt := range stmts {
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		_, err = tx.ExecContext(ctx, stmt)
		if err != nil {
			log.Error().Err(err).Msg("failed to execute statement")
			return err
		}
	}

	return tx.Commit()
}
