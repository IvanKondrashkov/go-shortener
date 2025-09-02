package logger

import (
	"go.uber.org/zap"
)

// ZapLogger оборачивает zap.Logger для удобного использования в приложении.
// Предоставляет методы для логирования и управления логгером.
type ZapLogger struct {
	Log *zap.Logger // Экземпляр zap.Logger
}

// NewZapLogger создает и возвращает новый экземпляр ZapLogger с указанным уровнем логирования.
// Уровень логирования должен быть одним из: "debug", "info", "warn", "error", "dpanic", "panic", "fatal".
// В случае ошибки парсинга уровня или создания логгера возвращает ошибку.
func NewZapLogger(level string) (*ZapLogger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		Log: zl,
	}, nil
}

// Sync синхронизирует буферизованные логи с их назначением (например, записывает на диск).
// Игнорирует ошибки синхронизации, так как они не критичны для работы приложения.
func (z *ZapLogger) Sync() {
	_ = z.Log.Sync()
}
