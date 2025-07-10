# Prometheus Installation Guide

## 1. Extract the zip package
```bash
unzip v0.9.0.zip
```

## 2. Navigate to directory
```bash
cd kube-prometheus-0.9.0
```

## 3. Execute installation commands
```bash
# Create namespace and CRD
kubectl create -f manifests/setup

# Wait for creation to complete
until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done

# Create monitoring components
kubectl create -f manifests/

# View pods in monitoring namespace
kubectl get pods -n monitoring
```

## 4. Add external access routing
Add grafana.caplost.com domain for the monitoring namespace through the route functionality developed in the route management service.

First login default credentials: admin / admin
