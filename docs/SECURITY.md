# Security Documentation

This document outlines the security practices, policies, and implementation details for the APP system log management dashboard.

## 🛡️ Overview

The APP system log management dashboard implements comprehensive security measures to protect sensitive log data, ensure user privacy, and maintain system integrity. Security is a fundamental aspect of the application's design and implementation.

## 🔒 Security Architecture

### 1. Authentication and Authorization

#### JWT-based Authentication
- **Token Generation**: Secure JWT tokens with expiration
- **Token Storage**: Tokens stored securely in browser (localStorage/sessionStorage)
- **Token Refresh**: Secure token refresh mechanism
- **Token Validation**: Middleware validation on all protected routes

#### Role-Based Access Control (RBAC)
- **User Roles**: Admin, User, Guest
- **Access Levels**: Different permissions based on role
- **Route Protection**: All routes secured with appropriate middleware
- **Resource Access**: Fine-grained access control

### 2. Data Protection

#### Log Data Security
- **PDPA Field Masking**: Sensitive fields automatically masked in exports
- **Data Encryption**: HTTPS for data in transit
- **Access Logging**: Audit trails of all data access
- **Data Retention**: Configurable data retention policies

#### Sensitive Field Masking
```go
// Example PDPA field masking
func maskSensitiveFields(logEntry *LogEntry) {
    if logEntry.IPAddress != "" {
        logEntry.IPAddress = maskIP(logEntry.IPAddress)
    }
    if logEntry.UserAgent != "" {
        logEntry.UserAgent = maskUserAgent(logEntry.UserAgent)
    }
    // Mask other sensitive fields...
}
```

### 3. Network Security

#### CORS Configuration
- **Allowed Origins**: Strict control over cross-origin requests
- **Allowed Methods**: Whitelisted HTTP methods
- **Allowed Headers**: Controlled request headers
- **Credentials**: Secure credential handling

#### HTTPS Enforcement
- **TLS 1.3**: Modern encryption standards
- **Certificate Management**: Automated certificate handling
- **Secure Headers**: HTTP security headers
- **Redirects**: Automatic HTTP to HTTPS redirects

## 🛡️ Security Implementation

### Backend Security

#### Authentication Middleware
```go
// JWT authentication middleware
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }
        
        // Validate JWT token
        claims, err := validateJWT(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

#### Input Validation
```go
// Input validation for search queries
func validateSearchQuery(query string) error {
    if len(query) > 1000 {
        return errors.New("query too long")
    }
    
    // Check for malicious patterns
    if strings.Contains(query, "DROP") || strings.Contains(query, "DELETE") {
        return errors.New("invalid query characters")
    }
    
    return nil
}
```

#### Rate Limiting
```go
// Rate limiting middleware
func RateLimit() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(10), 100) // 10 requests per second
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### Frontend Security

#### Secure Storage
```javascript
// Secure token handling
const setToken = (token) => {
    // Store in secure HTTP-only cookie or localStorage
    localStorage.setItem('auth_token', token);
};

const getToken = () => {
    return localStorage.getItem('auth_token');
};
```

#### XSS Protection
```javascript
// XSS prevention in frontend
const sanitizeHTML = (html) => {
    const div = document.createElement('div');
    div.textContent = html;
    return div.innerHTML;
};
```

## 🛡️ Security Features

### 1. Authentication Security

#### JWT Security
- **Secret Key**: Strong, randomly generated secret key
- **Expiration**: Short-lived tokens with refresh mechanism
- **Token Revocation**: Support for token invalidation
- **Secure Flags**: Secure, HttpOnly, SameSite flags on tokens

#### Session Management
- **Secure Sessions**: Session-based authentication
- **Timeout Handling**: Automatic session timeout
- **Multi-device**: Support for multiple simultaneous sessions
- **Logout**: Proper session cleanup on logout

### 2. Data Security

#### Export Security
- **PDPA Compliance**: Automatic sensitive field masking
- **Access Control**: Export access restricted to authorized users
- **Audit Logging**: Logging of all export activities
- **File Permissions**: Secure file handling

#### Search Security
- **Query Sanitization**: Prevention of injection attacks
- **Result Filtering**: Filtering of sensitive information
- **Access Logs**: Logging of search activities
- **Rate Limiting**: Protection against abuse

### 3. Network Security

#### CORS Protection
```go
// CORS configuration
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "https://yourdomain.com")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

#### Security Headers
```go
// Security headers
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}
```

## 🛡️ Security Testing

### 1. Security Testing Framework

#### Automated Security Tests
```go
// Security test examples
func TestJWTSecurity(t *testing.T) {
    // Test token expiration
    // Test token manipulation
    // Test role-based access control
    // Test session management
}
```

#### Vulnerability Scanning
```bash
# Security scanning
cd backend
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...

# Dependency scanning
cd frontend
npm audit
```

### 2. Penetration Testing

#### Manual Testing
- **Authentication bypass**: Test for vulnerabilities
- **Authorization testing**: Verify access controls
- **Input validation**: Test for injection attacks
- **Session management**: Test session handling

#### Automated Testing
```go
func TestSecurityEndpoints(t *testing.T) {
    // Test for SQL injection
    // Test for XSS
    // Test for path traversal
    // Test for CSRF
}
```

## 🛡️ Security Compliance

### PDPA Compliance

#### Data Protection Measures
- **Field Masking**: Automatic masking of sensitive fields
- **Access Control**: Strict access controls on data
- **Audit Logging**: Comprehensive logging of data access
- **Data Retention**: Configurable retention policies

#### Sensitive Fields
```go
// Fields that are automatically masked
sensitiveFields := []string{
    "ip_address",
    "user_agent",
    "email",
    "phone",
    "credit_card",
    "ssn",
}
```

### GDPR Compliance

#### Data Handling
- **User Consent**: Clear consent mechanisms
- **Data Minimization**: Collect only necessary data
- **Right to Erasure**: Support for data deletion
- **Data Portability**: Export functionality

## 🛡️ Security Best Practices

### 1. Development Security

#### Code Security
- **Input Validation**: Always validate user inputs
- **Output Encoding**: Encode outputs to prevent XSS
- **Error Handling**: Secure error handling without exposing sensitive data
- **Dependency Management**: Regular dependency updates

#### Secure Coding Practices
```go
// Secure error handling
func secureHandler(c *gin.Context) {
    defer func() {
        if r := recover(); r != nil {
            // Log error securely
            log.Error("Internal server error", "error", r)
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
    }()
    
    // Handle request
}
```

### 2. Deployment Security

#### Environment Security
- **Secret Management**: Secure handling of secrets
- **Configuration**: Environment-specific configurations
- **Network**: Secure network configurations
- **Updates**: Regular security updates

#### Container Security
```yaml
# Secure container configuration
services:
  backend:
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
    user: 1000:1000
    cap_drop:
      - ALL
    cap_add:
      - CHOWN
      - SETGID
      - SETUID
```

## 🛡️ Incident Response

### Security Incident Handling

#### Incident Response Process
1. **Detection**: Monitor for security incidents
2. **Assessment**: Evaluate incident severity
3. **Containment**: Isolate affected systems
4. **Eradication**: Remove security threats
5. **Recovery**: Restore systems
6. **Lessons Learned**: Improve security measures

#### Logging and Monitoring
```go
// Security event logging
func logSecurityEvent(eventType string, details map[string]interface{}) {
    log := logger.GetLogger()
    log.Info("Security event", 
        "event_type", eventType,
        "details", details,
        "timestamp", time.Now())
}
```

## 🛡️ Security Updates

### Vulnerability Management

#### Regular Updates
- **Dependency Updates**: Regular updates of dependencies
- **Framework Updates**: Keep frameworks updated
- **Security Patches**: Apply security patches promptly
- **Version Control**: Maintain secure version control

#### Security Monitoring
- **Threat Intelligence**: Monitor for new threats
- **Vulnerability Scanning**: Regular vulnerability scans
- **Security Audits**: Periodic security audits
- **Compliance Checks**: Regular compliance verification

## 🛡️ Resources

### Security Tools and Libraries

#### Backend Security
- **Gosec**: Go security scanner
- **GolangCI-Lint**: Linting with security checks
- **JWT**: JSON Web Token library
- **GORM**: ORM with security features

#### Frontend Security
- **ESLint**: JavaScript linting
- **Snyk**: Dependency security scanning
- **CSP**: Content Security Policy
- **Helmet**: Security headers

### Security Standards
- **OWASP Top 10**: Web application security risks
- **NIST Cybersecurity Framework**: Security guidelines
- **ISO 27001**: Information security management
- **GDPR**: Data protection regulations

### Security Documentation
- [OWASP Security Guidelines](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://golang.org/doc/effective_go.html#security)
- [JWT Security](https://jwt.io/)
- [QuickWit Security](https://quickwit.io/docs/security/)