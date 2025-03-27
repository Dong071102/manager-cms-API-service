package controllers

import (
	"cms-backend/config"
	"cms-backend/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"cms-backend/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	pgvector "github.com/pgvector/pgvector-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func GenerateFakeEmbedding(dim int) pgvector.Vector {
	vec := pgvector.NewVector(make([]float32, dim))
	return vec
}

func RegisterUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var existing models.User
	if err := config.DB.Where("username = ? OR email = ?", user.Username, user.Email).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Username or Email already exists"})
	}

	hashed, err := HashPassword(user.PasswordHash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Password hashing failed"})
	}
	user.PasswordHash = hashed

	config.DB.Create(&user)

	switch user.Role {
	case "student":
		var studentCode string
		for {
			code := models.GenerateStudentCode()
			var count int64
			config.DB.Model(&models.Student{}).Where("student_code = ?", code).Count(&count)
			if count == 0 {
				studentCode = code
				break
			}
		}
		s := models.Student{
			StudentID:     user.UserID,
			StudentCode:   studentCode,
			FaceEmbedding: GenerateFakeEmbedding(512)}
		config.DB.Create(&s)
	case "lecturer":
		l := models.Lecturer{LecturerID: user.UserID, LectainerCode: uuid.New().String()}
		config.DB.Create(&l)
	case "admin":
		a := models.Admin{AdminID: user.UserID, AdminCode: uuid.New().String()}
		config.DB.Create(&a)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid role"})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
		"user_id": user.UserID.String(),
	})
}
func LoginUser(c echo.Context) error {
	var input struct {
		UsernameOrEmail string `json:"username_or_email"`
		Password        string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	var user models.User
	if err := config.DB.Where("username = ? OR email = ?", input.UsernameOrEmail, input.UsernameOrEmail).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}
	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate access token"})
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
func ChangePassword(c echo.Context) error {
	claims := c.Get("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var user models.User
	if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Old password incorrect"})
	}

	hashed, err := HashPassword(input.NewPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}

	user.PasswordHash = hashed
	config.DB.Save(&user)

	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

func ForgotPassword(c echo.Context) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&input); err != nil || input.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid email"})
	}

	var user models.User
	if err := config.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Normally here you would generate a reset token and send it via email
	// For demo, we just return a fake token (insecure in real app)
	resetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID.String(),
		"exp":     jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 minutes expiry
	})
	tokenStr, err := resetToken.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate reset token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Reset token generated (normally sent via email)",
		"reset_token": tokenStr,
	})
}
func RefreshToken(c echo.Context) error {
	// Lấy refresh_token từ header Authorization
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("Authorization header:", authHeader)

	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing Authorization header"})
	}

	// Kiểm tra xem header có chứa Bearer token không
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid Authorization header format"})
	}

	// Lấy refresh_token từ Authorization header
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Kiểm tra và xử lý refresh_token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired refresh token"})
	}

	// Tiến hành xử lý khi token hợp lệ
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil || claims["role"] == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
	}

	// Tạo access_token mới
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims["user_id"],
		"role":    claims["role"],
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	tokenStr, err := newAccessToken.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate new access token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": tokenStr,
	})
}

func UpdateProfile(c echo.Context) error {
	claimsRaw := c.Get("user")
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID in token"})
	}

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		ImageURL  string `json:"image_url"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var user models.User
	if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.ImageURL != "" {
		user.ImageURL = input.ImageURL
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update profile"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

func GetCurrentUser(c echo.Context) error {
	// Lấy claims từ JWT
	claimsRaw := c.Get("user")
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid token claims"})
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid user ID in token"})
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid UUID format"})
	}

	var user models.User
	if err := config.DB.First(&user, "user_id = ?", userUUID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
	}

	// Trả thông tin user đã xác thực
	return c.JSON(http.StatusOK, echo.Map{
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"role":       user.Role,
		"image_url":  user.ImageURL,
	})
}
