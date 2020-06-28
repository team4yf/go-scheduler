package model

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var keys = strings.Split("ID,Code,Log,StartAt,EndAt,Status,Cost,URL", ",")

func TestString(t *testing.T) {
	task := &Task{
		Code:    "abc",
		StartAt: time.Now(),
		Log:     "log",
		Status:  1,
		Cost:    100,
		URL:     "http://baidu.com/",
	}
	values := task.String(keys)
	assert.Equal(t, values[1], "abc")
	assert.Equal(t, values[2], "log")
}
