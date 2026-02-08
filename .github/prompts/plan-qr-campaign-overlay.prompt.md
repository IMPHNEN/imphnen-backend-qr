# Plan: QR Campaign Overlay MVP

**TL;DR:** Tambah fitur QR campaign overlay ke backend yang sudah ada. Admin membuat campaign (nama + URL), system generate QR code PNG & cache di memory. User upload image, system merge QR ke bottom-right corner dari active campaign. Hanya 1 campaign aktif pada satu waktu. Seeder berjalan otomatis saat server start + tersedia sebagai CLI command terpisah.

---

## Phase 1: Database Schema

**Step 1** — Buat migration file `db/migrations/000002_create_qr_campaigns_table.up.sql` dan `.down.sql`

Schema `qr_campaigns`:

| Column | Type | Constraint |
|---|---|---|
| `id` | UUID | PK, DEFAULT `uuid_generate_v4()` |
| `name` | VARCHAR(255) | NOT NULL |
| `url` | TEXT | NOT NULL |
| `qr_code_data` | BYTEA | NOT NULL (PNG bytes) |
| `is_active` | BOOLEAN | NOT NULL DEFAULT `false` |
| `created_by` | UUID | NOT NULL, FK → `users(id)` ON DELETE CASCADE |
| `expires_at` | TIMESTAMP | NOT NULL (default: created_at + 7 days) |
| `created_at` | TIMESTAMP | NOT NULL DEFAULT NOW() |
| `updated_at` | TIMESTAMP | NOT NULL DEFAULT NOW() |

Index pada `is_active` untuk fast lookup active campaign. Partial unique index `WHERE is_active = true` untuk enforce max 1 active campaign di DB level.

Down migration: `DROP TABLE IF EXISTS qr_campaigns;`

---

## Phase 2: Backend Development

**Step 2** — Tambah dependencies di `go.mod`
- `github.com/skip2/go-qrcode` — pure Go QR code generator (generate PNG bytes langsung)
- Tidak perlu library tambahan untuk image merge — gunakan standard library `image`, `image/png`, `image/jpeg`, `image/draw`

**Step 3** — Domain layer: buat `internal/domain/qr_campaign.go`
- Struct `QRCampaign` dengan field: `ID`, `Name`, `URL`, `QRCodeData` ([]byte, `json:"-"`), `IsActive`, `CreatedBy`, `ExpiresAt`, `CreatedAt`, `UpdatedAt`
- Interface `QRCampaignRepository` dengan method:
  - `Create(campaign *QRCampaign) error`
  - `FindByID(id string) (*QRCampaign, error)`
  - `FindActive() (*QRCampaign, error)` — return campaign with `is_active = true`
  - `FindAll() ([]*QRCampaign, error)`
  - `SetActive(id string) error` — deactivate all lalu activate 1 (dalam transaksi)
  - `Delete(id string) error`

**Step 4** — Repository layer: buat `internal/repository/qr_campaign_repository.go`
- Implementasi `domain.QRCampaignRepository` mengikuti pattern di `internal/repository/user_repository.go`
- Raw SQL, parameterized queries, `uuid.New().String()` untuk ID, `(nil, nil)` pattern untuk not-found
- `SetActive`: gunakan SQL transaction — `UPDATE qr_campaigns SET is_active = false` lalu `UPDATE qr_campaigns SET is_active = true WHERE id = $1`

**Step 5** — Service layer: buat `internal/service/qr_campaign_service.go`
- Struct `QRCampaignService` dengan field:
  - `repo domain.QRCampaignRepository`
  - `cache` — simple in-memory cache (`sync.RWMutex` + `cachedQR []byte` + `cachedCampaignID string`)
- Input DTO: `CreateCampaignInput{Name, URL}`
- Method `CreateCampaign(input, createdBy) (*QRCampaign, error)`:
  1. Generate QR PNG bytes via `go-qrcode` (256x256px, medium recovery)
  2. Set `expires_at` = now + 7 hari
  3. Set `is_active = true` (auto-activate newly created campaign)
  4. Simpan ke DB via repository
  5. Update in-memory cache
- Method `GetActiveCampaign() (*QRCampaign, error)` — cek cache dulu, fallback ke DB, update cache
- Method `GetAllCampaigns() ([]*QRCampaign, error)`
- Method `SetActiveCampaign(id) error` — call repo `SetActive`, invalidate & refresh cache
- Method `ProcessImage(uploadedImage io.Reader) ([]byte, error)`:
  1. Load active campaign QR dari cache (via `GetActiveCampaign`)
  2. Decode uploaded image (support PNG & JPEG via `image.Decode` after registering decoders)
  3. Decode QR PNG bytes ke `image.Image`
  4. Buat canvas baru, draw original image, lalu draw QR di bottom-right (dengan padding ~10px, QR di-resize ke ~1/5 dari dimensi terkecil image jika perlu)
  5. Encode hasil ke PNG bytes, return
- Sentinel errors: `ErrCampaignNotFound`, `ErrNoActiveCampaign`, `ErrInvalidImage`

**Step 6** — Handler layer: buat `internal/handler/qr_campaign_handler.go`
- Struct `QRCampaignHandler` dengan `*service.QRCampaignService`
- **Admin endpoints:**
  - `CreateCampaign(c echo.Context) error` — bind JSON `{name, url}`, validate, call service, return `SuccessResponse` dengan campaign data
  - `GetAllCampaigns(c echo.Context) error` — return list campaigns
  - `SetActiveCampaign(c echo.Context) error` — param `:id`, call service, return success
- **User endpoint:**
  - `ProcessImage(c echo.Context) error` — read multipart file upload (field `image`), call service `ProcessImage`, return binary PNG response dengan `Content-Type: image/png`

**Step 7** — Route registration di `cmd/server/main.go`
- Instantiate `qrCampaignRepo`, `qrCampaignService`, `qrCampaignHandler` mengikuti DI pattern yang ada
- Admin routes (JWT + RBAC "admin"):
  - `POST /api/v1/campaigns` → `CreateCampaign`
  - `GET /api/v1/campaigns` → `GetAllCampaigns`
  - `PUT /api/v1/campaigns/:id/activate` → `SetActiveCampaign`
- User routes (JWT only, semua role bisa akses):
  - `POST /api/v1/campaigns/process-image` → `ProcessImage`

**Step 8** — Auto-seeder di `cmd/server/main.go` + CLI command di `cmd/seeder/main.go`
- Logika seeder (shared function di `internal/seeder/seeder.go`):
  1. Cek apakah user dengan email `admin@imphnen.id` sudah ada
  2. Jika belum, create admin: `{email: "admin@imphnen.id", password: "admin123", name: "Admin Demo", role: "admin"}`
  3. Cek apakah user dengan email `user@imphnen.id` sudah ada
  4. Jika belum, create user: `{email: "user@imphnen.id", password: "user123", name: "User Demo", role: "user"}`
  5. Log hasil (created / already exists)
- Di `main.go`: panggil seeder setelah DB connection established, sebelum server start
- Di `cmd/seeder/main.go`: standalone command yang load config, connect DB, run seeder, lalu exit

---

## Verification

1. **Migration**: Jalankan migration, verify tabel `qr_campaigns` terbuat dengan schema yang benar
2. **Seeder**: Start server, cek log bahwa 2 akun demo terbuat. Jalankan ulang, pastikan tidak duplicate
3. **Auth flow**: Login sebagai admin (`POST /api/v1/auth/login` dengan `admin@imphnen.id`), dapat JWT
4. **Create campaign**: `POST /api/v1/campaigns` dengan JWT admin, kirim `{name, url}`, verify response berisi campaign data
5. **Role check**: Coba create campaign dengan JWT user biasa, verify mendapat 403
6. **Process image**: `POST /api/v1/campaigns/process-image` dengan JWT (admin atau user), kirim multipart image, verify response adalah PNG binary dengan QR di bottom-right
7. **Cache**: Process image kedua kali, verify tidak ada query DB tambahan (observable via log atau timing)
8. **No active campaign**: Process image tanpa campaign aktif, verify mendapat error yang jelas

## Decisions

- **QR Library**: `skip2/go-qrcode` — pure Go, simple API, langsung generate PNG bytes
- **Image processing**: Standard library only (`image`, `image/draw`, `image/png`, `image/jpeg`) — no external dependency needed
- **Cache strategy**: Simple single-value in-memory cache (bukan map) karena hanya 1 active campaign — `sync.RWMutex` guarded
- **Active campaign constraint**: Partial unique index di DB level + transactional toggle di repository untuk mencegah race condition
- **QR size on image**: ~1/5 dimensi terkecil image, minimum 100px, dengan 10px padding dari bottom-right edge
- **Seeder**: Shared logic di `internal/seeder/`, dipanggil dari main.go (auto) dan `cmd/seeder/main.go` (CLI)
