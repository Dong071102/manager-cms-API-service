package controllers

import (
	"cms-backend/config"
	"cms-backend/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ScheduleResult struct {
	ScheduleID  string `json:"schedule_id"`
	ClassID     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	ClassroomID string `json:"classroom_id"`
	LecturerID  string `json:"lecturer_id"`
	RoomName    string `json:"room_name"`
	CourseID    string `json:"course_id"`
	CourseName  string `json:"course_name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
}

func GetSchedules(c echo.Context) error {

	classroomId := c.QueryParam("classroom_id")
	lecturerID := c.QueryParam("lecturer_id")
	weekStartStr := c.QueryParam("week_start")

	var results []ScheduleResult

	query := config.DB.Table("schedules AS s").
		Select(`s.schedule_id, s.class_id, c.class_name, 
		        s.classroom_id, cs.room_name, 
		        c.course_id, cr.course_name, c.lecturer_id,
		        s.start_time, s.end_time, s.topic, s.description`).
		Joins("JOIN classes c ON c.class_id = s.class_id").
		Joins("JOIN classrooms cs ON cs.classroom_id = s.classroom_id").
		Joins("JOIN courses cr ON cr.course_id = c.course_id")

	// --- Filter theo room_name nếu có ---
	if classroomId != "" {
		query = query.Where("cs.classroom_id = ?", classroomId)
	}

	// --- Filter theo lecturer_id nếu có ---
	if lecturerID != "" {
		query = query.Where("c.lecturer_id = ?", lecturerID)
	}

	// --- Filter theo tuần nếu có ---
	if weekStartStr != "" {
		weekStart, err := time.Parse("2006-01-02", weekStartStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid week_start format (must be YYYY-MM-DD)"})
		}
		weekEnd := weekStart.AddDate(0, 0, 7) // cộng 7 ngày
		query = query.Where("s.start_time >= ? AND s.start_time < ?", weekStart, weekEnd)
	}

	if err := query.Find(&results).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

func AddSubject(c echo.Context) error {
	var input models.Schedule

	// Parse JSON -> struct
	if err := c.Bind(&input); err != nil {
		log.Println("[Bind error]", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	log.Printf("Schedule input: %+v\n", input)
	// Gán ScheduleID nếu client không truyền
	if input.ScheduleID == uuid.Nil {
		input.ScheduleID = uuid.New()
		fmt.Println("📦 New schedule:", input.ScheduleID)

	}

	// Validate thời gian (optional)
	if input.StartTime.IsZero() || input.EndTime.IsZero() || input.EndTime.Before(input.StartTime) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid start/end time"})
	}

	// Lưu vào DB
	if err := config.DB.Create(&input).Error; err != nil {
		log.Println("Bind error:", err) // 👈 log ra lỗi thực sự
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, input)
}

func UpdateSubject(c echo.Context) error {
	idStr := c.Param("id")
	scheduleID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid UUID"})
	}

	var input models.Schedule
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	var existing models.Schedule
	if err := config.DB.First(&existing, "schedule_id = ?", scheduleID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Schedule not found"})
	}

	// Cập nhật các trường
	existing.ClassID = input.ClassID
	existing.ClassroomID = input.ClassroomID
	existing.StartTime = input.StartTime
	existing.EndTime = input.EndTime
	existing.Topic = input.Topic
	existing.Description = input.Description

	if err := config.DB.Save(&existing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, existing)
}
func DeleteSchedule(c echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Missing schedule ID",
		})
	}

	scheduleID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid schedule ID format",
		})
	}

	var schedule models.Schedule
	if err := config.DB.First(&schedule, "schedule_id = ?", scheduleID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Schedule not found",
		})
	}

	if err := config.DB.Delete(&schedule).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to delete schedule",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Schedule deleted successfully",
	})
}

func GetScheduleStartTimes(c echo.Context) error {
	classID := c.QueryParam("class_id")
	lecturerID := c.QueryParam("lecturer_id")
	scheduleID := c.QueryParam("schedule_id")

	var results []struct {
		StartTime  time.Time `json:"start_time"`
		ScheduleID uuid.UUID `json:"schedule_id"`
	}

	query := config.DB.
		Table("schedules AS s").
		Select("s.start_time,s.schedule_id").
		Joins("JOIN classes c ON c.class_id = s.class_id")

	if classID != "" {
		query = query.Where("c.class_id = ?", classID)
	}
	if lecturerID != "" {
		query = query.Where("c.lecturer_id = ?", lecturerID)
	}
	if scheduleID != "" {

	}
	if err := query.Order("s.start_time DESC").Find(&results).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

func GetScheduTimes(c echo.Context) error {
	scheduleID := c.QueryParam("schedule_id")

	if scheduleID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "schedule_id is required",
		})
	}

	var result struct {
		ScheduleID uuid.UUID `json:"schedule_id"`
		StartTime  time.Time `json:"start_time"`
		EndTime    time.Time `json:"end_time"`
	}

	query := config.DB.
		Table("schedules AS s").
		Select("s.schedule_id, s.start_time, s.end_time").
		Where("s.schedule_id = ?", scheduleID)

	if err := query.First(&result).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Schedule not found",
		})
	}

	return c.JSON(http.StatusOK, result)
}
