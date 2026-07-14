# Architecture Documentation

This document outlines the architectural design and components of the APP system log management dashboard.

## 🏗️ Overview

The APP system log management dashboard is a modern, secure, and scalable solution for viewing, filtering, and analyzing system logs. It provides a comprehensive log management experience with authentication, search capabilities, and export functionality, built using a Go backend with Gin framework and an Astro frontend.

## 📊 System Architecture

### High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Data Layer    │
│   (Astro)       │    │   (Go + Gin)    │    │   (QuickWit)    │
│                 │    │                 │    │                 │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │   UI      │  │    │  │   API     │  │    │  │  Logs     │  │
│  │ Components│  │    │  │  Handlers │  │    │  │  Storage  │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
│                 │    │                 │    │                 │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │   Pages   │  │    │  │  Routing  │  │    │  │  Search   │  │
│  │  Layouts  │  │    │  │  Middleware  │    │  │  Indexing │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        └───────────────────────┼───────────────────────┘
                                │
                        ┌─────────────────┐
                        │   Authentication│
                        │   (JWT)         │
                        └─────────────────┘
```

### Component Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              APP Dashboard                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  Frontend   │  │  Backend    │  │  Security   │  │  Data       │         │
│  │  (Astro)    │  │  (Go)       │  │  (JWT)      │  │  Layer      │         │
│  │             │  │             │  │             │  │  (QuickWit) │         │
│  │  UI         │  │  API        │  │  Auth       │  │  Search     │         │
│  │  Components │  │  Controllers│  │  Middleware │  │  Indexing   │         │
│  │  Pages      │  │  Services   │  │  RBAC       │  │  Storage    │         │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 📦 Components

### 1. Frontend (Astro)

#### Technology Stack
- **Framework**: Astro 4.x
- **Language**: TypeScript/JavaScript
- **Styling**: Tailwind CSS
- **Build Tool**: Vite

#### Key Features
- Responsive UI design
- Component-based architecture
- Server-side rendering
- Client-side interactivity
- TypeScript type safety

#### Component Structure
```
frontend/
├── src/
│   ├── components/     # Reusable UI components
│   ├── layouts/        # Page layouts
│   └── pages/          # Route pages
├── styles/             # CSS files
└── public/             # Static assets
```

#### Key Components
- **Dashboard Layout**: Main application layout
- **Search Component**: Log search interface
- **Log Viewer**: Log display and filtering
- **Export Component**: Log export functionality
- **Authentication**: Login and session management

### 2. Backend (Go + Gin)

#### Technology Stack
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: QuickWit (search engine)
- **Authentication**: JWT
- **Security**: CORS, Rate Limiting

#### Key Features
- RESTful API endpoints
- JWT-based authentication
- Search and filtering capabilities
- Export functionality
- Security middleware

#### Component Structure
```
backend/
├── main.go             # Entry point
├── routers/            # Route definitions
├── controllers/        # HTTP handlers
├── services/           # Business logic
├── middleware/         # Authentication and authorization
├── models/             # Data structures
├── config/             # Configuration loading
└── pages/              # Static HTML pages
```

#### Key Services
- **Authentication Service**: User authentication and JWT management
- **Search Service**: Log search and filtering
- **Export Service**: Log export functionality
- **Health Service**: System health monitoring

### 3. Data Layer (QuickWit)

#### Technology Stack
- **Search Engine**: QuickWit (Elasticsearch-compatible)
- **Storage**: File-based storage
- **Indexing**: Real-time indexing

#### Key Features
- Fast log search capabilities
- Real-time indexing
- Scalable storage
- Full-text search
- Aggregation support

## 🔄 Data Flow

### 1. User Authentication Flow
```
1. User accesses login page
2. User submits credentials
3. Backend validates credentials
4. JWT token generated and returned
5. Token stored in browser (localStorage/sessionStorage)
6. Subsequent requests include JWT in Authorization header
7. Middleware validates token on each request
```

### 2. Log Search Flow
```
1. User submits search query
2. Frontend sends request to backend
3. Backend validates authentication
4. Backend forwards request to QuickWit
5. QuickWit performs search
6. Results returned to backend
7. Backend processes results
8. Results returned to frontend
9. Frontend displays results
```

### 3. Log Export Flow
```
1. User initiates export
2. Frontend sends export request to backend
3. Backend validates authentication
4. Backend queries QuickWit for logs
5. Backend processes logs with PDPA field masking
6. Export file created
7. File served to frontend
8. User downloads export file
```

## 🛡️ Security Architecture

### Authentication and Authorization
- **JWT-based Authentication**: Secure token-based authentication
- **Role-Based Access Control (RBAC)**: Different access levels for users
- **Middleware Protection**: All protected routes secured with middleware
- **Secure Token Storage**: Tokens stored securely in browser

### Data Protection
- **PDPA Field Masking**: Sensitive fields masked in exports
- **Data Encryption**: HTTPS for data in transit
- **Access Control**: Strict access control on all endpoints
- **Audit Trail**: Logging of access and modification activities

### Network Security
- **CORS Configuration**: Strict cross-origin resource sharing
- **HTTPS Support**: All production deployments use HTTPS
- **Firewall Rules**: Network segmentation and access control
- **Session Management**: Secure session handling

## 📈 Performance Architecture

### Scalability
- **Horizontal Scaling**: Support for multiple backend instances
- **Load Balancing**: Distribute traffic across instances
- **Caching**: Implement caching strategies where appropriate
- **Database Optimization**: Efficient indexing and query optimization

### Monitoring
- **Health Checks**: API endpoints for system monitoring
- **Metrics Collection**: Performance metrics collection
- **Logging**: Structured logging for debugging and monitoring
- **Alerting**: Automated alerting for system issues

## 📦 Deployment Architecture

### Containerization
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   QuickWit      │
│   (Docker)      │    │   (Docker)      │    │   (Docker)      │
│                 │    │                 │    │                 │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │   Astro   │  │    │  │   Go      │  │    │  │  QuickWit │  │
│  │  Build    │  │    │  │  Server   │  │    │  │  Engine   │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        └───────────────────────┼───────────────────────┘
                                │
                    ┌─────────────────┐
                    │   Docker Compose│
                    │   Orchestration │
                    └─────────────────┘
```

### Environment Configuration
- **Development**: Local development with hot-reloading
- **Staging**: Pre-production environment
- **Production**: Live production environment

## 📋 API Architecture

### RESTful API Endpoints

#### Authentication
```
POST /api/v1/auth/login          # User login
POST /api/v1/auth/refresh        # Token refresh
```

#### Dashboard
```
GET /api/v1/dashboard            # Get dashboard data
GET /api/v1/dashboard/health     # Health check
```

#### Search
```
GET /api/v1/task/search          # Search logs
GET /api/v1/task/autocomplete    # Search autocomplete
```

#### Export
```
GET /api/v1/task/export          # Export logs
GET /api/v1/task/exports         # Get export list
```

### Request/Response Structure
```json
{
  "status": "success",
  "data": {},
  "message": "Operation completed successfully"
}
```

## 📊 Monitoring and Logging

### Logging Structure
```go
// Structured logging example
log := logger.GetLogger()
log.Info("User logged in", 
    "user_id", user.ID, 
    "ip", ctx.ClientIP(),
    "user_agent", ctx.GetHeader("User-Agent"))
```

### Metrics Collection
- Request count and timing
- Error rates
- Response codes
- System resource usage

## 📋 Technology Stack Summary

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: QuickWit (search engine)
- **Authentication**: JWT
- **Security**: CORS, Rate Limiting

### Frontend
- **Framework**: Astro 4.x
- **Language**: TypeScript/JavaScript
- **Styling**: Tailwind CSS
- **Build Tool**: Vite

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Monitoring**: Prometheus/Grafana
- **Security**: HTTPS, JWT, CORS

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Astro Documentation](https://docs.astro.build/)
- [Gin Framework](https://gin-gonic.com/)
- [QuickWit Documentation](https://quickwit.io/docs/)
- [Docker Documentation](https://docs.docker.com/)
