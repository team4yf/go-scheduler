package email

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/team4yf/go-scheduler/model"
)

// NewNotifyEmail
func NewNotifyEmail(tmpPath, title string, task model.Task) (subject string, body string) {
	mailTplContent := getEmailHTMLContent(tmpPath, task)
	return title, mailTplContent
}

type NotifyMailData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// getEmailHTMLContent 获取邮件模板
func getEmailHTMLContent(tplPath string, mailData interface{}) string {
	b, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return ""
	}
	mailTpl := string(b)
	tpl, err := template.New("email tpl").Parse(mailTpl)
	if err != nil {
		return ""
	}
	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, mailData)
	if err != nil {
		fmt.Println("exec err", err)
	}
	return buffer.String()
}
