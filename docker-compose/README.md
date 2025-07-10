# Docker Compose Setups

This directory contains different Docker Compose configurations for various deployment scenarios.

## Available Setups

### 1. Basic Setup (`basic-setup/`)

A minimal configuration for development and testing:

- **base-service**: Single Go service with hot reload
- **Ports**: 5000 (API), 8080 (Health check)

**Usage:**
```bash
cd basic-setup
docker-compose up -d
```

### 2. Production Stack (`production-stack/`)

A comprehensive production-ready setup with all supporting services:

#### Core Services
- **Consul Cluster**: 4-node service discovery cluster (ports 8500 for UI)
- **MySQL**: Database with persistent storage (port 3306)
- **Jaeger**: Distributed tracing (ports 6831/UDP, 16686 for UI)

#### Monitoring & Observability
- **Prometheus**: Metrics collection (port 9090)
- **Grafana**: Monitoring dashboards (port 3000, admin/admin)
- **ELK Stack**: Log aggregation and analysis
  - Elasticsearch (ports 9200, 9300)
  - Logstash (ports 5044, 5000, 9600)
  - Kibana (port 5601)

**Usage:**
```bash
cd production-stack
docker-compose up -d
```

## Quick Start

1. **For Development**: Use `basic-setup/`
2. **For Production**: Use `production-stack/`

## Service URLs

When running the production stack:

- **Consul UI**: http://localhost:8500
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger UI**: http://localhost:16686
- **Kibana**: http://localhost:5601
- **Prometheus**: http://localhost:9090
- **Elasticsearch**: http://localhost:9200

## Notes

- All data is persisted in the `./mysql` directory for the production stack
- The Hystrix dashboard is commented out but can be enabled if needed
- Memory limits are set for Elasticsearch and Logstash to run on development machines 