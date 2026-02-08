# IMPHNEN QR API Collection

Complete API documentation and testing collection for IMPHNEN QR Backend.

## 📁 Collection Structure

```
IMPHNEN-QR-API/
├── Check-Health.bru              # Health check endpoint
├── Auth-Public/                  # Authentication endpoints (public)
│   ├── Register-User.bru
│   ├── Login-User.bru
│   ├── Redirect-OAuth.bru
│   ├── Callback-OAuth.bru
│   └── Refresh-Token.bru
├── Users-Protected/              # User endpoints (auth required)
│   ├── Get-Profile.bru
│   └── Update-Profile.bru
├── Admin/                        # Admin endpoints (admin role required)
│   ├── List-All-Users.bru
│   └── Update-User-Role.bru
└── environments/                 # Environment configurations
    ├── Development.bru           # localhost:8080
    ├── Local-Alternative.bru     # localhost:3000
    └── Production.bru           # Production server
```

## 🎯 Features

### ✅ Complete API Coverage
All backend endpoints are documented and testable.

### 🤖 Auto-Save Tokens
Access tokens are automatically saved after login/register and used in protected endpoints.

### 🧪 Automated Tests
Every endpoint includes response validation tests:
- Status code checks
- Response structure validation
- Data integrity verification
- Business logic assertions

### 📖 Comprehensive Documentation
Each endpoint includes:
- Detailed descriptions
- Request/response examples
- Error scenarios
- Usage instructions

### 🔐 Environment-based Configuration
Switch between Development/Production environments easily.

## 🚀 Quick Start

See [QUICK-START.md](QUICK-START.md) for 5-minute setup guide.

## 📋 Endpoint Categories

### Public Endpoints
No authentication required:
- Health Check
- User Registration
- User Login
- Google OAuth
- Token Refresh

### Protected Endpoints
Requires Bearer token:
- Get User Profile
- Update User Profile

### Admin Endpoints
Requires Bearer token + Admin role:
- List All Users
- Update User Role

## 🔑 Authentication

### Token Flow
```
Register/Login → Access Token (auto-saved) → Use in Protected Endpoints
```

### Token Variables
- `accessToken`: Automatically saved after login/register
- Used automatically in all protected endpoints
- Refresh using dedicated refresh endpoint

### Authorization Levels
1. **None**: Public endpoints
2. **Bearer Token**: Protected user endpoints
3. **Bearer Token + Admin Role**: Admin endpoints

## 🧪 Testing Strategies

### 1. Smoke Test
Quick test to verify API is working:
```
Health Check → Register → Get Profile
```

### 2. Full User Journey
Complete user workflow:
```
Register → Login → Get Profile → Update Profile
```

### 3. Admin Workflow
Admin operations:
```
Login (admin) → List Users → Update User Role
```

### 4. Token Lifecycle
Token management:
```
Login → Save Token → Use Token → Refresh Token → Use New Token
```

### 5. OAuth Flow
Google authentication:
```
Redirect to Google → User Consent → Callback → Token Saved
```

## 📊 Test Results

All tests include assertions for:
- ✅ Correct HTTP status codes
- ✅ Valid response structure
- ✅ Required fields present
- ✅ Data type validation
- ✅ Business logic correctness

## 🔧 Configuration

### Environment Variables

| Variable | Description | Development | Production |
|----------|-------------|-------------|------------|
| `baseUrl` | API base URL | `http://localhost:8080` | `https://api.imphnen.com` |
| `accessToken` | JWT token | Auto-filled | Auto-filled |

### Adding New Environment
1. Duplicate existing environment file
2. Update `baseUrl`
3. Select new environment in Bruno

## 📝 Request Format

All requests follow consistent patterns:

### Headers
```
Content-Type: application/json
Authorization: Bearer {{accessToken}}  // For protected endpoints
```

### Body (POST/PUT)
```json
{
  "field": "value"
}
```

### Response
```json
{
  "success": true,
  "message": "operation successful",
  "data": { ... }
}
```

## 🎓 Best Practices

### 1. Environment Management
- Use Development for local testing
- Use Production for integration testing
- Never commit sensitive tokens

### 2. Token Management
- Let Bruno auto-save tokens
- Refresh tokens before expiry
- Re-login if refresh fails

### 3. Test Execution
- Run Health Check first
- Follow logical flow (register → login → use)
- Check test results after each request

### 4. Data Cleanup
- Use unique emails for testing
- Clean up test data periodically
- Don't use production data in development

## 🐛 Common Issues

### Issue: Token not saving
**Cause**: Login/register failed  
**Fix**: Check response status and fix request body

### Issue: 401 Unauthorized
**Cause**: Token expired or missing  
**Fix**: Login again to get new token

### Issue: 403 Forbidden
**Cause**: Missing admin role  
**Fix**: Update user role in database

### Issue: Connection refused
**Cause**: Backend server not running  
**Fix**: Start backend server on correct port

## 📚 Additional Resources

- [Full Documentation](../README.md)
- [Quick Start Guide](QUICK-START.md)
- [Backend Repository](https://github.com/IMPHNEN/imphnen-backend-qr)
- [Bruno Documentation](https://docs.usebruno.com)

## 🤝 Contributing

To add new endpoints:

1. Create `.bru` file in appropriate folder
2. Follow existing file structure:
   - `meta`: Metadata
   - `get/post/put/delete`: HTTP method and URL
   - `headers`: Request headers
   - `body:json`: Request body (if applicable)
   - `tests`: Response assertions
   - `docs`: Documentation
3. Add automated tests
4. Update this README if needed

## 📄 License

MIT License - See main repository

---

**Version**: 1.0.0  
**Last Updated**: February 2026  
**Maintainer**: IMPHNEN Team
