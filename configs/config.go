package configs

import "github.com/spf13/viper"

type Config struct {
	Host     string
	Port     string
	DBPort   string
	User     string
	Password string
	DBName   string
}

func NewConfig() (*Config, error) {
	viper.SetConfigFile("configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("server.port"),
		DBPort:   viper.GetString("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.name"),
	}

	return cfg, nil
}
