package job

import (
	"github.com/jinzhu/gorm"
	"github.com/team4yf/go-scheduler/constant"
	"github.com/team4yf/go-scheduler/model"
)

//JobRepo the rep for job
type JobRepo interface {
	Create(*model.Job) error
	List(status int) ([]*model.Job, error)
	Get(code string) (*model.Job, error)
	Update(name string, c model.CommonMap) error
}

type jobRepo struct {
	db *gorm.DB
}

//NewJobRepo create a new job repo
func NewJobRepo() JobRepo {
	return &jobRepo{
		db: model.Db,
	}
}

//List by status,
func (r *jobRepo) List(status int) (jobs []*model.Job, err error) {
	if status == constant.StatusAll {
		err = r.db.Model(model.Job{}).Find(&jobs).Error
		return
	}
	err = r.db.Model(model.Job{}).Where("status = ?", status).Find(&jobs).Error
	return
}

//Create new a job.
func (r *jobRepo) Create(job *model.Job) (err error) {
	err = r.db.Create(job).Error
	return
}

//Get by code.
func (r *jobRepo) Get(code string) (*model.Job, error) {
	job := model.Job{}
	err := r.db.Where("code = ?", code).First(&job).Error
	return &job, err
}

//Update information.
func (r *jobRepo) Update(code string, c model.CommonMap) error {
	return r.db.Model(model.Job{}).Where("code = ?", code).Updates(c).Error
}
