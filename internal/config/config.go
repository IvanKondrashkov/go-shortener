package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config содержит конфигурационные параметры приложения,
// которые могут быть установлены через переменные окружения.
type Config struct {
	Config          string `env:"CONFIG"`                                     // Конфигурационный файл в формате JSON
	ServerAddress   string `env:"SERVER_ADDRESS" json:"server_address"`       // Адрес сервера в формате host:port
	URL             string `env:"URL" json:"url"`                             // Базовый URL сервиса
	LogLevel        string `env:"LOG_LEVEL" json:"log_level"`                 // Уровень логирования (DEBUG, INFO, WARN, ERROR)
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"` // Путь к файловому хранилищу URL
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`           // DSN для подключения к БД
	AuthKey         string `env:"AUTH_KEY" json:"auth_key"`                   // Ключ для аутентификации

	TerminationTimeout int  `env:"TERMINATION_TIMEOUT" json:"termination_timeout"` // Таймаут завершения работы (в секундах)
	WorkerCount        int  `env:"WORKER_COUNT" json:"worker_count"`               // Количество воркеров
	EnableHTTPS        bool `env:"ENABLE_HTTPS" json:"enable_https"`               // Включение защищенного протокола
}

// Глобальные переменные конфигурации со значениями по умолчанию
var (
	ServerAddress   = "localhost:8080"
	URL             = "http://localhost:8080/"
	SecureURL       = "https://localhost:8080/"
	LogLevel        = "INFO"
	FileStoragePath = "internal/storage/urls.json"
	DatabaseDSN     = ""
	AuthKey         = []byte("6368616e676520746869732070617373776f726420746f206120736563726574")

	TerminationTimeout = time.Second * 30
	WorkerCount        = 10
	EnableHTTPS        = false
	FileConfigPath     = "internal/config/config.json"
)

// ParseConfig загружает конфигурацию приложения из:
// 1. Аргументов командной строки (имеют наивысший приоритет)
// 2. Переменных окружения
// 3. JSON файла конфигурации
// 4. Значений по умолчанию
//
// Возвращает ошибку если не удалось распарсить конфигурацию.
func ParseConfig() error {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "Base host host:port")
	flag.StringVar(&URL, "b", URL, "Base url protocol://host:port/")
	flag.StringVar(&LogLevel, "l", LogLevel, "Base log level info")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "Base file storage path")
	flag.StringVar(&DatabaseDSN, "d", DatabaseDSN, "Base url db connection")
	flag.BoolVar(&EnableHTTPS, "s", EnableHTTPS, "Enable secure protocol")
	flag.StringVar(&FileConfigPath, "c", FileConfigPath, "Configuration JSON file")
	flag.Parse()

	var envCfg Config
	err := env.Parse(&envCfg)
	if err != nil {
		return fmt.Errorf("config parse error: %w", err)
	}

	var jsonCfg *Config
	if envConfig := envCfg.Config; envConfig != "" {
		jsonCfg, err = parseJSONConfig(envCfg.Config)
		if err != nil {
			return fmt.Errorf("failed to parse JSON config: %w", err)
		}
		applyJSONConfig(envCfg, jsonCfg)
	}

	if envServerAddress := envCfg.ServerAddress; envServerAddress != "" {
		ServerAddress = envServerAddress
	}

	if envBaseURL := envCfg.URL; envBaseURL != "" {
		URL = envBaseURL
	}

	if envLogLevel := envCfg.LogLevel; envLogLevel != "" {
		LogLevel = envLogLevel
	}

	if envFileStoragePath := envCfg.FileStoragePath; envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}

	if envDatabaseDsn := envCfg.DatabaseDSN; envDatabaseDsn != "" {
		DatabaseDSN = envDatabaseDsn
	}

	if envAuthKey := envCfg.AuthKey; envAuthKey != "" {
		AuthKey = []byte(envAuthKey)
	}

	if envTerminationTimeout := envCfg.TerminationTimeout; envTerminationTimeout != 0 {
		TerminationTimeout = time.Duration(envTerminationTimeout)
	}

	if envWorkerCount := envCfg.WorkerCount; envWorkerCount != 0 {
		WorkerCount = envWorkerCount
	}

	if envEnableHTTPS := envCfg.EnableHTTPS; envEnableHTTPS {
		EnableHTTPS = true
	}

	if EnableHTTPS {
		URL = SecureURL
	}

	if !strings.HasSuffix(URL, "/") {
		URL += "/"
	}
	return nil
}

// parseJSONConfig загружает конфигурацию из JSON файла
func parseJSONConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// applyJSONConfig применяет значения из JSON конфигурации (низший приоритет)
func applyJSONConfig(envCfg Config, jsonCfg *Config) {
	if jsonCfg == nil {
		return
	}

	if envCfg.ServerAddress == "" {
		ServerAddress = jsonCfg.ServerAddress
	}

	if envCfg.URL == "" {
		URL = jsonCfg.URL
	}

	if envCfg.LogLevel == "" {
		LogLevel = jsonCfg.LogLevel
	}

	if envCfg.FileStoragePath == "" {
		FileStoragePath = jsonCfg.FileStoragePath
	}

	if envCfg.DatabaseDSN == "" {
		DatabaseDSN = jsonCfg.DatabaseDSN
	}

	if envCfg.AuthKey == "" {
		AuthKey = []byte(jsonCfg.AuthKey)
	}

	if envCfg.TerminationTimeout == 0 {
		TerminationTimeout = time.Duration(jsonCfg.TerminationTimeout) * time.Second
	}

	if envCfg.WorkerCount == 0 {
		WorkerCount = jsonCfg.WorkerCount
	}

	if !envCfg.EnableHTTPS {
		EnableHTTPS = jsonCfg.EnableHTTPS
	}
}
