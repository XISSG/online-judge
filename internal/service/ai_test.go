package service

import (
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/ai"
	"testing"
)

var (
	hostUrl = "wss://spark-api.xf-yun.com/v1.1/chat"
	//hostUrl   = "wss://spark-api.xf-yun.com/v3.5/chat"
	appid     = "3fd9659c"
	apiSecret = "NDM4MjZjY2VkYzE0MzE5NWJlMGEyOTc1"
	apiKey    = "030ebc99d5951548bf7991edd9d83a2c"
)

func TestAI(t *testing.T) {
	appConfig := config.LoadConfig()
	client := ai.NewAIClient(appConfig.AI.HostUrl, appConfig.AI.ApiKey, appConfig.AI.ApiSecret)
	//roleStr := "现在你是位偶像练习生，会唱跳rap,篮球技术精湛，接下来你将扮演该角色进行对话"
	service := NewAIService(client, appid, 4, 0.6)
	service.SendMessage("来一首rap,以鸡，篮球，背带裤，中分为关键词生成")
	service.ReceiveMessage()
}
