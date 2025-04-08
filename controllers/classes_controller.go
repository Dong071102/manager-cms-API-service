package controllers

import (
	"cms-backend/config"
	"cms-backend/models"
	"log"
	"net/http"
	"strconv"
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
	scheduleID := c.QueryParam("schedule_id")
	// Định nghĩa cấu trúc dữ liệu cho bản ghi trả về
	var records []struct {
		AttendanceId   uuid.UUID `json:"attendance_id"`
		StudentID      uuid.UUID `json:"student_id"`
		StudentCode    string    `json:"student_code"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		ScheduleId     uuid.UUID `json:"schedule_id"`
		ClassID        uuid.UUID `json:"class_id"`
		ClassName      string    `json:"class_name"`
		CourseID       uuid.UUID `json:"course_id"`
		CourseName     string    `json:"course_name"`
		Status         string    `json:"status"`
		AttendanceTime time.Time `json:"attendance_time"`
		StartTime      time.Time `json:"start_time"`
		Note           string    `json:"note"`
		// ImageUrl         string    `json:"evidence_image_url"`
		EvidenceImageUrl string `json:"evidence_image_url"`
	}

	// Xây dựng query cơ bản với điều kiện lecturer_id
	dbQuery := config.DB.Table("attendance a").
		Select(`a.attendance_id,a.student_id, st.student_code, u.first_name, u.last_name,a.schedule_id,
		        c.class_id, c.class_name, cs.course_id, cs.course_name,a.attendance_time,
		        a.status, s.start_time, s.start_time, a.note,a.evidence_image_url`).
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
	if scheduleID != "" {
		dbQuery = dbQuery.Where("a.schedule_id = ?", scheduleID)
	}

	// Thực hiện query và sắp xếp theo s.start_time giảm dần
	if err := dbQuery.Order("s.start_time DESC").Scan(&records).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve data"})
	}

	// Trả về kết quả dạng JSON
	return c.JSON(http.StatusOK, records)
}

func UpdateAttendance(c echo.Context) error {
	type AttendanceUpdate struct {
		AttendanceID     uuid.UUID `json:"attendance_id"`      // Khóa chính để cập nhật
		StudentID        uuid.UUID `json:"student_id"`         // UUID của học sinh
		ScheduleID       uuid.UUID `json:"schedule_id"`        // UUID của lịch học
		AttendanceTime   string    `json:"attendance_time"`    // Dạng string theo RFC3339, ví dụ "2025-03-29T14:50:43.590Z"
		Status           string    `json:"status"`             // "present", "absent", "late", ...
		EvidenceImageURL *string   `json:"evidence_image_url"` // Có thể là null
		Note             *string   `json:"note"`               // Có thể là null
	}
	var att AttendanceUpdate
	if err := c.Bind(&att); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Parse trường attendance_time sang kiểu time.Time, nếu không có hoặc lỗi thì lấy thời gian hiện tại
	var attTime time.Time
	if att.AttendanceTime != "" {
		t, err := time.Parse(time.RFC3339, att.AttendanceTime)
		if err != nil {
			log.Printf("Error parsing attendance_time: %v", err)
			attTime = time.Now().UTC()
		} else {
			attTime = t
		}
	} else {
		attTime = time.Now().UTC()
	}

	// Câu lệnh SQL cập nhật bản ghi dựa trên attendance_id
	query := `
		UPDATE attendance
		SET schedule_id = $1,
			attendance_time = $2,
			student_id = $3,
			status = $4,
			evidence_image_url = $5,
			note = $6
		WHERE attendance_id = $7
	`
	err := config.DB.Exec(query,
		att.ScheduleID,
		attTime,
		att.StudentID,
		att.Status,
		att.EvidenceImageURL,
		att.Note,
		att.AttendanceID,
	).Error

	if err != nil {
		log.Printf("Error updating attendance: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update attendance"})
	}

	// Trả về phản hồi thành công
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Attendance updated successfully",
	})
}
func GetAttendanceReport(c echo.Context) error {
	lecturerID := c.Param("lecturer_id")
	filter := c.QueryParam("filter") // "week", "month", or "year"
	year := c.QueryParam("year")
	month := c.QueryParam("month")
	week := c.QueryParam("week") // dùng khi filter là "week"
	classID := c.QueryParam("class_id")

	if filter == "" {
		filter = "month" // mặc định
	}

	type AttendanceReport struct {
		Period  string `json:"period"` // week: ngày, month: tuần, year: tháng
		Present int    `json:"present"`
		Late    int    `json:"late"`
		Absent  int    `json:"absent"`
	}

	var reports []AttendanceReport
	var query string
	var params []interface{}

	switch filter {
	case "year":
		// Báo cáo theo năm: trả về đầy đủ 12 tháng
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid year"})
		}
		// Tính ngày bắt đầu của năm (đầu năm) theo múi giờ địa phương
		startOfYear := time.Date(yearInt, time.January, 1, 0, 0, 0, 0, time.Local)
		// Cách tính endOfYear: thêm 1 năm vào startOfYear
		// (generate_series sẽ tạo ra các mốc tháng từ startOfYear đến endOfYear - 1 day)
		endOfYear := startOfYear.AddDate(1, 0, 0)

		// Truy vấn: Sử dụng CTE để tạo ra dãy 12 tháng, sau đó LEFT JOIN với dữ liệu điểm danh đã được tổng hợp
		if classID == "" {
			query = `
			WITH months AS (
				SELECT generate_series(?::timestamp, ?::timestamp - interval '1 day', interval '1 month') AS month_start
			),
			data AS (
				SELECT s.start_time, a.status
				FROM schedules s
				JOIN attendance a ON a.schedule_id = s.schedule_id
				JOIN classes c ON s.class_id = c.class_id
				WHERE c.lecturer_id = ?
			)
			SELECT 
				TO_CHAR(m.month_start, 'MM') AS period,
				COALESCE(SUM(CASE WHEN d.status = 'present' THEN 1 ELSE 0 END), 0) AS present,
				COALESCE(SUM(CASE WHEN d.status = 'late' THEN 1 ELSE 0 END), 0) AS late,
				COALESCE(SUM(CASE WHEN d.status = 'absent' THEN 1 ELSE 0 END), 0) AS absent
			FROM months m
			LEFT JOIN data d ON d.start_time >= m.month_start
				AND d.start_time < m.month_start + interval '1 month'
			GROUP BY m.month_start
			ORDER BY m.month_start;
		`
			// Tham số: startOfYear, endOfYear, lecturerID
			params = []interface{}{startOfYear, endOfYear, lecturerID}
		} else {
			query = `
			WITH months AS (
				SELECT generate_series(?::timestamp, ?::timestamp - interval '1 day', interval '1 month') AS month_start
			),
			data AS (
				SELECT s.start_time, a.status
				FROM schedules s
				JOIN attendance a ON a.schedule_id = s.schedule_id
				JOIN classes c ON s.class_id = c.class_id
				WHERE c.lecturer_id = ? AND c.class_id = ?
			)
			SELECT 
				TO_CHAR(m.month_start, 'MM') AS period,
				COALESCE(SUM(CASE WHEN d.status = 'present' THEN 1 ELSE 0 END), 0) AS present,
				COALESCE(SUM(CASE WHEN d.status = 'late' THEN 1 ELSE 0 END), 0) AS late,
				COALESCE(SUM(CASE WHEN d.status = 'absent' THEN 1 ELSE 0 END), 0) AS absent
			FROM months m
			LEFT JOIN data d ON d.start_time >= m.month_start
				AND d.start_time < m.month_start + interval '1 month'
			GROUP BY m.month_start
			ORDER BY m.month_start;
		`
			// Tham số: startOfYear, endOfYear, lecturerID, classID
			params = []interface{}{startOfYear, endOfYear, lecturerID, classID}
		}

	case "month":
		// Báo cáo theo tháng: chia thành 4 khoảng (tuần)
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid year"})
		}
		monthInt, err := strconv.Atoi(month)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid month"})
		}
		baseDate := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.Local)
		week1Start := baseDate
		week1End := baseDate.AddDate(0, 0, 7)
		week2Start := week1End
		week2End := baseDate.AddDate(0, 0, 14)
		week3Start := week2End
		week3End := baseDate.AddDate(0, 0, 21)
		week4Start := week3End
		week4End := baseDate.AddDate(0, 1, 0) // ngày đầu tháng kế tiếp

		baseCondition := "c.lecturer_id = ?"
		if classID != "" {
			baseCondition += " AND c.class_id = ?"
		}
		query = `
			SELECT 'Tuần 1' AS period,
				COALESCE(present,0) AS present,
				COALESCE(late,0) AS late,
				COALESCE(absent,0) AS absent
			FROM (
				SELECT COUNT(*) FILTER (WHERE a.status = 'present') AS present,
					   COUNT(*) FILTER (WHERE a.status = 'late') AS late,
					   COUNT(*) FILTER (WHERE a.status = 'absent') AS absent
				FROM attendance a
				JOIN schedules s ON a.schedule_id = s.schedule_id
				JOIN classes c ON s.class_id = c.class_id
				WHERE ` + baseCondition + ` AND s.start_time >= ? AND s.start_time < ?
			) sub
			UNION ALL
			SELECT 'Tuần 2' AS period,
				COALESCE(present,0),
				COALESCE(late,0),
				COALESCE(absent,0)
			FROM (
				SELECT COUNT(*) FILTER (WHERE a.status = 'present') AS present,
					   COUNT(*) FILTER (WHERE a.status = 'late') AS late,
					   COUNT(*) FILTER (WHERE a.status = 'absent') AS absent
				FROM attendance a
				JOIN schedules s ON a.schedule_id = s.schedule_id
				JOIN classes c ON s.class_id = c.class_id
				WHERE ` + baseCondition + ` AND s.start_time >= ? AND s.start_time < ?
			) sub
			UNION ALL
			SELECT 'Tuần 3' AS period,
				COALESCE(present,0),
				COALESCE(late,0),
				COALESCE(absent,0)
			FROM (
				SELECT COUNT(*) FILTER (WHERE a.status = 'present') AS present,
					   COUNT(*) FILTER (WHERE a.status = 'late') AS late,
					   COUNT(*) FILTER (WHERE a.status = 'absent') AS absent
				FROM attendance a
				JOIN schedules s ON a.schedule_id = s.schedule_id
				JOIN classes c ON s.class_id = c.class_id
				WHERE ` + baseCondition + ` AND s.start_time >= ? AND s.start_time < ?
			) sub
			UNION ALL
			SELECT 'Tuần 4' AS period,
				COALESCE(present,0),
				COALESCE(late,0),
				COALESCE(absent,0)
			FROM (
				SELECT COUNT(*) FILTER (WHERE a.status = 'present') AS present,
					   COUNT(*) FILTER (WHERE a.status = 'late') AS late,
					   COUNT(*) FILTER (WHERE a.status = 'absent') AS absent
				FROM attendance a
				JOIN schedules s ON a.schedule_id = s.schedule_id
				JOIN classes c ON s.class_id = c.class_id
				WHERE ` + baseCondition + ` AND s.start_time >= ? AND s.start_time < ?
			) sub;
		`
		params = []interface{}{}
		// Tuần 1
		params = append(params, lecturerID)
		if classID != "" {
			params = append(params, classID)
		}
		params = append(params, week1Start, week1End)
		// Tuần 2
		params = append(params, lecturerID)
		if classID != "" {
			params = append(params, classID)
		}
		params = append(params, week2Start, week2End)
		// Tuần 3
		params = append(params, lecturerID)
		if classID != "" {
			params = append(params, classID)
		}
		params = append(params, week3Start, week3End)
		// Tuần 4
		params = append(params, lecturerID)
		if classID != "" {
			params = append(params, classID)
		}
		params = append(params, week4Start, week4End)
	case "week":
		// Báo cáo theo tuần: trả về 7 ngày của tuần được chọn sử dụng generate_series
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid year"})
		}
		monthInt, err := strconv.Atoi(month)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid month"})
		}
		weekInt, err := strconv.Atoi(week)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid week"})
		}
		// Tính ngày đầu tiên của tháng theo múi giờ địa phương
		firstOfMonth := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.Local)
		// Tính ngày bắt đầu của tuần thứ weekInt trong tháng
		startDate := firstOfMonth.AddDate(0, 0, (weekInt-1)*7)
		// Tính endDate: cộng thêm 7 ngày, để bao gồm đủ 7 ngày của tuần
		endDate := startDate.AddDate(0, 0, 7)

		// Sử dụng generate_series để tạo dãy 7 ngày
		// Đưa các điều kiện của lecturer (và class nếu có) vào trong phần LEFT JOIN của bảng classes,
		// để không loại bỏ các ngày không có dữ liệu (NULL) từ generate_series.
		if classID == "" {
			query = `
			WITH days AS (
				SELECT generate_series(?::timestamp, ?::timestamp - interval '1 second', interval '1 day') AS day_date
			)
			SELECT 
				TO_CHAR(d.day_date, 'DD') AS period,
				COALESCE(COUNT(*) FILTER (WHERE a.status = 'present'), 0) AS present,
				COALESCE(COUNT(*) FILTER (WHERE a.status = 'late'), 0) AS late,
				COALESCE(COUNT(*) FILTER (WHERE a.status = 'absent'), 0) AS absent
			FROM days d
			LEFT JOIN schedules s 
				ON s.start_time >= d.day_date 
			   AND s.start_time < d.day_date + interval '1 day'
			LEFT JOIN attendance a 
				ON a.schedule_id = s.schedule_id
			LEFT JOIN classes c 
				ON s.class_id = c.class_id AND c.lecturer_id = ?
			WHERE d.day_date >= ? AND d.day_date < ?
			GROUP BY d.day_date
			ORDER BY d.day_date;
		`
			// Thứ tự tham số: startDate, endDate, lecturerID, startDate, endDate
			params = []interface{}{startDate, endDate, lecturerID, startDate, endDate}
		} else {
			query = `
			WITH days AS (
				SELECT generate_series(?::timestamp, ?::timestamp - interval '1 second', interval '1 day') AS day_date
			)
			SELECT 
				TO_CHAR(d.day_date, 'DD') AS period,
				COALESCE(COUNT(*) FILTER (WHERE a.status = 'present'), 0) AS present,
				COALESCE(COUNT(*) FILTER (WHERE a.status = 'late'), 0) AS late,
				COALESCE(COUNT(*) FILTER (WHERE a.status = 'absent'), 0) AS absent
			FROM days d
			LEFT JOIN schedules s 
				ON s.start_time >= d.day_date 
			   AND s.start_time < d.day_date + interval '1 day'
			LEFT JOIN attendance a 
				ON a.schedule_id = s.schedule_id
			LEFT JOIN classes c 
				ON s.class_id = c.class_id AND c.lecturer_id = ? AND c.class_id = ?
			WHERE d.day_date >= ? AND d.day_date < ?
			GROUP BY d.day_date
			ORDER BY d.day_date;
		`
			// Thứ tự tham số: startDate, endDate, lecturerID, classID, startDate, endDate
			params = []interface{}{startDate, endDate, lecturerID, classID, startDate, endDate}
		}

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid filter"})
	}

	//====================
	err := config.DB.Raw(query, params...).Scan(&reports).Error
	if err != nil {
		c.Logger().Error("Error executing query:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve data"})
	}

	return c.JSON(http.StatusOK, reports)
}

func GetStudentsInClass(c echo.Context) error {
	lecturerID := c.Param("lecturer_id")
	if lecturerID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Thiếu lecturer_id",
		})
	}

	type StudentInClassResponse struct {
		ClassID     uuid.UUID `json:"class_id"`
		ClassName   string    `json:"class_name"`
		StudentID   uuid.UUID `json:"student_id"`
		StudentCode string    `json:"student_code"`
		FirstName   string    `json:"first_name"`
		LastName    string    `json:"last_name"`
		Status      string    `json:"status"`
	}

	classID := c.QueryParam("class_id")

	var students []StudentInClassResponse

	query := config.DB.Table("class_students AS cs").
		Select("cs.class_id, c.class_name, cs.student_id, s.student_code, u.first_name, u.last_name, cs.status").
		Joins("JOIN students s ON s.student_id = cs.student_id").
		Joins("JOIN users u ON u.user_id = s.student_id").
		Joins("JOIN classes c ON c.class_id = cs.class_id").
		Where("c.lecturer_id = ?", lecturerID)

	if classID != "" && classID != "0" {
		query = query.Where("cs.class_id = ?", classID)
	}

	if err := query.Scan(&students).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "Không thể truy vấn danh sách sinh viên",
		})
	}

	return c.JSON(http.StatusOK, students)
}
func AddStudentToClass(c echo.Context) error {
	// Định nghĩa kiểu yêu cầu
	type AddStudentRequest struct {
		StudentCode string    `json:"studentCode"`
		ClassID     uuid.UUID `json:"classId"`
		Status      string    `json:"status"`
	}

	// Lấy dữ liệu từ yêu cầu
	var req []AddStudentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	// Danh sách các sinh viên đã được thêm thành công
	var successfullyAdded []string

	// Lặp qua từng sinh viên trong yêu cầu
	for _, student := range req {
		// Kiểm tra xem sinh viên có tồn tại trong bảng `students` không
		var studentRecord models.Student
		if err := config.DB.Where("student_code = ?", student.StudentCode).First(&studentRecord).Error; err != nil {
			// Nếu sinh viên không tồn tại, bỏ qua sinh viên đó
			continue
		}

		// Kiểm tra xem sinh viên có đã có trong lớp này chưa
		var existingEntry models.ClassStudent
		if err := config.DB.Where("student_id = ? AND class_id = ?", studentRecord.StudentID, student.ClassID).First(&existingEntry).Error; err == nil {
			// Nếu sinh viên đã có trong lớp, chỉ cần cập nhật trạng thái
			if err := config.DB.Model(&models.ClassStudent{}).
				Where("student_id = ? AND class_id = ?", studentRecord.StudentID, student.ClassID).
				Update("status", student.Status).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "Không thể cập nhật trạng thái sinh viên",
					"error":   err.Error(),
				})
			}
			successfullyAdded = append(successfullyAdded, student.StudentCode) // Ghi nhận sinh viên đã được cập nhật
			continue
		}

		// Nếu sinh viên chưa có trong lớp, thêm mới sinh viên vào bảng `class_students`
		newEntry := models.ClassStudent{
			ClassID:   student.ClassID,
			StudentID: studentRecord.StudentID,
			Status:    student.Status,
		}

		if err := config.DB.Create(&newEntry).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Không thể thêm sinh viên vào lớp",
				"error":   err.Error(),
			})
		}

		// Ghi nhận sinh viên đã thêm thành công
		successfullyAdded = append(successfullyAdded, student.StudentCode)
	}

	// Trả về danh sách sinh viên đã được thêm hoặc cập nhật
	return c.JSON(http.StatusOK, echo.Map{
		"message":        "Sinh viên đã được thêm hoặc cập nhật trạng thái thành công",
		"added_students": successfullyAdded,
	})
}

// API endpoint to get attendance summary by lecturer and class_id
func GetStudentAttendanceSummary(c echo.Context) error {
	// Lấy lecturer_id từ URL parameter
	lecturerID := c.Param("lecturer_id")
	// Lấy class_id và course_id từ query parameters
	classID := c.QueryParam("class_id")
	courseId := c.QueryParam("course_id")

	if lecturerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lecturer ID is required"})
	}

	// Truy vấn chính: lấy dữ liệu điểm danh nếu có
	query := `
		SELECT 
			s.student_id AS "studentId",
			s.student_code AS "studentCode",
			u.first_name || ' ' || u.last_name AS "fullName",
			COUNT(a.attendance_id) AS "attendanceDays",
			COALESCE(SUM(CASE WHEN a.status = 'present' THEN 1 ELSE 0 END), 0) AS "presentDays",  
			COALESCE(SUM(CASE WHEN a.status = 'late' THEN 1 ELSE 0 END), 0) AS "lateDays",  
			COALESCE(SUM(CASE WHEN a.status = 'absent' THEN 1 ELSE 0 END), 0) AS "absentDays"  
		FROM 
			students s
		JOIN 
			attendance a ON s.student_id = a.student_id
		JOIN 
			users u ON s.student_id = u.user_id
		JOIN 
			schedules sc ON sc.schedule_id = a.schedule_id
		JOIN 
			classes c ON c.class_id = sc.class_id
		WHERE 
			c.lecturer_id = ?`
	// Nếu có điều kiện bổ sung về class_id và course_id thì thêm
	if classID != "" && courseId != "" {
		query += ` AND c.class_id = ? AND c.course_id = ?`
	} else if classID != "" {
		query += ` AND c.class_id = ?`
	} else if courseId != "" {
		query += ` AND c.course_id = ?`
	}
	query += ` GROUP BY s.student_id, s.student_code, u.first_name, u.last_name ORDER BY s.student_id`

	var results []map[string]interface{}
	var err error

	if classID != "" && courseId != "" {
		err = config.DB.Raw(query, lecturerID, classID, courseId).Scan(&results).Error
	} else if classID != "" {
		err = config.DB.Raw(query, lecturerID, classID).Scan(&results).Error
	} else if courseId != "" {
		err = config.DB.Raw(query, lecturerID, courseId).Scan(&results).Error
	} else {
		err = config.DB.Raw(query, lecturerID).Scan(&results).Error
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch data"})
	}

	// Nếu truy vấn chính không trả về dữ liệu (tức là không có thông tin attendance/schedule),
	// thực hiện fallback query để lấy danh sách học sinh từ bảng ánh xạ (class_students)
	if len(results) == 0 {
		fallbackQuery := `
		SELECT 
			s.student_id AS "studentId",
			s.student_code AS "studentCode",
			u.first_name || ' ' || u.last_name AS "fullName",
			0 AS "attendanceDays",
			0 AS "presentDays",
			0 AS "lateDays",
			0 AS "absentDays"
		FROM 
			class_students cs
		JOIN 
			students s ON cs.student_id = s.student_id
		JOIN 
			users u ON s.student_id = u.user_id
		JOIN 
			classes c ON cs.class_id = c.class_id
		WHERE 
			cs.class_id = ? 
			AND c.lecturer_id = ?`
		// Nếu có courseId, thêm điều kiện cho course
		if courseId != "" {
			fallbackQuery += ` AND c.course_id = ?`
			err = config.DB.Raw(fallbackQuery, classID, lecturerID, courseId).Scan(&results).Error
		} else {
			err = config.DB.Raw(fallbackQuery, classID, lecturerID).Scan(&results).Error
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch fallback data"})
		}
	}

	return c.JSON(http.StatusOK, results)
}

func GetClassesByCourses(c echo.Context) error {
	type ClassResponse struct {
		ClassID       uuid.UUID `json:"class_id"`
		ClassName     string    `json:"class_name"`
		LecturerID    uuid.UUID `json:"lecturer_id"`
		CreatedAt     time.Time `json:"created_at"`
		CurrentLesson int       `json:"current_lessons"`
		CourseID      uuid.UUID `json:"course_id"`
	}

	// Lấy tham số course_id từ query string
	courseIDStr := c.QueryParam("course_id")
	if courseIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing course_id parameter"})
	}

	// Chuyển đổi course_id từ chuỗi sang uuid.UUID
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid course_id format"})
	}

	// Truy vấn cơ sở dữ liệu để lấy danh sách các lớp theo course_id
	var classes []models.Class
	result := config.DB.Where("course_id = ?", courseID).Find(&classes)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	// Chuyển đổi dữ liệu từ model Class sang response ClassResponse để loại bỏ trường Course
	var response []ClassResponse
	for _, class := range classes {
		response = append(response, ClassResponse{
			ClassID:       class.ClassID,
			ClassName:     class.ClassName,
			LecturerID:    class.LecturerID,
			CreatedAt:     class.CreatedAt,
			CurrentLesson: class.CurrentLesson,
			CourseID:      class.CourseID,
		})
	}

	// Trả về dữ liệu đã chuyển đổi
	return c.JSON(http.StatusOK, response)
}
