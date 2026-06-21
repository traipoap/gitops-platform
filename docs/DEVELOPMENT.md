# Development Guide

This document provides comprehensive instructions for setting up, developing, and maintaining the APP system log management dashboard.

## 🏗️ Project Structure

```
app/
├── backend/           # Go backend application
│   ├── main.go        # Entry point
│   ├── routers/       # Route definitions
│   ├── controllers/   # HTTP handlers
│   ├── services/      # Business logic
│   ├── middleware/    # Authentication and authorization
│   ├── models/        # Data structures
│   ├── config/        # Configuration loading
│   └── docs/          # API documentation
├── frontend/          # Astro frontend application
│   ├── src/           # Source code
│   │   ├── components/ # Reusable UI components
│   │   ├── layouts/    # Page layouts
│   │   └── pages/      # Route pages
│   ├── styles/        # CSS files
│   └── public/        # Static assets
├── docker/            # Docker configurations
│   ├── backend/       # Backend Dockerfile
│   ├── frontend/      # Frontend Dockerfile
│   └── docker-compose.yml # Docker Compose configuration
├── docs/              # Documentation files
├── scripts/           # Utility scripts
└── README.md          # Main documentation
```

## 🛠️ Prerequisites

### Required Software
- **Go 1.21+** - Programming language
- **Node.js 18+** - JavaScript runtime
- **Docker 20+** - Containerization platform
- **Docker Compose** - Multi-container Docker applications
- **Git** - Version control system

### Installation Commands
```bash
# Install Go (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# Install Node.js (Ubuntu/Debian)
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install Docker
sudo apt install docker.io
sudo apt install docker-compose

# Install Git
sudo apt install git
```

## 🚀 Setup Development Environment

### 1. Clone Repository
```bash
git clone <repository-url>
cd app
```

### 2. Backend Setup
```bash
cd backend
go mod tidy
```

### 3. Frontend Setup
```bash
cd frontend
npm install
```

### 4. Database Setup (QuickWit)
```bash
# QuickWit will be automatically set up via Docker Compose
# No manual setup required for development
```

## 🐳 Running the Application

### Using Docker Compose (Recommended)
```bash
# Start all services
docker-compose up

# Start services in detached mode
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs
```

### Manual Setup
#### Backend
```bash
cd backend
go run main.go
```

#### Frontend
```bash
cd frontend
npm run dev
```

## 📁 Development Workflow

### 1. Branching Strategy
```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Create fix branch
git checkout -b fix/your-fix-name

# Create hotfix branch
git checkout -b hotfix/your-hotfix-name
```

### 2. Code Development
```bash
# Make changes to code
# Test changes locally
# Commit changes
git add .
git commit -m "feat: add new feature"
git push origin feature/your-feature-name
```

### 3. Testing
#### Backend Tests
```bash
cd backend
go test ./...
go test -v ./...
go test -race ./...
```

#### Frontend Tests
```bash
cd frontend
npm test
npm run test:coverage
```

### 4. Code Quality
#### Backend
```bash
# Run linters
golangci-lint run

# Run security checks
gosec ./...
```

#### Frontend
```bash
# Run linters
npm run lint
npm run type-check

# Run security checks
npm audit
```

## 📊 Development Tools

### Backend Tools
- **GoLand/VS Code** - IDE for Go development
- **Delve** - Go debugger
- **GolangCI-Lint** - Linter for Go code
- **GoSec** - Security scanner for Go code
- **Testify** - Testing tools for Go

### Frontend Tools
- **VS Code** - IDE for frontend development
- **ESLint** - JavaScript linter
- **TypeScript** - Type checking
- **Tailwind CSS** - CSS framework
- **Jest** - Testing framework

## 🛡️ Development Security

### Security Practices
1. **Input Validation**: All inputs are validated and sanitized
2. **Authentication**: JWT-based authentication with secure tokens
3. **Authorization**: Role-based access control
4. **CORS**: Strict cross-origin resource sharing configuration
5. **Rate Limiting**: Protection against abuse
6. **Error Handling**: Secure error handling without exposing sensitive information

### Security Testing
```bash
# Security checks for backend
cd backend
gosec ./...
golangci-lint run

# Security checks for frontend
cd frontend
npm audit
```

## 📦 Building for Production

### Backend Build
```bash
cd backend
go build -o app-backend
```

### Frontend Build
```bash
cd frontend
npm run build
```

### Docker Build
```bash
# Build all services
docker-compose build

# Build specific service
docker-compose build backend
docker-compose build frontend
```

## 🧪 Testing Strategy

### Backend Testing
```go
// Example test structure
func TestSearchLogs(t *testing.T) {
    // Test setup
    // Execute
    // Verify results
    // Cleanup
}
```

### Frontend Testing
```typescript
// Example test structure
describe('LogSearch', () => {
    test('should render search input', () => {
        // Test implementation
    });
});
```

### Test Coverage
- **Unit Tests**: Individual function testing
- **Integration Tests**: Component interaction testing
- **End-to-End Tests**: Full user flow testing
- **Security Tests**: Vulnerability and compliance testing

## 📋 Code Standards

### Backend (Go)
```go
// Use clear, descriptive variable names
func ProcessLogEntry(entry *LogEntry) error {
    if entry == nil {
        return errors.New("entry cannot be nil")
    }
    
    // Process log entry
    return nil
}

// Add comments for exported functions
// ProcessLogEntry processes a log entry and stores it
func ProcessLogEntry(entry *LogEntry) error {
    // Implementation here
}
```

### Frontend (TypeScript)
```typescript
// Use TypeScript for type safety
interface LogEntry {
    id: string;
    timestamp: string;
    level: string;
    message: string;
    source: string;
}

// Use descriptive function names
const getLogEntries = async (query: string): Promise<LogEntry[]> => {
    // Implementation here
};
```

## 📊 Monitoring and Debugging

### Logging
```go
// Structured logging example
log := logger.GetLogger()
log.Info("User logged in", 
    "user_id", user.ID, 
    "ip", ctx.ClientIP(),
    "user_agent", ctx.GetHeader("User-Agent"))
```

### Debugging
```bash
# Backend debugging
dlv debug main.go

# Frontend debugging
# Use browser developer tools
# Set breakpoints in VS Code
```

## 🔄 Version Control

### Git Workflow
```bash
# Fetch latest changes
git fetch origin

# Rebase on main
git checkout main
git pull origin main
git checkout feature/your-feature
git rebase main

# Create pull request
git push origin feature/your-feature
```

### Commit Message Format
```bash
# Example commit messages
feat(auth): add JWT authentication
fix(search): resolve issue with query parsing
docs(readme): update installation instructions
style: format code according to style guide
refactor: restructure authentication module
test: add tests for export functionality
```

## 📚 Documentation

### Documentation Updates
```bash
# Update documentation files
# Follow existing style and format
# Ensure consistency across documents
```

### API Documentation
- Auto-generated API documentation
- Examples for all endpoints
- Clear parameter descriptions
- Error response examples

## 🛠️ Troubleshooting

### Common Issues

#### Docker Issues
```bash
# Clear Docker cache
docker system prune -a

# Rebuild images
docker-compose build --no-cache

# Check Docker logs
docker-compose logs
```

#### Backend Issues
```bash
# Check dependencies
cd backend
go mod tidy

# Run with verbose logging
go run main.go -v
```

#### Frontend Issues
```bash
# Clear npm cache
npm cache clean --force

# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install

# Check for TypeScript errors
npm run type-check
```

## 📋 Performance Optimization

### Backend Optimization
- Efficient database queries
- Caching strategies
- Memory management
- Concurrency handling

### Frontend Optimization
- Code splitting
- Lazy loading
- Asset optimization
- Bundle size reduction

## 📋 Environment Configuration

### Development Environment
```bash
# Environment variables for development
export APP_ENV=development
export DATABASE_URL=postgresql://localhost:5432/app_dev
export JWT_SECRET=your-secret-key
```

### Production Environment
```bash
# Environment variables for production
export APP_ENV=production
export DATABASE_URL=postgresql://prod-db:5432/app_prod
export JWT_SECRET=your-production-secret-key
```

## 📚 Resources

### Development Resources
- [Go Documentation](https://golang.org/doc/)
- [Astro Documentation](https://docs.astro.build/)
- [Gin Framework](https://gin-gonic.com/)
- [QuickWit Documentation](https://quickwit.io/docs/)
- [Docker Documentation](https://docs.docker.com/)

### Development Tools
- **IDE**: VS Code with Go and TypeScript extensions
- **Linter**: GolangCI-Lint, ESLint
- **Security**: GoSec, npm audit
- **Testing**: Go test, Jest, Cypress

### Best Practices
- Follow semantic versioning
- Maintain clean commit history
- Write comprehensive tests
- Keep documentation updated
- Follow security best practices