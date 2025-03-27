package main

import (
	"cms-backend/config"
	"cms-backend/models"
	"cms-backend/routes"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Tải tệp .env và xử lý lỗi nếu không thể tải
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Kiểm tra giá trị của JWT_SECRET để chắc chắn rằng nó đã được tải
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required but not set in the environment variables")
	}
	fmt.Println("Loaded JWT_SECRET:", jwtSecret)

	// Kết nối đến cơ sở dữ liệu
	config.ConnectDB()

	// Thực hiện AutoMigration cho các model
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Lecturer{},
		&models.Admin{},
		// &models.Class{},
		// &models.Course{},
	); err != nil {
		log.Fatalf("Error during database migration: %v", err)
	}

	// Khởi tạo một instance của Echo
	e := echo.New()

	// Cấu hình middleware của Echo
	e.Use(echomiddleware.Logger())  // Log mỗi request
	e.Use(echomiddleware.Recover()) // Giúp ứng dụng không bị sập khi có lỗi

	// Cấu hình CORS
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // Địa chỉ frontend của bạn (có thể thay đổi nếu khác)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Cho phép gửi cookies hoặc thông tin xác thực (nếu cần)
	}))

	// Setup các route từ tệp routes
	routes.SetupRoutes(e)

	// Lấy cổng từ biến môi trường hoặc mặc định là 9000
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	// Khởi chạy server trên cổng đã cấu hình
	e.Logger.Fatal(e.Start(":" + port))
}
