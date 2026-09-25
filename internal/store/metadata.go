package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NASA_API/internal/config"
	"github.com/NASA_API/internal/entity"
	_ "github.com/lib/pq"
)

type MetaDataStore struct {
	db *sql.DB
}

func NewMetaDataStore(cfg config.PsqlConfig) (*MetaDataStore, error) {
	connStr := fmt.Sprintf(
		"host =%s port =%s user =%s password =%s dbname =%s sslmode=disable",
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

func (m *MetaDataStore) SaveApod(apod *entity.Apod) error {
	_, err := m.db.Exec(`INSERT INTO metadata (date,explanation,hdurl,media_type,title,url)
						VALUES ($1, $2, $3, $4, $5, $6)`,
		apod.Date,
		apod.Explanation,
		apod.HdURL,
		apod.MediaType,
		apod.Title,
		apod.Url)
	return err
}

func (m *MetaDataStore) GetApodByDate(date string) (entity.Apod, error) {
	apod := entity.Apod{}
	err := m.db.QueryRow(`SELECT date,explanation,hdurl,media_type,title,url 
							FROM metadata
							WHERE date = $1`, date).Scan(
		apod.Date,
		apod.Explanation,
		apod.HdURL,
		apod.MediaType,
		apod.Title,
		apod.Url,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return apod, fmt.Errorf("no such apod")
	}

	return apod, err
}

func (m *MetaDataStore) GetAllApods() ([]entity.Apod, error) {
	rows, err := m.db.Query(`SELECT date,explanation,hdurl,media_type,title,url 
							FROM metadata`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apods := []entity.Apod{}

	for rows.Next() {
		apod := entity.Apod{}
		err := rows.Scan(
			&apod.Date,
			&apod.Explanation,
			&apod.HdURL,
			&apod.MediaType,
			&apod.Title,
			&apod.Url,
		)
		if err != nil {
			continue
		}
		apods = append(apods, apod)
	}
	return apods, nil
}
