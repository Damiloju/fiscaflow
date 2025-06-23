# OpenTelemetry Integration Specification

## 📊 **Overview**

FiscaFlow implements comprehensive observability using OpenTelemetry (OTel) to provide distributed tracing, metrics collection, and structured logging across all system components, built with Go for high performance and reliability.

## 🎯 **Objectives**

- **Distributed Tracing**: Track request flows across microservices
- **Performance Monitoring**: Identify bottlenecks and optimize performance
- **Error Tracking**: Rapidly identify and resolve issues
- **Business Metrics**: Track financial and user engagement metrics
- **Compliance**: Audit trails for financial data operations

## 🏗️ **Architecture**

### **OTel Components**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Application   │    │   OTel SDK      │    │   OTel Collector│
│   (Go + Gin)    │───▶│   (Instrumentation)│───▶│   (OTel Agent)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Jaeger        │    │   Prometheus    │    │   Elasticsearch │
│   (Tracing)     │◀───│   (Metrics)     │◀───│   (Logging)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🔍 **Distributed Tracing**

### **Trace Context Propagation**
- **Trace ID**: Unique identifier for each request
- **Span ID**: Individual operation within a trace
- **Parent Span ID**: Links child operations to parent
- **Baggage**: Custom key-value pairs for context

### **Instrumented Operations**

#### **HTTP Requests**
```go
// Example: Transaction creation trace
type TransactionTrace struct {
    TraceID    string                 `json:"trace_id"`
    SpanID     string                 `json:"span_id"`
    Operation  string                 `json:"operation_name"`
    Attributes map[string]interface{} `json:"attributes"`
    Events     []TraceEvent           `json:"events"`
}

type TraceEvent struct {
    Name       string                 `json:"name"`
    Timestamp  time.Time              `json:"timestamp"`
    Attributes map[string]interface{} `json:"attributes"`
}

// Example trace data
{
    "trace_id": "1-632c7c4b-5c3d2e1f-8a9b0c1d2e3f4a5b",
    "span_id": "2a3b4c5d6e7f8a9b",
    "operation_name": "POST /api/v1/transactions",
    "attributes": {
        "http.method": "POST",
        "http.url": "/api/v1/transactions",
        "http.status_code": 201,
        "user.id": "user_123",
        "transaction.amount": 150.00,
        "transaction.category": "groceries"
    },
    "events": [
        {
            "name": "transaction.created",
            "timestamp": "2024-01-15T10:30:00Z",
            "attributes": {
                "transaction.id": "txn_456",
                "processing_time_ms": 45
            }
        }
    ]
}
```

#### **Database Operations**
```go
// Example: Database query trace
{
    "trace_id": "1-632c7c4b-5c3d2e1f-8a9b0c1d2e3f4a5b",
    "span_id": "3b4c5d6e7f8a9b0c",
    "parent_span_id": "2a3b4c5d6e7f8a9b",
    "operation_name": "SELECT transactions",
    "attributes": {
        "db.system": "postgresql",
        "db.name": "fiscaflow",
        "db.statement": "SELECT * FROM transactions WHERE user_id = $1",
        "db.operation": "SELECT",
        "db.rows_affected": 25
    }
}
```

#### **External API Calls**
```go
// Example: Plaid API integration trace
{
    "trace_id": "1-632c7c4b-5c3d2e1f-8a9b0c1d2e3f4a5b",
    "span_id": "4c5d6e7f8a9b0c1d",
    "parent_span_id": "2a3b4c5d6e7f8a9b",
    "operation_name": "GET /plaid/transactions",
    "attributes": {
        "http.method": "GET",
        "http.url": "https://api.plaid.com/transactions/get",
        "http.status_code": 200,
        "plaid.account_id": "acc_789",
        "plaid.transactions_count": 15
    }
}
```

### **Business Operation Tracing**

#### **Transaction Processing Flow**
1. **User Authentication** → **Account Validation** → **Transaction Creation** → **Categorization** → **Budget Update** → **Notification**

#### **Budget Calculation Flow**
1. **Data Aggregation** → **Category Grouping** → **Budget Comparison** → **Alert Generation** → **Dashboard Update**

## 📈 **Metrics Collection**

### **Technical Metrics**

#### **Application Metrics**
```go
// HTTP Metrics
http_requests_total{method="POST", endpoint="/api/v1/transactions", status="201"}
http_request_duration_seconds{method="POST", endpoint="/api/v1/transactions", quantile="0.95"}

// Database Metrics
db_connections_active{database="fiscaflow"}
db_query_duration_seconds{operation="SELECT", table="transactions", quantile="0.95"}

// Memory and CPU
process_memory_usage_bytes
process_cpu_usage_percent
```

#### **Business Metrics**
```go
// Financial Metrics
transactions_created_total{category="groceries", user_type="premium"}
budget_alerts_triggered_total{severity="warning", category="entertainment"}
goals_completed_total{goal_type="savings"}

// User Engagement Metrics
daily_active_users
feature_usage_total{feature="categorization", action="auto_categorize"}
api_calls_total{endpoint="analytics", user_tier="premium"}
```

### **Custom Metrics**

#### **Categorization Accuracy**
```go
categorization_accuracy_percent{category="groceries", model_version="v2.1"}
categorization_confidence_score{category="groceries", confidence_level="high"}
manual_categorization_corrections_total{original_category="groceries", corrected_category="dining"}
```

#### **Financial Health Indicators**
```go
average_monthly_spending{category="groceries", user_segment="millennial"}
budget_compliance_percent{category="entertainment", compliance_threshold="90"}
savings_rate_percent{user_tier="premium"}
```

## 📝 **Structured Logging**

### **Log Format**
```go
// Log structure in Go
type LogEntry struct {
    Timestamp   time.Time              `json:"timestamp"`
    Level       string                 `json:"level"`
    TraceID     string                 `json:"trace_id"`
    SpanID      string                 `json:"span_id"`
    Service     string                 `json:"service"`
    Version     string                 `json:"version"`
    Environment string                 `json:"environment"`
    Message     string                 `json:"message"`
    Attributes  map[string]interface{} `json:"attributes"`
    Context     LogContext             `json:"context"`
}

type LogContext struct {
    RequestID  string `json:"request_id"`
    IPAddress  string `json:"ip_address"`
    UserAgent  string `json:"user_agent"`
}

// Example log output
{
    "timestamp": "2024-01-15T10:30:00.123Z",
    "level": "INFO",
    "trace_id": "1-632c7c4b-5c3d2e1f-8a9b0c1d2e3f4a5b",
    "span_id": "2a3b4c5d6e7f8a9b",
    "service": "fiscaflow-api",
    "version": "1.0.0",
    "environment": "production",
    "message": "Transaction created successfully",
    "attributes": {
        "user_id": "user_123",
        "transaction_id": "txn_456",
        "amount": 150.00,
        "category": "groceries",
        "processing_time_ms": 45
    },
    "context": {
        "request_id": "req_789",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0..."
    }
}
```

### **Log Levels and Usage**

#### **ERROR** - System failures and exceptions
```go
// Error logging in Go
logger.Error("Failed to process transaction",
    "error", err,
    "transaction_id", transactionID,
    "retry_count", retryCount,
)
```

#### **WARN** - Potential issues and degraded performance
```go
// Warning logging in Go
logger.Warn("High response time detected",
    "endpoint", "/api/v1/transactions",
    "response_time_ms", 850,
    "threshold_ms", 500,
)
```

#### **INFO** - Business operations and state changes
```go
// Info logging in Go
logger.Info("User upgraded to premium plan",
    "user_id", userID,
    "plan", "premium",
    "payment_method", "credit_card",
)
```

#### **DEBUG** - Detailed debugging information
```go
// Debug logging in Go
logger.Debug("Processing categorization request",
    "transaction_description", "WALMART GROCERY",
    "ml_model_version", "v2.1",
    "confidence_score", 0.95,
)
```

## 🔧 **Implementation Details**

### **Go SDK Configuration**
```go
// otel/config.go
package otel

import (
    "context"
    "fmt"
    "os"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
    "go.opentelemetry.io/otel/trace"
)

func InitTracer(serviceName, serviceVersion string) (trace.TracerProvider, error) {
    ctx := context.Background()
    
    // Create OTLP exporter
    exporter, err := otlptracehttp.New(ctx,
        otlptracehttp.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
        otlptracehttp.WithInsecure(),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
    }
    
    // Create resource
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(serviceName),
            semconv.ServiceVersion(serviceVersion),
            semconv.DeploymentEnvironment(os.Getenv("ENVIRONMENT")),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }
    
    // Create trace provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )
    
    // Set global trace provider
    otel.SetTracerProvider(tp)
    
    return tp, nil
}

// main.go
func main() {
    // Initialize OpenTelemetry
    tp, err := otel.InitTracer("fiscaflow-user-service", "1.0.0")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            log.Printf("Error shutting down tracer provider: %v", err)
        }
    }()
    
    // Initialize metrics
    mp, err := otel.InitMeter("fiscaflow-user-service")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := mp.Shutdown(context.Background()); err != nil {
            log.Printf("Error shutting down meter provider: %v", err)
        }
    }()
    
    // Start server
    server := gin.Default()
    server.Use(middleware.OpenTelemetry())
    // ... rest of server setup
}
```

### **Custom Instrumentation**
```go
// tracing/transaction.go
package tracing

import (
    "context"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

type TransactionTracer struct {
    tracer trace.Tracer
}

func NewTransactionTracer() *TransactionTracer {
    return &TransactionTracer{
        tracer: otel.Tracer("fiscaflow-transactions"),
    }
}

func (t *TransactionTracer) TraceTransactionCreation(ctx context.Context, userID string, transactionData map[string]interface{}) (interface{}, error) {
    ctx, span := t.tracer.Start(ctx, "transaction.creation")
    defer span.End()
    
    startTime := time.Now()
    
    // Add business attributes
    span.SetAttributes(
        attribute.String("user.id", userID),
        attribute.Float64("transaction.amount", transactionData["amount"].(float64)),
        attribute.String("transaction.category", transactionData["category"].(string)),
    )
    
    // Create transaction
    transaction, err := t.createTransaction(ctx, transactionData)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    // Add success event
    span.AddEvent("transaction.created",
        trace.WithAttributes(
            attribute.String("transaction.id", transaction.ID),
            attribute.Int64("processing_time_ms", time.Since(startTime).Milliseconds()),
        ),
    )
    
    return transaction, nil
}

// middleware/opentelemetry.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func OpenTelemetry() gin.HandlerFunc {
    return func(c *gin.Context) {
        tracer := otel.Tracer("fiscaflow-http")
        
        ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path,
            trace.WithAttributes(
                attribute.String("http.method", c.Request.Method),
                attribute.String("http.url", c.Request.URL.String()),
                attribute.String("http.user_agent", c.Request.UserAgent()),
                attribute.String("http.remote_addr", c.ClientIP()),
            ),
        )
        defer span.End()
        
        // Add trace context to request
        c.Request = c.Request.WithContext(ctx)
        
        // Process request
        c.Next()
        
        // Add response attributes
        span.SetAttributes(
            attribute.Int("http.status_code", c.Writer.Status()),
            attribute.Int("http.response_size", c.Writer.Size()),
        )
        
        // Set span status based on HTTP status code
        if c.Writer.Status() >= 400 {
            span.SetStatus(codes.Error, http.StatusText(c.Writer.Status()))
        }
    }
}
```

### **Custom Metrics**
```go
// metrics/business.go
package metrics

import (
    "context"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/metric"
)

type BusinessMetrics struct {
    transactionsCreated    metric.Int64Counter
    categorizationAccuracy metric.Float64Histogram
    meter                  metric.Meter
}

func NewBusinessMetrics() (*BusinessMetrics, error) {
    meter := otel.Meter("fiscaflow-business")
    
    transactionsCreated, err := meter.Int64Counter("transactions_created_total",
        metric.WithDescription("Total number of transactions created"),
    )
    if err != nil {
        return nil, err
    }
    
    categorizationAccuracy, err := meter.Float64Histogram("categorization_accuracy_percent",
        metric.WithDescription("Accuracy of transaction categorization"),
        metric.WithUnit("percent"),
    )
    if err != nil {
        return nil, err
    }
    
    return &BusinessMetrics{
        transactionsCreated:    transactionsCreated,
        categorizationAccuracy: categorizationAccuracy,
        meter:                  meter,
    }, nil
}

func (m *BusinessMetrics) RecordTransactionCreated(ctx context.Context, category, userType string) {
    m.transactionsCreated.Add(ctx, 1,
        metric.WithAttributes(
            attribute.String("category", category),
            attribute.String("user_type", userType),
        ),
    )
}

func (m *BusinessMetrics) RecordCategorizationAccuracy(ctx context.Context, category string, accuracy float64, modelVersion string) {
    m.categorizationAccuracy.Record(ctx, accuracy,
        metric.WithAttributes(
            attribute.String("category", category),
            attribute.String("model_version", modelVersion),
        ),
    )
}

// Usage in handlers
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    ctx := c.Request.Context()
    
    // ... transaction creation logic
    
    // Record metrics
    h.metrics.RecordTransactionCreated(ctx, transaction.Category.Name, user.Role)
    
    c.JSON(http.StatusCreated, transaction)
}
```

### **Structured Logging with Zap**
```go
// logging/logger.go
package logging

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "go.opentelemetry.io/otel/trace"
)

func NewLogger(serviceName, version, environment string) (*zap.Logger, error) {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    config.EncoderConfig.MessageKey = "message"
    config.EncoderConfig.LevelKey = "level"
    config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
    
    logger, err := config.Build()
    if err != nil {
        return nil, err
    }
    
    return logger.With(
        zap.String("service", serviceName),
        zap.String("version", version),
        zap.String("environment", environment),
    ), nil
}

// Middleware to inject trace context into logs
func TraceLogger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := c.Request.Context()
        
        // Get trace context
        span := trace.SpanFromContext(ctx)
        traceID := span.SpanContext().TraceID().String()
        spanID := span.SpanContext().SpanID().String()
        
        // Create logger with trace context
        traceLogger := logger.With(
            zap.String("trace_id", traceID),
            zap.String("span_id", spanID),
        )
        
        // Add to context
        c.Set("logger", traceLogger)
        
        c.Next()
    }
}

// Usage in handlers
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    logger := c.MustGet("logger").(*zap.Logger)
    
    logger.Info("Creating transaction",
        zap.String("user_id", userID),
        zap.Float64("amount", req.Amount),
        zap.String("category", req.Category),
    )
    
    // ... transaction creation logic
    
    logger.Info("Transaction created successfully",
        zap.String("transaction_id", transaction.ID.String()),
        zap.Int64("processing_time_ms", time.Since(startTime).Milliseconds()),
    )
}
```

## 📊 **Observability Backends**

### **Jaeger (Tracing)**
- **URL**: `http://jaeger:16686`
- **Storage**: Elasticsearch
- **Retention**: 30 days
- **Sampling**: 100% for errors, 10% for successful requests

### **Prometheus (Metrics)**
- **URL**: `http://prometheus:9090`
- **Scrape Interval**: 15 seconds
- **Retention**: 15 days
- **Alerting**: Grafana AlertManager

### **Elasticsearch (Logging)**
- **URL**: `http://elasticsearch:9200`
- **Index Pattern**: `fiscaflow-logs-*`
- **Retention**: 90 days
- **Search**: Kibana interface

## 🚨 **Alerting Rules**

### **Performance Alerts**
```yaml
# High Response Time
- alert: HighAPIResponseTime
  expr: histogram_quantile(0.95, http_request_duration_seconds) > 0.5
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "High API response time detected"
    description: "95th percentile response time is {{ $value }}s"

# High Error Rate
- alert: HighErrorRate
  expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
  for: 2m
  labels:
    severity: critical
  annotations:
    summary: "High error rate detected"
    description: "Error rate is {{ $value | humanizePercentage }}"
```

### **Business Alerts**
```yaml
# Categorization Accuracy Drop
- alert: LowCategorizationAccuracy
  expr: avg_over_time(categorization_accuracy_percent[1h]) < 85
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "Low categorization accuracy"
    description: "Average accuracy is {{ $value }}%"

# System Uptime
- alert: SystemDown
  expr: up == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "Service is down"
    description: "Service {{ $labels.instance }} is down"
```

## 🔍 **Troubleshooting Guide**

### **Common Issues**

#### **Missing Traces**
1. Check OTel collector connectivity
2. Verify sampling configuration
3. Check trace context propagation

#### **High Memory Usage**
1. Monitor span creation rate
2. Check for memory leaks in custom instrumentation
3. Adjust buffer sizes and batch processing

#### **Slow Query Performance**
1. Use database query traces to identify slow queries
2. Check database connection pool metrics
3. Monitor query execution plans

### **Debugging Commands**
```bash
# Check OTel collector status
curl http://otel-collector:13133/health

# View recent traces
curl "http://jaeger:16686/api/traces?service=fiscaflow-api&limit=10"

# Check metrics
curl "http://prometheus:9090/api/v1/query?query=up"

# Search logs
curl -X GET "http://elasticsearch:9200/fiscaflow-logs-*/_search" \
  -H "Content-Type: application/json" \
  -d '{"query":{"match":{"message":"error"}}}'
```

## 📋 **Compliance & Security**

### **Data Privacy**
- **PII Masking**: Automatically mask sensitive data in logs
- **Data Retention**: Configurable retention policies
- **Access Control**: Role-based access to observability data

### **Audit Trail**
- **Financial Operations**: Complete audit trail for all financial transactions
- **User Actions**: Track all user actions for compliance
- **System Changes**: Log all configuration and deployment changes

## 🐳 **Docker Configuration**

### **OTel Collector Configuration**
```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
  memory_limiter:
    check_interval: 1s
    limit_mib: 1500

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true
  prometheus:
    endpoint: 0.0.0.0:9464
  logging:
    loglevel: debug

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch, memory_limiter]
      exporters: [jaeger, logging]
    metrics:
      receivers: [otlp]
      processors: [batch, memory_limiter]
      exporters: [prometheus, logging]
    logs:
      receivers: [otlp]
      processors: [batch, memory_limiter]
      exporters: [logging]
```

### **Docker Compose with Observability**
```yaml
# docker-compose.observability.yml
version: '3.8'

services:
  otel-collector:
    image: otel/opentelemetry-collector:latest
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP
      - "9464:9464"   # Prometheus metrics
    depends_on:
      - jaeger
      - prometheus

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # UI
      - "14250:14250"  # gRPC
    environment:
      - COLLECTOR_OTLP_ENABLED=true

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  grafana_data:
```

---

*This specification ensures comprehensive observability for FiscaFlow, enabling rapid troubleshooting and performance optimization with Go-based services.* 