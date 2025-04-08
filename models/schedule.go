// models/course.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ScheduleID  uuid.UUID `json:"schedule_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ClassID     uuid.UUID `json:"class_id" gorm:"type:uuid"`
	ClassroomID uuid.UUID `json:"classroom_id" gorm:"type:uuid"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Topic       string    `json:"topic"`
	Description string    `json:"description"`
}
