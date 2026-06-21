# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

APP is a self-hosted system log management dashboard for viewing, filtering, and analyzing server logs. It's a full-stack application with:
- **Backend**: Go (Gin framework) serving as the API gateway
- **Frontend**: Astro (TypeScript) for the dashboard UI
- **External**: QuickWit for log search and Elasticsearch

## Commands

### Frontend
- Start dev server: `npm run dev` (from frontend/ directory)
- Build production: `npm run build`
- Preview build: `npm run preview`

### Backend
- Start dev server: `go run main.go` (from backend/ directory)
- Default port: 8080
- Build: `go build -o exporter main.go`

### Common Workflows
- Run both servers: `npm run dev` in frontend/, `go run main.go` in backend/
- Check backend API: http://localhost:8080/api
- Test authentication: http://localhost:8080/api/auth/login (basic auth with admin@example.com/admin123)

## Architecture

```
app
├── backend
│   ├── backend.md
│   ├── components
│   ├── config
│   ├── controllers
│   ├── exports
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── middleware
│   ├── models
│   ├── pages
│   ├── routers
│   ├── services
│   └── static
├── chart
│   ├── nfs-subdir-external-provisioner
│   ├── quickwit
│   └── vector
├── CLAUDE.md
├── config.json
├── docker
│   ├── backend
│   └── frontend
├── docker-compose.yml
├── frontend
│   ├── astro.config.mjs
│   ├── frontend.md
│   ├── package.json
│   ├── package-lock.json
│   ├── public
│   ├── README.md
│   ├── src
│   └── tsconfig.json
└── k8s
    ├── base
    ├── overlays
    └── README.md
```

### Frontend (Astro)
```
frontend/
├── src/
│   ├── components/     # Astro components
│   ├── layouts/        # Layout templates
│   └── pages/          # Route pages (dashboard, signin, signup)
├── styles/             # CSS files
├── assets/             # Static assets
└── package.json
```

### Backend (Go/Gin)
```
backend/
├── main.go             # Entry point using gin.Default()
├── routers/            # Route setup (SetupRoutes function)
│   └── routers.go      # Route definitions
├── controllers/        # HTTP handlers
│   ├── api_controller.go
│   ├── export_controller.go
│   ├── exports_list_controller.go
│   ├── protected_controller.go
│   └── search_controller.go
├── services/           # Business logic (JWTService, etc.)
├── middleware/         # Auth middleware (JWTAuth, RoleAuth)
├── models/             # JWT claims structures
├── config/             # Configuration loading
└── pages/              # Static HTML/JS/CSS (optional)
```

### External Integration
- Backend proxies requests to QuickWit via `AppConfig.QuickwitURL`
- CORS configured for frontend domains: localhost:3000, localhost:4321, 192.168.1.107:3000

## Key Patterns

1. **JWT Authentication**: Token-based auth using golang-jwt/v5
2. **Role-based Access**: Admin role checks via middleware
3. **Export Functionality**: Large log exports with file generation
4. **Dark Theme UI**: CSS variables in main.css
5. **Progressive Enhancement**: Vanilla JS, no heavy framework requirements

## File Locations

- Environment variables: backend/.env
- Build output: frontend/dist/, backend/exporter/
- Static assets: frontend/src/assets/, backend/static/

## Development Notes

- Frontend uses `type: module` for ES modules
- Backend uses Go modules (go.mod/go.sum)
- Frontend is Astro v6 with TypeScript support
- Backend CORS allows localhost:4321, 192.168.1.107:3000 for development
- Authentication credentials are hardcoded for development (admin@example.com/admin123)
- The application has both API endpoints and static HTML pages for search functionality
- The backend uses a modular structure with controllers, services, middleware, and routers
- The frontend uses Astro's component-based architecture with island architecture patterns

## Testing and Quality

The application follows a modular architecture where:
- Controllers handle HTTP requests and responses
- Services contain business logic
- Middleware provides authentication and authorization
- Routers define the API endpoints
- Models define data structures

## Important Endpoints

- `/api/auth/login` - Login endpoint (basic auth with admin@example.com/admin123)
- `/api/dashboard` - Dashboard data endpoint
- `/task/search` - Search logs endpoint
- `/task/export` - Export logs endpoint
- `/task/exports` - List exports endpoint
- `/task/exports/:filename` - Download export endpoint