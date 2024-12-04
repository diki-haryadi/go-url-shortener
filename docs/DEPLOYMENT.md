# Deployment Guide

## Deployment Options

1. Docker Compose (Simple)
2. Kubernetes (Recommended for production)
3. Bare Metal

## Docker Compose Deployment

### Prerequisites
- Docker
- Docker Compose
- Access to container registry

### Steps

1. Build the image
```bash
make docker-build
```

2. Configure environment
```bash
cp .env.example .env
# Edit .env with your production values
```

3. Deploy
```bash
make docker-up
```

## Kubernetes Deployment

### Prerequisites
- Kubernetes cluster
- kubectl configured
- Helm (optional)

### Configuration Files

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: url-shortener
spec:
  replicas: 3
  selector:
    matchLabels:
      app: url-shortener
  template:
    metadata:
      labels:
        app: url-shortener
    spec:
      containers:
      - name: url-shortener
        image: your-registry/url-shortener:latest
        ports:
        - containerPort: 8080
        env:
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: url-shortener-config
              key: REDIS_URL
```

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: url-shortener
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: url-shortener
```

### Deployment Steps

1. Apply configurations
```bash
kubectl apply -f k8s/
```

2. Verify deployment
```bash
kubectl get pods
kubectl get services
```

3. Configure ingress (if needed)
```bash
kubectl apply -f k8s/ingress.yaml
```

## Production Considerations

### Security
1. Enable TLS/SSL
2. Use secure Redis connection
3. Implement rate limiting
4. Set up monitoring
5. Configure backups

### Scaling
1. Configure horizontal pod autoscaling
2. Use Redis cluster
3. Implement caching strategies
4. Configure load balancing

### Monitoring
1. Set up Prometheus
2. Configure Grafana dashboards
3. Implement logging
4. Set up alerts

### Backup
1. Configure Redis persistence
2. Set up regular backups
3. Test recovery procedures

## Environment Variables

```env
# App Configuration
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info

# Redis Configuration
REDIS_URL=redis://redis-cluster:6379
REDIS_PASSWORD=secure_password
REDIS_DB=0

# Rate Limiting
RATE_LIMIT=100
RATE_WINDOW=60

# Metrics
ENABLE_METRICS=true
METRICS_PORT=9090
```

## Health Checks

### Liveness Probe
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 3
```

### Readiness Probe
```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

## Rollback Procedure

1. Identify the issue
```bash
kubectl logs deployment/url-shortener
```

2. Rollback deployment
```bash
kubectl rollout undo deployment/url-shortener
```

3. Verify rollback
```bash
kubectl rollout status deployment/url-shortener
```

## Performance Tuning

### Redis Configuration
```conf
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### Application Settings
```env
GOMAXPROCS=8
```

## Troubleshooting

### Common Issues

1. Redis Connection
```bash
kubectl exec -it pod/redis-0 -- redis-cli ping
```

2. Service Discovery
```bash
kubectl describe service url-shortener
```

3. Log Analysis
```bash
kubectl logs -l app=url-shortener --tail=100
```

### Monitoring Checklist

- [ ] CPU Usage
- [ ] Memory Usage
- [ ] Redis Connections
- [ ] Request Latency
- [ ] Error Rates
- [ ] Cache Hit Ratio