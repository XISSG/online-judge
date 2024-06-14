package ai

import (
	"errors"
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/request"
	"testing"
)

func Test(t *testing.T) {
	appConfig := config.LoadConfig()
	client := NewAIClient(appConfig.AI)

	if client == nil {
		panic(errors.New("client is nil"))
	}
	data := &request.AI{
		Header: &request.Header{
			AppID: appConfig.AI.AppId,
		},
		Parameter: &request.Parameter{
			Chat: &request.Chat{
				Domain:      "general",
				Temperature: float64(0.8),
				TopK:        6,
				MaxTokens:   2048,
			},
		},
		Payload: &request.Payload{
			Message: &request.Message{
				Text: []*request.Text{
					{Role: constant.USER_ROLE, Content: "hello "},
				},
			},
		},
	}
	err := client.SendMessage(data)
	if err != nil {
		panic(err)
	}
	res, err := client.ReadMessage()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
