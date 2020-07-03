package model

import (
	"github.com/jinzhu/gorm"
)

//Job the job detail
type Job struct {
	gorm.Model     `json:"-"`
	Code           string `json:"code" gorm:"index"` // the code of the job, it's unique
	Name           string `json:"name" gorm:"index"`
	Cron           string `json:"cron"`   // the cron of the job
	Status         int    `json:"status"` // the status of the job, running/stoped/shutdown
	Title          string `json:"title"`
	Remark         string `json:"remark"`
	RetryMax       int    `json:"retryMax"`
	Timeout        int    `json:"timeout"`
	Argument       string `json:"argument"`       //the body/param of the request
	Delay          int    `json:"delay"`          // the delay of failue, it should be second
	ExecuteType    string `json:"executeType"`    // http invoke, grpc, internal
	Auth           string `json:"auth"`           //auth method, None,Basic,Other,now only support None/Basic
	AuthProperties string `json:"authProperties"` // properties, normally json here
	NotifyTopic    string `json:"notifyTopic"`    // the topic of the notify
	URL            string `json:"url"`            //the http url
}

//TableName the table name
func (Job) TableName() string {
	return "schedule_job"
}
