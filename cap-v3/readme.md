# Code Generation Tool

The executable tool is built for Linux environment with Go version 1.16.

## Installation of protoc and protoc-gen-go

### 1. Download the corresponding version of protoc
```bash
# In Dockerfile, use the following command in the image
apk add protoc
```

### 2. Download and install the Go plugin protoc-gen-go
```bash
go install github.com/golang/protobuf/protoc-gen-go@v1.27
```
Note: Use protoc-gen-go@v1.27 version

### 3. Install gen-micro
Choose the v3 version here
During development, go-micro v4 was just released and tested, but there are still compatibility issues with proto. Switch when stable.
```bash
go get -u github.com/asim/go-micro/cmd/protoc-gen-micro/v3
```

### 4. After installation, run on the installed machine
```bash
protoc --version
``` 