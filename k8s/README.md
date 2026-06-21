# Kustomize Configuration for Kubernetes

This directory contains Kustomize configurations for deploying the log-management application to Kubernetes clusters.

## Structure

```
k8s/
├── base/                          # Base configuration for all environments
│   ├── kustomization.yaml         # Root kustomization
│   ├── namespace.yaml             # Namespace definition
│   ├── backend-deployment.yaml    # Backend deployment
│   ├── backend-service.yaml       # Backend service
│   ├── frontend-deployment.yaml   # Frontend deployment
│   ├── frontend-service.yaml      # Frontend service
│   ├── configmap.yaml             # Application config
│   └── ingress.yaml               # Ingress configuration
├── overlays/                      # Environment-specific overlays
│   ├── development/               # Local development
│   ├── staging/                   # Staging environment
│   └── production/               # Production environment
└── README.md                      # This file
```

## Quick Start

### Development

Apply the development configuration:

```bash
kubectl apply -k k8s/overlays/development
```

### Staging

Apply the staging configuration:

```bash
kubectl apply -k k8s/overlays/staging
```

### Production

Apply the production configuration:

```bash
kubectl apply -k k8s/overlays/production
```

## Environment Differences

| Setting | Development | Staging | Production |
|---------|-------------|---------|------------|
| Replicas | 1 | 2 | 2+ |
| NodeEnv | development | staging | production |
| Log Level | debug | info | info |
| Image Tag | dev-0.1.0 | staging-0.1.0 | v1.2.3 |
| Secrets | None | ConfigMap | Secrets |

## Available Services

- **frontend**: React/Vite frontend (port 4321)
- **backend**: Go API server (port 8080)
- **quickwit**: Log ingestion/search (port 8000)

## Resources

View deployed resources:

```bash
# All resources
kubectl get all -n log-management

# Specific service
kubectl get svc -n log-management
```

## Cleanup

Remove all resources:

```bash
kubectl delete -k k8s/overlays/development
```
