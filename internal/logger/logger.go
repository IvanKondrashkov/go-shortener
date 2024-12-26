package logger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	Log *zap.Logger
}

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

func (z *ZapLogger) Sync() {
	_ = z.Log.Sync()
}
