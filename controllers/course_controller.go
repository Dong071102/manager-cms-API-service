package controllers

import (
	"cms-backend/config"
	"cms-backend/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetCoursesByLecturerID(c echo.Context) error {
	lecturerIDStr := c.QueryParam("lecturer_id")
	if lecturerIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing lecturer_id parameter",
		})
	}

	lecturerID, err := uuid.Parse(lecturerIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid lecturer_id format (expect UUID)",
		})
	}

	var courses []models.Course
	result := config.DB.Where("main_lecturer_id = ?", lecturerID).Find(&courses)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	// Chuẩn bị dữ liệu trả về
	type CourseResponse struct {
		CourseID       uuid.UUID `json:"course_id"`
		CourseName     string    `json:"course_name"`
		MainLecturerID uuid.UUID `json:"main_lecturer"`
		CreatedAt      string    `json:"created_at"`
		TotalLesson    int       `json:"total_lesson"`
	}

	var response []CourseResponse
	for _, course := range courses {
		response = append(response, CourseResponse{
			CourseID:       course.CourseID,
			CourseName:     course.CourseName,
			MainLecturerID: course.MainLecturerID,
			CreatedAt:      course.CreatedAt.Format("2006-01-02"),
			TotalLesson:    course.TotalLesson,
		})
	}

	return c.JSON(http.StatusOK, response)
}
