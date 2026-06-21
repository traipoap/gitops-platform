# Contributing Guide

This document provides guidelines and instructions for contributing to the APP system log management dashboard project.

## 🤝 Welcome!

We welcome contributions from the community to help improve and extend the APP system log management dashboard. Whether you're reporting bugs, suggesting features, or submitting code changes, your contributions are valuable to the project.

## 📋 Before You Contribute

### 1. Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming and inclusive environment for all contributors.

### 2. Project Overview

The APP system log management dashboard is built with:
- **Backend**: Go 1.21+ with Gin framework
- **Frontend**: Astro 4.x with TypeScript and Tailwind CSS
- **Data Layer**: QuickWit search engine
- **Deployment**: Docker and Docker Compose

### 3. Development Environment

Make sure you have the following prerequisites installed:
- Go 1.21+
- Node.js 18+
- Docker 20+
- Docker Compose
- Git

## 🚀 Getting Started

### 1. Fork and Clone

```bash
# Fork the repository on GitHub
git clone <your-forked-repo-url>
cd app
```

### 2. Setup Development Environment

#### Backend Setup
```bash
cd backend
go mod tidy
```

#### Frontend Setup
```bash
cd frontend
npm install
```

### 3. Run Development Server

#### Option 1: Using Docker Compose
```bash
docker-compose up
```

#### Option 2: Manual Setup
```bash
# Backend
cd backend
go run main.go

# Frontend
cd frontend
npm run dev
```

## 📦 Contribution Process

### 1. Create a Branch

```bash
# Create a new branch for your feature or fix
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-fix-name
```

### 2. Make Changes

Follow the project's coding standards and conventions:
- Write clean, readable code
- Add comments for complex logic
- Include documentation for new features
- Write tests for new functionality

### 3. Test Your Changes

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

### 4. Commit Your Changes

Follow the commit message conventions:
```bash
# Example commit messages
feat(auth): add JWT authentication
fix(search): resolve issue with query parsing
docs(readme): update installation instructions
style: format code according to style guide
refactor: restructure authentication module
test: add tests for export functionality
```

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:
- Clear title and description
- Reference to related issues
- Testing instructions
- Screenshots (if applicable)

## 📋 Code Standards

### Backend (Go)

#### Code Style
```go
// Use clear, descriptive variable names
func ProcessLogEntry(entry *LogEntry) error {
    if entry == nil {
        return errors.New("entry cannot be nil")
    }
    
    // Process log entry
    return nil
}
```

#### Documentation
```go
// Add comments for exported functions
// ProcessLogEntry processes a log entry and stores it
func ProcessLogEntry(entry *LogEntry) error {
    // Implementation here
}
```

### Frontend (TypeScript/Astro)

#### Code Style
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

#### Documentation
```typescript
/**
 * Search logs with the given query
 * @param query - Search query string
 * @returns Promise resolving to array of log entries
 */
const searchLogs = async (query: string): Promise<LogEntry[]> => {
    // Implementation here
};
```

## 🧪 Testing

### Backend Testing
```go
// Example test
func TestProcessLogEntry(t *testing.T) {
    entry := &LogEntry{
        ID: "test-123",
        Message: "test message",
    }
    
    err := ProcessLogEntry(entry)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

### Frontend Testing
```typescript
// Example test
describe('LogSearch', () => {
    test('should render search input', () => {
        const { getByPlaceholderText } = render(<LogSearch />);
        expect(getByPlaceholderText('Search logs...')).toBeInTheDocument();
    });
});
```

## 📋 Pull Request Guidelines

### Before Submitting
- Ensure all tests pass
- Check code formatting
- Update documentation if needed
- Add comments for complex changes
- Follow the project's coding standards

### What to Include
- Clear description of changes
- Reference to related issues
- Testing instructions
- Screenshots (for UI changes)
- Performance impact (if any)

## 📚 Documentation

### Update Documentation
When adding new features or changing existing functionality, update the relevant documentation:
- README.md
- DEVELOPMENT.md
- ARCHITECTURE.md
- API documentation

### Documentation Style
- Use consistent formatting
- Include examples where appropriate
- Keep documentation up-to-date with code changes
- Use clear, concise language

## 🛡️ Security Considerations

### Security Best Practices
- Validate all inputs
- Sanitize outputs
- Use secure authentication
- Follow security guidelines
- Regular security testing

### Security Testing
```bash
# Run security checks
cd backend
gosec ./...

# Check dependencies
cd frontend
npm audit
```

## 📊 Code Review Process

### Review Checklist
- [ ] Code follows project standards
- [ ] Tests pass
- [ ] Documentation updated
- [ ] Security considerations addressed
- [ ] Performance impact considered

### Review Comments
- Be constructive and helpful
- Focus on code quality and maintainability
- Suggest improvements where needed
- Acknowledge good practices

## 📋 Issue Reporting

### Bug Reports
When reporting bugs, include:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment information
- Screenshots (if applicable)

### Feature Requests
When requesting features, include:
- Clear description of the feature
- Use cases for the feature
- Implementation suggestions (if any)
- Priority level

## 📋 Community Guidelines

### Communication
- Be respectful and professional
- Follow the code of conduct
- Provide helpful feedback
- Ask questions when unsure

### Recognition
Contributors will be recognized in:
- CHANGELOG.md
- GitHub contributors list
- Project documentation
- Community announcements

## 📚 Resources

### Learning Resources
- [Go Documentation](https://golang.org/doc/)
- [Astro Documentation](https://docs.astro.build/)
- [Gin Framework](https://gin-gonic.com/)
- [QuickWit Documentation](https://quickwit.io/docs/)
- [Docker Documentation](https://docs.docker.com/)

### Tools and Libraries
- Go tools: `golangci-lint`, `gosec`, `air`
- Frontend tools: `Vite`, `TypeScript`, `Tailwind CSS`
- Testing: `Go test`, `Jest`, `Cypress`
- Security: `npm audit`, `gosec`

## 📋 Frequently Asked Questions

### How do I set up the development environment?
Follow the instructions in [DEVELOPMENT.md](DEVELOPMENT.md).

### What are the coding standards?
Follow the existing code style and conventions in the project.

### How do I run tests?
Run backend tests with `go test ./...` and frontend tests with `npm test`.

### How do I submit a pull request?
Fork the repository, create a branch, make changes, and submit a pull request.

### What if I have questions?
Join our community discussions or open an issue for help.

## 📋 License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project. See [LICENSE](LICENSE) for details.

## 🙏 Thank You!

Thank you for contributing to the APP system log management dashboard. Your efforts help make this project better for everyone!