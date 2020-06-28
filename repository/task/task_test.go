package task

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/team4yf/go-scheduler/config"
	"github.com/team4yf/go-scheduler/model"
)

var (
	rep taskRepo
)

func init() {
	viper.SetConfigFile("../../conf/config.test.yaml")
	viper.SetConfigType("yaml") // 设置配置文件格式为yaml

	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		panic(err)
	}
	if db, err := model.CreateTmpDb(&config.DBSetting{
		Engine:   viper.GetString("db.engine"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Database: viper.GetString("db.database"),
		Charset:  viper.GetString("db.charset"),
		ShowSQL:  viper.GetBool("db.showSql"),
	}); err != nil {
		panic(err)
	} else {
		rep.db = db
	}

	// migration.Install()
}
func TestCreate(t *testing.T) {
	code := "test3"
	task := &model.Task{
		Code:    code,
		StartAt: time.Now(),
		Log:     "log",
		Status:  1,
		Cost:    100,
		URL:     "http://baidu.com",
	}
	rep.Create(task)
	fmt.Printf("Create-Task-> id: %d \n", task.ID)

	tk, err := rep.Get(task.ID)
	assert.Equal(t, 100, tk.Cost)
	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, tk, "task should not nil")

	rep.Clear(code)
}
func TestList(t *testing.T) {

	code := "test4"
	rep.Clear(code)
	for i := 0; i < 10; i++ {
		task := &model.Task{
			Code:    code,
			StartAt: time.Now(),
			Log:     "log",
			Status:  1,
			Cost:    100,
			URL:     "http://baidu.com/" + strconv.Itoa(i),
		}
		rep.Create(task)
	}

	list, err := rep.List(code, 1, -1)
	assert.Nil(t, err, "err should be nil")
	assert.NotNil(t, list, "task list should not nil")
	assert.Equal(t, len(list), 10)
	rep.Clear(code)
}

func TestUpdate(t *testing.T) {
	code := "test5"
	rep.Clear(code)
	task := &model.Task{
		Code:    code,
		StartAt: time.Now(),
		Log:     "log",
		Status:  1,
		Cost:    100,
		URL:     "http://baidu.com",
	}
	rep.Create(task)
	fmt.Printf("Create-Task-> id: %d \n", task.ID)

	task.EndAt = time.Now()
	task.Status = 2
	err := rep.Update(task)
	assert.Nil(t, err, "err should be nil")
	tk, _ := rep.Get(task.ID)
	assert.Equal(t, 2, tk.Status)
	rep.Clear(code)
}
