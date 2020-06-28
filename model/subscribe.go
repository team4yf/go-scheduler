package model

import (
	"github.com/jinzhu/gorm"
)

//Subscribe the subscribe detail
type Subscribe struct {
	gorm.Model  `json:"-"`
	Code        string `json:"code"`        // the code of the job
	Subscriber  string `json:"subscriber"`  //
	NotifyTopic string `json:"notifyTopic"` // the topic of the notify
	NotifyEvent string `json:"notifyEvent"` // the event of the notify
	NotifyType  string `json:"notifyType"`  // the event notify type : webhook? email?
}

//TableName the table name
func (Subscribe) TableName() string {
	return "schedule_job_subscribe"
}
