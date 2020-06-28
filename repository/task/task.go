package task

import (
	"github.com/jinzhu/gorm"
	"github.com/team4yf/go-scheduler/model"
)

//TaskRepo the rep for task
type TaskRepo interface {
	Create(*model.Task) error
	List(code string, p, l int) ([]*model.Task, error)
	Get(id uint) (*model.Task, error)
	Update(*model.Task) error
	Clear(code string) error
}

type taskRepo struct {
	db *gorm.DB
}

//NewJobRepo create a new job repo
func NewTaskRepo() TaskRepo {
	return &taskRepo{
		db: model.Db,
	}
}

//Create new a task
func (r *taskRepo) Create(task *model.Task) (err error) {
	err = r.db.Create(&task).Error
	return
}

//List by code
func (r *taskRepo) List(code string, p, l int) (tasks []*model.Task, err error) {
	err = r.db.Where("code = ?", code).Limit(l).Offset((p - 1) * l).Find(&tasks).Error
	return
}

//Get by code.
func (r *taskRepo) Get(id uint) (*model.Task, error) {
	var tk model.Task
	err := r.db.Where("id = ?", id).First(&tk).Error
	return &tk, err
}

//Update information.
func (r *taskRepo) Update(task *model.Task) (err error) {
	err = r.db.Save(&task).Error
	return
}
func (r *taskRepo) Clear(code string) (err error) {
	err = r.db.Where("code = ?", code).Delete(model.Task{}).Error
	return
}
