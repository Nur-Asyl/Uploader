package storages

import (
	"Manual_Parser/configs"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"sync"
)

type Storage struct {
	db *sql.DB
}

var lock = &sync.Mutex{}

var singleInstance *Storage

func Connect(cfg *configs.Config) (*Storage, error) {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		slog.Info("Creating single instance of postgreSQL")
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.DBPort, cfg.User, cfg.Password, cfg.DBName)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			return nil, err
		}

		return &Storage{db: db}, nil
	} else {
		slog.Info("Connecting to single instance of postgreSQL")
	}

	return singleInstance, nil
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}
