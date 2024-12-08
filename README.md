# URL Shortener Service

A high-performance URL shortening service built with Go-kit, featuring Redis caching and RESTful API endpoints.

## Documentation

- [Features](docs/FEATURES.md)
- [Getting Started](docs/GETTING_STARTED.md)
- [API Documentation](docs/API.md)
- [Development Guide](docs/DEVELOPMENT.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [Contributing Guide](docs/CONTRIBUTING.md)

## Quick Start

1. Clone the repository
```bash
git clone https://github.com/diki-haryadi/go-url-shortener
cd go-url-shortener
```

2. Start the services
```bash
make docker-up
```

3. Access the services
- API: http://localhost:8080
- Redis Commander: http://localhost:8081 (admin/admin)

## Tech Stack

- Go 1.21+
- Fiber
- Redis (Primary database and caching)
- Docker & Docker Compose
- Make (Build automation)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.