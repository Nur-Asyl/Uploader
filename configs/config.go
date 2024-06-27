package configs

import (
	"github.com/spf13/viper"
	"log/slog"
	"strings"
	"sync"
)

type Config struct {
	Server   *server
	Database *database
}

type server struct {
	Port string
}

type database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigFile("configs/config.yaml")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			slog.Error("Failed to read in config", "error", err)
			panic(err)
		}

		if err := viper.Unmarshal(&configInstance); err != nil {
			slog.Error("Failed to unmarshal config", "error", err)
			panic(err)
		}

	})

	return configInstance
}
