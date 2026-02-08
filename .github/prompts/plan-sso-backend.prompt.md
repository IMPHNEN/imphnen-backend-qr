## Plan: Go Echo SSO Backend dengan RBAC

**TL;DR** — Membangun backend SSO menggunakan Go + Echo dengan autentikasi JWT, dual login (email/password + Google OAuth), dan RBAC sederhana untuk dua role (Admin, User). Arsitektur menggunakan clean architecture pattern dengan layer handler, service, dan repository. Database menggunakan PostgreSQL dengan sqlc (raw SQL, bukan ORM) lengkap dengan migrasi. Environment variable dikelola dengan Viper (prioritas env OS, fallback .env).

---

### Struktur Proyek

```
imphnen-backend-qr/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── db/
│   ├── migrations/             # SQL migration files
│   └── query.sql               # sqlc query definitions
├── internal/
│   ├── config/
│   │   └── config.go            # Environment & app config
│   ├── domain/
│   │   ├── user.go              # User entity & interfaces
│   │   └── token.go             # Token entity
│   ├── handler/
│   │   ├── auth_handler.go      # Login, register, OAuth endpoints
│   │   └── user_handler.go      # User management endpoints
│   ├── middleware/
│   │   ├── jwt_middleware.go    # JWT validation middleware
│   │   └── rbac_middleware.go   # Role-based access control
│   ├── repository/
│   │   └── user_repository.go   # Database operations
│   ├── service/
│   │   ├── auth_service.go      # Auth business logic
│   │   └── user_service.go      # User business logic
│   └── utils/
│       ├── jwt.go               # JWT helper functions
│       ├── password.go          # Password hashing
│       └── response.go          # Standard API response
├── pkg/
│   └── database/
│       └── postgres.go          # Database connection (sqlc, raw SQL)
├── go.mod
├── go.sum
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

### Steps

**1. Inisialisasi Project**
- Buat `go.mod` dengan module name `github.com/IMPHNEN/imphnen-backend-qr`
- Install dependencies: `echo/v4`, `oauth2`, `jwt/v5`, `bcrypt`, `sqlc`, `viper`, `pq` (Postgres driver)

**2. Setup Config & Database**
- Buat `config.go` untuk load environment variables (JWT_SECRET, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, DATABASE_URL) menggunakan Viper (prioritas env OS, fallback .env)
- Buat `postgres.go` untuk koneksi database menggunakan sqlc (raw SQL, bukan ORM)
- Setup migrasi database dengan SQL migration files (misal pakai golang-migrate atau migrate CLI)

**3. Domain Layer**
- Definisikan `User` entity dengan field: ID, Email, Password (nullable untuk OAuth), Name, Role (enum: admin/user), Provider (local/google), ProviderID, CreatedAt, UpdatedAt
- Definisikan interface `UserRepository` dan `AuthService`

**4. Repository Layer**
- Implementasi `UserRepository` dengan method: Create, FindByEmail, FindByID, FindByProviderID, Update menggunakan sqlc (raw SQL)

**5. Service Layer - Auth**
- `Register`: Hash password dengan bcrypt, create user dengan role default "user"
- `Login`: Validate credentials, generate JWT dengan claims (user_id, email, role)
- `GoogleLogin`: Exchange authorization code, get user info dari Google, create/update user, generate JWT
- `RefreshToken`: Validate refresh token, generate new access token

**6. Service Layer - User**
- `GetProfile`: Get current user info
- `UpdateProfile`: Update user data
- `GetAllUsers`: Admin only - list semua users
- `UpdateUserRole`: Admin only - ubah role user

**7. JWT & Password Utils**
- `GenerateToken`: Buat access token (15 menit) dan refresh token (7 hari)
- `ValidateToken`: Validate dan parse JWT claims
- `HashPassword`: Hash password dengan bcrypt cost 10
- `ComparePassword`: Compare password dengan hash

**8. Middleware**
- `JWTMiddleware`: Extract dan validate JWT dari header `Authorization: Bearer <token>`
- `RBACMiddleware`: Check role dari JWT claims, restrict berdasarkan allowed roles

**9. Handler Layer - Auth**
```
POST /api/v1/auth/register     → Register dengan email/password
POST /api/v1/auth/login        → Login dengan email/password
GET  /api/v1/auth/google       → Redirect ke Google OAuth consent
GET  /api/v1/auth/google/callback → Handle OAuth callback
POST /api/v1/auth/refresh      → Refresh access token
POST /api/v1/auth/logout       → Logout (invalidate token - optional)
```

**10. Handler Layer - User (Protected)**
```
GET    /api/v1/users/me        → Get current user profile [User, Admin]
PUT    /api/v1/users/me        → Update profile [User, Admin]
GET    /api/v1/users           → List all users [Admin only]
PUT    /api/v1/users/:id/role  → Update user role [Admin only]
```

**11. Entry Point & Routes**
- Setup Echo instance di `main.go`
- Register middleware: Logger, Recover, CORS
- Mount auth routes (public)
- Mount user routes (protected dengan JWT + RBAC middleware)

**12. Environment & Documentation**
- Buat `.env.example` dengan template config
- Update `README.md` dengan setup instructions

**13. Docker & Deployment**
- Buat `Dockerfile` untuk build Go app
- Buat `docker-compose.yml` untuk menjalankan app dan Postgres (gunakan volume agar data Postgres persisten di server deploy)
- Contoh volume:
  ```yaml
  services:
    db:
      image: postgres:16
      restart: always
      environment:
        POSTGRES_USER: youruser
        POSTGRES_PASSWORD: yourpassword
        POSTGRES_DB: yourdb
      ports:
        - "5432:5432"
      volumes:
        - pgdata:/var/lib/postgresql/data
  volumes:
    pgdata:
  ```

---

### API Response Format

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

Error response:
```json
{
  "success": false,
  "message": "Error description",
  "error": "error_code"
}
```

---

### JWT Claims Structure

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "user",
  "exp": 1234567890,
  "iat": 1234567890
}
```

---

### Database Schema (User)

| Column      | Type         | Notes                          |
|-------------|--------------|--------------------------------|
| id          | UUID         | Primary key                    |
| email       | VARCHAR(255) | Unique, required               |
| password    | VARCHAR(255) | Nullable (for OAuth users)     |
| name        | VARCHAR(255) | Required                       |
| role        | VARCHAR(20)  | "admin" or "user"              |
| provider    | VARCHAR(20)  | "local" or "google"            |
| provider_id | VARCHAR(255) | Google user ID (nullable)      |
| created_at  | TIMESTAMP    |                                |
| updated_at  | TIMESTAMP    |                                |

---

### Verification

1. **Unit Test**: Test service layer dengan mock repository
2. **Integration Test**: Test endpoints dengan database test
3. **Manual Testing**:
   - Register user baru → expect JWT
   - Login dengan credentials → expect JWT
   - Access protected route tanpa token → expect 401
   - Access admin route dengan user token → expect 403
   - Google OAuth flow → expect redirect dan JWT

---

### Decisions

- **Database**: PostgreSQL + sqlc (raw SQL, bukan ORM) — type safety, migration manual, query optimal
- **Token Strategy**: Access token (15 min) + Refresh token (7 days) — balance security & UX
- **Password Hash**: bcrypt cost 10 — industry standard, good security/performance balance
- **RBAC**: Simple role check di middleware — cukup untuk 2 roles, mudah di-extend
- **Env Management**: Viper (prioritas env OS, fallback .env)
