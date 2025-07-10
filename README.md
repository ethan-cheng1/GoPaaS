# GoPaaS - Cloud-Native Platform as a Service

A comprehensive cloud-native Platform as a Service (PaaS) built with Go microservices and Kubernetes.

## Overview

GoPaaS is a modern cloud platform that provides container orchestration, service management, and infrastructure automation. Built with microservices architecture using Go and go-micro v3, it offers a complete solution for managing cloud-native applications.

## Architecture

### Core Services

- **Pod Management** (`pod/`, `podapi/`) - Container orchestration and lifecycle management
- **Service Management** (`svc/`, `svcapi/`) - Microservice deployment and management
- **Route Management** (`route/`, `routeapi/`) - API gateway and routing
- **Volume Management** (`volume/`, `volumeapi/`) - Storage and persistent volume management
- **Base Service** (`base/`) - Core platform functionality

### Infrastructure Components

- **API Gateway** (`cap-api-gateway/`) - Unified API access point
- **Service Registry** - Consul-based service discovery
- **Monitoring** - Prometheus and Grafana integration
- **Tracing** - Jaeger distributed tracing
- **Database** - MySQL for persistent storage

### Frontend

- **Web UI** (`go-paas-front/`, `go-paas-html/`) - Management interface

## Quick Start

### Prerequisites

- Kubernetes cluster (1.21.5+)
- Docker
- Go 1.16+

### 1. Infrastructure Setup

Start the required infrastructure services:

```bash
cd docker-compose/production-stack/
docker-compose up -d
```

This starts:
- Consul (Service Discovery) - Port 8500
- MySQL (Database) - Port 3306
- Jaeger (Tracing) - Port 16686
- Prometheus (Monitoring) - Port 9090
- Grafana (Dashboards) - Port 3000

### 2. Build and Deploy Services

For each service:

```bash
cd <service-directory>
make proto    # Generate Go code from proto files
make build    # Compile the service
make docker   # Build Docker image
```

### 3. Start API Gateway

```bash
docker run -d -p 8080:8080 cap1573/cap-api-gateway \
  --registry=consul \
  --registry_address=<consul-address>:8500 \
  api --handler=api
```

### 4. Access the Platform

- **API Gateway**: http://localhost:8080
- **Consul UI**: http://localhost:8500
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686
- **Prometheus**: http://localhost:9090

## Development

### Technology Stack

- **Backend**: Go with go-micro v3
- **Service Discovery**: Consul
- **Container Orchestration**: Kubernetes
- **Monitoring**: Prometheus + Grafana
- **Tracing**: Jaeger
- **Database**: MySQL
- **Frontend**: HTML/CSS/JavaScript

### Development Workflow

1. **Service Development**
   - Domain-driven design (Model, Repository, Service layers)
   - Protocol Buffers for API definitions
   - Code generation with cap-v3 tool

2. **API Development**
   - RESTful API interfaces
   - gRPC communication
   - API gateway integration

3. **Deployment**
   - Docker containerization
   - Kubernetes manifests
   - Service mesh integration

### Tools

- **cap-tool** - Project scaffolding
- **cap-v3** - Code generation for proto and go-micro
- **cap-api-gateway** - Unified API gateway

## Features

### Container Management
- Pod lifecycle management
- Container orchestration
- Resource allocation
- Health monitoring

### Service Management
- Microservice deployment
- Service discovery
- Load balancing
- Auto-scaling

### Storage Management
- Persistent volume provisioning
- Storage class management
- Volume lifecycle

### Monitoring & Observability
- Metrics collection
- Distributed tracing
- Log aggregation
- Alert management

### API Management
- Unified API gateway
- Rate limiting
- Authentication/Authorization
- API versioning

## Configuration

The platform uses Consul for configuration management:

- Database connections
- Service registry settings
- Tracing endpoints
- Monitoring configurations

## Production Deployment

For production environments:

1. **Security**
   - Enable TLS/SSL
   - Implement authentication
   - Configure RBAC

2. **Scalability**
   - Horizontal pod autoscaling
   - Load balancing
   - Resource quotas

3. **Reliability**
   - High availability setup
   - Backup and disaster recovery
   - Monitoring and alerting

4. **Networking**
   - Ingress controllers
   - Service mesh (optional)
   - Network policies

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For issues and questions:
- Check the documentation in each service directory
- Review the development workflow guides
- Examine the configuration examples 