# APP System Log Management Dashboard

[![License](https://img.shields.io/github/license/your-username/app)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/your-username/app)](go.mod)
[![Docker](https://img.shields.io/docker/v/your-username/app-backend?label=docker%20backend)](https://hub.docker.com/repository/docker/your-username/app-backend)
[![Docker](https://img.shields.io/docker/v/your-username/app-frontend?label=docker%20frontend)](https://hub.docker.com/repository/docker/your-username/app-frontend)

## 📋 Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Getting Started](#getting-started)
5. [Installation](#installation)
6. [Usage](#usage)
7. [Security](#security)
8. [Development](#development)
9. [API Documentation](#api-documentation)
10. [Contributing](#contributing)
11. [License](#license)
12. [Resources](#resources)

## 📖 Overview

The APP system log management dashboard is a comprehensive solution for managing and analyzing system logs. Built with modern technologies, it provides real-time log monitoring, powerful search capabilities, and secure export functionality with PDPA compliance.

## 🚀 Features

### 🔍 Log Management
- Real-time log monitoring and visualization
- Advanced search with filtering capabilities
- Log entry details and context
- Multi-source log aggregation

### 📊 Export Functionality
- Export logs in JSON, CSV, and TXT formats
- PDPA-compliant field masking
- Secure export processing
- Export status tracking

### 🔐 Security
- JWT-based authentication
- Role-based access control
- CORS configuration
- Input validation and sanitization
- Rate limiting
- Secure token handling

### 🛠️ System Management
- Health check endpoints
- System statistics and monitoring
- User management and permissions
- Configurable data retention

### 🌐 Technology Stack
- **Backend**: Go 1.21+ with Gin framework
- **Frontend**: Astro 4.x with TypeScript and Tailwind CSS
- **Data Layer**: QuickWit search engine
- **Deployment**: Docker and Docker Compose

## 🏗️ Architecture

### System Components
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   QuickWit      │
│   (Astro)       │    │   (Go)          │    │   (Search)      │
│                 │    │                 │    │                 │
│  - UI           │    │  - REST API     │    │  - Log Storage  │
│  - Search       │    │  - Authentication│   │  - Search       │
│  - Export       │    │  - Business     │    │  - Indexing     │
│  - Dashboard    │    │  - Security     │    │                 │
└─────────────────┘    │  - Middleware   │    └─────────────────┘
                         │  - Controllers  │
                         │  - Services     │
                         └─────────────────┘
```

### Data Flow
1. **User Request**: Frontend sends request to backend API
2. **Authentication**: JWT validation and authorization
3. **Processing**: Backend processes request and interacts with QuickWit
4. **Response**: Data returned to frontend for display
5. **Export**: Secure export processing with field masking

### Security Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   QuickWit      │
│   (Secure)      │    │   (Secure)      │    │   (Secure)      │
│                 │    │                 │    │                 │
│  - HTTPS        │    │  - JWT Auth     │    │  - Encrypted    │
│  - CORS         │    │  - Input Sanit  │    │  - Access       │
│  - CSP          │    │  - Rate Limit   │    │  - Search       │
│  - Security     │    │  - Error Log    │    │  - Audit        │
│  - Tokens       │    │  - SSL/TLS      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker 20+
- Docker Compose

### Quick Start
```bash
# Clone the repository
git clone <repository-url>
cd app

# Start the application
docker-compose up

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

## 📦 Installation

### Docker Installation
```bash
# Clone repository
git clone <repository-url>
cd app

# Build and start services
docker-compose up --build

# Stop services
docker-compose down
```

### Manual Installation
#### Backend Setup
```bash
cd backend
go mod tidy
go run main.go
```

#### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## 📋 Usage

### Authentication
1. Access the application at `http://localhost:3000`
2. Login with credentials (default admin user)
3. Use the JWT token for API access

### Searching Logs
1. Enter search query in the search bar
2. Apply filters for date range, log level, or source
3. View results in the log table
4. Click on any log entry for detailed view

### Exporting Logs
1. Perform a search to get log results
2. Click the export button
3. Select export format (JSON, CSV, TXT)
4. Download the exported file (sensitive fields will be masked)

### System Monitoring
1. Navigate to the system dashboard
2. View real-time statistics
3. Check system health status
4. Monitor resource usage

## 🛡️ Security

### Authentication
The application implements JWT-based authentication with:
- Secure token generation and validation
- Token refresh mechanism
- Secure token storage
- Token expiration handling

### Data Protection
- PDPA-compliant field masking in exports
- HTTPS encryption for data in transit
- Input validation and sanitization
- Access logging and audit trails

### Access Control
- Role-based access control (RBAC)
- Different permissions for Admin, User, and Guest roles
- Route protection with middleware
- Fine-grained resource access control

### Network Security
- CORS configuration with strict controls
- HTTPS enforcement
- Security headers implementation
- Redirects from HTTP to HTTPS

## 🛠️ Development

### Development Environment
```bash
# Clone repository
git clone <repository-url>
cd app

# Setup development environment
docker-compose up -d

# Access development services
# Backend: http://localhost:8080
# Frontend: http://localhost:3000
```

### Testing
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### Documentation
Comprehensive documentation is available in the `/docs` directory:
- [DEVELOPMENT.md](DEVELOPMENT.md) - Development guide
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [SECURITY.md](SECURITY.md) - Security practices
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [API.md](API.md) - API documentation

## 📚 API Documentation

Full API documentation is available at [API.md](API.md).

### Key Endpoints
- `POST /api/auth/login` - Authenticate user
- `GET /api/logs/search` - Search log entries
- `POST /api/logs/export` - Export log entries
- `GET /api/health` - Health check
- `GET /api/system/stats` - System statistics

## 🤝 Contributing

We welcome contributions from the community! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📚 Resources

### Documentation
- [DEVELOPMENT.md](DEVELOPMENT.md)
- [ARCHITECTURE.md](ARCHITECTURE.md)
- [SECURITY.md](SECURITY.md)
- [CONTRIBUTING.md](CONTRIBUTING.md)
- [API.md](API.md)

### Technology Stack
- **Go**: Backend development
- **Astro**: Frontend framework
- **QuickWit**: Search engine
- **Docker**: Containerization
- **Gin**: Web framework

### Security Resources
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://golang.org/doc/effective_go.html#security)
- [JWT Security](https://jwt.io/)
- [QuickWit Security](https://quickwit.io/docs/security/)

### Development Tools
- Go tools: `golangci-lint`, `gosec`, `air`
- Frontend tools: `Vite`, `TypeScript`, `Tailwind CSS`
- Testing: `Go test`, `Jest`, `Cypress`
- Security: `npm audit`, `gosec`

### Community
- GitHub Issues
- Pull Requests
- Discussion Forums