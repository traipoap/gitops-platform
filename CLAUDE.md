# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Backend (Go/Gin)
- **Run server**: `cd backend && go run main.go`
- **Build binary**: `cd backend && go build -o server main.go`
- **Tidy dependencies**: `cd backend && go mod tidy`
- **Run tests**: `cd backend && go test ./...`

### Frontend (Astro/TypeScript)
- **Install dependencies**: `cd frontend && npm install`
- **Run development server**: `cd frontend && npm run dev`
	- Port: 4321
- **Build production assets**: `cd frontend && npm run build`

### Local Orchestration
- **Start both services (Go + Astro)**: `./start-app.sh`
- **Run with Docker Compose**: `docker compose up --build`

### Infrastructure & Deployment
- **Helm deployment**: `helm install app ./chart`
- **Kustomize deployment (Dev environment)**: `kubectl apply -k k8s/overlays/development`
- **Build Docker images**: 
    - Backend: `docker build -t app-backend ./docker/backend`
    - Frontend: `docker build -t app-frontend ./docker/frontend`

## Architecture Overview

This repository contains a full-stack log management system composed of several layers:

### 1. Core Application Layers
- **Backend (Go/Gin)**: A RESTful API that handles business logic, authentication (JWT), and interface with search engines.
    - `controllers/`: Request handling and response orchestration.
     	- `services/`: Core business logic and integration (e.g., Quickwit client).
    	- `middleware/`: Security and utility layers (JWT Auth, Role checks).
    	- `models/`: Data structures and JWT claims.
    	- `routers/`: API endpoint definitions.
- **Frontend (Astro)**: A responsive dashboard for interacting with the backend.
    - Uses Astro components, layouts, and pages to build a modern UI for searching and exporting logs.

### 2. Search & Data Engine
- **Quickwit**: The primary search engine used for log indexing and searching via Lucene syntax.
- **Vector/NFS**: Infrastructure components for log ingestion and persistent storage.

### 3. Deployment & Infrastructure
- **Containerization**: Docker and Docker Compose for local development and standardizing runtime environments.
- **Kubernetes (K8s)**: Managed deployments using Helm charts and Kustomize overlays.
- **Infrastructure as Code (IaC)**:
    - **Terraform**: Provisioning of Proxmox-based infrastructure.
    - **Ansible**: Configuration management for server setup and role deployment.

## Key Concepts & Security
- **Authentication**: Uses JWT (JSON Web Tokens) with middleware enforcing authorization at the router level.
- **Search Syntax**: Supports Lucene query syntax through integration with Quickwit.
- **Security Features**: Includes path traversal protection, PDPA field masking in exports, and CORS management.
