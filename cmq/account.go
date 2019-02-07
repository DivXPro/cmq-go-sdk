package cmq

// Account CMQ服务账号
type Account struct {
	host      string
	secretID  string
	secretKey string
	client    *Client
}

// NewAccount 创建Account实例
func NewAccount(host string, secretID string, secretKey string) *Account {
	var account = &Account{host, secretID, secretKey, NewClient(host, secretID, secretKey)}
	return account
}

// GetQueue 获取队列实例
func (account *Account) GetQueue(queueName string) *Queue {
	queue := &Queue{queueName, account.client, true}
	return queue
}
