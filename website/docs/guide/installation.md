# Installation Guide

Detailed installation instructions for different environments.

## Docker Installation (Recommended)

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+

### Steps

1. Clone the repository:
```bash
git clone https://github.com/zjxzjw/Carrot-Agent.git
cd carrot-agent
```

2. Create environment file:
```bash
cp .env.example .env
```

3. Edit `.env` and add your API key

4. Start the service:
```bash
docker-compose up -d
```

5. Verify:
```bash
curl http://localhost:8080/health
```

## Local Installation

### Prerequisites

- Go 1.22+
- GCC (for go-sqlite3)
- Node.js 18+ (for web UI)

### Build from Source

```bash
# Clone repository
git clone https://github.com/zjxzjw/Carrot-Agent.git
cd carrot-agent

# Build backend
make build

# Build frontend
cd ui
npm install
npm run build

# Run API server
./bin/carrot-agent-api
```

## Kubernetes Deployment

Create `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: carrot-agent
spec:
  replicas: 1
  selector:
    matchLabels:
      app: carrot-agent
  template:
    metadata:
      labels:
        app: carrot-agent
    spec:
      containers:
      - name: carrot-agent
        image: carrotagent/carrot-agent:latest
        ports:
        - containerPort: 8080
        env:
        - name: CARROT_API_KEY
          valueFrom:
            secretKeyRef:
              name: carrot-secrets
              key: api-key
        volumeMounts:
        - name: data
          mountPath: /home/carrot/.carrot
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: carrot-data-pvc
```

## Configuration

See [Configuration Guide](/guide/configuration) for detailed options.

## Verification

After installation, test the API:

```bash
# Health check
curl http://localhost:8080/health

# Test chat
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello!", "session_id": "test"}'
```

## Next Steps

- [Quick Start](/guide/quick-start)
- [Configuration](/guide/configuration)
- [Docker Deployment](/guide/docker)
