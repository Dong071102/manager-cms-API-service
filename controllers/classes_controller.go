package controllers

import (
	"cms-backend/config"
	"net/http"
	"time"

	"github.com/google/uuid"

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

func GetNearestClassByLecturer(c echo.Context) error {
	lecturerID := c.Param("lecturer_id") // Lấy lecturer_id từ tham số URL

	var nearesClasses []struct {
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
		Scan(&nearesClasses).
		Error

	if err != nil {
		// Nếu có lỗi xảy ra khi truy vấn
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve data"})
	}

	// Trả về kết quả dưới dạng JSON
	return c.JSON(http.StatusOK, nearesClasses)
}

func AttendanceSummaryHandler(c echo.Context) error {
	lecturerID := c.QueryParam("lecturer_id")
	if lecturerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Thiếu tham số lecturer_id"})
	}

	classID := c.QueryParam("class_id")

	// Định nghĩa struct cho kết quả trả về
	type AttendanceSummary struct {
		TotalStudents int `json:"total_students"`
		CountAbsent   int `json:"count_absent"`
		CountPresent  int `json:"count_present"`
		CountLate     int `json:"count_late"`
	}

	var summary AttendanceSummary

	// Xây dựng truy vấn với GORM
	query := config.DB.Table("attendance a").
		Select(`
			COUNT(DISTINCT a.student_id) AS total_students,
			SUM(CASE WHEN a.status = 'absent' THEN 1 ELSE 0 END) AS count_absent,
			SUM(CASE WHEN a.status = 'present' THEN 1 ELSE 0 END) AS count_present,
			SUM(CASE WHEN a.status = 'late' THEN 1 ELSE 0 END) AS count_late
		`).
		Joins("JOIN schedules cs ON a.schedule_id = cs.schedule_id").
		Joins("JOIN classes c ON cs.class_id = c.class_id").
		Where("c.lecturer_id = ?", lecturerID)

	// Nếu có class_id, thêm điều kiện lọc
	if classID != "" {
		query = query.Where("c.class_id = ?", classID)
	}

	if err := query.Scan(&summary).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve attendance summary"})
	}

	return c.JSON(http.StatusOK, summary)
}
func GetAttendanceDetails(c echo.Context) error {
	// Lấy lecturer_id từ tham số URL
	lecturerID := c.QueryParam("lecturer_id")
	if lecturerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Thiếu lecturer_id"})
	}

	// Lấy class_id từ query string (có thể không có)
	classID := c.QueryParam("class_id")
	// Định nghĩa cấu trúc dữ liệu cho bản ghi trả về
	var records []struct {
		StudentID      uuid.UUID `json:"student_id"`
		StudentCode    string    `json:"student_code"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		ClassID        uuid.UUID `json:"class_id"`
		ClassName      string    `json:"class_name"`
		CourseID       uuid.UUID `json:"course_id"`
		CourseName     string    `json:"course_name"`
		Status         string    `json:"status"`
		AttendanceTime time.Time `json:"attendance_time"`
		StartTime      time.Time `json:"start_time"`
		Note           string    `json:"note"`
	}

	// Xây dựng query cơ bản với điều kiện lecturer_id
	dbQuery := config.DB.Table("attendance a").
		Select(`a.student_id, st.student_code, u.first_name, u.last_name,
		        c.class_id, c.class_name, cs.course_id, cs.course_name,
		        a.status, a.attendance_time, s.start_time, a.note`).
		Joins("JOIN schedules s ON a.schedule_id = s.schedule_id").
		Joins("JOIN students st ON st.student_id = a.student_id").
		Joins("JOIN users u ON u.user_id = a.student_id").
		Joins("JOIN classes c ON s.class_id = c.class_id").
		Joins("JOIN courses cs ON cs.course_id = c.course_id").
		Where("c.lecturer_id = ?", lecturerID)

	// Nếu có class_id thì thêm điều kiện
	if classID != "" {
		dbQuery = dbQuery.Where("c.class_id = ?", classID)
	}

	// Thực hiện query và sắp xếp theo s.start_time giảm dần
	if err := dbQuery.Order("s.start_time DESC").Scan(&records).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve data"})
	}

	// Trả về kết quả dạng JSON
	return c.JSON(http.StatusOK, records)
}
