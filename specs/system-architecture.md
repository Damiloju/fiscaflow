# System Architecture Specification

## 🏗️ **Overview**

FiscaFlow follows a well-structured monolith architecture pattern with clean separation of concerns, ensuring maintainability, testability, and easy migration to microservices when needed. The system is designed to handle financial data with enterprise-grade security and observability, built with Go for high performance and reliability.

## 🎯 **Architecture Principles**

- **Monolith First**: Start with a well-structured monolith, migrate to microservices when needed
- **Clean Architecture**: Separation of concerns with clear boundaries
- **Domain-Driven Design**: Business logic organized by domain
- **API-First**: RESTful APIs with comprehensive documentation
- **Security by Design**: Zero-trust security model
- **Observability**: Full OpenTelemetry integration
- **Scalability**: Horizontal scaling with auto-scaling capabilities
- **Resilience**: Circuit breakers, retries, and graceful degradation
- **Performance**: Go-based services for high throughput and low latency

## 🏛️ **High-Level Architecture**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Load Balancer (NGINX)                          │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API Gateway (Kong)                             │
│  - Rate Limiting  - Authentication  - Request Routing  - CORS Handling     │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        FiscaFlow Monolith (Go + Gin)                       │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌─────────────┐ │
│  │   User Domain   │ │ Transaction     │ │   Analytics     │ │   Budget    │ │
│  │   - Auth        │ │   Domain        │ │   Domain        │ │   Domain    │ │
│  │   - Profile     │ │   - CRUD        │ │   - Reports     │ │   - Goals   │ │
│  │   - Permissions │ │   - Categories  │ │   - Insights    │ │   - Alerts  │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └─────────────┘ │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌─────────────┐ │
│  │ Notification    │ │   ML Service    │ │   File Storage  │ │   Audit     │ │
│  │   Domain        │ │   - Categorize  │ │   - Exports     │ │   - Logs    │ │
│  │   - Email       │ │   - Predict     │ │   - Imports     │ │   - Events  │ │
│  │   - SMS         │ │   - Anomaly     │ │   - Attachments │ │   - Events  │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Message Queue (RabbitMQ)                       │
│  - Event Publishing  - Async Processing  - Background Jobs                 │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Data Layer                                     │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ ┌─────────────┐ │
│  │   PostgreSQL    │ │     Redis       │ │   Elasticsearch │ │   MinIO     │ │
│  │  (Primary DB)   │ │   (Cache/Queue) │ │   (Search/Logs) │ │ (File Store)│ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 🏗️ **Monolith Architecture**

### **Project Structure**
```
fiscaflow/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   │   ├── config.go
│   │   └── environment.go
│   ├── domain/                     # Domain models and business logic
│   │   ├── user/
│   │   │   ├── model.go
│   │   │   ├── service.go
│   │   │   └── repository.go
│   │   ├── transaction/
│   │   │   ├── model.go
│   │   │   ├── service.go
│   │   │   └── repository.go
│   │   ├── budget/
│   │   ├── analytics/
│   │   ├── notification/
│   │   └── audit/
│   ├── infrastructure/             # External dependencies
│   │   ├── database/
│   │   │   ├── postgres.go
│   │   │   └── migrations/
│   │   ├── cache/
│   │   │   └── redis.go
│   │   ├── search/
│   │   │   └── elasticsearch.go
│   │   ├── storage/
│   │   │   └── minio.go
│   │   └── messaging/
│   │       └── rabbitmq.go
│   ├── api/                        # HTTP API layer
│   │   ├── handlers/               # HTTP handlers
│   │   │   ├── user.go
│   │   │   ├── transaction.go
│   │   │   ├── budget.go
│   │   │   └── analytics.go
│   │   ├── middleware/             # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logging.go
│   │   │   └── opentelemetry.go
│   │   ├── routes/                 # Route definitions
│   │   │   └── routes.go
│   │   └── server.go               # HTTP server setup
│   ├── observability/              # Observability components
│   │   ├── tracing/
│   │   ├── metrics/
│   │   └── logging/
│   └── shared/                     # Shared utilities
│       ├── errors/
│       ├── utils/
│       └── constants/
├── pkg/                            # Public packages
│   ├── auth/
│   ├── validator/
│   └── crypto/
├── migrations/                     # Database migrations
├── scripts/                        # Build and deployment scripts
├── docs/                           # Documentation
├── tests/                          # Integration tests
├── docker/                         # Docker configurations
├── k8s/                            # Kubernetes manifests
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### **Domain Organization**

#### **1. User Domain**
**Purpose**: User management and authentication

**Core Features**:
- User registration and authentication
- Profile management
- Role-based access control (RBAC)
- Multi-tenancy support
- Session management

**API Endpoints**:
```
POST   /api/v1/users/register
POST   /api/v1/users/login
POST   /api/v1/users/logout
GET    /api/v1/users/profile
PUT    /api/v1/users/profile
POST   /api/v1/users/refresh-token
DELETE /api/v1/users/:id
```

**Go Implementation Structure**:
```go
// internal/domain/user/model.go
package user

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email             string    `json:"email" gorm:"unique;not null"`
    PasswordHash      string    `json:"-" gorm:"not null"`
    FirstName         string    `json:"first_name"`
    LastName          string    `json:"last_name"`
    Phone             string    `json:"phone"`
    DateOfBirth       *time.Time `json:"date_of_birth"`
    Timezone          string    `json:"timezone" gorm:"default:'UTC'"`
    Locale            string    `json:"locale" gorm:"default:'en-US'"`
    Role              UserRole  `json:"role" gorm:"default:'user'"`
    Status            UserStatus `json:"status" gorm:"default:'active'"`
    EmailVerified     bool      `json:"email_verified" gorm:"default:false"`
    PhoneVerified     bool      `json:"phone_verified" gorm:"default:false"`
    TwoFactorEnabled  bool      `json:"two_factor_enabled" gorm:"default:false"`
    LastLoginAt       *time.Time `json:"last_login_at"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

type UserRole string

const (
    UserRoleUser        UserRole = "user"
    UserRolePremium     UserRole = "premium"
    UserRoleAdmin       UserRole = "admin"
    UserRoleFamilyOwner UserRole = "family_owner"
    UserRoleFamilyMember UserRole = "family_member"
)

type UserStatus string

const (
    UserStatusActive    UserStatus = "active"
    UserStatusInactive  UserStatus = "inactive"
    UserStatusSuspended UserStatus = "suspended"
    UserStatusDeleted   UserStatus = "deleted"
)

// internal/domain/user/service.go
package user

import (
    "context"
    "errors"
    "time"
    
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type Service struct {
    repo   Repository
    cache  Cache
    logger Logger
}

type Repository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type Cache interface {
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Get(ctx context.Context, key string, dest interface{}) error
    Delete(ctx context.Context, key string) error
}

type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Check if user already exists
    existing, err := s.repo.GetByEmail(ctx, req.Email)
    if err == nil && existing != nil {
        return nil, errors.New("user already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &User{
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        Role:         UserRoleUser,
        Status:       UserStatusActive,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    s.logger.Info("User registered successfully", "user_id", user.ID)
    return user, nil
}

// internal/domain/user/repository.go
package user

import (
    "context"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

#### **2. Transaction Domain**
**Purpose**: Core financial transaction management

**Core Features**:
- Transaction CRUD operations
- Categorization (manual and ML-powered)
- Reconciliation
- Import/export functionality
- Transaction history

**API Endpoints**:
```
GET    /api/v1/transactions
POST   /api/v1/transactions
GET    /api/v1/transactions/:id
PUT    /api/v1/transactions/:id
DELETE /api/v1/transactions/:id
POST   /api/v1/transactions/import
GET    /api/v1/transactions/export
POST   /api/v1/transactions/categorize
```

**Go Implementation Structure**:
```go
// internal/domain/transaction/model.go
package transaction

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
)

type Transaction struct {
    ID                    uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID                uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
    FamilyID              *uuid.UUID      `json:"family_id" gorm:"type:uuid"`
    AccountID             uuid.UUID       `json:"account_id" gorm:"type:uuid;not null"`
    CategoryID            *uuid.UUID      `json:"category_id" gorm:"type:uuid"`
    
    Amount                decimal.Decimal  `json:"amount" gorm:"type:decimal(15,2);not null"`
    Currency              string          `json:"currency" gorm:"default:'USD'"`
    Description           string          `json:"description" gorm:"not null"`
    Merchant              string          `json:"merchant"`
    Location              datatypes.JSON  `json:"location"`
    
    TransactionDate       time.Time       `json:"transaction_date" gorm:"not null"`
    PostedDate            *time.Time      `json:"posted_date"`
    Status                TransactionStatus `json:"status" gorm:"default:'pending'"`
    
    CategorizationSource  string          `json:"categorization_source" gorm:"default:'manual'"`
    CategorizationConfidence *float64     `json:"categorization_confidence"`
    
    Tags                  pq.StringArray  `json:"tags" gorm:"type:text[]"`
    Notes                 string          `json:"notes"`
    ReceiptURL            string          `json:"receipt_url"`
    
    PlaidTransactionID    string          `json:"plaid_transaction_id"`
    ExternalID            string          `json:"external_id"`
    
    CreatedAt             time.Time       `json:"created_at"`
    UpdatedAt             time.Time       `json:"updated_at"`
    
    // Relationships
    User                  User            `json:"user" gorm:"foreignKey:UserID"`
    Account               Account         `json:"account" gorm:"foreignKey:AccountID"`
    Category              *Category       `json:"category" gorm:"foreignKey:CategoryID"`
}

type TransactionStatus string

const (
    TransactionStatusPending   TransactionStatus = "pending"
    TransactionStatusPosted    TransactionStatus = "posted"
    TransactionStatusCancelled TransactionStatus = "cancelled"
    TransactionStatusDisputed  TransactionStatus = "disputed"
)

// internal/domain/transaction/service.go
package transaction

import (
    "context"
    "time"
    
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type Service struct {
    repo        Repository
    cache       Cache
    mlService   MLService
    eventBus    EventBus
    logger      Logger
}

type Repository interface {
    Create(ctx context.Context, transaction *Transaction) error
    GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
    GetByUserID(ctx context.Context, userID uuid.UUID, filters TransactionFilters) ([]*Transaction, error)
    Update(ctx context.Context, transaction *Transaction) error
    Delete(ctx context.Context, id uuid.UUID) error
    BulkCreate(ctx context.Context, transactions []*Transaction) error
}

type MLService interface {
    CategorizeTransaction(ctx context.Context, description string) (*CategorizationResult, error)
}

type EventBus interface {
    Publish(ctx context.Context, event Event) error
}

func (s *Service) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*Transaction, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Auto-categorize if not provided
    if req.CategoryID == nil {
        categorization, err := s.mlService.CategorizeTransaction(ctx, req.Description)
        if err == nil && categorization != nil {
            req.CategoryID = &categorization.CategoryID
            req.CategorizationSource = "ml"
            req.CategorizationConfidence = &categorization.Confidence
        }
    }
    
    // Create transaction
    transaction := &Transaction{
        UserID:                req.UserID,
        FamilyID:              req.FamilyID,
        AccountID:             req.AccountID,
        CategoryID:            req.CategoryID,
        Amount:                req.Amount,
        Currency:              req.Currency,
        Description:           req.Description,
        Merchant:              req.Merchant,
        Location:              req.Location,
        TransactionDate:       req.TransactionDate,
        Status:                TransactionStatusPending,
        CategorizationSource:  req.CategorizationSource,
        CategorizationConfidence: req.CategorizationConfidence,
        Tags:                  req.Tags,
        Notes:                 req.Notes,
    }
    
    if err := s.repo.Create(ctx, transaction); err != nil {
        return nil, err
    }
    
    // Publish event
    s.eventBus.Publish(ctx, TransactionCreatedEvent{
        TransactionID: transaction.ID,
        UserID:        transaction.UserID,
        Amount:        transaction.Amount,
        Timestamp:     time.Now(),
    })
    
    s.logger.Info("Transaction created successfully", 
        "transaction_id", transaction.ID,
        "user_id", transaction.UserID,
        "amount", transaction.Amount,
    )
    
    return transaction, nil
}
```

#### **3. Analytics Domain**
**Purpose**: Financial analytics and reporting

**Core Features**:
- Spending analysis and trends
- Budget vs actual reporting
- Financial goal tracking
- Custom report generation
- Data visualization APIs

**API Endpoints**:
```
GET    /api/v1/analytics/spending
GET    /api/v1/analytics/budget
GET    /api/v1/analytics/goals
GET    /api/v1/analytics/trends
POST   /api/v1/analytics/reports
GET    /api/v1/analytics/dashboard
```

#### **4. Budget Domain**
**Purpose**: Budget management and goal tracking

**Core Features**:
- Budget creation and management
- Goal setting and tracking
- Alert generation
- Budget recommendations

**API Endpoints**:
```
GET    /api/v1/budgets
POST   /api/v1/budgets
GET    /api/v1/budgets/:id
PUT    /api/v1/budgets/:id
DELETE /api/v1/budgets/:id
GET    /api/v1/goals
POST   /api/v1/goals
PUT    /api/v1/goals/:id
```

#### **5. Notification Domain**
**Purpose**: Multi-channel notification delivery

**Core Features**:
- Email notifications
- SMS alerts
- Push notifications
- In-app notifications
- Notification preferences

#### **6. Audit Domain**
**Purpose**: Audit trail and compliance

**Core Features**:
- User action logging
- Financial transaction audit
- System event tracking
- Compliance reporting

## 📊 **Data Architecture**

### **Primary Database (PostgreSQL)**
**Purpose**: ACID-compliant transactional data

**Key Tables**:
- `users` - User accounts and profiles
- `transactions` - Financial transactions
- `accounts` - Bank and investment accounts
- `categories` - Transaction categories
- `budgets` - Budget definitions
- `goals` - Financial goals
- `notifications` - Notification history
- `audit_logs` - Audit trail

### **Cache Layer (Redis)**
**Purpose**: High-performance caching and session storage

**Use Cases**:
- User session storage
- API response caching
- Rate limiting counters
- Real-time data (dashboards)
- Job queue management

### **Search Engine (Elasticsearch)**
**Purpose**: Full-text search and analytics

**Indices**:
- `transactions` - Transaction search
- `categories` - Category search
- `logs` - Application logs
- `metrics` - Business metrics

### **File Storage (MinIO)**
**Purpose**: Document and file storage

**Buckets**:
- `exports` - CSV/PDF exports
- `imports` - Transaction imports
- `attachments` - Receipt images
- `backups` - Database backups

## 🔄 **Event-Driven Communication**

### **Event Types**
```go
// internal/shared/events/events.go
package events

import (
    "time"
    "github.com/google/uuid"
)

type Event interface {
    Type() string
    Timestamp() time.Time
}

type UserCreatedEvent struct {
    UserID    uuid.UUID `json:"user_id"`
    Email     string    `json:"email"`
    Timestamp time.Time `json:"timestamp"`
}

func (e UserCreatedEvent) Type() string { return "user.created" }
func (e UserCreatedEvent) Timestamp() time.Time { return e.Timestamp }

type TransactionCreatedEvent struct {
    TransactionID uuid.UUID `json:"transaction_id"`
    UserID        uuid.UUID `json:"user_id"`
    Amount        decimal.Decimal `json:"amount"`
    Timestamp     time.Time `json:"timestamp"`
}

func (e TransactionCreatedEvent) Type() string { return "transaction.created" }
func (e TransactionCreatedEvent) Timestamp() time.Time { return e.Timestamp }

type BudgetAlertEvent struct {
    UserID   uuid.UUID `json:"user_id"`
    BudgetID uuid.UUID `json:"budget_id"`
    Severity string    `json:"severity"`
    Timestamp time.Time `json:"timestamp"`
}

func (e BudgetAlertEvent) Type() string { return "budget.alert" }
func (e BudgetAlertEvent) Timestamp() time.Time { return e.Timestamp }
```

### **Message Queue (RabbitMQ)**
**Configuration**:
```yaml
# rabbitmq.yml
exchanges:
  - name: fiscaflow.events
    type: topic
    durable: true

queues:
  - name: fiscaflow.events
    exchange: fiscaflow.events
    routing_key: "#"
```

**Go Event Publisher**:
```go
// internal/infrastructure/messaging/rabbitmq.go
package messaging

import (
    "context"
    "encoding/json"
    
    "github.com/streadway/amqp"
    "github.com/fiscaflow/internal/shared/events"
)

type RabbitMQEventBus struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func (r *RabbitMQEventBus) Publish(ctx context.Context, event events.Event) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    return r.channel.Publish(
        "fiscaflow.events",
        event.Type(),
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        data,
        },
    )
}
```

## 🔒 **Security Architecture**

### **Authentication Flow**
```
1. User Login Request
   ↓
2. Validate Credentials
   ↓
3. Generate JWT Access Token (15min)
   ↓
4. Generate Refresh Token (7 days)
   ↓
5. Store Refresh Token in Redis
   ↓
6. Return Tokens to Client
```

### **Authorization Model**
```go
// internal/api/middleware/auth.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/fiscaflow/internal/domain/user"
)

func JWTAuthMiddleware(userService user.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Remove "Bearer " prefix
        token = strings.TrimPrefix(token, "Bearer ")
        
        // Validate token
        claims, err := validateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Get user from database
        user, err := userService.GetByID(c.Request.Context(), claims.UserID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            c.Abort()
            return
        }
        
        // Set user context
        c.Set("user", user)
        c.Set("user_id", user.ID)
        c.Set("user_role", user.Role)
        
        c.Next()
    }
}
```

### **Data Encryption**
- **At Rest**: AES-256 encryption for sensitive data
- **In Transit**: TLS 1.3 for all communications
- **PII Masking**: Automatic masking in logs and exports

## 🚀 **Deployment Architecture**

### **Container Orchestration (Kubernetes)**
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fiscaflow-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: fiscaflow-api
  template:
    metadata:
      labels:
        app: fiscaflow-api
    spec:
      containers:
      - name: api
        image: fiscaflow/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: fiscaflow-secrets
              key: database-url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: fiscaflow-secrets
              key: redis-url
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### **Auto-scaling**
```yaml
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: fiscaflow-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: fiscaflow-api
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## 📈 **Performance & Scalability**

### **Performance Targets**
- **API Response Time**: < 200ms (95th percentile)
- **Database Queries**: < 50ms (95th percentile)
- **Cache Hit Rate**: > 95%
- **Uptime**: 99.9%

### **Scaling Strategies**
- **Horizontal Scaling**: Auto-scaling based on CPU/memory
- **Database Scaling**: Read replicas for analytics queries
- **Cache Scaling**: Redis cluster for high availability
- **CDN**: Static assets and API responses

### **Load Balancing**
- **Application Level**: Kubernetes service load balancing
- **Network Level**: NGINX load balancer
- **Database Level**: Connection pooling with PgBouncer

## 🔍 **Monitoring & Observability**

### **Application Monitoring**
- **APM**: New Relic / DataDog integration
- **Health Checks**: Kubernetes liveness/readiness probes
- **Custom Metrics**: Business metrics via OpenTelemetry

### **Infrastructure Monitoring**
- **Kubernetes**: Prometheus + Grafana
- **Database**: PostgreSQL monitoring
- **Cache**: Redis monitoring
- **Message Queue**: RabbitMQ monitoring

### **Logging Strategy**
- **Centralized Logging**: ELK stack (Elasticsearch, Logstash, Kibana)
- **Structured Logging**: JSON format with correlation IDs
- **Log Levels**: ERROR, WARN, INFO, DEBUG

## 🔄 **CI/CD Pipeline**

### **Build Pipeline**
```
1. Code Commit
   ↓
2. Automated Testing
   - Unit Tests (go test)
   - Integration Tests
   - Security Scans
   ↓
3. Build Docker Images
   ↓
4. Push to Registry
   ↓
5. Deploy to Staging
   ↓
6. Run E2E Tests
   ↓
7. Deploy to Production
```

### **Go-specific CI/CD**
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run tests
      run: |
        go test -v ./...
        go test -race ./...
        go test -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
    
    - name: Build and push
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: fiscaflow/api:latest
```

### **Deployment Strategy**
- **Blue-Green Deployment**: Zero-downtime deployments
- **Rolling Updates**: Gradual service updates
- **Canary Deployments**: Risk mitigation for major changes

## 🛡️ **Disaster Recovery**

### **Backup Strategy**
- **Database**: Daily automated backups with point-in-time recovery
- **File Storage**: Cross-region replication
- **Configuration**: Version-controlled infrastructure as code

### **Recovery Procedures**
- **RTO (Recovery Time Objective)**: 4 hours
- **RPO (Recovery Point Objective)**: 1 hour
- **Failover**: Automated failover to secondary region

## 🐳 **Docker Configuration**

### **Multi-stage Build**
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
```

### **Docker Compose for Development**
```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:password@postgres:5432/fiscaflow
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=your-secret-key
      - ELASTICSEARCH_URL=http://elasticsearch:9200
      - MINIO_ENDPOINT=minio:9000
      - RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
    depends_on:
      - postgres
      - redis
      - elasticsearch
      - minio
      - rabbitmq

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=fiscaflow
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  elasticsearch:
    image: elasticsearch:8.8.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin
    volumes:
      - minio_data:/data

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"

volumes:
  postgres_data:
  elasticsearch_data:
  minio_data:
```

## 🔄 **Migration Path to Microservices**

### **Phase 1: Domain Separation (Current)**
- Clean domain boundaries
- Event-driven communication
- Shared database with domain-specific schemas

### **Phase 2: Database Separation**
- Separate databases per domain
- Data synchronization via events
- API composition layer

### **Phase 3: Service Extraction**
- Extract domains into separate services
- Implement service mesh
- Maintain API compatibility

### **Phase 4: Full Microservices**
- Independent deployment
- Service-specific scaling
- Advanced monitoring and tracing

---

*This monolith architecture provides a robust, scalable, and maintainable foundation for FiscaFlow's personal finance management platform, with a clear migration path to microservices when needed.* 