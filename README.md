# Go Clean Architecture REST API Boilerplate

Dự án này là một bản boilerplate hoàn chỉnh cho Go REST API, được xây dựng theo kiến trúc **Clean Architecture** (Domain, Usecase, Repository, Delivery) nhằm đảm bảo tính bảo trì, mở rộng và dễ dàng kiểm thử.

## 🚀 Công nghệ sử dụng
- **Languague:** Go (Golang)
- **Framework:** [Gin-gonic](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/) (với driver MySQL)
- **Log:** [Uber-zap](https://github.com/uber-go/zap) (custom middleware)
- **Auth:** JWT (JSON Web Token)
- **Validation:** Go Playground Validator (tích hợp trong Gin)
- **Database:** MySQL 8.0

## 📂 Cấu trúc thư mục
```bash
.
├── cmd/api/            # Điểm khởi đầu (main.go), Dependency Injection
├── config/             # Quản lý cấu hình tập trung từ .env
├── internal/
│   ├── domain/         # Entities (Thực thể), Interfaces (không phụ thuộc tag)
│   ├── repository/     # Triển khai tầng lưu trữ (GORM models & operations)
│   ├── usecase/        # Logic nghiệp vụ (Business Logic)
│   └── delivery/http/  # Tầng vận chuyển (Handlers, Routes, Middlewares, DTOs)
├── pkg/                # Các thư viện dùng chung (Logger, Database connection)
├── .env.example        # Mẫu file cấu hình môi trường
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yaml # Docker setup cho MySQL và App
└── Makefile            # Lệnh tắt cho build, run, test
```

## 🛠 Hướng dẫn chạy nhanh

### 1. Chuẩn bị môi trường
Copy file mẫu cấu hình và điều chỉnh các thông số (Database, JWT Secret):
```bash
cp .env.example .env
```

### 2. Chạy với Docker (Khuyên dùng)
Lệnh này sẽ khởi chạy cả MySQL và API App trong container:
```bash
make docker-up
```

### 3. Chạy local
Nếu bạn đã có MySQL chạy sẵn, hãy cập nhật thông tin trong `.env` và chạy trực tiếp:
```bash
make run
```

## 📝 Danh sách API chính

| Group | Method | Endpoint | Auth |
| :--- | :--- | :--- | :--- |
| **Auth** | POST | `/api/v1/auth/register` | ❌ |
| | POST | `/api/v1/auth/login` | ❌ |
| **User** | GET | `/api/v1/users/me` | ✅ |
| | PUT | `/api/v1/users/me` | ✅ |
| | GET | `/api/v1/users` | ✅ |
| **Task** | POST | `/api/v1/tasks` | ✅ |
| | GET | `/api/v1/tasks` | ✅ |
| | GET | `/api/v1/tasks/:id` | ✅ |
| | PUT | `/api/v1/tasks/:id` | ✅ |
| | DELETE | `/api/v1/tasks/:id` | ✅ |

## 🛡 Đặc điểm nổi bật
- **Strict Decoupling:** Tách biệt hoàn toàn DTO (giao tiếp), Domain (lõi) và Model (database).
- **Security:** Tự động check quyền sở hữu (ownership) đối với các CRUD Task.
- **Graceful Shutdown:** Đảm bảo server đóng các kết nối an toàn khi nhận tín hiệu tắt.
- **Centralized Error Handling:** Quản lý mã lỗi và thông báo lỗi tập trung tại tầng Domain.
