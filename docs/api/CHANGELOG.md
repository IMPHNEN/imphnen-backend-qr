# Changelog

All notable changes to the API documentation will be documented in this file.

## [1.1.0] - 2026-02-09

### ✨ QR Campaign Overlay

New campaign management and QR code overlay feature.

### ✨ Added

#### Campaign Endpoints
- **Create Campaign** (Admin) - Create QR campaign with name + URL
- **List All Campaigns** (Admin) - View all campaigns
- **Activate Campaign** (Admin) - Set active campaign
- **Process Image** (All users) - Upload image, get QR overlay PNG

#### Testing Scenarios
- Campaign Flow (Steps 5-8): Login Admin → Create Campaign → List → Process Image
- Error - Campaign Forbidden: Non-admin creating campaign
- Error - No Active Campaign: Processing image without active campaign
- Error - Missing Image: Processing without image upload

#### Features
- 🎯 Auto-seeder creates demo admin + user accounts on startup
- 📷 QR code overlay on uploaded images (bottom-right, ~1/5 size)
- 💾 In-memory QR cache for active campaign
- 🔒 Partial unique index ensures max 1 active campaign
- 🧪 Campaign ID auto-saved to environment variables

#### Documentation
- Bruno endpoint docs for all 4 campaign endpoints
- Campaigns-Admin folder with 3 admin endpoints
- Campaigns-User folder with process-image endpoint
- Updated QUICK-START, README, OVERVIEW with campaign info
- Updated all environments with campaignId variable

---

## [1.0.0] - 2026-02-09

### 🎉 Initial Release

Complete API documentation and testing collection for IMPHNEN QR Backend.

### ✨ Added

#### Endpoints Documentation
- **Health Check** - API health verification
- **Auth Endpoints (Public)**
  - Register User
  - Login User
  - Google OAuth Redirect
  - Google OAuth Callback
  - Refresh Token
- **User Endpoints (Protected)**
  - Get User Profile
  - Update User Profile
- **Admin Endpoints**
  - List All Users
  - Update User Role

#### Features
- 🤖 Auto-save tokens after authentication
- 🧪 Automated tests for all endpoints
- 📖 Comprehensive documentation with examples
- 🔄 Multiple environment support
- 🎯 Pre-configured testing scenarios
- 📝 Collection-level scripts for logging
- ✅ Response validation tests
- 🔐 Proper authentication handling

#### Environments
- Development (localhost:8080)
- Production (configurable)
- Local Alternative (localhost:3000)

#### Testing Scenarios
- Happy Path - New User (4-step workflow)
- Error Scenarios:
  - Invalid Login
  - Duplicate Email
  - Missing Authorization
  - Invalid Request Body

#### Documentation Files
- Main README with full documentation
- Quick Start Guide for 5-minute setup
- Collection README with detailed structure
- Contributing Guide for contributors
- Testing Scenarios Guide

#### Configuration
- Collection-level pre-request scripts
- Collection-level post-response scripts
- Collection-level tests
- Auto-logging for all requests
- Token management automation

### 🎨 Structure

```
docs/api/
├── README.md                      # Main documentation
├── CONTRIBUTING.md                # Contribution guidelines
├── .gitignore                     # Git ignore rules
└── IMPHNEN-QR-API/
    ├── README.md                  # Collection overview
    ├── QUICK-START.md             # Quick start guide
    ├── bruno.json                 # Collection config
    ├── collection.bru             # Collection scripts
    ├── Check-Health.bru           # Health check
    ├── environments/
    │   ├── Development.bru        # Dev environment
    │   ├── Production.bru         # Prod environment
    │   └── Local-Alternative.bru  # Alt local env
    ├── Auth-Public/
    │   ├── folder.bru
    │   ├── Register-User.bru
    │   ├── Login-User.bru
    │   ├── Redirect-OAuth.bru
    │   ├── Callback-OAuth.bru
    │   └── Refresh-Token.bru
    ├── Users-Protected/
    │   ├── folder.bru
    │   ├── Get-Profile.bru
    │   └── Update-Profile.bru
    ├── Admin/
    │   ├── folder.bru
    │   ├── List-All-Users.bru
    │   └── Update-User-Role.bru
    ├── Campaigns-Admin/
    │   ├── folder.bru
    │   ├── Create-Campaign.bru
    │   ├── List-All-Campaigns.bru
    │   └── Activate-Campaign.bru
    ├── Campaigns-User/
    │   ├── folder.bru
    │   └── Process-Image.bru
    └── Testing-Scenarios/
        ├── folder.bru
        ├── README.bru
        ├── 1-Register-New-User.bru
        ├── 2-Get-New-User-Profile.bru
        ├── 3-Update-User-Profile.bru
        ├── 4-Verify-Profile-Updated.bru
        ├── Error-Invalid-Login.bru
        ├── Error-Duplicate-Email.bru
        ├── Error-Missing-Authorization.bru
        └── Error-Invalid-Request-Body.bru
```

### 📊 Statistics

- **Total Endpoints**: 14
- **Total Test Requests**: 21 (including scenarios)
- **Total Environments**: 3
- **Test Scenarios**: 8 happy path/campaign flow + 7 error cases
- **Documentation Files**: 5 markdown files
- **Lines of Documentation**: ~3000+ lines

### 🔒 Security

- No sensitive data in repository
- Tokens stored in environment variables
- Example data only (no real credentials)
- Secrets marked as secret in environments

### 🎯 Coverage

API Endpoint Coverage: **100%**
- ✅ All public endpoints documented
- ✅ All protected endpoints documented
- ✅ All admin endpoints documented
- ✅ All authentication flows covered
- ✅ Error scenarios included

### 📝 Documentation Quality

Each endpoint includes:
- ✅ Clear description
- ✅ Authentication requirements
- ✅ Request body examples
- ✅ Response examples
- ✅ Error scenarios
- ✅ Usage notes
- ✅ Automated tests
- ✅ Post-response scripts

### 🧪 Testing Coverage

Each endpoint has tests for:
- ✅ Status code validation
- ✅ Response structure
- ✅ Required fields
- ✅ Data types
- ✅ Business logic
- ✅ Token management

---

## Future Enhancements

<!-- ### Planned for v1.2.0
- [ ] Event management endpoints
- [ ] Analytics endpoints
- [ ] Webhook endpoints
- [ ] Bulk operations
- [ ] Search and filtering
- [ ] Pagination examples
- [ ] Rate limiting documentation
- [ ] WebSocket documentation (if applicable)

### Ideas for v2.0.0
- [ ] GraphQL documentation (if implemented)
- [ ] API versioning examples
- [ ] Advanced testing scenarios
- [ ] Performance testing requests
- [ ] Load testing scenarios
- [ ] CI/CD integration examples
- [ ] Mock server configuration
- [ ] Contract testing -->

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on adding new endpoints.

## Version Format

We follow [Semantic Versioning](https://semver.org/):
- MAJOR: Breaking changes in API
- MINOR: New endpoints or features
- PATCH: Bug fixes and documentation updates

## Tags

- `Added` - New endpoints or features
- `Changed` - Changes in existing endpoints
- `Deprecated` - Soon-to-be removed features
- `Removed` - Removed features
- `Fixed` - Bug fixes
- `Security` - Security improvements
