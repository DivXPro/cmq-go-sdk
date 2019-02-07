package cmq

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/resty.v1"
)

// PATH 接口路径
const Path = "/v2/index.php"

const ActionSendMessage = "SendMessage"
const ActionReceiveMessage = "ReceiveMessage"

// Client CMQ客户端
type Client struct {
	host      string
	secretID  string
	secretKey string
	version   string
	method    string
}

// Params CMQ接口信息的数据结构
type Params struct {
	Action          string `json:"action"`
	Nonce           int64  `json:"Nonce"`
	RequestClient   string `json:"RequestClient"`
	SecretID        string `json:"secretId"`
	SignatureMethod string `json:"signatureMethod"`
	Signature       string `json:"signature"`
	Timestamp       int64  `json:"Timestamp"`
}

type SendMessageParam struct {
	Params
	QueueMessage
}

type ReceiveMessageParam struct {
	Params
	ReceiveMessage
}

// SendMessageResp 发送消息返回的数据解雇
type SendMessageResp struct {
	Code      int64  `json:"code"`
	Message   string `json:"message"`
	MsgID     string `json:"msgId"`
	RequestID string `json:"requestId"`
}

// ReceiveMessageResp 接受消息返回的数据结构
type ReceiveMessageResp struct {
	SendMessageResp
	MsgBody          string `json:"msgBody"`
	ReceiptHandle    string `json:"receiptHandle"`
	EnqueueTime      int64  `json:"enqueueTime"`
	FirstDequeueTime int64  `json:"firstDequeueTime"`
	NextVisibleTime  int64  `json:"nextVisibleTime"`
	DequeueCount     int64  `json:"dequeueCount"`
}

// NewClient 创建新的Client实例
func NewClient(host string, secretID string, secretKey string) *Client {
	var client = &Client{host, secretID, secretKey, "CMQ_CLIENT_GO_1.0", "POST"}
	return client
}

func (client *Client) buildParams(action string, message interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())
	switch action {
	case ActionSendMessage:
		params := SendMessageParam{
			Params{action, rand.Int63(), client.version, client.secretID, "HmacSHA1", "", time.Now().Unix()},
			message.(QueueMessage),
		}
		plainText := MakeSignPlainText(params, client.method, client.host, Path)
		params.Params.Signature = Sign("", client.secretID, []byte(plainText))
		return params
	case ActionReceiveMessage:
		params := ReceiveMessageParam{
			Params{action, rand.Int63(), client.version, client.secretID, "HmacSHA1", "", time.Now().Unix()},
			message.(ReceiveMessage),
		}
		return params
	default:
		return nil
	}
}

func (client *Client) request(action string, params interface{}) (interface{}, error) {
	sendParams := client.buildParams(action, params)
	if resp, err := resty.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetBody(sendParams).Post(client.host + Path); err != nil {
		return resp, err
	} else {
		if status := resp.StatusCode(); status != 200 {
			return resp, errors.New("Request fail " + strconv.Itoa(status))
		} else {
			switch action {
			case ActionSendMessage:
				var sendMessageResp SendMessageResp
				if err = json.Unmarshal(resp.Body(), &sendMessageResp); err != nil {
					return resp, errors.New("SendMessageResp Unmarshal fail")
				} else {
					return sendMessageResp, nil
				}
			case ActionReceiveMessage:
				var receiveMessageResp ReceiveMessageResp
				if err = json.Unmarshal(resp.Body(), &receiveMessageResp); err != nil {
					return resp, errors.New("SendMessageResp Unmarshal fail")
				} else {
					return receiveMessageResp, nil
				}
			default:
				return resp, errors.New("unknown action type")
			}
		}
	}
}

// SendMessage 发送消息到队列
func (client *Client) SendMessage(params interface{}) (SendMessageResp, error) {
	if resp, err := client.request(ActionSendMessage, params); err != nil {
		return resp.(SendMessageResp), err
	} else {
		return resp.(SendMessageResp), nil
	}
}

// ReceiveMessage 发送消息到队列
func (client *Client) ReceiveMessage(params interface{}) (ReceiveMessageResp, error) {
	if resp, err := client.request(ActionReceiveMessage, params); err != nil {
		return ReceiveMessageResp{}, err
	} else {
		return resp.(ReceiveMessageResp), nil
	}
}
