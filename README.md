# 📚 CMS Backend Services (Golang API)

**CMS Backend Services** là hệ thống API được viết bằng **Golang (Echo framework)**, phục vụ cho ứng dụng quản lý điểm danh sinh viên, lớp học, giảng viên và camera AI.

---

## 🚀 Tính năng chính

Hệ thống cung cấp các API RESTful hỗ trợ các chức năng:

### 👨‍🏫 Dành cho giảng viên:
- Lấy danh sách lớp theo `lecturer_id`
- Lấy tổng quan và chi tiết điểm danh theo lớp/học sinh
- Cập nhật trạng thái điểm danh
- Lấy báo cáo điểm danh tổng hợp

### 👨‍🎓 Dành cho sinh viên:
- Lấy danh sách sinh viên trong lớp
- Kiểm tra mã sinh viên đã tồn tại chưa
- Xem lịch sử điểm danh theo sinh viên

### 🏫 Dành cho hệ thống:
- Thêm / xóa / cập nhật lịch học (schedule)
- Lấy danh sách phòng học, lớp học, khóa học
- Lấy thông tin socket stream camera AI (Face / Human Counter)
- Lấy chi tiết ảnh snapshot từ hệ thống AI

---

## 📁 Các API chính (Router)

| Method | Endpoint | Chức năng |
|--------|----------|-----------|
| GET | `/classes/:lecturer_id` | Danh sách lớp theo giảng viên |
| GET | `/attendance-summary` | Tổng hợp điểm danh |
| GET | `/attendance-detail` | Chi tiết điểm danh |
| POST | `/update-attendance` | Cập nhật trạng thái điểm danh |
| GET | `/attendance-report/:lecturer_id` | Báo cáo điểm danh |
| GET | `/students-in-class/:lecturer_id` | Danh sách sinh viên trong lớp |
| PUT | `/update/student/:id` | Cập nhật thông tin sinh viên |
| DELETE | `/del-student-from-class/:student_id/:class_id` | Xóa sinh viên khỏi lớp |
| POST | `/add-student-to-class` | Thêm sinh viên vào lớp |
| GET | `/student-attendance-summary/:lecturer_id` | Tổng hợp điểm danh sinh viên |
| GET | `/get-student-attendances/:student_id/:lecturer_id` | Lịch sử điểm danh sinh viên |
| GET | `/get-classrooms` | Danh sách phòng học |
| GET | `/get-schedules` | Danh sách lịch học |
| GET | `/get-courses-by-lecturerID` | Khóa học theo giảng viên |
| GET | `/get-class-by-course-id` | Lớp học theo khóa |
| POST | `/add-schedule` | Thêm lịch học |
| PUT | `/update-schedule/:id` | Cập nhật lịch học |
| DELETE | `/delete-schedule/:id` | Xóa lịch học |
| GET | `/get-schedule-start-times` | Các giờ bắt đầu lịch học |
| GET | `/get-schedule-times` | Danh sách giờ học |
| GET | `/get-attendance-socket-path` | URL stream camera nhận diện khuôn mặt |
| GET | `/get-human-couter-socket-path` | URL stream đếm người |
| GET | `/get-snapshot-details` | Thông tin ảnh snapshot |

---

## ⚙️ Cài đặt & Chạy API

### 1. Clone project

```bash
git clone https://github.com/yourusername/cms-backend.git
cd cms-backend
```

### 2. Khởi tạo module

```bash
go mod tidy
```

### 3. Chạy server

```bash
go run main.go
```

Server sẽ chạy tại: `http://localhost:8080`

---

## 🔐 Ghi chú bảo mật

- Bạn có thể thêm middleware xác thực JWT/Token cho từng route nếu cần.
- Dữ liệu `userId` ví dụ: `"2d536da8-fdf3-437b-a812-fb4e08aad955"` sẽ được client gửi kèm trong request header/body.

---

## 🛠️ Công nghệ sử dụng

- Language: **Go**
- Web Framework: **Echo**
- DB support: (PostgreSQL, MySQL, SQLite tùy bạn kết nối ở tầng `repository`)
- Kiểu trả về: **JSON API**

---

## 📧 Tác giả

- **Tên**: Vũ Bá Đông  
- 📩 Email: [vubadong071102@gmail.com](mailto:vubadong071102@gmail.com)

---

## 📄 License

MIT License