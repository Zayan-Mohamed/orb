# Docker Deployment

Deploy Orb relay server using Docker.

## Quick Start

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o orb

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/orb .

EXPOSE 8080
CMD ["./orb", "relay", "--host", "0.0.0.0", "--port", "8080"]
```

### Build Image

```bash
docker build -t orb-relay .
```

### Run Container

```bash
docker run -d \
  --name orb-relay \
  -p 8080:8080 \
  --restart unless-stopped \
  orb-relay
```

## Docker Compose

```yaml
version: "3.8"

services:
  relay:
    build: .
    ports:
      - "8080:8080"
    restart: unless-stopped
    environment:
      - ORB_LOG_LEVEL=info
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    depends_on:
      - relay
```

Start:

```bash
docker-compose up -d
```

## Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: orb-relay
spec:
  replicas: 3
  selector:
    matchLabels:
      app: orb-relay
  template:
    metadata:
      labels:
        app: orb-relay
    spec:
      containers:
        - name: orb-relay
          image: orb-relay:latest
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: orb-relay
spec:
  selector:
    app: orb-relay
  ports:
    - port: 8080
      targetPort: 8080
```

Deploy:

```bash
kubectl apply -f orb-deployment.yaml
```

## Management

```bash
# View logs
docker logs -f orb-relay

# Restart container
docker restart orb-relay

# Stop container
docker stop orb-relay

# Remove container
docker rm orb-relay
```

## Next Steps

- [Production Deployment](production.md)
- [Relay Server Setup](relay-server.md)
