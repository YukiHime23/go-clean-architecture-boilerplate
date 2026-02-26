package logger

import (
	"go-clean-architecture-boilerplate/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a Zap logger based on the application config.
// In release mode it uses structured JSON; in debug mode it uses colored console output.
func New(cfg config.AppConfig) (*zap.Logger, error) {
	var zapCfg zap.Config
	if cfg.GinMode == "release" {
		zapCfg = zap.NewProductionConfig()
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return zapCfg.Build()
}

// NewNop returns a no-op logger suitable for testing.
func NewNop() *zap.Logger {
	return zap.NewNop()
}
