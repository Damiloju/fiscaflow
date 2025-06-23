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
- Go 1.21+
- Docker & Docker Compose

### Setup
1. Clone the repository
2. Copy `.env.example` to `.env` and configure as needed
3. Start services:
   ```bash
   docker-compose up --build
   ```
4. Run the app locally:
   ```bash
   go run cmd/server/main.go
   ```

### Development
- All prompts and decisions are tracked in `Prompt.md`
- See `specs/` for architecture and technical details

## 📄 License
Proprietary. All rights reserved. 