package ai

import (
	"errors"
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"testing"
)

func Test(t *testing.T) {
	appConfig := config.LoadConfig()

	client := NewAIClient(appConfig.AI)

	if client == nil {
		panic(errors.New("client is nil"))
	}
	err := client.SendMessage("hello")
	if err != nil {
		panic(err)
	}
	res, err := client.ReadMessage()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
