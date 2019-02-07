package cmq

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
)

type Queue struct {
	queueName string
	client    *Client
	encoding  bool
}

type QueueMessage struct {
	QueueName    string      `json:"queueName"`
	MsgBody      interface{} `json:"msgBody"`
	DelaySeconds int64       `json:"delaySeconds"`
}

type ReceiveMessage struct {
	QueueName          string `json:"queueName"`
	PollingWaitSeconds int64  `json:"pollingWaitSeconds"`
}

// NewQueue 构建队列实例
func NewQueue(queueName string, client *Client, encoding bool) *Queue {
	var queue = &Queue{queueName, client, encoding}
	return queue
}

// SendMessage 发送信息到队列
func (queue *Queue) SendMessage(msg interface{}, delay int64) (SendMessageResp, error) {
	var msgBody string
	if reflect.ValueOf(msg).Kind() == reflect.String {
		msgBody = msg.(string)
	} else {
		if data, err := json.Marshal(msg); err != nil {
			msgBody = ""
			return SendMessageResp{}, err
		} else {
			msgBody = string(data)
		}
	}
	encodeMsg := base64.URLEncoding.EncodeToString([]byte(msgBody))
	queueMessage := QueueMessage{queue.queueName, encodeMsg, delay}
	return queue.client.SendMessage(queueMessage)
}

// ReceiveMessage 从队列消费消息
func (queue *Queue) ReceiveMessage(pollingWaitSeconds int64) (ReceiveMessageResp, error) {
	if pollingWaitSeconds > 30 || pollingWaitSeconds < 0 {
		pollingWaitSeconds = 30
	}
	receiveMessage := ReceiveMessage{queue.queueName, pollingWaitSeconds}
	return queue.client.ReceiveMessage(receiveMessage)
}
