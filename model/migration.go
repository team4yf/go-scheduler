package model

//Migration auto create or upgrade table
type Migration struct{}

//Install create the tables
func (migration *Migration) Install() error {

	tables := []interface{}{
		&Job{}, &Task{}, &Subscribe{},
	}
	Db.AutoMigrate(tables...)
	return nil
}
