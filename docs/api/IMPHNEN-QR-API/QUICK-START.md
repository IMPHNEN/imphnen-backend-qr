# Quick Start Guide

## 🚀 Getting Started in 5 Minutes

### 1. Install Bruno
Download from: https://www.usebruno.com/downloads

### 2. Open This Collection
```
File → Open Collection → Select this folder
```

### 3. Select Environment
Choose "Development" from the environment dropdown (top-right corner)

### 4. Start Your Server
Make sure your backend server is running at `http://localhost:8080`

### 5. Test the API

#### First Request: Health Check ✅
- Open: `Check Health`
- Click: **Send**
- Expected: Status 200, `{"status": "ok"}`

#### Register a User 👤
- Open: `Auth-Public/Register User`
- Update request body with your details:
```json
{
  "name": "Your Name",
  "email": "your@email.com",
  "password": "YourSecurePassword123!"
}
```
- Click: **Send**
- ✨ Access token automatically saved!

#### Get Your Profile 📋
- Open: `Users-Protected/Get Profile`
- Click: **Send**
- See your user data! (Token is automatically used)

## 🎯 Common Workflows

### Testing as Regular User
1. ✅ Health Check
2. 📝 Register User (or Login)
3. 👤 Get Profile
4. ✏️ Update Profile

### Testing as Admin
1. ✅ Health Check
2. 🔐 Login as admin
3. 📋 List All Users
4. ⚙️ Update User Role

### Testing Google OAuth
1. ✅ Health Check
2. 🔗 Open Google OAuth Redirect in browser
3. 🔑 Complete Google authentication
4. ✨ Token automatically saved

### Testing Token Refresh
1. 🔐 Login/Register
2. 📋 Copy `refresh_token` from response
3. 🔄 Use Refresh Token endpoint
4. ✨ New access token automatically saved

## 📌 Important Notes

### Auto-Save Feature
After successful login/register, the access token is **automatically saved** to the environment variable `accessToken`. You don't need to copy-paste tokens manually!

### Authorization
- 🟢 **Public endpoints**: No token needed
- 🟡 **Protected endpoints**: Auto-use saved token
- 🔴 **Admin endpoints**: Need admin role + token

### Request Body
All request bodies are pre-filled with examples. Just update the values before sending.

### Tests
Each request includes automated tests. Check the "Tests" tab after sending to see validation results.

## 🆘 Troubleshooting

❌ **401 Unauthorized?**
→ Login first! Token might be expired.

❌ **403 Forbidden on admin endpoints?**
→ Your user needs admin role. Update via database:
```sql
UPDATE users SET role = 'admin' WHERE email = 'your@email.com';
```

❌ **Connection Error?**
→ Check if backend server is running on `http://localhost:8080`

## 📚 Next Steps

- Read full documentation: `../README.md`
- Explore each endpoint's documentation (click "Docs" tab)
- Check automated tests (click "Tests" tab)
- Try all request variations

## 💡 Pro Tips

1. **Environment Switching**: Create Production environment for testing against production server
2. **Pre-request Scripts**: Each request has scripts that run automatically
3. **Variables**: Use `{{variableName}}` anywhere in the request
4. **Collections**: Organize related requests in folders

Happy Testing! 🎉
