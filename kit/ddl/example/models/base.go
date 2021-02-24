package models

import "time"

type Base struct {
	CreateAt *time.Time `papi:"default:CURRENT_TIMESTAMP"`
	UpdateAt *time.Time `papi:"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP"`
	DeleteAt *time.Time
}
