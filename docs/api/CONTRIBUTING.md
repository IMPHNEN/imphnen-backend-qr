# Contributing to API Documentation

Thank you for contributing to IMPHNEN QR API documentation! This guide will help you add or update API endpoints in our Bruno collection.

## 📋 Prerequisites

- Bruno API Client installed
- IMPHNEN QR Backend running locally
- Understanding of the API endpoint you're documenting
- Git and GitHub knowledge

## 🎯 Adding a New Endpoint

### 1. Determine the Category

Place your endpoint in the appropriate folder:
- `Auth-Public/` - Authentication endpoints (no auth required)
- `Users-Protected/` - User endpoints (bearer token required)
- `Admin/` - Admin endpoints (bearer token + admin role)
- Create new folder if needed for new features

### 2. Create the .bru File

**Naming Convention**: Use kebab-case with descriptive names
```
Good: Create-Event.bru, Delete-QR-Code.bru
Bad: endpoint1.bru, test.bru
```

### 3. File Structure Template

```plaintext
meta {
  name: Your Endpoint Name
  type: http
  seq: 1
}

[method] {
  url: {{baseUrl}}/api/v1/your/endpoint
  body: json|none
  auth: bearer|none
}

auth:bearer {
  token: {{accessToken}}
}

headers {
  Content-Type: application/json
  Authorization: Bearer {{accessToken}}
}

body:json {
  {
    "field": "value"
  }
}

script:post-response {
  // Auto-save important values
  if (res.getStatus() === 200) {
    const data = res.getBody();
    // Save to variable if needed
    // bru.setVar("variableName", data.value);
  }
}

tests {
  test("Status code is correct", function() {
    expect(res.getStatus()).to.equal(200);
  });
  
  test("Response has success structure", function() {
    const body = res.getBody();
    expect(body).to.have.property("success", true);
    expect(body).to.have.property("message");
    expect(body).to.have.property("data");
  });
  
  // Add more specific tests
}

docs {
  # Endpoint Title
  
  Brief description of what this endpoint does.
  
  ## Authentication
  [Required/Not Required] - Explain auth requirements
  
  ## Request Body
  - `field1` (type, required/optional): Description
  - `field2` (type, required/optional): Description
  
  ## Success Response (200)
  ```json
  {
    "success": true,
    "message": "operation successful",
    "data": {
      // expected response structure
    }
  }
  ```
  
  ## Error Responses
  - 400: Bad Request - Invalid input
  - 401: Unauthorized - Missing/invalid token
  - 404: Not Found - Resource not found
  
  ## Notes
  Any additional information or edge cases
}
```

### 4. Essential Components

#### Meta Section
```plaintext
meta {
  name: Descriptive Name (Title Case)
  type: http
  seq: [number] # Order in folder
}
```

#### HTTP Method Section
Choose appropriate method: `get`, `post`, `put`, `delete`, `patch`

```plaintext
post {
  url: {{baseUrl}}/api/v1/endpoint
  body: json  # or 'none' for GET/DELETE
  auth: bearer  # or 'none' for public endpoints
}
```

#### Headers (for authenticated endpoints)
```plaintext
headers {
  Content-Type: application/json
  Authorization: Bearer {{accessToken}}
}
```

#### Request Body (for POST/PUT/PATCH)
```plaintext
body:json {
  {
    "field": "example_value"
  }
}
```

#### Tests (REQUIRED)
At minimum, include these tests:

```javascript
tests {
  test("Status code is correct", function() {
    expect(res.getStatus()).to.equal(200);
  });
  
  test("Response structure is valid", function() {
    const body = res.getBody();
    expect(body).to.have.property("success");
    expect(body).to.have.property("message");
  });
  
  test("Response data exists", function() {
    const data = res.getBody().data;
    expect(data).to.exist;
  });
}
```

#### Documentation (REQUIRED)
```plaintext
docs {
  # Clear Title
  
  Description of endpoint functionality.
  
  ## Required sections:
  - Authentication requirements
  - Request parameters/body
  - Success response example
  - Error responses
  - Usage notes
}
```

## 🧪 Testing Your Endpoint

### Before Committing

1. **Test the request**
   - Ensure it works with actual backend
   - Verify all fields in request body
   - Check response matches documentation

2. **Validate all tests**
   - All tests should pass (green checkmarks)
   - Add edge case tests
   - Test error scenarios

3. **Check auto-save scripts**
   - Verify tokens/IDs are saved correctly
   - Test variable usage in dependent requests

4. **Review documentation**
   - Clear and complete
   - Examples are accurate
   - No sensitive data (passwords, real tokens)

## 📝 Documentation Standards

### Writing Style
- Use clear, concise language
- Write in present tense
- Use active voice
- Be specific about requirements

### Code Examples
- Use realistic but safe example data
- Include all required fields
- Show expected response structure
- Format JSON properly

### Parameter Documentation
Format: `name` (type, required/optional): Description

Example:
```
- `email` (string, required): User's email address
- `name` (string, optional): User's display name
```

## 🔄 Updating Existing Endpoints

### When Backend Changes
1. Update request body/parameters
2. Update response examples
3. Update tests if response structure changed
4. Update documentation
5. Test thoroughly

### Version Changes
If endpoint behavior changes significantly:
1. Consider versioning in URL
2. Update documentation with version notes
3. Mark deprecated fields clearly

## 🎨 Best Practices

### Request Bodies
- Use realistic example values
- Include all available fields
- Comment optional fields
- Use safe test data (no real emails/passwords)

### Tests
- Test status codes
- Validate response structure
- Check data types
- Verify business logic
- Test edge cases

### Scripts
- Keep post-response scripts simple
- Only save variables that will be reused
- Add console.log for debugging
- Handle errors gracefully

### Documentation
- Start with a clear summary
- List all parameters with types
- Show complete response examples
- Document all error codes
- Add usage notes for complex endpoints

## 🚫 Common Mistakes to Avoid

1. **Don't hard-code URLs** - Always use `{{baseUrl}}`
2. **Don't commit tokens** - Use variables, not actual tokens
3. **Don't skip tests** - Every endpoint needs tests
4. **Don't skip documentation** - Docs are just as important as code
5. **Don't use real user data** - Use test data only
6. **Don't forget error scenarios** - Document all possible errors

## ✅ Checklist Before Submitting

- [ ] File named with clear, descriptive name
- [ ] Placed in correct folder
- [ ] URL uses `{{baseUrl}}` variable
- [ ] Auth configured correctly
- [ ] Request body has example data
- [ ] All tests pass
- [ ] Tests cover main scenarios
- [ ] Documentation is complete
- [ ] Response examples are accurate
- [ ] Error scenarios documented
- [ ] No sensitive data in file
- [ ] Tested against running backend
- [ ] Variables auto-save where needed
- [ ] Sequence number set appropriately

## 📤 Submitting Your Changes

### Git Workflow

1. **Create a branch**
   ```bash
   git checkout -b docs/add-endpoint-name
   ```

2. **Add your changes**
   ```bash
   git add docs/api/IMPHNEN-QR-API/
   ```

3. **Commit with clear message**
   ```bash
   git commit -m "docs: Add/Update [Endpoint Name] documentation"
   ```

4. **Push and create PR**
   ```bash
   git push origin docs/add-endpoint-name
   ```

### PR Description Template

```markdown
## Changes
- Added/Updated [Endpoint Name] documentation

## Endpoint Details
- Method: [GET/POST/PUT/DELETE]
- Path: `/api/v1/...`
- Auth Required: [Yes/No]

## Testing
- [ ] Tested against running backend
- [ ] All tests pass
- [ ] Documentation reviewed
- [ ] Examples validated

## Checklist
- [ ] Follows naming conventions
- [ ] Includes automated tests
- [ ] Has complete documentation
- [ ] No sensitive data included
```

## 🔍 Review Process

Your PR will be reviewed for:
- Correctness of endpoint documentation
- Quality of tests
- Clarity of documentation
- Code style consistency
- No sensitive data

## 💡 Getting Help

- Check existing endpoints for examples
- Read [Bruno Documentation](https://docs.usebruno.com)
- Ask in team chat or GitHub discussions
- Review [Quick Start Guide](IMPHNEN-QR-API/QUICK-START.md)

## 🎓 Learning Resources

- [Bruno Official Docs](https://docs.usebruno.com)
- [API Testing Best Practices](https://docs.usebruno.com/testing/introduction)
- [JavaScript Test Assertions](https://www.chaijs.com/api/bdd/)
- [HTTP Status Codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)

Thank you for contributing! 🙏
