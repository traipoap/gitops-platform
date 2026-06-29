# AGENTS.md

## Development Workflows

### Backend (Go/Gin)
- **Run server**: `cd backend && go run main.go`
- **Build**: `cd backend && go build -o exporter main.go`
- **Environment**: Uses `.env` in `backend/` directory.

### Frontend (Astro)
- **Run dev server**: `cd frontend && npm run dev`
- **Build**: `cd frontend && npm run build`
- **Preview**: `cd frontend && npm run preview`

## Key Project Context

- **Architecture**: Full-stack system. Backend (Go) acts as an API gateway proxying to QuickWit and Elasticsearch.
- **Authentication**: JWT-based. Development credentials: `admin@example.com` / `admin1mu123`.
- **CORS**: Backend allows `localhost:3000`, `localhost:4321`, and `192.168.1.107:3000`.
- **External Services**: Depends on QuickWit for log search and Elasticsearch for storage.
- **Endpoints**:
  - `/api/auth/login` (Auth)
  - `/api/dashboard` (Dashboard data)
  - `/task/search` (Log search)
  - `/task/export` (Log export)
  - `/task/exports` (List exports)
