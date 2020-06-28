package email

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/team4yf/go-scheduler/config"
	"github.com/team4yf/go-scheduler/model"
)

func TestSend(t *testing.T) {
	// init config

	config.Init("../../conf/config.test.yaml")

	Init()
	time.Sleep(3 * time.Second)
	type args struct {
		to      string
		subject string
		body    string
	}

	task := model.Task{
		Code:   "Test Code",
		Log:    "ok",
		Status: 200,
	}
	subject, body := NewNotifyEmail("../../templates/task-email.html", "test", task)

	err := Send("fwang@evolveconsulting.com.hk", subject, body)
	time.Sleep(20 * time.Second)
	assert.Nil(t, err, "should not be err")

}
