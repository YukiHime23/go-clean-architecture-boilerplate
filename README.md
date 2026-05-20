# Go Clean Architecture REST API Boilerplate

REST API boilerplate được xây dựng theo **Clean Architecture** với Go, Gin, GORM và JWT.

## Tech Stack

- **Go** 1.22+
- **Gin** — HTTP framework
- **GORM** + MySQL driver
- **golang-jwt/jwt v5** — JWT authentication
- **godotenv** — `.env` loader
- **bcrypt** — password hashing

## Project Structure

```
.
├── cmd/api/                        # Entry point (main.go)
├── config/                         # Config loader từ env vars
├── internal/
│   ├── domain/
│   │   └── entity/                 # Core domain structs (User, Task)
│   ├── repository/mysql/           # GORM implementations
│   ├── usecase/
│   │   ├── auth/                   # DTOs + repo interfaces + business logic
│   │   ├── user/
│   │   └── task/
│   └── delivery/
│       ├── handler/                # Gin handlers + usecase interfaces
│       ├── middleware/             # JWT auth middleware
│       └── router/                 # Route registration
└── pkg/
    ├── apperror/                   # Typed app errors
    ├── jwt/                        # JWT generate/parse helpers
    └── response/                   # Unified JSON response helpers
```

## Getting Started

### 1. Clone và cài đặt

```bash
git clone <repo>
cd <repo>
go mod download
```

### 2. Cấu hình environment

```bash
cp .env.example .env
# Chỉnh sửa .env với thông tin DB và JWT secret của bạn
```

### 3. Tạo database MySQL

```sql
CREATE DATABASE go_clean_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. Chạy server

```bash
go run ./cmd/api
```

Server khởi động tại `http://localhost:8080`. GORM sẽ tự động `AutoMigrate` tạo bảng.

## API Endpoints

| Group  | Method | Endpoint                | Auth |
| :----- | :----- | :---------------------- | :--- |
| Auth   | POST   | `/api/v1/auth/register` | ❌   |
|        | POST   | `/api/v1/auth/login`    | ❌   |
| User   | GET    | `/api/v1/users/me`      | ✅   |
|        | PUT    | `/api/v1/users/me`      | ✅   |
|        | GET    | `/api/v1/users`         | ✅   |
| Task   | POST   | `/api/v1/tasks`         | ✅   |
|        | GET    | `/api/v1/tasks`         | ✅   |
|        | GET    | `/api/v1/tasks/:id`     | ✅   |
|        | PUT    | `/api/v1/tasks/:id`     | ✅   |
|        | DELETE | `/api/v1/tasks/:id`     | ✅   |

## Authentication

Các endpoint có `✅` yêu cầu header:

```
Authorization: Bearer <token>
```

Token nhận được từ response của `/auth/register` hoặc `/auth/login`.

## Request Examples

### Register
```json
POST /api/v1/auth/register
{
  "name": "Nguyen Van A",
  "email": "user@example.com",
  "password": "secret123"
}
```

### Login
```json
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "secret123"
}
```

### Create Task
```json
POST /api/v1/tasks
Authorization: Bearer <token>
{
  "title": "Implement feature X",
  "description": "Details here",
  "status": "pending"
}
```

## Task Status Values

- `pending`
- `in_progress`
- `done`
