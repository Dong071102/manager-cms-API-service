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
	e.POST("/update-attendance", controllers.UpdateAttendance)
	e.GET("/attendance-report/:lecturer_id", controllers.GetAttendanceReport)
	e.GET("/students-in-class/:lecturer_id", controllers.GetStudentsInClass)
	e.GET("/students-in-class/:lecturer_id", controllers.GetStudentsInClass)
	e.PUT("/update/student/:id", controllers.UpdateStudent)
	e.DELETE("/del-student-from-class/:student_id/:class_id", controllers.DeleteStudentFromClass)
	e.GET("/check-student-existence/:studentCode", controllers.CheckStudentExistence)
	// Thêm route cho API add-student-to-class
	e.POST("/add-student-to-class", controllers.AddStudentToClass)
	e.GET("/student-attendance-summary/:lecturer_id", controllers.GetStudentAttendanceSummary)
	e.GET("/get-student-attendances/:student_id/:lecturer_id", controllers.GetStudentAttendances)
	e.GET("/get-classrooms", controllers.GetClassrooms)
	e.GET("/get-schedules", controllers.GetSchedules)
	e.GET("/get-courses-by-lecturerID", controllers.GetCoursesByLecturerID)
	e.GET("/get-class-by-course-id", controllers.GetClassesByCourses)

	e.POST("/add-schedule", controllers.AddSubject)
	e.PUT("update-schedule/:id", controllers.UpdateSubject)
	e.DELETE("/delete-schedule/:id", controllers.DeleteSchedule)
	e.GET("/get-schedule-start-times", controllers.GetScheduleStartTimes)
	e.GET("/get-schedule-times", controllers.GetScheduTimes)
	e.GET("/get-attendance-socket-path", controllers.GetCameraSocketPath)
	e.GET("/get-human-couter-socket-path", controllers.GetHumanCouterSocketPath)
	e.GET("/get-snapshot-details", controllers.GetSnapshotDetails)
}

// userId:"2d536da8-fdf3-437b-a812-fb4e08aad955"
