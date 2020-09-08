package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

const TimeFormatter = "2006-01-02 15:04:05"

//Task task detail
type Task struct {
	gorm.Model `json:"-"`
	Code       string    `json:"code"`    // the job code
	Log        string    `json:"log"`     // the result of run the task
	StartAt    time.Time `json:"startAt"` // the start time of the task
	EndAt      time.Time `json:"endAt"`   // the end time of the task
	Status     int       `json:"status"`  // the status of the job, running/stoped/shutdown
	Cost       int64     `json:"cost"`    // the cost of the execute task
	URL        string    `json:"url"`     // the http url
}

//TableName the table name
func (Task) TableName() string {
	return "schedule_task"
}

func (t *Task) String(headers []string) []string {
	values := []string{}

	for _, key := range headers {
		switch strings.ToLower(strings.ReplaceAll(key, "_", "")) {
		case "id":
			values = append(values, strconv.FormatUint((uint64)(t.ID), 10))
		case "code":
			values = append(values, t.Code)
		case "log":
			values = append(values, t.Log)
		case "startat":
			values = append(values, t.StartAt.Format(TimeFormatter))
		case "endat":
			values = append(values, t.EndAt.Format(TimeFormatter))
		case "status":
			values = append(values, strconv.FormatInt((int64)(t.Status), 10))
		case "cost":
			values = append(values, strconv.FormatInt((int64)(t.Cost), 10))
		case "url":
			values = append(values, t.URL)

		}
	}

	return values
}
