# K3s GitOps Platform on Proxmox

[![Frontend Build & Deploy](https://github.com/traipoap/gitops-platform/actions/workflows/frontend-build.yml/badge.svg)](https://github.com/traipoap/gitops-platform/actions/workflows/frontend-build.yml)
[![Backend Build & Deploy](https://github.com/traipoap/gitops-platform/actions/workflows/backend-build.yml/badge.svg)](https://github.com/traipoap/gitops-platform/actions/workflows/backend-build.yml)
![Kubernetes](https://img.shields.io/badge/Kubernetes-K3s-326ce5)
![GitOps](https://img.shields.io/badge/GitOps-FluxCD-5468ff)
![IaC](https://img.shields.io/badge/IaC-Terraform%20%7C%20Ansible-7b42bc)
![License](https://img.shields.io/github/license/traipoap/gitops-platform)

A production-oriented home lab / portfolio platform for provisioning and managing a Kubernetes cluster using Infrastructure as Code, GitOps, CI/CD, observability, and security best practices.

This project demonstrates how to build a repeatable, declarative, and automated Kubernetes platform from bare-metal/virtual infrastructure to application deployment.

---

## Table of Contents
- [Overview](#-overview)
- [Lab Environment](#lab-environment)
- [Architecture](#architecture)
- [Network Topology Diagram](#network-topology-diagram)
- [Deployment Time Estimates](#deployment-time-estimates)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Tech Stack](#tech-stack)
- [Repository Structure](#repository-structure)
- [Quickstart](#quickstart)
- [GitOps Workflow](#gitops-workflow)
- [CI/CD Pipeline](#cicd-pipeline)
- [Observability](#observability)
- [Storage](#storage)
- [Security](#security)
- [Backup and Restore](#backup-and-restore)
- [Troubleshooting](#troubleshooting)
- [Results / Impact](#results--impact)
- [Lessons Learned](#lessons-learned)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## 🏗 Overview

This project automates the deployment and management of a K3s Kubernetes cluster on Proxmox virtual machines using Terraform and Ansible.

The cluster is managed using a GitOps workflow with FluxCD. All infrastructure and Kubernetes workloads are declared in Git, and FluxCD continuously reconciles the cluster state with the desired state defined in the repository.

The platform includes:

- Automated VM and infrastructure provisioning
- Automated Kubernetes cluster installation
- GitOps-based application and infrastructure management
- CI/CD pipeline using GitHub Actions
- Monitoring and logging stack
- Service mesh visibility with Istio and Kiali
- Object storage using Garage
- TLS automation using cert-manager

---

## Lab Environment

This platform is deployed on a single-node Proxmox lab environment.

| Resource | Specification |
|---|---:|
| CPU | 4 x Intel(R) Core(TM) i5-3470S CPU @ 2.90GHz (1 Socket) |
| RAM | 16 GB |
| Storage | 2 x 1 TB HDD, 1 x 500 GB HDD |
| Hypervisor | Proxmox VE 9.1.5 |

### VM Allocation

| VM | Role | vCPU | RAM | Disk | Network |
|---|---:|---:|---:|---:|---:|
| super-node-1 | LB / NFS / S3 | 1 | 2 GB | 32 GB / 100 GB | 10.10.16.4 / vLan 16 |
| k3s-master-1 | K3s control plane | 1 | 4 GB | 32 GB | 10.10.16.11 / vLan 16 |
| k3s-worker-1 | K3s worker node | 1 | 8 GB | 32 GB | 10.10.16.21 / vLan 16 |

> Note: The lab environment has limited resources. Workloads are configured with resource requests/limits, and non-essential components can be disabled to reduce memory usage.

### Traffic Flow
```mermaid
flowchart LR
    EndUser[End User] -->|HTTPS| CloudFlare[Cloudflare]
    CloudFlare -->|Cloudflare Tunnel| LB[Load Balancer]
    LB -->|HTTP| GatewayAPI[Gateway API]
    GatewayAPI -->|HTTP| HTTPRoute[HTTPRoute]
    HTTPRoute -->|HTTP| Service[Kubernetes Service]
    Service -->|HTTP| Apps[Application Pods]
```

---

## Architecture

```mermaid
flowchart LR
    Dev[Developer]
    AppRepo[Application Repository]
    GHActions[GitHub Actions]
    Registry[Container Registry]
    GitOpsRepo[GitOps Repository]
    FluxCD[FluxCD]
    K3s[K3s Cluster]
    Apps[Applications]
    Monitoring[Prometheus / Grafana / Kiali]
    Logging[Quickwit / Vector]
    Storage[NFS / Garage Object Storage]

    Dev -->|git push| AppRepo
    AppRepo -->|trigger pipeline| GHActions
    GHActions -->|build and push image| Registry
    GHActions -->|update image tag| GitOpsRepo
    FluxCD -->|sync manifests| GitOpsRepo
    FluxCD -->|deploy workloads| K3s
    K3s --> Apps
    K3s --> Monitoring
    K3s --> Logging
    K3s --> Storage
```

---

## Network Topology Diagram
```mermaid
flowchart TD
    subgraph Proxmox ["Proxmox Host (vLan 16)"]
        SN[super-node-1<br/>10.10.16.4<br/>LB/NFS/S3]
        M[k3s-master-1<br/>10.10.16.11<br/>Control Plane]
        W[k3s-worker-1<br/>10.10.16.21<br/>Worker]
    end
    Internet[Internet] -->|via Cloudflare Tunnel| SN
    SN --> M
    SN --> W
    M -->|mounted| SN
    W -->|mounted| SN
```

---

## Deployment Time Estimates

| Phase | Estimated Time |
|---|---|
| Terraform provision (VM creation) | ~3 min |
| Ansible prerequisites + K3s install | ~6 min |
| Istio install | ~1 min |
| Storage networking setup | ~1 min |
| S3 bucket setup | ~1 min |
| FluxCD bootstrap | ~1 min |
| First full reconciliation | ~6 min |
| **Total** | **~20-25 min** |

---

## Features
- Provision virtual machines on Proxmox using Infrastructure as Code
- Automate system prerequisites and kernel tuning with Ansible
- Install and configure a K3s Kubernetes cluster
- Manage Kubernetes workloads using GitOps with FluxCD
- Deploy infrastructure components using Kustomize and Helm
- Build and publish container images using GitHub Actions
- Automatically update GitOps repository when new images are built
- Monitor cluster and workloads with Prometheus and Grafana
- Visualize service mesh traffic with Kiali
- Collect and query logs using Vector and Quickwit
- Provide S3-compatible object storage using Garage
- Automate TLS certificates with cert-manager
- Store cluster configuration declaratively in Git

---

## Prerequisites

Before you begin, ensure the following tools are installed on your workstation:

| Tool | Minimum Version | Install Command |
|---|---|---|
| Terraform | ≥ 1.5 | `brew install hashicorp/tap/terraform` |
| Ansible | ≥ 2.15 | `pip install ansible-core` |
| kubectl | — | `kubectl version --client` |
| flux | ≥ 2.3 | `brew install fluxcd/tap/flux` |
| helm | ≥ 3.12 | `brew install helm` |
| Proxmox Terraform Provider | latest | Configured in `terraform/provider.tf` |

### Version Compatibility

| Component | Version |
|---|---|
| K3s | v1.31.x |
| FluxCD | ≥ 2.3 |
| Istio | 1.24.x |
| Proxmox VE | 8.x / 9.x |
| Terraform | ≥ 1.5 |
| Ansible | ≥ 2.15 |

---

## Tech Stack
### Infrastructure
- Proxmox VE
- Terraform
- Ansible
- Linux
- NFS
### Kubernetes
- CRUN Runtime
- Youki Runtime
- K3s
- kubectl
- Helm
- Kustomize
- FluxCD
- Istio
- cert-manager
### CI/CD
- GitHub Actions
- Docker
- Container Registry
### Observability
- Prometheus
- Grafana
- Kiali
- Vector
- Quickwit

### Storage
- NFS
- Garage Object Storage, S3-compatible

---

## Repository Structure

```
.
├── ansible
│   ├── ansible.cfg
│   ├── inventory
│   │   └── hosts
│   ├── playbooks
│   │   ├── 00-prerequisites.yml
│   │   ├── 01-cluster-setup.yml
│   │   ├── 02-servicemesh.yml
│   │   ├── 03-gitops-bootstrap.yml
│   │   ├── 04-storage-networking.yml
│   │   ├── 05-garage-deploy.yml
│   │   ├── files
│   │   │   └── k3s-kubeconfig.yaml
│   │   └── pipeline.sh
│   └── roles
│       ├── common
│       │   └── tasks
│       │       └── main.yml
│       ├── fluxcd
│       │   └── tasks
│       │       └── main.yml
│       ├── garage
│       │   ├── tasks
│       │   │   └── main.yml
│       │   └── templates
│       │       ├── garage.toml.j2
│       │       └── systemd
│       │           └── garage.service.j2
│       ├── istio
│       │   └── tasks
│       │       └── main.yml
│       ├── k3s-agent
│       │   └── tasks
│       │       └── main.yml
│       ├── k3s-server
│       │   └── tasks
│       │       └── main.yml
│       ├── load-balance
│       │   ├── tasks
│       │   │   └── main.yml
│       │   └── templates
│       │       └── haproxy.cfg.j2
│       ├── nfs
│       │   ├── tasks
│       │   │   └── main.yml
│       │   └── templates
│       │       └── exports.j2
│       └── runtime
│           └── tasks
│               └── main.yml
├── backend
│   ├── config
│   │   ├── config.go
│   │   └── jwt.go
│   ├── controllers
│   │   ├── export_controller.go
│   │   ├── exports_list_controller.go
│   │   ├── protected_controller.go
│   │   ├── register_controller.go
│   │   ├── search_controller.go
│   │   └── signin_controller.go
│   ├── exports
│   │   ├── 192.168.1.1_20260401_213305.csv
│   │   └── 192.168.1.1_20260401_213305.zip
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── middleware
│   │   └── jwt.go
│   ├── models
│   │   ├── export.go
│   │   ├── jwt.go
│   │   ├── register.go
│   │   └── search.go
│   ├── routers
│   │   └── routers.go
│   ├── services
│   │   ├── download_service.go
│   │   ├── export_service.go
│   │   ├── jwt.go
│   │   └── quickwit_client.go
│   └── users.db
├── docker
│   ├── backend
│   │   └── Dockerfile
│   └── frontend
│       └── Dockerfile
├── docker-compose.yml
├── frontend
│   ├── astro.config.mjs
│   ├── frontend.md
│   ├── package.json
│   ├── package-lock.json
│   ├── public
│   │   ├── favicon.ico
│   │   └── favicon.svg
│   ├── src
│   │   ├── assets
│   │   │   ├── astro.svg
│   │   │   ├── background.svg
│   │   │   └── logout.svg
│   │   ├── components
│   │   │   └── Welcome.astro
│   │   ├── layouts
│   │   │   └── Layout.astro
│   │   ├── pages
│   │   │   ├── 404.astro
│   │   │   ├── dashboard.astro
│   │   │   ├── index.astro
│   │   │   ├── signin.astro
│   │   │   └── signup.astro
│   │   ├── scripts
│   │   │   └── dashboard.js
│   │   └── styles
│   │       ├── logs.css
│   │       ├── main.css
│   │       ├── signin.css
│   │       └── signup.css
│   └── tsconfig.json
├── README.md
├── start-app.sh
└── terraform
    ├── backend.tf
    ├── hosts.tpl
    ├── locals.tf
    ├── main.tf
    ├── provider.tf
    ├── secrets.auto.tfvars
    └── variables.tf

49 directories, 81 files
```

---

## Quickstart

### Prerequisites Checklist

Before running the setup, ensure you have completed the following:

- [ ] Terraform ≥ 1.5 installed locally
- [ ] Ansible ≥ 2.15 installed locally
- [ ] kubectl and flux CLI installed
- [ ] Proxmox VE 9.1.5 server with VMs configured
- [ ] GitHub repository created with GitHub Actions enabled
- [ ] `.env` file created with required environment variables
- [ ] SSH access configured to Proxmox VMs

### Clone the repository
```bash
git clone https://github.com/traipoap/gitops-platform.git
cd gitops-platform
```

### Copy the example environment file
```bash
cp .env.example .env
```

### Edit .env with your environment values
```bash
# Proxmox API token
echo 'proxmox_password = "xxx"' > terraform/secrets.auto.tfvars
# RPC secret between nodes 
export RPC_SECRET="$(openssl rand -hex 32)"
# Admin token for the application backend 
export ADMIN_TOKEN="$(openssl rand -base64 32)"
# Garage S3-compatible object storage credentials
export GARAGE_DEFAULT_ACCESS_KEY="GK$(openssl rand -hex 16)"
export GARAGE_DEFAULT_SECRET_KEY="$(openssl rand -hex 32)"
# GitHub PAT for GitOps repository access (FluxCD bootstrap) GitHub PAT
export APP_GIT_SECRET="xxx"
# Container registry URL (e.g. ghcr.io/traipoap/gitops-platform) GitHub PAT
export GITHUB_REGISTRY="xxx"
```

### Provision infrastructure
```bash
cd terraform
terraform init
terraform plan
terraform apply
```

### Run Ansible playbooks

#### Option A: Environment variables (recommended — no interactive prompts)

**bash / zsh:**
```bash
source .env
cd ansible
ansible-playbook -i inventory/hosts playbooks/00-prerequisites.yml
ansible-playbook -i inventory/hosts playbooks/01-cluster-setup.yml
ansible-playbook -i inventory/hosts playbooks/02-servicemesh.yml
ansible-playbook -i inventory/hosts playbooks/03-storage-networking.yml
ansible-playbook -i inventory/hosts playbooks/04-garage-deploy.yml
ansible-playbook -i inventory/hosts playbooks/05-gitops-bootstrap.yml
```

**fish shell:**
```bash
set -x RPC_SECRET (cat .env | grep RPC_SECRET | cut -d'"' -f2)
set -x ADMIN_TOKEN (cat .env | grep ADMIN_TOKEN | cut -d'"' -f2)
set -x GARAGE_DEFAULT_ACCESS_KEY (cat .env | grep GARAGE_DEFAULT_ACCESS_KEY | cut -d'"' -f2)
set -x GARAGE_DEFAULT_SECRET_KEY (cat .env | grep GARAGE_DEFAULT_SECRET_KEY | cut -d'"' -f2)
set -x APP_GIT_SECRET (cat .env | grep APP_GIT_SECRET | cut -d'"' -f2)
set -x GITHUB_REGISTRY (cat .env | grep GITHUB_REGISTRY | cut -d'"' -f2)
cd ansible
ansible-playbook -i inventory/hosts playbooks/00-prerequisites.yml
ansible-playbook -i inventory/hosts playbooks/01-cluster-setup.yml
ansible-playbook -i inventory/hosts playbooks/02-servicemesh.yml
ansible-playbook -i inventory/hosts playbooks/03-storage-networking.yml
ansible-playbook -i inventory/hosts playbooks/04-garage-deploy.yml
ansible-playbook -i inventory/hosts playbooks/05-gitops-bootstrap.yml
```

> When environment variables are set, Ansible uses them directly and skips interactive prompts.

#### Option B: Interactive mode (no env vars — prompts for input)
```bash
cd ansible
ansible-playbook -i inventory/hosts playbooks/00-prerequisites.yml
ansible-playbook -i inventory/hosts playbooks/01-cluster-setup.yml
ansible-playbook -i inventory/hosts playbooks/02-servicemesh.yml
ansible-playbook -i inventory/hosts playbooks/03-storage-networking.yml
ansible-playbook -i inventory/hosts playbooks/04-garage-deploy.yml # Prompts for garage credentials (blank to auto-generate)
ansible-playbook -i inventory/hosts playbooks/05-gitops-bootstrap.yml # Prompts for FluxCD and GitHub PAT
```

---

## GitOps Workflow
This project uses FluxCD as the GitOps operator.
FluxCD watches this Git repository and reconciles the cluster state using:
- GitRepository
- Kustomization
- HelmRepository
- HelmRelease

Bootstrap FluxCD:
```bash
flux bootstrap github \
      --owner=traipoap \
      --repository=fleet-infra \
      --branch=main \
      --path=./clusters/dev \
      --personal
```
Check FluxCD status:
```
flux get all -A
```
Check Kubernetes resources:
```
kubectl get nodes
kubectl get pods -A
kubectl get gitrepositories -A
kubectl get kustomizations -A
kubectl get helmreleases -A
```
---

## CI/CD Pipeline
The CI/CD pipeline uses GitHub Actions.

When code is pushed to the repository, the pipeline performs:
1. Lint and test the application
2. Build a Docker image
3. Push the image to a container registry
4. Update the GitOps repository with the new image tag
5. FluxCD detects the change and deploys the new version

### Example workflow
```mermaid
flowchart LR
    A["git push"] --> B["GitHub Actions"]
    B --> C["Build image"]
    C --> D["Push image to registry"]
    D --> E["Update GitOps repo"]
    E --> F["FluxCD sync"]
    F --> G["Application updated"]
```

---

## Observability
### Prometheus
Prometheus is used to collect metrics from the cluster and workloads.

**Access locally:**
```bash
kubectl -n istio-system port-forward svc/prometheus-server 9090:9090
```
Open: <http://localhost:9090>

### Grafana
Grafana is used for dashboards and visualization.

**Access locally:**
```bash
kubectl -n istio-system port-forward svc/grafana 3000:3000
```
Open: <http://localhost:3000>

### Kiali
Kiali provides visibility into Istio service mesh traffic.

**Access locally:**
```bash
kubectl -n istio-system port-forward svc/kiali 20001:20001
```
Open: <http://localhost:20001>

### Logging
Logs are collected using Vector and stored/queryable in Quickwit.

**View log pods:**
```bash
kubectl -n logging get pods
kubectl -n logging logs -l app.kubernetes.io/name=vector
```
**Access Quickwit UI:**
```bash
kubectl -n logging port-forward svc/quickwit 7280:7280
```
Open: <http://localhost:7280>

---

## Storage
### NFS
NFS is used for persistent storage in this lab environment.

**Example StorageClass:**
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: nfs-client
provisioner: cluster.local/nfs-subdir-external-provisioner
parameters:
  server: 10.10.16.4
  path: /nfs
reclaimPolicy: Retain
volumeBindingMode: Immediate
```

### Garage Object Storage
Garage provides S3-compatible object storage.

**Example use cases:**
- Backup storage
- Application assets
- Registry storage backend
- S3-compatible testing

---

## Security
### Implemented practices
- TLS certificates managed by cert-manager
- Declarative configuration managed through Git
- Secrets separated from application manifests
- Namespace isolation
- RBAC for access control
- Pod security best practices
- Automated reconciliation to reduce configuration drift

### Recommended additional improvements
- SOPS + age for secret encryption
- External Secrets Operator
- Kyverno or OPA Gatekeeper policies
- NetworkPolicies
- Container image scanning
- Signed container images
- Private registry authentication

---

## Backup and Restore
### Recommended backup targets
- etcd snapshots
- Kubernetes resource manifests
- Persistent volumes
- GitOps repository
- Secrets
- NFS data
- Garage object storage data

### Example tools
- Velero
- Kopia
- Restic
- etcd snapshot
- Longhorn snapshots, if using Longhorn

> 📝 Detailed backup and restore procedures are documented in [docs/backup-restore.md](docs/backup-restore.md). If this file does not exist yet, it is planned as part of the Roadmap.

---

## Troubleshooting
Check node status
```
kubectl get nodes -o wide
kubectl describe node <node-name>
```

Check pod status
```
kubectl get pods -A
kubectl describe pod <pod-name> -n <namespace>
kubectl logs <pod-name> -n <namespace>
kubectl logs <pod-name> -n <namespace> --previous
```

Check FluxCD status
```
flux get all -A
flux get kustomizations -A
flux get helmreleases -A
```

Force FluxCD reconciliation
```
flux reconcile source git flux-system -n flux-system
flux reconcile kustomization <name> -n <namespace> --with-source
```

Check HelmRelease status
```
kubectl describe helmrelease <name> -n <namespace>
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
```

---

## Results / Impact
- Reduced infrastructure provisioning time from several hours to approximately 20 minutes
- Eliminated manual Kubernetes deployment steps using GitOps
- Improved environment consistency by defining all workloads declaratively
- Enabled repeatable cluster rebuild from Git and automation scripts
- Improved observability using Prometheus, Grafana, Kiali, and centralized logging
- Reduced configuration drift through continuous reconciliation by FluxCD

---

## Lessons Learned
- GitOps improves consistency and auditability compared to manual kubectl apply
- Infrastructure automation requires careful handling of secrets and state files
- Observability should be installed early, not after problems occur
- Backup and restore testing is as important as deployment automation
- Kubernetes troubleshooting requires strong Linux and networking fundamentals

---

## Roadmap
- Add Velero backup and restore testing
- Add Kyverno policy enforcement
- Add NetworkPolicies for namespace isolation
- Add SOPS or External Secrets for secret management
- Add image vulnerability scanning in CI
- Add automated testing for Kubernetes manifests
- Add multi-environment promotion: dev, staging, production
- Add disaster recovery runbook with measured RTO/RPO
- Add Longhorn or Ceph for more production-grade storage
- Add OpenTelemetry tracing

---

## Contributing
Contributions are welcome! Please feel free to:

- Open an [issue](https://github.com/traipoap/gitops-platform/issues) for bugs or feature requests
- Submit a [pull request](https://github.com/traipoap/gitops-platform/pulls) for improvements

When contributing, please ensure your changes follow the existing style and include updated documentation.

---

## License
This project is licensed under the MIT License.
