package ai

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/utils"
	"time"
)

type AIClient struct {
	conn *websocket.Conn
}

const NORMAL_RESPONSE_CODE = 101

func NewAIClient(cfg config.AIConfig) *AIClient {
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	authStr := utils.AssembleAuthUrl(cfg.HostUrl, cfg.ApiKey, cfg.ApiSecret)
	conn, resp, err := dialer.Dial(authStr, nil)
	if err != nil {
		return nil
	}
	if resp.StatusCode != NORMAL_RESPONSE_CODE {
		return nil
	}
	return &AIClient{
		conn: conn,
	}
}

// topk灵活度，1-6默认为4,temperature随机性0-1, roleSetting设置ai扮演的角色
func (c *AIClient) AISetting(appId string, topK int, temperature float64) *request.AI {
	if topK <= 0 || topK > 6 {
		topK = 4
	}
	if temperature <= 0.0 || temperature > 1 {
		temperature = 0.8
	}
	data := &request.AI{
		Header: &request.Header{
			AppID: appId,
		},
		Parameter: &request.Parameter{
			Chat: &request.Chat{
				Domain:      "generalv3.5",
				Temperature: temperature,
				TopK:        topK,
				MaxTokens:   4096,
			},
		},
		Payload: &request.Payload{
			Message: &request.Message{
				Text: []*request.Text{},
			},
		},
	}
	return data
}

func (c *AIClient) SendMessage(aiRequest *request.AI) error {
	return c.conn.WriteJSON(aiRequest)
}

// staus字段为2时代表数据发送完毕
func (c *AIClient) ReadMessage() ([]*response.AI, error) {
	var dataStream []*response.AI
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return nil, err
		}
		var dataResponse *response.AI
		err = json.Unmarshal(data, &dataResponse)
		if err != nil {
			return nil, err
		}
		dataStream = append(dataStream, dataResponse)
		if dataResponse.Header.Status == 2 {
			c.conn.Close()
			break
		}
	}

	return dataStream, nil
}
