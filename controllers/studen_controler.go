package controllers

import (
	"cms-backend/config"
	"cms-backend/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UpdateStudent(c echo.Context) error {
	type UpdateStudentRequest struct {
		ClassID     uuid.UUID `json:"classId"`
		StudentCode string    `json:"studentCode"`
		LastName    string    `json:"lastName"`
		FirstName   string    `json:"firstName"`
		Status      string    `json:"status"`
	}

	studentID := c.Param("id")

	var req UpdateStudentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	// Cập nhật thông tin trong bảng `users`
	if err := config.DB.Model(&models.User{}).
		Where("user_id = ?", studentID).
		Updates(map[string]interface{}{
			"first_name": req.FirstName,
			"last_name":  req.LastName,
		}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	// Cập nhật thông tin trong bảng `students`
	if err := config.DB.Model(&models.Student{}).
		Where("student_id = ?", studentID).
		Update("student_code", req.StudentCode).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to update student",
			"error":   err.Error(),
		})
	}

	// Cập nhật thông tin trong bảng `class_students`
	if err := config.DB.Model(&models.ClassStudent{}).
		Where("student_id = ? AND class_id = ?", studentID, req.ClassID).
		Update("status", req.Status).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to update class student status",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Cập nhật sinh viên thành công",
	})
}

func DeleteStudentFromClass(c echo.Context) error {
	studentID := c.Param("student_id")
	classID := c.Param("class_id")

	// Xóa sinh viên khỏi lớp trong bảng `class_students`
	if err := config.DB.Where("student_id = ? AND class_id = ?", studentID, classID).Delete(&models.ClassStudent{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to delete student from class",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sinh viên đã được xoá khỏi lớp thành công",
	})
}

// Kiểm tra sự tồn tại của sinh viên trong cơ sở dữ liệu
func CheckStudentExistence(c echo.Context) error {
	studentCode := c.Param("studentCode")

	if studentCode == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Mã sinh viên không được để trống",
		})
	}

	var student models.Student
	if err := config.DB.Where("student_code = ?", studentCode).First(&student).Error; err != nil {
		// Nếu không tìm thấy sinh viên, trả về false
		return c.JSON(http.StatusOK, echo.Map{
			"exists": false,
		})
	}

	// Nếu tìm thấy sinh viên, trả về true
	return c.JSON(http.StatusOK, echo.Map{
		"exists": true,
	})
}
func GetStudentAttendances(c echo.Context) error {
	type AttendanceResponse struct {
		AttendanceID     string `json:"attendance_id"`
		StudentCode      string `json:"student_code"`
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		ClassName        string `json:"class_name"`
		CourseName       string `json:"course_name"`
		StartTime        string `json:"start_time"`
		AttendanceTime   string `json:"attendance_time"`
		Status           string `json:"status"`
		Note             string `json:"note"`
		EvidenceImageURL string `json:"evidence_image_url"`
	}
	// Lấy student_id và lecturer_id từ URL parameter
	studentID := c.Param("student_id")
	lecturerID := c.Param("lecturer_id")
	// Lấy class_id từ query parameter (tùy chọn)
	classID := c.QueryParam("class_id")

	if studentID == "" || lecturerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "student_id and lecturer_id are required"})
	}

	var results []AttendanceResponse

	// Truy vấn cơ bản
	query := `
		SELECT 
			a.attendance_id,
			s.student_code,
			u.first_name,
			u.last_name,
			courses.course_name,
			c.class_name,
			sc.start_time,
			a.attendance_time,
			a.status  "status",
			a.note,
			a.evidence_image_url
		FROM attendance a
		JOIN schedules sc ON a.schedule_id = sc.schedule_id
		JOIN students s ON s.student_id = a.student_id
		JOIN users u ON u.user_id = s.student_id
		JOIN classes c ON sc.class_id = c.class_id
		join courses on courses.course_id=c.course_id
		WHERE a.student_id = ? AND c.lecturer_id = ?
	`

	// Nếu có class_id, thêm điều kiện
	if classID != "" {
		query += " AND c.class_id = ?"
		if err := config.DB.Raw(query, studentID, lecturerID, classID).Scan(&results).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch data"})
		}
	} else {
		if err := config.DB.Raw(query, studentID, lecturerID).Scan(&results).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch data"})
		}
	}

	return c.JSON(http.StatusOK, results)
}
