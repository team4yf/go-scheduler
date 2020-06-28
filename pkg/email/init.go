package email

import (
	"errors"
	"sync"

	"github.com/spf13/viper"

	"github.com/team4yf/go-scheduler/pkg/log"
)

// Client 邮件发送客户端
var Client Driver

// Lock 读写锁
var Lock sync.RWMutex

var (
	// ErrChanNotOpen 邮件队列没有开启
	ErrChanNotOpen = errors.New("email queue does not open")
)

// Init 初始化客户端
func Init() {
	log.Info("email init")
	Lock.Lock()
	defer Lock.Unlock()

	// 确保是已经关闭的
	if Client != nil {
		Client.Close()
	}

	client := NewSMTPClient(SMTPConfig{
		Name:      viper.GetString("email.name"),
		Address:   viper.GetString("email.address"),
		ReplyTo:   viper.GetString("email.reply"),
		Host:      viper.GetString("email.host"),
		Port:      viper.GetInt("email.port"),
		Username:  viper.GetString("email.username"),
		Password:  viper.GetString("email.password"), // 这里读取的是系统的环境变量 #GS_EMAIL_PASSWORD
		Keepalive: viper.GetInt("email.keepalive"),
	})
	// log.Info("client: %+v\n", client)
	client.Init()
	Client = client

}
