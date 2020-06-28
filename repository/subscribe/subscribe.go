package subscribe

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/team4yf/go-scheduler/model"
)

//SubscribeRepo the rep for Subscribe
type SubscribeRepo interface {
	Subscribe([]*model.Subscribe) error
	List(code string) ([]*model.Subscribe, error)
	UnSubscribe(*model.Subscribe) error
}

type subscribeRepo struct {
	db *gorm.DB
}

//NewJobRepo create a new job repo
func NewSubscribeRepo() SubscribeRepo {
	return &subscribeRepo{
		db: model.Db,
	}
}

//Create new a Subscribe
func (r *subscribeRepo) Subscribe(subs []*model.Subscribe) (err error) {
	tx := r.db.Begin()
	txOK := false
	defer func() {
		if !txOK {
			tx.Rollback()
		}
	}()
	for _, sub := range subs {
		fmt.Printf("%+v\n", &sub)
		err = tx.Create(&sub).Error
		if err != nil {
			return
		}
	}
	txOK = true

	tx.Commit()
	return
}

//List by code
func (r *subscribeRepo) List(code string) (Subscribes []*model.Subscribe, err error) {
	err = r.db.Where("code = ?", code).Find(&Subscribes).Error
	return
}

func (r *subscribeRepo) UnSubscribe(sub *model.Subscribe) (err error) {
	err = r.db.Where("code = ? and subscriber = ? and notify_type = ? and notify_event = ?", sub.Code, sub.Subscriber, sub.NotifyType, sub.NotifyEvent).Delete(model.Subscribe{}).Error
	return
}
