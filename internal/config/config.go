package config

import (
	"flag"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURL          string `env:"DATABASE_URL"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
	SecretKey            string `env:"SECRET_KEY"`
	UpdateInterval       int    `env:"UPDATE_INTERVAL"`
	WorkerCount          int    `env:"WORKER_COUNT"`
}

func GetConfig() (Config, error) {
	var config Config

	flag.StringVar(&config.RunAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&config.DatabaseURL, "d", "", "database dsn")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8081", "accrual system address")
	flag.StringVar(&config.LogLevel, "l", "debug", "log level")
	flag.StringVar(&config.SecretKey, "s", "", "secret key")
	flag.IntVar(&config.UpdateInterval, "i", 10, "update interval in seconds")
	flag.IntVar(&config.WorkerCount, "w", 5, "number of workers")

	flag.Parse()

	/*if err := env.Parse(&config); err != nil {
		return Config{}, err
	}*/

	return config, nil
}
