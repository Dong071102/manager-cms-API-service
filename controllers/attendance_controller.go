package controllers

import (
	"cms-backend/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCameraSocketPath(c echo.Context) error {
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
		WHERE s.schedule_id = ? AND c.camera_type = 'recognition'
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
