# Volume API

Volume management API for the GoPaaS platform.

## Quick Start

### 1. Generate Go code from proto files
```bash
make proto
```

### 2. Build the service
```bash
make build
```
This will generate a volumeApi binary file

### 3. Build Docker image
```bash
make docker
```
This will automatically generate a volumeApi:latest image
You can verify with: `docker images | grep volumeApi`

## Technology Stack

This service uses go-micro v3 as the microservice development framework.
Framework: https://github.com/asim/go-micro

## Development Workflow

### Prerequisites
* cap1573/cap-tool - Project scaffolding tool
* cap1573/cap-v3 - Code generation tool for proto and go-micro
* cap1573/cap-api-gateway - Unified API gateway for go-micro v3

### 1. Backend Service Development
* 1.1 Generate project structure using cap-tool
* 1.2 Develop domain - model layer
* 1.3 Develop domain - repository layer
* 1.4 Develop domain - service layer
* 1.5 Create proto files and generate code using cap1573/cap-v3
* 1.6 Develop exposed services
* 1.7 Write main.go
* 1.8.1 Package Docker image, write Dockerfile (for k8s operations, copy or mount .kube/config file)
* 1.8.2 When packaging Docker, ensure registry and tracing addresses use non-internal network addresses to avoid access failures
* 1.8.3 Expose and map circuit breaker and monitoring system ports to collect data
* 1.8.4 Fix service external port: micro.Address(":8081")
* 1.8.5 Correct custom service addresses, use service names in containers
* 1.8.6 Add MySQL connection information
* 1.8.7 Update MySQL connection address in Consul
* 1.8.8 Cross-compile

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o <service-name> *.go
```
* 1.8.9 Build

```bash
sudo docker build cap1573/<service-name>
```
* 1.8.10 Run Docker

```bash
sudo docker run -p 8081:8081 -p 9092:9092 -p 9192:9192 \
  -v /absolute/path/.kube/config:/root/.kube/config \
  cap1573/<service-name>
```

### 2. API Interface Development
* 2.1 Generate project structure using cap-tool
* 2.2 Write API proto files and generate code using cap1573/cap-v3
* 2.3 Develop exposed API interfaces
* 2.4 Write main.go
* 2.5 Package Docker image

### 3. Start Gateway
* 3.1 Use cap-api-gateway to establish gateway

```bash
sudo docker run -d -p 8080:8080 cap1573/cap-api-gateway \
  --registry=consul \
  --registry_address=<registry-address>:8500 \
  api --handler=api
```
Note: Use host addresses that are accessible, not internal network addresses.

### 4. Frontend Application Development


       
