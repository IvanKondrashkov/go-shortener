package config

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config содержит конфигурационные параметры приложения,
// которые могут быть установлены через переменные окружения.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`    // Адрес сервера в формате host:port
	URL             string `env:"URL"`               // Базовый URL сервиса
	LogLevel        string `env:"LOG_LEVEL"`         // Уровень логирования (DEBUG, INFO, WARN, ERROR)
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // Путь к файловому хранилищу URL
	DatabaseDSN     string `env:"DATABASE_DSN"`      // DSN для подключения к БД
	AuthKey         string `env:"AUTH_KEY"`          // Ключ для аутентификации

	TerminationTimeout int  `env:"TERMINATION_TIMEOUT"` // Таймаут завершения работы (в секундах)
	WorkerCount        int  `env:"WORKER_COUNT"`        // Количество воркеров
	EnableHTTPS        bool `env:"ENABLE_HTTPS"`        // Включение защищенного протокола
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
)

// ParseConfig загружает конфигурацию приложения из:
// 1. Аргументов командной строки (имеют наивысший приоритет)
// 2. Переменных окружения
// 3. Значений по умолчанию
//
// Возвращает ошибку если не удалось распарсить конфигурацию.
func ParseConfig() error {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "Base host host:port")
	flag.StringVar(&URL, "b", URL, "Base url protocol://host:port/")
	flag.StringVar(&LogLevel, "l", LogLevel, "Base log level info")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "Base file storage path")
	flag.StringVar(&DatabaseDSN, "d", DatabaseDSN, "Base url db connection")
	flag.BoolVar(&EnableHTTPS, "s", EnableHTTPS, "Enable secure protocol")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("config parse error: %w", err)
	}

	if envServerAddress := cfg.ServerAddress; envServerAddress != "" {
		ServerAddress = envServerAddress
	}

	if envBaseURL := cfg.URL; envBaseURL != "" {
		URL = envBaseURL
	}

	if envLogLevel := cfg.LogLevel; envLogLevel != "" {
		LogLevel = envLogLevel
	}

	if envFileStoragePath := cfg.FileStoragePath; envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}

	if envDatabaseDsn := cfg.DatabaseDSN; envDatabaseDsn != "" {
		DatabaseDSN = envDatabaseDsn
	}

	if envAuthKey := cfg.AuthKey; envAuthKey != "" {
		AuthKey = []byte(envAuthKey)
	}

	if envTerminationTimeout := cfg.TerminationTimeout; envTerminationTimeout != 0 {
		TerminationTimeout = time.Duration(envTerminationTimeout)
	}

	if envWorkerCount := cfg.WorkerCount; envWorkerCount != 0 {
		WorkerCount = envWorkerCount
	}

	if envEnableHTTPS := cfg.EnableHTTPS; envEnableHTTPS {
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
