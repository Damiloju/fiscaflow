# FiscaFlow

FiscaFlow is a modern, high-performance personal finance tracker backend built with Go. It features robust domain-driven design, full observability with OpenTelemetry, and a clean monolith architecture that can be easily migrated to microservices.

## 🚀 Features
- User management and authentication
- Transaction tracking and categorization
- Budgeting and analytics
- Multi-channel notifications
- Audit logging and compliance
- Full observability (tracing, metrics, logging)
- Scalable, production-ready deployment (Docker, Kubernetes)

## 🛠️ Tech Stack
- **Language:** Go
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL
- **Cache:** Redis
- **Search:** Elasticsearch
- **File Storage:** MinIO
- **Queue:** RabbitMQ
- **Observability:** OpenTelemetry, Zap
- **Containerization:** Docker, Kubernetes

## 📦 Project Structure
See `specs/system-architecture.md` for a detailed breakdown.

## 🏁 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)

### Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd fiscaflow
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start all services**
   ```bash
   docker-compose up --build
   ```

4. **Access the application**
   - **FiscaFlow API**: http://localhost:8080
   - **Health Check**: http://localhost:8080/health
   - **Jaeger UI (Tracing)**: http://localhost:16686
   - **MinIO Console**: http://localhost:9001 (admin/minioadmin)
   - **RabbitMQ Management**: http://localhost:15672 (guest/guest)
   - **Elasticsearch**: http://localhost:9200

### Development Setup

For local development without Docker:

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up environment**
   ```bash
   cp .env.example .env
   # Configure your local database and other services
   ```

3. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

### API Endpoints

#### Authentication
- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/login` - User login
- `POST /api/v1/users/logout` - User logout
- `POST /api/v1/users/refresh` - Refresh access token

#### User Management
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile

#### Health & Monitoring
- `GET /health` - Health check
- `GET /ready` - Readiness check

### Docker Services

The Docker Compose setup includes:

| Service | Port | Description |
|---------|------|-------------|
| fiscaflow | 8080 | Main application |
| postgres | 5432 | Primary database |
| redis | 6379 | Cache and sessions |
| elasticsearch | 9200 | Search and analytics |
| minio | 9000/9001 | File storage (API/Console) |
| rabbitmq | 5672/15672 | Message queue (AMQP/Management) |
| jaeger | 16686/4317 | Distributed tracing (UI/OTLP) |

### Environment Variables

Key environment variables (see `.env.example` for complete list):

```bash
# Database
DATABASE_URL=postgres://postgres:postgres@postgres:5432/fiscaflow?sslmode=disable

# Redis
REDIS_URL=redis:6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=168h

# OpenTelemetry
OTEL_ENDPOINT=http://jaeger:4317
```

### Development Workflow

1. **Start services for development**
   ```bash
   docker-compose up -d postgres redis elasticsearch minio rabbitmq jaeger
   ```

2. **Run application locally**
   ```bash
   go run cmd/server/main.go
   ```

3. **Run tests**
   ```bash
   make test
   ```

4. **Build and test Docker image**
   ```bash
   docker-compose build fiscaflow
   docker-compose up fiscaflow
   ```

### Production Deployment

For production deployment:

1. **Update environment variables**
   - Use strong secrets for JWT_SECRET
   - Configure production database credentials
   - Set up proper SSL certificates

2. **Build production image**
   ```bash
   docker build -t fiscaflow:latest .
   ```

3. **Deploy with Kubernetes**
   ```bash
   kubectl apply -f k8s/
   ```

### Monitoring & Observability

- **Jaeger**: Distributed tracing and request flow analysis
- **Health Checks**: Built-in health endpoints for monitoring
- **Logging**: Structured JSON logging with correlation IDs
- **Metrics**: OpenTelemetry metrics collection

### Troubleshooting

**Common Issues:**

1. **Port conflicts**: Ensure ports 8080, 5432, 6379, 9200, 9000, 5672, 16686 are available
2. **Database connection**: Wait for postgres to be ready before starting the app
3. **Memory issues**: Elasticsearch requires at least 2GB RAM

**Useful Commands:**

```bash
# View logs
docker-compose logs fiscaflow

# Restart specific service
docker-compose restart fiscaflow

# Clean up volumes
docker-compose down -v

# Check service health
docker-compose ps
```

## 📄 License
Proprietary. All rights reserved. 