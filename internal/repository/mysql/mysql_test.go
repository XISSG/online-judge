package mysql

import (
	"github.com/xissg/online-judge/internal/config"
	"testing"
)

func TestMysqlClient(t *testing.T) {
	appConfig := config.LoadConfig()
	NewMysqlClient(appConfig.Database)
}
