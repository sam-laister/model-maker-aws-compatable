package model

import "database/sql/driver"

type TaskStatus string

const (
	SUCCESS    TaskStatus = "SUCCESS"
	INPROGRESS TaskStatus = "INPROGRESS"
	FAILED     TaskStatus = "FAILED"
	INITIAL    TaskStatus = "INITIAL"
)

func (self *TaskStatus) Scan(value interface{}) error {
	*self = TaskStatus(value.([]byte))
	return nil
}

func (self TaskStatus) Value() (driver.Value, error) {
	return string(self), nil
}
