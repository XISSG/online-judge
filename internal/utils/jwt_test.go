package utils

import (
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"testing"
)

func TestJwt(t *testing.T) {
	appConfig := config.LoadConfig()
	tokenString, err := generate(1, appConfig.Jwt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tokenString)
	res, err := parse(tokenString, appConfig.Jwt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
