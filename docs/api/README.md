# IMPHNEN QR API Documentation

Dokumentasi lengkap API IMPHNEN-QR menggunakan Bruno API Client.

## 📋 Daftar Isi

- [Tentang](#tentang)
- [Instalasi Bruno](#instalasi-bruno)
- [Cara Menggunakan](#cara-menggunakan)
- [Struktur Endpoint](#struktur-endpoint)
- [Environment Variables](#environment-variables)
- [Testing Flow](#testing-flow)

## 🎯 Tentang

Dokumentasi ini menggunakan [Bruno](https://www.usebruno.com/), sebuah API client yang open-source dan Git-friendly. Semua koleksi API disimpan dalam format text files yang mudah di-version control.

## 🚀 Instalasi Bruno

### Desktop App
Download dan install Bruno dari [website resmi](https://www.usebruno.com/downloads).

### Via Package Manager

**macOS (Homebrew):**
```bash
brew install bruno
```

**Linux (Snap):**
```bash
snap install bruno
```

**Windows (Chocolatey):**
```bash
choco install bruno
```

## 📖 Cara Menggunakan

### 1. Open Collection
1. Buka Bruno
2. Click "Open Collection"
3. Navigate ke folder: `docs/api/IMPHNEN-QR-API`
4. Collection akan ter-load dengan semua endpoints

### 2. Setup Environment
1. Pilih environment "Development" di pojok kanan atas
2. Pastikan `baseUrl` sudah diset ke: `http://localhost:8080`
3. Variable `accessToken` akan otomatis terisi setelah login/register

### 3. Testing Workflow
Ikuti urutan berikut untuk testing lengkap:

#### Step 1: Health Check
```
GET /health
```
Memastikan server berjalan dengan baik.

#### Step 2: Register User
```
POST /api/v1/auth/register
```
- Input email, password, dan name
- Access token akan otomatis tersimpan
- Gunakan untuk testing sebagai user biasa

#### Step 3: Login User
```
POST /api/v1/auth/login
```
- Login dengan credentials yang sudah didaftarkan
- Access token akan otomatis diupdate

#### Step 4: Test Protected Endpoints
Setelah login, test endpoint protected:
- `GET /api/v1/users/me` - View profile
- `PUT /api/v1/users/me` - Update profile

#### Step 5: Test Admin Endpoints
⚠️ **Perhatian**: Endpoint admin memerlukan user dengan role 'admin'

Untuk testing admin endpoints:
1. Login sebagai admin, atau
2. Update role user melalui database/migration

Admin endpoints:
- `GET /api/v1/users` - List semua users
- `PUT /api/v1/users/:id/role` - Update role user

## 🗂️ Struktur Endpoint

### Public Endpoints (Auth-Public)
Tidak memerlukan authentication:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register user baru |
| POST | `/api/v1/auth/login` | Login dengan email/password |
| GET | `/api/v1/auth/google` | Redirect ke Google OAuth |
| GET | `/api/v1/auth/google/callback` | Handle Google OAuth callback |
| POST | `/api/v1/auth/refresh` | Refresh access token |

### Protected Endpoints (Users-Protected)
Memerlukan authentication (Bearer token):

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get user profile |
| PUT | `/api/v1/users/me` | Update user profile |

### Admin Endpoints (Admin)
Memerlukan authentication + role admin:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List semua users |
| PUT | `/api/v1/users/:id/role` | Update user role |

## 🔐 Environment Variables

### Development Environment

| Variable | Description | Auto-filled |
|----------|-------------|-------------|
| `baseUrl` | Base API URL | ❌ Manual |
| `accessToken` | JWT access token | ✅ Auto (after login) |

### Cara Kerja Auto-fill Token

Setiap request di folder **Auth-Public** yang berhasil (login/register) akan otomatis menyimpan access token ke variable `accessToken`. Token ini kemudian digunakan otomatis oleh semua request di folder **Users-Protected** dan **Admin**.

Script yang dijalankan setelah login/register:
```javascript
if (res.getStatus() === 200 || res.getStatus() === 201) {
  const data = res.getBody();
  if (data.data && data.data.tokens) {
    bru.setVar("accessToken", data.data.tokens.access_token);
  }
}
```

## 🧪 Testing Flow

### Basic User Flow
```
1. Register → Auto-save token
2. Get Profile → Uses saved token
3. Update Profile → Uses saved token
```

### Admin Flow
```
1. Login as admin → Auto-save token
2. List All Users → Uses saved token
3. Update User Role → Uses saved token
```

### Token Refresh Flow
```
1. Save refresh_token from login response
2. Use Refresh Token endpoint
3. New access token auto-saved
4. Continue with protected endpoints
```

## ✅ Automated Tests

Setiap endpoint dilengkapi dengan automated assertions:

### Response Structure Tests
- Status code validation
- Response body structure
- Required fields validation

### Business Logic Tests
- Token auto-save verification
- Data integrity checks
- Authorization checks

### Example Test Output
```javascript
✓ Status code is 200
✓ Response has success structure
✓ Response contains user and tokens
✓ Access token is saved
```

## 📝 Response Format

Semua endpoint mengikuti format response yang konsisten:

### Success Response
```json
{
  "success": true,
  "message": "operation successful",
  "data": {
    // response data
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "error message",
  "error": "error_code",
  "data": null
}
```

## 🔧 Troubleshooting

### Problem: "401 Unauthorized" pada protected endpoints
**Solution**: 
- Pastikan sudah login/register terlebih dahulu
- Check variable `accessToken` terisi dengan benar
- Token mungkin expired, coba login ulang

### Problem: "403 Forbidden" pada admin endpoints  
**Solution**:
- Pastikan user memiliki role 'admin'
- Update role melalui database jika perlu:
```sql
UPDATE users SET role = 'admin' WHERE email = 'your@email.com';
```

### Problem: Collection tidak muncul
**Solution**:
- Pastikan membuka folder yang benar: `docs/api/IMPHNEN-QR-API`
- Check file `bruno.json` ada di folder tersebut

## 📚 Resources

- [Bruno Documentation](https://docs.usebruno.com/)
- [Bruno GitHub](https://github.com/usebruno/bruno)
- [API Backend Repository](https://github.com/IMPHNEN/imphnen-backend-qr)

## 🤝 Contributing

Untuk menambah atau mengupdate endpoint:

1. Buat/edit file `.bru` di folder yang sesuai
2. Ikuti struktur yang ada (meta, headers, body, tests, docs)
3. Tambahkan automated tests
4. Update README ini jika perlu

## 📄 License

MIT License - see main repository for details
