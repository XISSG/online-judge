package logger

import (
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	curPath, _ := os.Getwd()
	outputPath := curPath + "/../../public/log/stdout.log"
	errorPath := curPath + "/../../public/log/error.log"
	outputPaths := []string{"stdout", outputPath}
	errorPaths := []string{errorPath}
	logger, _ := newCustomLogger(
		outputPaths,
		errorPaths,
	)
	// 增加一个 skip 选项，触发 zap 内部 error，将错误输出到 error.log
	logger.WithOptions(zap.AddCallerSkip(100))
	sugar := logger.Sugar()
	defer func(sugar *zap.SugaredLogger) {
		err := sugar.Sync()
		if err != nil {

		}
	}(sugar)

	sugar.Info("Info msg")
	sugar.Error("Error msg")
}
