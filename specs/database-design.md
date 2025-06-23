# Database Design Specification

## 📊 **Overview**

FiscaFlow uses a multi-database architecture with PostgreSQL as the primary transactional database, Redis for caching and sessions, Elasticsearch for search and analytics, and MinIO for file storage.

## 🎯 **Database Architecture**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Application Layer                               │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                    ┌───────────────────┼───────────────────┐
                    ▼                   ▼                   ▼
┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐
│     PostgreSQL          │ │        Redis            │ │     Elasticsearch       │
│   (Primary Database)    │ │     (Cache/Sessions)    │ │    (Search/Analytics)   │
│  - ACID Transactions    │ │  - Session Storage      │ │  - Full-text Search     │
│  - User Data            │ │  - API Response Cache   │ │  - Log Aggregation      │
│  - Financial Data       │ │  - Rate Limiting        │ │  - Business Metrics     │
│  - Audit Trail          │ │  - Job Queues           │ │  - Real-time Analytics  │
└─────────────────────────┘ └─────────────────────────┘ └─────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              MinIO (File Storage)                           │
│  - Document Storage  - File Uploads  - Exports  - Backups                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 🗄️ **PostgreSQL Schema Design**

### **Core Tables**

#### **1. Users Table**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    date_of_birth DATE,
    timezone VARCHAR(50) DEFAULT 'UTC',
    locale VARCHAR(10) DEFAULT 'en-US',
    role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('user', 'premium', 'admin', 'family_owner', 'family_member')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended', 'deleted')),
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT users_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);
```

#### **2. User Sessions Table**
```sql
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    access_token_hash VARCHAR(255),
    device_info JSONB,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT user_sessions_expires_at_future CHECK (expires_at > created_at)
);

-- Indexes
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
```

#### **3. Family Groups Table**
```sql
CREATE TABLE family_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id),
    description TEXT,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_family_groups_owner_id ON family_groups(owner_id);
```

#### **4. Family Members Table**
```sql
CREATE TABLE family_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID NOT NULL REFERENCES family_groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    permissions JSONB DEFAULT '{}',
    joined_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(family_id, user_id)
);

-- Indexes
CREATE INDEX idx_family_members_family_id ON family_members(family_id);
CREATE INDEX idx_family_members_user_id ON family_members(user_id);
```

#### **5. Accounts Table**
```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_id UUID REFERENCES family_groups(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('checking', 'savings', 'credit_card', 'investment', 'loan', 'other')),
    institution VARCHAR(100),
    account_number_hash VARCHAR(255),
    balance DECIMAL(15,2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'USD',
    is_active BOOLEAN DEFAULT TRUE,
    plaid_account_id VARCHAR(255),
    last_sync_at TIMESTAMP,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT accounts_balance_non_negative CHECK (balance >= 0)
);

-- Indexes
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_family_id ON accounts(family_id);
CREATE INDEX idx_accounts_type ON accounts(type);
CREATE INDEX idx_accounts_plaid_account_id ON accounts(plaid_account_id);
```

#### **6. Categories Table**
```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(7),
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT categories_color_format CHECK (color ~* '^#[0-9A-Fa-f]{6}$')
);

-- Indexes
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_is_default ON categories(is_default);
CREATE INDEX idx_categories_sort_order ON categories(sort_order);
```

#### **7. Transactions Table**
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_id UUID REFERENCES family_groups(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    description TEXT NOT NULL,
    merchant VARCHAR(255),
    location JSONB,
    
    transaction_date DATE NOT NULL,
    posted_date DATE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'posted', 'cancelled', 'disputed')),
    
    categorization_source VARCHAR(20) DEFAULT 'manual' CHECK (categorization_source IN ('manual', 'ml', 'plaid', 'user_correction')),
    categorization_confidence DECIMAL(3,2) CHECK (categorization_confidence >= 0 AND categorization_confidence <= 1),
    
    tags TEXT[],
    notes TEXT,
    receipt_url VARCHAR(500),
    
    plaid_transaction_id VARCHAR(255),
    external_id VARCHAR(255),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT transactions_amount_non_zero CHECK (amount != 0)
);

-- Indexes
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_family_id ON transactions(family_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_category_id ON transactions(category_id);
CREATE INDEX idx_transactions_transaction_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_amount ON transactions(amount);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_plaid_transaction_id ON transactions(plaid_transaction_id);
CREATE INDEX idx_transactions_external_id ON transactions(external_id);
```

#### **8. Budgets Table**
```sql
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_id UUID REFERENCES family_groups(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    period_type VARCHAR(20) NOT NULL CHECK (period_type IN ('monthly', 'quarterly', 'yearly', 'custom')),
    start_date DATE NOT NULL,
    end_date DATE,
    
    total_amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    
    is_active BOOLEAN DEFAULT TRUE,
    settings JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT budgets_total_amount_positive CHECK (total_amount > 0),
    CONSTRAINT budgets_end_date_after_start CHECK (end_date IS NULL OR end_date > start_date)
);

-- Indexes
CREATE INDEX idx_budgets_user_id ON budgets(user_id);
CREATE INDEX idx_budgets_family_id ON budgets(family_id);
CREATE INDEX idx_budgets_period_type ON budgets(period_type);
CREATE INDEX idx_budgets_start_date ON budgets(start_date);
```

#### **9. Budget Categories Table**
```sql
CREATE TABLE budget_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    
    allocated_amount DECIMAL(15,2) NOT NULL,
    spent_amount DECIMAL(15,2) DEFAULT 0.00,
    
    alert_threshold DECIMAL(3,2) DEFAULT 0.80 CHECK (alert_threshold >= 0 AND alert_threshold <= 1),
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(budget_id, category_id),
    CONSTRAINT budget_categories_allocated_positive CHECK (allocated_amount > 0)
);

-- Indexes
CREATE INDEX idx_budget_categories_budget_id ON budget_categories(budget_id);
CREATE INDEX idx_budget_categories_category_id ON budget_categories(category_id);
```

#### **10. Goals Table**
```sql
CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_id UUID REFERENCES family_groups(id) ON DELETE CASCADE,
    
    name VARCHAR(100) NOT NULL,
    description TEXT,
    goal_type VARCHAR(50) NOT NULL CHECK (goal_type IN ('savings', 'debt_payoff', 'investment', 'emergency_fund', 'vacation', 'purchase', 'other')),
    
    target_amount DECIMAL(15,2) NOT NULL,
    current_amount DECIMAL(15,2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'USD',
    
    target_date DATE,
    start_date DATE DEFAULT CURRENT_DATE,
    
    icon VARCHAR(50),
    color VARCHAR(7),
    
    is_active BOOLEAN DEFAULT TRUE,
    settings JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT goals_target_amount_positive CHECK (target_amount > 0),
    CONSTRAINT goals_current_amount_non_negative CHECK (current_amount >= 0),
    CONSTRAINT goals_target_date_future CHECK (target_date IS NULL OR target_date >= CURRENT_DATE)
);

-- Indexes
CREATE INDEX idx_goals_user_id ON goals(user_id);
CREATE INDEX idx_goals_family_id ON goals(family_id);
CREATE INDEX idx_goals_goal_type ON goals(goal_type);
CREATE INDEX idx_goals_target_date ON goals(target_date);
```

#### **11. Notifications Table**
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    type VARCHAR(50) NOT NULL CHECK (type IN ('budget_alert', 'goal_milestone', 'unusual_spending', 'bill_reminder', 'system', 'security')),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    
    data JSONB DEFAULT '{}',
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    
    read_at TIMESTAMP,
    sent_at TIMESTAMP,
    expires_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT notifications_expires_at_future CHECK (expires_at IS NULL OR expires_at > created_at)
);

-- Indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_read_at ON notifications(read_at);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
```

### **Audit and Logging Tables**

#### **12. Audit Logs Table**
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id UUID REFERENCES user_sessions(id) ON DELETE SET NULL,
    
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    resource_data JSONB,
    
    ip_address INET,
    user_agent TEXT,
    
    success BOOLEAN NOT NULL,
    error_message TEXT,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

#### **13. Data Sync Logs Table**
```sql
CREATE TABLE data_sync_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    
    sync_type VARCHAR(50) NOT NULL CHECK (sync_type IN ('plaid', 'manual', 'csv_import')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'in_progress', 'completed', 'failed')),
    
    transactions_added INTEGER DEFAULT 0,
    transactions_updated INTEGER DEFAULT 0,
    transactions_deleted INTEGER DEFAULT 0,
    
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    error_message TEXT,
    
    metadata JSONB DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_data_sync_logs_user_id ON data_sync_logs(user_id);
CREATE INDEX idx_data_sync_logs_account_id ON data_sync_logs(account_id);
CREATE INDEX idx_data_sync_logs_status ON data_sync_logs(status);
CREATE INDEX idx_data_sync_logs_started_at ON data_sync_logs(started_at);
```

## 🔄 **Database Triggers and Functions**

### **Updated At Trigger**
```sql
-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to all tables with updated_at column
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ... (apply to all relevant tables)
```

### **Transaction Balance Update Trigger**
```sql
-- Function to update account balance when transactions are added/modified
CREATE OR REPLACE FUNCTION update_account_balance()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE accounts 
        SET balance = balance + NEW.amount,
            updated_at = NOW()
        WHERE id = NEW.account_id;
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE accounts 
        SET balance = balance - OLD.amount + NEW.amount,
            updated_at = NOW()
        WHERE id = NEW.account_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE accounts 
        SET balance = balance - OLD.amount,
            updated_at = NOW()
        WHERE id = OLD.account_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_account_balance_trigger
    AFTER INSERT OR UPDATE OR DELETE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_account_balance();
```

## 📊 **Redis Data Structures**

### **Session Storage**
```redis
# User sessions
SET session:{session_id} "{user_id, expires_at, device_info}"
EXPIRE session:{session_id} {ttl_seconds}

# User session mapping
SADD user_sessions:{user_id} {session_id}
```

### **API Response Cache**
```redis
# Cached API responses
SET cache:api:{endpoint}:{user_id}:{params_hash} "{response_data}"
EXPIRE cache:api:{endpoint}:{user_id}:{params_hash} 300

# Cache invalidation patterns
DEL cache:api:transactions:{user_id}:*
DEL cache:api:analytics:{user_id}:*
```

### **Rate Limiting**
```redis
# Rate limiting counters
INCR rate_limit:{user_id}:{endpoint}:{window}
EXPIRE rate_limit:{user_id}:{endpoint}:{window} {window_seconds}

# Check rate limit
GET rate_limit:{user_id}:{endpoint}:{window}
```

### **Real-time Data**
```redis
# Live dashboard data
SET dashboard:{user_id}:spending_today "{amount, categories}"
SET dashboard:{user_id}:budget_status "{alerts, progress}"

# WebSocket connections
SADD ws_connections:{user_id} {connection_id}
```

## 🔍 **Elasticsearch Indices**

### **Transactions Index**
```json
{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "user_id": { "type": "keyword" },
      "family_id": { "type": "keyword" },
      "account_id": { "type": "keyword" },
      "category_id": { "type": "keyword" },
      "amount": { "type": "float" },
      "currency": { "type": "keyword" },
      "description": { 
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "merchant": { 
        "type": "text",
        "analyzer": "standard"
      },
      "transaction_date": { "type": "date" },
      "status": { "type": "keyword" },
      "tags": { "type": "keyword" },
      "location": {
        "type": "geo_point"
      },
      "categorization_confidence": { "type": "float" }
    }
  }
}
```

### **Categories Index**
```json
{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "name": { 
        "type": "text",
        "analyzer": "standard"
      },
      "description": { "type": "text" },
      "icon": { "type": "keyword" },
      "color": { "type": "keyword" },
      "parent_id": { "type": "keyword" },
      "is_default": { "type": "boolean" },
      "is_active": { "type": "boolean" },
      "sort_order": { "type": "integer" }
    }
  }
}
```

### **Logs Index**
```json
{
  "mappings": {
    "properties": {
      "timestamp": { "type": "date" },
      "level": { "type": "keyword" },
      "service": { "type": "keyword" },
      "trace_id": { "type": "keyword" },
      "span_id": { "type": "keyword" },
      "message": { "type": "text" },
      "user_id": { "type": "keyword" },
      "attributes": { "type": "object" },
      "error": {
        "properties": {
          "type": { "type": "keyword" },
          "message": { "type": "text" },
          "stack_trace": { "type": "text" }
        }
      }
    }
  }
}
```

## 📁 **MinIO Bucket Structure**

### **File Organization**
```
fiscaflow-storage/
├── exports/
│   ├── {user_id}/
│   │   ├── transactions_{date}.csv
│   │   ├── budget_report_{date}.pdf
│   │   └── analytics_{date}.xlsx
├── imports/
│   ├── {user_id}/
│   │   ├── bank_statement_{date}.csv
│   │   └── credit_card_{date}.csv
├── attachments/
│   ├── {user_id}/
│   │   ├── receipts/
│   │   │   └── {transaction_id}.jpg
│   │   └── documents/
│   │       └── {document_id}.pdf
└── backups/
    ├── database/
    │   └── {date}_{time}.sql
    └── files/
        └── {date}_{time}.tar.gz
```

## 🔒 **Data Security and Privacy**

### **Encryption**
- **At Rest**: AES-256 encryption for sensitive fields
- **In Transit**: TLS 1.3 for all database connections
- **Field-level Encryption**: Sensitive data encrypted before storage

### **Data Masking**
```sql
-- Function to mask sensitive data in logs
CREATE OR REPLACE FUNCTION mask_sensitive_data(data JSONB)
RETURNS JSONB AS $$
BEGIN
    -- Mask email addresses
    data = jsonb_set(data, '{email}', '"***@***.***"');
    
    -- Mask account numbers
    data = jsonb_set(data, '{account_number}', '"****-****-****-****"');
    
    -- Mask phone numbers
    data = jsonb_set(data, '{phone}', '"***-***-****"');
    
    RETURN data;
END;
$$ language 'plpgsql';
```

### **Data Retention**
```sql
-- Cleanup old audit logs (retain for 7 years)
DELETE FROM audit_logs 
WHERE created_at < NOW() - INTERVAL '7 years';

-- Cleanup old sessions (retain for 30 days)
DELETE FROM user_sessions 
WHERE expires_at < NOW() - INTERVAL '30 days';

-- Archive old transactions (move to archive table after 5 years)
INSERT INTO transactions_archive 
SELECT * FROM transactions 
WHERE transaction_date < NOW() - INTERVAL '5 years';
```

## 📈 **Performance Optimization**

### **Indexing Strategy**
- **Primary Keys**: UUID with gen_random_uuid()
- **Foreign Keys**: Indexed for join performance
- **Date Ranges**: B-tree indexes on date columns
- **Full-text Search**: GIN indexes on text columns
- **Composite Indexes**: For common query patterns

### **Partitioning**
```sql
-- Partition transactions by date
CREATE TABLE transactions_y2024m01 PARTITION OF transactions
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE transactions_y2024m02 PARTITION OF transactions
FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
```

### **Connection Pooling**
```yaml
# PgBouncer configuration
[databases]
fiscaflow = host=postgres port=5432 dbname=fiscaflow

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
```

---

*This database design provides a robust, scalable, and secure foundation for FiscaFlow's financial data management requirements.* 