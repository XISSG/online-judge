package ai

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/model/response"
	"github.com/xissg/online-judge/internal/utils"
	"time"
)

type AIClient struct {
	conn *websocket.Conn
	ctx  *request.AI
}

// option可以指定maxtokens,topk,temperature,domain
func NewAIClient(cfg config.AIConfig, option ...Option) *AIClient {
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	authStr := utils.AssembleAuthUrl(cfg.HostUrl, cfg.ApiKey, cfg.ApiSecret)
	conn, resp, err := dialer.Dial(authStr, nil)
	if err != nil {
		return nil
	}
	if resp.StatusCode != constant.NORMAL_RESPONSE_CODE {
		return nil
	}

	ctx := initAIParams(cfg.AppId, option...)
	return &AIClient{
		conn: conn,
		ctx:  ctx,
	}
}

func (c *AIClient) SendMessage(message string) error {
	text := &request.Text{
		Role:    constant.USER_ROLE,
		Content: message,
	}
	c.ctx.Payload.Message.Text = append(c.ctx.Payload.Message.Text, text)
	return c.conn.WriteJSON(c.ctx)
}

// ReadMessage status字段为2时代表数据发送完毕
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
		if dataResponse.Header.Status == constant.EOF_RESPONSE_STATUS {
			c.conn.Close()
			break
		}
	}

	return dataStream, nil
}
