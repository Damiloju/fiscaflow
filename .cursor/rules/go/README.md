# Go-Specific Cursor Rules

This directory contains comprehensive Cursor rules specifically designed for Go development best practices. These rules are created by Go experts and cover all aspects of modern Go development.

## 📋 **Rules Overview**

### **Code Organization & Structure**
- **`go-code-organization.mdc`**: Package structure, naming conventions, file organization, and import management

### **Error Handling & Resilience**
- **`go-error-handling.mdc`**: Error handling patterns, custom errors, error wrapping, and error documentation

### **Concurrency & Parallelism**
- **`go-concurrency.mdc`**: Goroutines, channels, context usage, synchronization, and concurrency patterns

### **Performance & Optimization**
- **`go-performance.mdc`**: Memory management, profiling, benchmarking, and performance optimization techniques

### **Testing & Quality Assurance**
- **`go-testing.mdc`**: Unit testing, integration testing, test organization, mocking, and test coverage

### **API Design & HTTP**
- **`go-api-design.mdc`**: RESTful API design, HTTP handlers, middleware, authentication, and API patterns

### **Database & Data Management**
- **`go-database.mdc`**: Database operations, migrations, connection management, transactions, and query optimization

## 🎯 **Usage**

These rules automatically activate when working with `.go` files and provide:

- **Real-time suggestions** for Go best practices
- **Code examples** showing before/after patterns
- **Comprehensive checklists** for each development area
- **Go-specific patterns** and idioms
- **Performance optimization** guidance
- **Security best practices** for Go applications

## 📊 **Rule Priorities**

- **Critical**: Error handling, concurrency, database security
- **High**: Performance, testing, API design
- **Medium**: Code organization, documentation

## 🔧 **Customization**

Each rule can be customized by modifying the `.mdc` files. Rules include:

- **Filters**: Target specific file types and content patterns
- **Actions**: Suggest improvements and provide examples
- **Examples**: Before/after code comparisons
- **Metadata**: Priority levels, versions, and tags

## 📚 **Best Practices Covered**

### **Code Organization**
- Package naming conventions
- Import organization
- File structure and size limits
- Package documentation

### **Error Handling**
- Always check errors
- Error wrapping with context
- Custom error types
- Error handling location

### **Concurrency**
- Goroutine management
- Channel usage patterns
- Context for cancellation
- Race condition prevention

### **Performance**
- Memory allocation optimization
- Slice and map pre-allocation
- String concatenation efficiency
- JSON handling optimization

### **Testing**
- Table-driven tests
- Test organization
- Mocking strategies
- Integration test separation

### **API Design**
- RESTful conventions
- HTTP status codes
- Request validation
- Middleware patterns

### **Database**
- Connection pooling
- Transaction management
- Query optimization
- Migration strategies

---

*These rules are designed to help maintain high-quality, performant, and maintainable Go code throughout the development lifecycle.* 