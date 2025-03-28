package routes

import (
	"cms-backend/controllers"
	custommiddleware "cms-backend/middleware"
	"fmt"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/auth/health", func(c echo.Context) error {
		fmt.Println("🔥 Hit /health endpoint")
		return c.JSON(200, map[string]string{"message": "API is running ✅"})
	})
	e.POST("/auth/register", controllers.RegisterUser)
	e.POST("/auth/login", controllers.LoginUser)

	// Protected route example:
	e.GET("/auth/me", controllers.GetCurrentUser, custommiddleware.JWTAuthMiddleware)
	// Admin-only route example:
	e.GET("/auth/admin-only", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Welcome, admin!"})
	}, custommiddleware.JWTAuthMiddleware, custommiddleware.RoleMiddleware("admin"))

	// Student-only route example:
	e.GET("/auth/student-only", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Welcome, student!"})
	}, custommiddleware.JWTAuthMiddleware, custommiddleware.RoleMiddleware("student"))

	// Lecturer-only route example:
	e.GET("/auth/lecturer-only", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Welcome, lecturer!"})
	}, custommiddleware.JWTAuthMiddleware, custommiddleware.RoleMiddleware("lecturer"))
	e.POST("/auth/change-password", controllers.ChangePassword, custommiddleware.JWTAuthMiddleware)
	e.POST("/auth/forgot-password", controllers.ForgotPassword)
	e.POST("/auth/refresh-token", controllers.RefreshToken)
	e.PATCH("/auth/update-profile", controllers.UpdateProfile, custommiddleware.JWTAuthMiddleware)

	// ================================================================
	e.GET("/classes/:lecturer_id", controllers.GetClassesByLecturer)
	e.GET("/attendance-summary", controllers.AttendanceSummaryHandler)
	e.GET("/attendance-detail", controllers.GetAttendanceDetails)
}
