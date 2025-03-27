package controllers

import (
	"cms-backend/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetClassesByLecturer(c echo.Context) error {
	lecturerID := c.Param("lecturer_id") // Lấy lecturer_id từ tham số URL 

	var classes []struct {
		ClassID    string `json:"class_id"`
		ClassName  string `json:"class_name"`
		CourseID   string `json:"course_id"`
		CourseName string `json:"course_name"`
	}

	// Thực hiện truy vấn SQL trực tiếp mà không dùng Preload tiếp 
	
	err := config.DB.Table("classes cl").
		Select("cl.class_id, cl.class_name, cl.course_id, cs.course_name").
		Joins("JOIN courses cs ON cs.course_id = cl.course_id").
		Where("cl.lecturer_id = ?", lecturerID).
		Scan(&classes).
		Error

	if err != nil {
		// Nếu có lỗi xảy ra khi truy vấn
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve data"})
	}

	// Trả về kết quả dưới dạng JSON
	return c.JSON(http.StatusOK, classes)
}
