package store

import (
	"database/sql"
	"fmt"

	"github.com/NASA_API/internal/config"
	_ "github.com/lib/pq"
)

type MetaDataStore struct {
	db *sql.DB
}

func NewMetaDataStore(cfg config.PsqlConfig) (*MetaDataStore, error) {
	connStr := fmt.Sprintf(
		"host =%s port =%s user =%s password =%s dbname =%s sslmode=allow",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &MetaDataStore{db: db}, nil
}
