package configuration

import (
	"flag"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Network *NetworkCfg `yaml:"network"`
	Logger  *LoggerCfg  `yaml:"logger"`
	Storage *Storage    `yaml:"storage"`
}
type LoggerCfg struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
type NetworkCfg struct {
	ServerAddress  string `yaml:"address"`
	BaseURL        string `yaml:"base-url"`
	AccuralAddress string `yaml:"accural-address"`
}

type Storage struct {
	FileStoragePath string `yaml:"address"`
	ConnString      string
	PGPoolCfg       *pgxpool.Config
}

func New() *Config {
	networkCfg := NetworkCfg{}
	loggerCfg := LoggerCfg{}
	storageCfg := Storage{}
	flag.StringVar(&networkCfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&networkCfg.AccuralAddress, "r", "http://localhost:8081", "system location address")
	flag.StringVar(&loggerCfg.Level, "l", "info", "logger level")
	flag.StringVar(&storageCfg.ConnString, "d", "", "database storage")

	flag.Parse()
	if serverAddress := os.Getenv("RUN_ADDRESS"); serverAddress != "" {
		networkCfg.ServerAddress = serverAddress
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		networkCfg.BaseURL = baseURL
	}
	if loggerLevel := os.Getenv("LOG_LEVEL"); loggerLevel != "" {
		loggerCfg.Level = loggerLevel
	}

	if accuralAddressStr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accuralAddressStr != "" {
		networkCfg.AccuralAddress = accuralAddressStr
	}
	return &Config{
		Network: &networkCfg,
		Logger:  &loggerCfg,
		Storage: &storageCfg,
	}
}
