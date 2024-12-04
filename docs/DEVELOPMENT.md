# Development Guide

## Project Structure

```
├── cmd/
│   └── main.go
├── pkg/
│   ├── service/
│   │   ├── service.go
│   │   ├── middleware.go
│   │   └── service_test.go
│   ├── endpoint/
│   │   ├── endpoint.go
│   │   └── middleware.go
│   ├── transport/
│   │   ├── http.go
│   │   └── metrics.go
│   └── database/
│       └── redis.go
├── docker/
│   ├── api/
│   └── redis/
├── docs/
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make
- Redis CLI (optional)

## Development Setup

1. Clone the repository
```bash
git clone https://github.com/yourusername/urlshortener.git
cd urlshortener
```

2. Install dependencies
```bash
make deps
```

3. Start development environment
```bash
make docker-up
```

4. Build the application
```bash
make build
```

## Available Make Commands

### Development
- `make dev`: Run application in development mode
- `make build`: Build the application
- `make clean`: Clean up built binaries
- `make deps`: Download dependencies
- `make lint`: Run linter

### Docker
- `make docker-build`: Build docker images
- `make docker-up`: Start docker containers
- `make docker-down`: Stop docker containers
- `make docker-logs`: View docker logs
- `make docker-rebuild`: Rebuild and restart containers

### Testing
- `make test`: Run tests
- `make test-coverage`: Run tests with coverage
- `make test-integration`: Run integration tests

### Database
- `make redis-cli`: Connect to Redis CLI

## Testing

### Unit Tests
```bash
make test
```

### Integration Tests
```bash
make test-integration
```

### Coverage Report
```bash
make test-coverage
```

## Code Style

### Go Formatting
```bash
# Format code
go fmt ./...

# Run linter
make lint
```

### Commit Messages
- Follow conventional commits
- Use present tense
- Keep it concise

## Development Workflow

1. Create a new branch
```bash
git checkout -b feature/your-feature
```

2. Make changes and test
```bash
make test
make lint
```

3. Commit changes
```bash
git add .
git commit -m "feat: add new feature"
```

4. Push and create PR
```bash
git push origin feature/your-feature
```

## Debugging

### Redis Commander
- Access at http://localhost:8081
- Username: admin
- Password: admin

### Logs
```bash
# View service logs
make docker-logs

# View Redis logs
docker-compose logs redis
```

## Environment Variables

```env
# App Settings
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=debug

# Redis Settings
REDIS_URL=redis:6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# Rate Limiting
RATE_LIMIT=100
RATE_WINDOW=60
```

## Common Issues

### Redis Connection
```bash
# Check Redis connection
redis-cli -h localhost -p 6379 ping
```

### Permission Issues
```bash
# Fix permissions
chmod +x scripts/*.sh
```

## Best Practices

1. Always write tests
2. Document your code
3. Follow Go idioms
4. Use meaningful variable names
5. Handle errors appropriately
6. Add logging for debugging
7. Keep functions small and focused