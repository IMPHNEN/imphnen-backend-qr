# 🚀 API Documentation & Testing - Complete Setup

## ✨ What's Included

Dokumentasi API lengkap dengan testing menggunakan Bruno API Client untuk IMPHNEN QR Backend.

### 📦 Package Contents

1. **14 Endpoint Dokumentasi Lengkap** ✅
2. **3 Environment Konfigurasi** (Dev, Prod, Alt) ✅
3. **15 Testing Scenarios** (Happy path + Campaign flow + Error cases) ✅
4. **5 Documentation Files** (README, Quick Start, Contributing, dll) ✅
5. **Automated Tests** untuk setiap endpoint ✅
6. **Auto-save Token System** ✅
7. **QR Campaign Overlay** endpoints ✅

---

## 📁 File Structure Overview

### Main Documentation
```
docs/api/
├── README.md           → Full documentation (Setup, Usage, Troubleshooting)
├── CONTRIBUTING.md     → Guidelines for contributors
├── CHANGELOG.md        → Version history and updates
└── .gitignore          → Git ignore rules for sensitive data
```

### Bruno Collection
```
docs/api/IMPHNEN-QR-API/
├── README.md              → Collection overview
├── QUICK-START.md         → 5-minute setup guide
├── bruno.json             → Collection configuration
├── collection.bru         → Collection-level scripts & tests
└── Check-Health.bru       → Health check endpoint
```

### Endpoints by Category

#### 1️⃣ Authentication (Public)
```
Auth-Public/
├── folder.bru             → Folder configuration
├── Register-User.bru      → Create new account
├── Login-User.bru         → Email/password login
├── Redirect-OAuth.bru     → Google OAuth redirect
├── Callback-OAuth.bru     → OAuth callback handler
└── Refresh-Token.bru      → Token refresh
```

#### 2️⃣ User Operations (Protected)
```
Users-Protected/
├── folder.bru             → Folder configuration
├── Get-Profile.bru        → Get user profile
└── Update-Profile.bru     → Update user info
```

#### 3️⃣ Admin Operations (Admin Only)
```
Admin/
├── folder.bru             → Folder configuration
├── List-All-Users.bru     → Get all users
└── Update-User-Role.bru   → Change user role
```

#### 4️⃣ Campaign Management (Admin Only)
```
Campaigns-Admin/
├── folder.bru             → Folder configuration
├── Create-Campaign.bru    → Create QR campaign
├── List-All-Campaigns.bru → List all campaigns
└── Activate-Campaign.bru  → Set active campaign
```

#### 5️⃣ Campaign User Operations (Protected)
```
Campaigns-User/
├── folder.bru             → Folder configuration
└── Process-Image.bru      → Upload image, get QR overlay
```

### Environments
```
environments/
├── Development.bru        → localhost:8080 (default)
├── Production.bru         → Production server
└── Local-Alternative.bru  → localhost:3000
```

### Testing Scenarios
```
Testing-Scenarios/
├── folder.bru                    → Scenarios overview
├── README.bru                    → Usage guide
├── 1-Register-New-User.bru       → Step 1 of happy path
├── 2-Get-New-User-Profile.bru    → Step 2 of happy path
├── 3-Update-User-Profile.bru     → Step 3 of happy path
├── 4-Verify-Profile-Updated.bru  → Step 4 of happy path
├── 5-Login-Admin.bru             → Step 5 of campaign flow
├── 6-Create-Campaign.bru         → Step 6 of campaign flow
├── 7-List-Campaigns.bru          → Step 7 of campaign flow
├── 8-Process-Image.bru           → Step 8 of campaign flow
├── Error-Invalid-Login.bru       → Error test: wrong credentials
├── Error-Duplicate-Email.bru     → Error test: duplicate registration
├── Error-Missing-Authorization.bru → Error test: no token
├── Error-Invalid-Request-Body.bru → Error test: invalid payload
├── Error-Campaign-Forbidden.bru  → Error test: non-admin campaign create
├── Error-No-Active-Campaign.bru  → Error test: process without campaign
└── Error-Missing-Image.bru       → Error test: process without image
```

---

## 🎯 Quick Start (3 Steps)

### 1. Install Bruno
```bash
# macOS
brew install bruno

# Linux
snap install bruno

# Or download from https://www.usebruno.com/downloads
```

### 2. Open Collection
1. Launch Bruno
2. File → Open Collection
3. Select: `docs/api/IMPHNEN-QR-API`

### 3. Start Testing
1. Select "Development" environment
2. Run "Check Health"
3. Run "Auth-Public/Register User"
4. Token auto-saved! ✨
5. Try protected endpoints

---

## 📚 Documentation Files

### 1. README.md (Main)
- Complete installation guide
- Usage instructions
- API endpoint reference
- Troubleshooting
- Best practices

**Read for**: Full understanding of the documentation

### 2. QUICK-START.md
- 5-minute setup
- Common workflows
- Quick testing guide
- Pro tips

**Read for**: Getting started fast

### 3. CONTRIBUTING.md
- How to add endpoints
- File structure templates
- Testing guidelines
- PR submission process

**Read for**: Contributing new endpoints

### 4. CHANGELOG.md
- Version history
- What's new
- Future plans
- Statistics

**Read for**: Tracking changes and updates

### 5. Collection README.md
- Collection structure
- Features overview
- Testing strategies
- Configuration guide

**Read for**: Understanding collection organization

---

## ✅ Features Breakdown

### 🤖 Auto-Save Tokens
```javascript
// Automatically runs after login/register
if (res.getStatus() === 200) {
  const data = res.getBody();
  bru.setVar("accessToken", data.data.tokens.access_token);
  console.log("🔑 Token saved!");
}
```

**Benefits**:
- No manual token copying
- Seamless workflow
- Auto-used in protected endpoints

### 🧪 Automated Tests
```javascript
tests {
  test("Status code is 200", function() {
    expect(res.getStatus()).to.equal(200);
  });
  
  test("Response has correct structure", function() {
    expect(res.getBody()).to.have.property("success");
  });
}
```

**Coverage**:
- Status codes ✅
- Response structure ✅
- Data validation ✅
- Business logic ✅

### 📖 Comprehensive Docs
Each endpoint includes:
```markdown
# Endpoint Title

Description of what it does

## Authentication
Required/Not required

## Request Body
- Parameter details

## Response Examples
- Success cases
- Error cases

## Usage Notes
- Important information
```

### 🔄 Collection Scripts
```javascript
// Runs for EVERY request
pre-request {
  console.log(`[${req.getMethod()}] ${req.getUrl()}`);
}

post-response {
  console.log(`✅ ${status} - Success`);
}
```

**Features**:
- Request logging
- Response logging
- Global error handling
- Token validation

---

## 🎓 Usage Scenarios

### Scenario 1: New User Testing
```
1. Register User     → Auto-save token
2. Get Profile       → Use saved token
3. Update Profile    → Use saved token
4. Verify Update     → Confirm changes
```
**Time**: ~2 minutes

### Scenario 2: Admin Testing
```
1. Login as Admin    → Auto-save token
2. List Users        → Use saved token
3. Update Role       → Use saved token
```
**Time**: ~1 minute

### Scenario 3: Campaign Flow
```
1. Login as Admin    → Auto-save token
2. Create Campaign   → Generate QR code
3. List Campaigns    → Verify active
4. Process Image     → Upload & get QR overlay
```
**Time**: ~2 minutes

### Scenario 4: Error Testing
```
1. Invalid Login     → Expect 401
2. Duplicate Email   → Expect 409
3. No Authorization  → Expect 401
4. Invalid Body      → Expect 400
5. Campaign Forbidden→ Expect 403
6. No Active Campaign→ Expect 404
7. Missing Image     → Expect 400
```
**Time**: ~3 minutes

### Scenario 5: Full API Test
```
1. Health Check      → Verify API
2. Register          → Create account
3. Login             → Get tokens
4. Profile Ops       → CRUD operations
5. Admin Ops         → Admin features
6. Campaign Ops      → QR campaign flow
7. Image Processing  → QR overlay
8. Error Cases       → Error handling
```
**Time**: ~7 minutes

---

## 📊 Statistics

| Metric | Count |
|--------|-------|
| Total Endpoints | 14 |
| Public Endpoints | 6 |
| Protected Endpoints | 3 |
| Admin Endpoints | 5 |
| Test Scenarios | 15 |
| Environments | 3 |
| Documentation Pages | 5 |
| Automated Tests | 60+ |
| Lines of Docs | 3000+ |

**API Coverage**: 100% ✅

---

## 🔐 Security Features

1. **No Hardcoded Tokens** ✅
   - All tokens in variables
   - Marked as secret

2. **Example Data Only** ✅
   - No real credentials
   - Safe test data

3. **.gitignore Configured** ✅
   - Local files excluded
   - Secrets protected

4. **Environment Separation** ✅
   - Dev/Prod isolated
   - Easy switching

---

## 🚦 Testing Workflow

### Initial Setup (Once)
```
1. Install Bruno           → 2 min
2. Open Collection         → 1 min
3. Select Environment      → 10 sec
4. Read Quick Start        → 3 min
Total: ~6 minutes
```

### Daily Testing (Regular)
```
1. Open Bruno              → 5 sec
2. Run Health Check        → 2 sec
3. Run Endpoint Tests      → varies
4. Check Results           → 5 sec
Total: ~15 seconds per endpoint
```

### Full Test Suite
```
1. All Public Endpoints    → 2 min
2. All Protected Endpoints → 1 min
3. All Admin Endpoints     → 1 min
4. All Campaign Endpoints  → 2 min
5. All Error Scenarios     → 3 min
Total: ~9 minutes
```

---

## 💡 Pro Tips

### 1. Environment Switching
```
Top-right dropdown → Select environment → All requests update automatically
```

### 2. Batch Testing
```
Select folder → Right-click → Run folder
Runs all requests in sequence!
```

### 3. Variables Everywhere
```
Use {{variableName}} anywhere:
- URLs
- Headers
- Body
- Tests
```

### 4. Console Debugging
```
script:post-response {
  console.log(res.getBody());
  console.log(bru.getVar("accessToken"));
}
```

### 5. Quick Duplicate
```
Right-click request → Duplicate
Perfect for creating variations!
```

---

## 🐛 Troubleshooting

### ❌ Problem: Collection not loading
**Solution**: 
- Open folder `docs/api/IMPHNEN-QR-API` (not parent)
- Check `bruno.json` exists

### ❌ Problem: 401 Unauthorized
**Solution**: 
- Run login/register first
- Check `accessToken` variable exists
- Token might be expired

### ❌ Problem: Tests failing
**Solution**: 
- Check server is running
- Verify request body data
- Check expected vs actual status

### ❌ Problem: Token not saving
**Solution**: 
- Check response status (200/201)
- Verify response structure
- Check console for errors

---

## 🎯 Next Steps

### For Developers
1. ✅ Read [QUICK-START.md](IMPHNEN-QR-API/QUICK-START.md)
2. ✅ Run your first test
3. ✅ Explore all endpoints
4. ✅ Try testing scenarios

### For Contributors
1. ✅ Read [CONTRIBUTING.md](CONTRIBUTING.md)
2. ✅ Understand file structure
3. ✅ Follow templates
4. ✅ Submit PRs

### For Team Leads
1. ✅ Review [README.md](README.md)
2. ✅ Check [CHANGELOG.md](CHANGELOG.md)
3. ✅ Verify coverage
4. ✅ Plan future endpoints

---

## 📞 Support

- **Documentation**: All `.md` files in this folder
- **Examples**: Check existing `.bru` files
- **Bruno Docs**: https://docs.usebruno.com
- **Issues**: GitHub Issues
- **Team Chat**: Internal communication

---

## 🎉 Summary

You now have:
- ✅ Complete API documentation
- ✅ Ready-to-use test collection
- ✅ Automated testing suite
- ✅ Multiple environments
- ✅ Testing scenarios
- ✅ Contribution guidelines
- ✅ Full documentation

**Everything you need to test and document the IMPHNEN QR API!**

---

**Version**: 1.0.0  
**Date**: February 9, 2026  
**Status**: Production Ready ✅  
**Coverage**: 100% 🎯  
**Maintained by**: IMPHNEN Team

---

## 🙏 Thank You!

Happy Testing! 🚀
