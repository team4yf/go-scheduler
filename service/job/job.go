package job

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/team4yf/go-scheduler/constant"
	"github.com/team4yf/go-scheduler/model"
	"github.com/team4yf/go-scheduler/pkg/email"
	"github.com/team4yf/go-scheduler/repository/job"
	"github.com/team4yf/go-scheduler/repository/subscribe"
	"github.com/team4yf/go-scheduler/repository/task"
	"google.golang.org/appengine/log"
)

var (
	jobRepo      job.JobRepo
	taskRepo     task.TaskRepo
	subRepo      subscribe.SubscribeRepo
	errNotInited = errors.New("Schedule Not Inited!")
	inited       = false //the flag of the service.Init()
)

type notifyBody struct {
	Event string     `json:"event"`
	Task  model.Task `json:"task"`
}

//Init init by the caller
func Init() {
	jobRepo = job.NewJobRepo()
	taskRepo = task.NewTaskRepo()
	subRepo = subscribe.NewSubscribeRepo()
}

func generateCallback(theJob *model.Job) func() {
	return func() {
		go runJob(theJob, func(data interface{}, err error) {
			if err != nil {
				log.Errorf("Run Job: %+v, Error: %+v\n", theJob, err)
				return
			}
		})
	}

}

//Callback the callback struct
type Callback func(data interface{}, err error)

//JobService service for job
type JobService interface {
	Init() error
	Start() error
	Execute(job *model.Job, callback Callback)
	Restart(job *model.Job, callback Callback) error
	Pause(job *model.Job) error
	Shutdown(job *model.Job) error
}

type JobWrapper struct {
	job *model.Job
	f   func()
	id  cron.EntryID
}
type simpleJobService struct {
	schedule *cron.Cron
	locker   sync.RWMutex
	handler  map[string]*JobWrapper
}

func NewSimpleJobService() JobService {
	serviec := &simpleJobService{
		handler: make(map[string]*JobWrapper),
	}
	Init()
	return serviec
}

func (s *simpleJobService) Init() (err error) {
	if inited {
		return nil
	}

	list, err := jobRepo.List(-1)
	if err != nil {
		return err
	}
	s.locker.RLock()
	defer s.locker.RUnlock()
	for _, j := range list {
		//define a callback
		//Important: 这里出现了闭包的问题

		//wrap the data and the callback
		s.handler[j.Code] = &JobWrapper{
			job: j,
			f:   generateCallback(j),
		}
	}
	s.schedule = cron.New()
	inited = true
	return nil
}

func (s *simpleJobService) Start() (err error) {
	if s.schedule == nil {
		return errNotInited
	}
	//add the func
	for _, wrapper := range s.handler {
		//ignore the not run job
		if wrapper.job.Status != constant.StatusRun {
			continue
		}
		id, err := s.schedule.AddFunc(wrapper.job.Cron, wrapper.f)
		if err != nil {
			return err
		}
		wrapper.id = id
		log.Infof("Start() Job: Code-> %v; Corn-> %v;\n", wrapper.job.Code, wrapper.job.Cron)
	}
	//startup
	s.schedule.Start()
	return nil
}

func (s *simpleJobService) Execute(job *model.Job, callback Callback) {
	go runJob(job, callback)
}

func (s *simpleJobService) Restart(job *model.Job, callback Callback) error {
	wrapper, ok := s.handler[job.Code]
	if !ok {
		//not exists
		return errors.New("job:" + job.Code + ", not exists")
	}

	if wrapper.id > 0 {
		//running
		return nil
	}
	id, err := s.schedule.AddFunc(wrapper.job.Cron, wrapper.f)
	if err != nil {
		return err
	}
	wrapper.id = id
	return nil
}

func (s *simpleJobService) Pause(job *model.Job) error {
	wrapper, ok := s.handler[job.Code]
	if !ok {
		//not exists
		return errors.New("job:" + job.Code + ", not exists")
	}
	if wrapper.id < 0 {
		//need not to pause
		return nil
	}
	s.schedule.Remove(wrapper.id)
	wrapper.id = -99
	return nil
}

func (s *simpleJobService) Shutdown(job *model.Job) error {
	s.schedule.Stop()
	s.schedule = nil
	return nil
}

//Actually run the http request
//Log the response for the task
func runJob(job *model.Job, callback Callback) {
	startAt := time.Now()
	task := &model.Task{
		Code:    job.Code,
		StartAt: startAt,
		URL:     job.URL,
		Status:  0,
	}
	if err := taskRepo.Create(task); err != nil {
		callback(nil, err)
		return
	}
	var rsp utils.ResponseWrapper
	//construct the auth data
	authProp := job.AuthProperties
	auth := &utils.HttpAuth{
		Type: utils.HTTPAuthType(job.Auth),
	}
	if authProp != "" {
		utils.StringToStruct(authProp, &auth.Data)
	} else {
		auth.Data = utils.HTTPAuthData(make(map[string]interface{}))
	}
	//perform the request
	switch job.ExecuteType {
	case "POST":
		rsp = utils.PostJsonWithAuth(job.URL, job.Argument, job.Timeout, auth)
	case "GET":
		rsp = utils.GetWithAuth(job.URL, job.Timeout, auth)
	case "FORM":
		rsp = utils.PostFormWithAuth(job.URL, job.Argument, job.Timeout, auth)
	}

	task.EndAt = time.Now()
	task.Cost = task.EndAt.UnixNano()/1e6 - task.StartAt.UnixNano()/1e6
	task.Status = rsp.StatusCode
	isSuccess := rsp.StatusCode == http.StatusOK
	if isSuccess {
		task.Log = rsp.Body
	}

	//update the task status
	if err := taskRepo.Update(task); err != nil {
		callback(nil, err)
		return
	}
	//notify the job's subscriber
	go func() {
		subs, err := subRepo.List(job.Code)
		if err != nil {
			log.Errorf("get subscriber err: %+v\n", err)
			return
		}
		if len(subs) < 1 {
			//no subscriber
			return
		}
		event := "fail"
		if isSuccess {
			event = "success"
		}
		body := notifyBody{
			Event: event,
			Task:  *task,
		}
		//用一个数组来存放通知的消费者，进行去重的操作，防止同一个事件被多次调用
		subscribers := []string{}

		for _, sub := range subs {
			if isSuccess == (sub.NotifyEvent == "success") {
				if utils.SliceIndexOf(subscribers, sub.Subscriber) > -1 {
					// contains
					continue
				}
				switch sub.NotifyType {
				case "webhook":
					utils.PostJson(sub.Subscriber, utils.JSON2String(body), 120)
				case "email":
					subject, content := email.NewNotifyEmail("./templates/task-email.html", "Scheduler-Notify", *task)
					email.Send(sub.Subscriber, subject, content)
				}
				subscribers = append(subscribers, sub.Subscriber)
			}
		}

	}()

	if !isSuccess {
		callback(nil, errors.New("status not ok!"+rsp.Body))
		return
	}
	callback(rsp.Body, nil)

}
