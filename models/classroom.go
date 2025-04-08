// models/course.go
package models

import (
	"github.com/google/uuid"
)

type Classroom struct {
	ClassroomID uuid.UUID `gorm:"column:classroom_id;primaryKey" json:"classroom_id"`
	RoomName    string    `gorm:"column:room_name" json:"room_name"`
	RoomType    string    `gorm:"column:room_type" json:"room_type"`
	Location    string    `gorm:"column:location" json:"location"`
	Description string    `gorm:"column:description" json:"description"`
}
