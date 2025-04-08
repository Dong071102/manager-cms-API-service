package models

import (
	"github.com/google/uuid"
)

type ClassStudent struct {
	ClassID   uuid.UUID `gorm:"primaryKey" json:"class_id"`
	StudentID uuid.UUID `gorm:"primaryKey" json:"student_id"`
	Status    string    `json:"status"`
}
