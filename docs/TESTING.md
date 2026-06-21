# Testing Guide

This document provides comprehensive information about testing practices for the APP system log management dashboard.

## 🧪 Overview

The APP system log management dashboard follows a comprehensive testing strategy that includes unit testing, integration testing, and end-to-end testing to ensure reliability, security, and performance.

## 📋 Testing Strategy

### Test Types

1. **Unit Tests** - Test individual functions and components
2. **Integration Tests** - Test interactions between components
3. **End-to-End Tests** - Test complete user workflows
4. **Security Tests** - Test authentication, authorization, and data protection
5. **Performance Tests** - Test system performance under load

## 🛠️ Testing Frameworks

### Backend Testing

#### Go Testing Framework

The backend uses Go's built-in testing framework with additional libraries:

```go
// Example test structure
func TestSearchLogs(t *testing.T) {
    // Setup
    // Execute
    // Assert
    // Teardown
}
```

#### Testing Libraries

- `testing` - Go's built-in testing package
- `testify` - Assertion and mocking library
- `httptest` - HTTP test helpers
- `gomega` - Behavior-driven testing

### Frontend Testing

#### Astro/TypeScript Testing

The frontend uses:

- `vitest` - Fast test runner
- `@testing-library/react` - React testing utilities
- `jsdom` - DOM environment for testing
- `cypress` - End-to-end testing framework

## 🧪 Backend Testing

### Unit Tests

#### Authentication Tests

```go
func TestLogin(t *testing.T) {
    // Test valid login
    // Test invalid login
    // Test token generation
    // Test token validation
}
```

#### Search Tests

```go
func TestSearchLogs(t *testing.T) {
    // Test valid search query
    // Test invalid search query
    // Test search with pagination
    // Test search with sorting
}
```

#### Export Tests

```go
func TestExportLogs(t *testing.T) {
    // Test CSV export functionality
    // Test PDPA field masking
    // Test export file creation
    // Test export file permissions
}
```

### Integration Tests

#### API Integration Tests

```go
func TestAPIEndpoints(t *testing.T) {
    // Test authenticated endpoints
    // Test unauthorized access
    // Test rate limiting
    // Test error handling
}
```

#### Database Integration Tests

```go
func TestDatabaseOperations(t *testing.T) {
    // Test QuickWit integration
    // Test Elasticsearch operations
    // Test data consistency
    // Test connection handling
}
```

### Security Tests

#### Authentication Security Tests

```go
func TestJWTSecurity(t *testing.T) {
    // Test token expiration
    // Test token manipulation
    // Test role-based access control
    // Test session management
}
```

#### Input Validation Tests

```go
func TestInputValidation(t *testing.T) {
    // Test SQL injection prevention
    // Test XSS protection
    // Test path traversal prevention
    // Test parameter validation
}
```

## 🧪 Frontend Testing

### Component Tests

#### Dashboard Components

```typescript
describe('Dashboard Component', () => {
    test('renders dashboard correctly', () => {
        // Test component rendering
        // Test data display
        // Test user interactions
    });
});
```

#### Search Components

```typescript
describe('Search Component', () => {
    test('executes search query', () => {
        // Test query execution
        // Test result display
        // Test error handling
    });
});
```

### Integration Tests

#### API Integration Tests

```typescript
describe('API Integration', () => {
    test('fetches dashboard data', async () => {
        // Test API calls
        // Test data handling
        // Test error responses
    });
});
```

### End-to-End Tests

#### Cypress Tests

```javascript
describe('Application Flow', () => {
    it('should login and search logs', () => {
        cy.visit('/login');
        cy.get('[data-testid="email"]').type('admin@example.com');
        cy.get('[data-testid="password"]').type('admin123');
        cy.get('[data-testid="login-button"]').click();
        
        cy.get('[data-testid="search-input"]').type('ERROR');
        cy.get('[data-testid="search-button"]').click();
        
        cy.contains('Log entries found');
    });
});
```

## 📊 Test Coverage

### Coverage Targets

- **Backend**: 80% code coverage
- **Frontend**: 70% code coverage
- **Security**: 100% test coverage for security-critical functions

### Coverage Tools

```bash
# Backend coverage
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Frontend coverage
cd frontend
npm run test:coverage
```

## 🚀 Test Execution

### Local Testing

#### Backend Tests

```bash
cd backend
go test ./...
go test -v ./...
go test -race ./...  # Run with race detector
```

#### Frontend Tests

```bash
cd frontend
npm test
npm run test:coverage
```

### Continuous Integration

#### GitHub Actions

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Setup Go
      uses: actions/setup-go@v2
      with:
        go-version: '1.21'
    - name: Run backend tests
      run: cd backend && go test ./...
    - name: Run frontend tests
      run: cd frontend && npm test
```

## 📈 Performance Testing

### Load Testing

#### Using Apache Bench

```bash
# Test API endpoints
ab -n 1000 -c 10 http://localhost:8080/api/dashboard
```

#### Using wrk

```bash
# Load testing with wrk
wrk -t12 -c400 -d30s http://localhost:8080/api/dashboard
```

### Stress Testing

```go
func TestStressSearch(t *testing.T) {
    // Test concurrent search operations
    // Test with large datasets
    // Test memory usage
}
```

## 🛡️ Security Testing

### Vulnerability Scanning

```bash
# Scan dependencies
cd backend
go list -m all | grep -E "(vuln|security)" | gosec

# Test for common vulnerabilities
cd frontend
npm audit
```

### Penetration Testing

#### Manual Testing

- Test authentication bypass
- Test authorization testing
- Test input validation
- Test session management

#### Automated Testing

```go
func TestSecurityEndpoints(t *testing.T) {
    // Test for SQL injection
    // Test for XSS
    // Test for path traversal
    // Test for CSRF
}
```

## 📋 Test Data Management

### Test Datasets

#### Sample Log Data

```json
{
  "timestamp": "2023-01-01T00:00:00Z",
  "level": "INFO",
  "message": "Application started successfully",
  "source": "app-server",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0..."
}
```

### Test Database Setup

```bash
# Setup test database
docker run -d --name test-quickwit \
  -p 7280:7280 \
  quickwit/quickwit:latest run --config /config/quickwit.yaml --data-dir /data
```

## 📋 Test Environment

### Development Environment

```bash
# Start test environment
docker-compose -f docker-compose.test.yml up
```

### Test Configuration

```env
# Test environment variables
TEST_MODE=true
JWT_SECRET=test-secret-key
QUICKWIT_URL=http://localhost:7280
DATABASE_URL=postgresql://test:test@localhost:5432/test_db
```

## 📊 Test Reporting

### Test Results

```bash
# Generate test reports
cd backend
go test -v ./... -json > test-report.json

# Run with coverage
cd frontend
npm run test:coverage -- --watchAll=false
```

### Code Coverage Reports

```bash
# Backend coverage
go tool cover -html=coverage.out -o coverage.html

# Frontend coverage
npm run test:coverage
```

## 📋 Testing Best Practices

### Test Organization

- Group tests by functionality
- Use descriptive test names
- Keep tests independent
- Use fixtures for test data

### Test Maintenance

- Regularly update tests with code changes
- Remove obsolete tests
- Refactor tests for clarity
- Monitor test performance

### Test Reliability

- Use deterministic test data
- Mock external dependencies
- Handle timeouts gracefully
- Ensure test isolation

## 📚 Resources

- [Go Testing Documentation](https://pkg.go.dev/testing)
- [Astro Testing Guide](https://docs.astro.build/en/guides/testing/)
- [Cypress Documentation](https://docs.cypress.io/)
- [Security Testing Best Practices](https://owasp.org/www-project-testing/)