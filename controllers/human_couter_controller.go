package controllers

import (
	"cms-backend/config"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetHumanCouterSocketPath(c echo.Context) error {

	scheduleID := c.QueryParam("schedule_id")

	if scheduleID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "schedule_id is required",
		})
	}

	var result struct {
		SocketPath string `json:"socket_path"`
	}

	query := `
		SELECT c.socket_path
		FROM schedules s
		JOIN classrooms cr ON cr.classroom_id = s.classroom_id
		JOIN cameras c ON cr.classroom_id = c.classroom_id
		WHERE s.schedule_id = ? AND c.camera_type = 'surveillance'
		LIMIT 1;
	`

	err := config.DB.Raw(query, scheduleID).Scan(&result).Error
	if err != nil || result.SocketPath == "" {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Socket path not found",
		})
	}

	return c.JSON(http.StatusOK, result)
}
func GetSnapshotDetails(c echo.Context) error {
	type SnapshotResponse struct {
		SnapshotID    string    `json:"snapshot_id"`
		ScheduleID    string    `json:"schedule_id"`
		PeopleCounter int       `json:"people_counter"`
		CapturedAt    time.Time `json:"captured_at"`
		StartTime     time.Time `json:"start_time"`
		EndTime       time.Time `json:"end_time"`
		ClassID       string    `json:"class_id"`
		ClassName     string    `json:"class_name"`
		CourseID      string    `json:"course_id"`
		CourseName    string    `json:"course_name"`
		RoomName      string    `json:"room_name"`
		ImagePath     string    `json:"image_path"`
	}

	scheduleID := c.QueryParam("schedule_id")
	classID := c.QueryParam("class_id")
	lecturerID := c.QueryParam("lecturer_id")

	// Base query
	baseQuery := `
		SELECT p.snapshot_id,
			   p.schedule_id,
			   p.people_counter,
			   p.captured_at,
			   s.start_time,
			   s.end_time,
			   c.class_id,
			   c.class_name,
			   c.course_id,
			   cs.course_name,
			   p.image_path,
			   cr.room_name
		FROM people_count_snapshots p
		JOIN schedules s ON s.schedule_id = p.schedule_id
		JOIN classes c ON c.class_id = s.class_id
		JOIN courses cs ON cs.course_id = c.course_id
		JOIN classrooms cr ON cr.classroom_id = s.classroom_id
		WHERE 1 = 1
	`

	params := map[string]interface{}{}

	if lecturerID != "" {
		baseQuery += " AND c.lecturer_id = @lecturer_id"
		params["lecturer_id"] = lecturerID
	}
	if scheduleID != "" {
		baseQuery += " AND p.schedule_id = @schedule_id"
		params["schedule_id"] = scheduleID
	}
	if classID != "" {
		baseQuery += " AND c.class_id = @class_id"
		params["class_id"] = classID
	}

	baseQuery += " ORDER BY s.start_time"

	var results []SnapshotResponse
	if err := config.DB.Raw(baseQuery, params).Scan(&results).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}
