package model

import "database/sql/driver"

type ReportType string

const (
	BUG      ReportType = "BUG"
	FEEDBACK ReportType = "FEEDBACK"
)

func (ts *ReportType) Scan(value interface{}) error {
	*ts = ReportType(value.(string))
	return nil
}

func (ts ReportType) Value() (driver.Value, error) {
	return string(ts), nil
}
