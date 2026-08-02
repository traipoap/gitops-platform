# Log Management System

A full-stack log management and search dashboard built with Go, Astro, and Quickwit.

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Frontend   │ ─── │    Backend   │ ─── │   Quickwit  │
│  (Astro/TS) │     │  (Go/Gin)    │     │  (Search)   │
└─────────────┘     └──────────────┘     └─────────────┘
```

- **Backend** (`backend/`) - Go/Gin RESTful API with JWT authentication, role-based access control, PD