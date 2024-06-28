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

var (
	once           sync.Once
	singleInstance *Storage
)

func Connect(cfg *configs.Config) *Storage {
	once.Do(func() {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			slog.Error("Failed to connect to the Database", "error", err)
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			slog.Error("Failed to Ping", "error", err)
			panic(err)
		}

		singleInstance = &Storage{db: db}
	})

	return singleInstance
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}
