# APP - System Log Management Dashboard

APP is a self-hosted system log management dashboard for viewing, filtering, and analyzing server logs. It provides a comprehensive solution for monitoring and analyzing system logs with a modern web interface and powerful backend services.

## 🚀 Features

- **Modern Dashboard UI**: Built with Astro (TypeScript) for a responsive and intuitive user experience
- **JWT Authentication**: Secure token-based authentication with role-based access control
- **Log Search**: Powerful search capabilities with Lucene query syntax support
- **Log Export**: Export logs to CSV format for further analysis
- **Multi-Index Support**: Support for multiple log indices in QuickWit
- **Dark Theme**: Built-in dark theme for comfortable viewing
- **Responsive Design**: Works on desktop, tablet, and mobile devices

## 🛠️ Architecture
## 📦 Prerequisites

- Go 1.21+
- Node.js 18+
- Docker (for containerized deployment)
- QuickWit (for log search functionality)
- Elasticsearch (for log storage)

## 🚀 Getting Started

### Development Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd app
```

2. **Backend Setup**
```bash
cd backend
go mod tidy
go run main.go
```

3. **Frontend Setup**
```bash
cd frontend
npm install
npm run dev
```

### Environment Configuration

The application uses environment variables for configuration. Create a `.env` file in the backend directory:

```env
# Backend Port
PORT=8080

# QuickWit URL
QUICKWIT_URL=http://localhost:7280

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h
```

## 🔐 Authentication

The application uses JWT (JSON Web Token) authentication:

- **Login Endpoint**: `POST /api/auth/login`
- **Credentials**: `admin@example.com` / `admin123` (hardcoded for development)
- **Protected Routes**: All routes under `/api` require a valid JWT token
- **Role-based Access**: Admin role check via middleware

## 📡 API Endpoints

### Authentication
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/profile` - Get user profile
- `GET /api/admin` - Admin-only endpoint
- `POST /api/logout` - Logout

### Dashboard
- `GET /api/dashboard` - Dashboard data
- `GET /api/data-inventory` - Data inventory
- `GET /api/logs` - Logs data
- `GET /api/compliance` - Compliance data
- `GET /api/settings` - Settings
- `GET /api/help` - Help information

### Search & Export
- `GET /task/indices` - Get available indices
- `GET /task/search` - Search logs with query parameters
- `GET /task/export` - Export logs to CSV
- `GET /task/exports` - List exported files
- `GET /task/exports/:filename` - Download exported file

## 🐳 Deployment

### Docker

#### Build Images
```bash
cd docker/backend
docker build -t app-backend .
cd ../frontend
docker build -t app-frontend .
```

#### Run with Docker Compose
```bash
docker-compose up
```

### Kubernetes

#### Deploy with Helm
```bash
helm install app ./chart
```

#### Deploy with Kustomize
```bash
kubectl apply -k k8s/overlays/development
```

## 📁 Directory Structure
### Root
```
app
├── backend
│   ├── main.go             # Entry point using gin.Default()
│   ├── routers/            # Route setup (SetupRoutes function)
│   │   └── routers.go      # Route definitions
│   ├── controllers/        # HTTP handlers
│   │   ├── api_controller.go
│   │   ├── export_controller.go
│   │   ├── exports_list_controller.go
│   │   ├── protected_controller.go
│   │   └── search_controller.go
│   ├── services/           # Business logic (JWTService, etc.)
│   ├── middleware/         # Auth middleware (JWTAuth, RoleAuth)
│   ├── models/             # JWT claims structures
│   ├── config/             # Configuration loading
│   └── pages/              # Static HTML/JS/CSS (optional)
├── frontend
│   ├── src/
│   │   ├── components/     # Astro components
│   │   ├── layouts/        # Layout templates
│   │   └── pages/          # Route pages (dashboard, signin, signup)
│   ├── styles/             # CSS files
│   └── package.json
├── chart
│   ├── nfs-subdir-external-provisioner
│   ├── quickwit
│   └── vector
├── docker
│   ├── backend
│   └── frontend
├── docker-compose.yml
├── k8s
│   ├── base
│   ├── overlays
│   └── README.md
└── config.json
```
### Backend
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

### Frontend
```
frontend/
├── src/
│   ├── components/     # Astro components
│   ├── layouts/        # Layout templates
│   └── pages/          # Route pages
├── styles/             # CSS files
└── package.json
```

## 🧪 Testing and Quality

The application follows a modular architecture where:
- Controllers handle HTTP requests and responses
- Services contain business logic
- Middleware provides authentication and authorization
- Routers define the API endpoints
- Models define data structures

## 🛡️ Security Considerations

- JWT tokens are used for authentication
- Role-based access control is implemented
- CORS is configured for development domains only
- Export functionality includes PDPA field masking
- Path traversal protection in file downloads

## 📚 Documentation

- **Frontend**: Astro documentation available at https://docs.astro.build
- **Backend**: Gin framework documentation at https://gin-gonic.com
- **QuickWit**: QuickWit documentation at https://quickwit.io/docs

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with Go and Astro
- Integrated with QuickWit for log search
- Uses Elasticsearch for log storage


```
[ Developer ] --(Push Code)--> [ GitHub ]
                                    │
                            (Trigger Pipeline)
                                    ▼
                            [ GitHub Actions ]
                            ┌───────┴───────┐
                    (Build & Push)   (Update YAML Tag)
                            ▼               ▼
                  [Image Registry]   [ k3s-flux-infra Repo ]
                            ^               ^
                            │               │ (Auto Detect)
                            └───────┬───────┘
                                    ▼
                              [ FluxCD ] (K3s) -> Deploy!
```
