package logger

import (
	"github.com/xissg/online-judge/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newCustomLogger(outputPaths, errorOutputPaths []string) (*zap.Logger, error) {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "", // 不记录日志调用位置
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
	}
	return cfg.Build()
}

func NewLogger(cfg config.LogConfig) (*zap.SugaredLogger, error) {
	logger, err := newCustomLogger(cfg.OutputPaths, cfg.ErrorOutputPaths)
	if err != nil {
		return nil, err
	}
	logger.WithOptions(zap.AddCallerSkip(100))
	sugar := logger.Sugar()
	defer func(sugar *zap.SugaredLogger) {
		sugar.Sync()
	}(sugar)
	return sugar, nil
}
