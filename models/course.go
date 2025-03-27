// models/course.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	CourseID       uuid.UUID `json:"course_id" gorm:"type:uuid;primaryKey"`
	CourseName     string    `json:"course_name"`
	MainLecturerID uuid.UUID `json:"main_lecturer_id"` // Mã giảng viên chính
	CreatedAt      time.Time `json:"created_at"`
	TotalLesson    int       `json:"total_lesson"` // Tổng số bài giảng
	MainLecturer   Lecturer  `json:"main_lecturer" gorm:"foreignKey:MainLecturerID;references:LecturerID"`
}
