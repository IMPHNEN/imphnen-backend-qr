# IMPHNEN Backend QR Code Generator

QR code overlay generator untuk campaign. Dengan fitur SSO dibangun menggunakan Go + Echo dengan JWT auth, Google OAuth, dan RBAC.

## Tech Stack

- **Go** + **Echo v4** — HTTP framework
- **PostgreSQL** — Database
- **JWT** — Access token (15 min) + Refresh token (7 days)
- **bcrypt** — Password hashing
- **Viper** — Environment config (prioritas env OS, fallback `.env`)
- **Google OAuth2** — Social login
- **go-qrcode** — QR code generation
- **image/draw** — Image overlay processing

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

### Dengan Docker Compose

```bash
docker compose up --build
```

Ini akan menjalankan app di `localhost:8080` dan Postgres di `localhost:5432`.

### Manual

1. Copy dan isi environment variables:
   ```bash
   cp .env.example .env
   # edit .env sesuai konfigurasi
   ```

2. Jalankan migration (manual atau pakai golang-migrate CLI):
   ```bash
   # Install golang-migrate
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

   # Run migration
   migrate -path db/migrations -database "$DATABASE_URL" up
   ```

3. Jalankan server:
   ```bash
   go run cmd/server/main.go
   ```
   Server akan otomatis menjalankan seeder saat start (membuat akun demo admin & user jika belum ada).

4. (Opsional) Jalankan seeder secara manual:
   ```bash
   go run cmd/seeder/main.go
   ```

### Demo Accounts (Auto-seeded)

| Email | Password | Role |
|---|---|---|
| `admin@imphnen.dev` | `admin123` | admin |
| `user@imphnen.dev` | `user123` | user |

## API Endpoints

### Auth (Public)

| Method | Path                            | Description              |
|--------|---------------------------------|--------------------------|
| POST   | `/api/v1/auth/register`         | Register email/password  |
| POST   | `/api/v1/auth/login`            | Login email/password     |
| GET    | `/api/v1/auth/google`           | Redirect ke Google OAuth |
| GET    | `/api/v1/auth/google/callback`  | OAuth callback           |
| POST   | `/api/v1/auth/refresh`          | Refresh access token     |

### Users (Protected — Bearer Token)

| Method | Path                       | Role         | Description         |
|--------|----------------------------|--------------|---------------------|
| GET    | `/api/v1/users/me`         | User, Admin  | Get profile         |
| PUT    | `/api/v1/users/me`         | User, Admin  | Update profile      |
| GET    | `/api/v1/users`            | Admin        | List all users      |
| PUT    | `/api/v1/users/:id/role`   | Admin        | Update user role    |

### QR Campaigns (Protected — Bearer Token)

| Method | Path                                    | Role         | Description                       |
|--------|-----------------------------------------|--------------|-----------------------------------|
| POST   | `/api/v1/campaigns`                     | Admin        | Create campaign (auto-activates)  |
| GET    | `/api/v1/campaigns`                     | Admin        | List all campaigns                |
| PUT    | `/api/v1/campaigns/:id/activate`        | Admin        | Set campaign as active            |
| POST   | `/api/v1/campaigns/process-image`       | User, Admin  | Upload image, get QR overlay PNG  |

### Health Check

```
GET /health → {"status": "ok"}
```

## QR Campaign Overlay

### Flow
1. **Admin** membuat campaign via `POST /api/v1/campaigns` dengan `name` dan `url`
2. System generate QR code PNG (256x256) dari URL campaign
3. Campaign baru otomatis menjadi active (hanya 1 active pada satu waktu)
4. **User** upload image via `POST /api/v1/campaigns/process-image` (multipart, field: `image`)
5. System merge QR code ke bottom-right corner dari image
6. Response berupa binary PNG

### Create Campaign Request
```json
{
  "name": "My Campaign",
  "url": "https://example.com/promo"
}
```

### Process Image Request
```bash
curl -X POST http://localhost:8080/api/v1/campaigns/process-image \
  -H "Authorization: Bearer <token>" \
  -F "image=@photo.jpg"
```
Response: binary `image/png` dengan QR overlay di bottom-right.

## API Response Format

**Success:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error:**
```json
{
  "success": false,
  "message": "Error description",
  "error": "error_code"
}
```

## Project Structure

```
cmd/server/main.go          — Entry point + auto-seeder
cmd/seeder/main.go          — Standalone seeder CLI
internal/config/             — Environment config (Viper)
internal/domain/             — Entities & interfaces
internal/handler/            — HTTP handlers
internal/middleware/          — JWT & RBAC middleware
internal/repository/         — Database operations (raw SQL)
internal/service/            — Business logic + QR generation + image processing
internal/seeder/             — Database seeder (demo accounts)
internal/utils/              — JWT, password, response helpers
pkg/database/                — Postgres connection
db/migrations/               — SQL migration files
```

## Environment Variables

| Variable              | Required | Default | Description                    |
|-----------------------|----------|---------|--------------------------------|
| `PORT`                | No       | 8080    | Server port                    |
| `DATABASE_URL`        | Yes      | —       | PostgreSQL connection string   |
| `JWT_SECRET`          | Yes      | —       | Secret for signing JWT         |
| `GOOGLE_CLIENT_ID`    | No       | —       | Google OAuth client ID         |
| `GOOGLE_CLIENT_SECRET`| No       | —       | Google OAuth client secret     |
| `GOOGLE_REDIRECT_URL` | No       | —       | Google OAuth redirect URL      |
