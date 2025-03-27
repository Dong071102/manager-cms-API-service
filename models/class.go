package models

import (
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ClassID       uuid.UUID `json:"class_id" gorm:"type:uuid;primaryKey"`
	ClassName     string    `json:"class_name"`
	LecturerID    uuid.UUID `json:"lecturer_id"` // Đổi thành uuid.UUID để khớp với dữ liệu giảng viên
	CreatedAt     time.Time `json:"created_at"`
	CurrentLesson int       `json:"current_lessons"`                                       // Số bài học hiện tại
	CourseID      uuid.UUID `json:"course_id"`                                             // Đổi thành uuid.UUID
	Course        Course    `json:"course" gorm:"foreignKey:CourseID;references:CourseID"` // Quan hệ với Course
}
