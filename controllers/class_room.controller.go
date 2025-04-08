package controllers

import (
	"cms-backend/config"
	"cms-backend/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetClassrooms(c echo.Context) error {
	var classrooms []models.Classroom
	if err := config.DB.Find(&classrooms).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to retrieve classrooms",
		})
	}
	return c.JSON(http.StatusOK, classrooms)
}

