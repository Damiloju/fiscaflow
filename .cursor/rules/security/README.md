# Security Rules

Rules for security best practices and vulnerability prevention.

## Rules

### security-best-practices.mdc
Prevents common security vulnerabilities and enforces security best practices.

**Features:**
- Environment variable usage for secrets
- Input validation and sanitization
- SQL injection prevention
- XSS prevention strategies
- CSRF protection requirements
- Authentication best practices
- Authorization and access control
- HTTPS enforcement
- Dependency vulnerability management
- Secure error handling

**Usage:**
- Applied when working with security-sensitive code
- Detects potential security vulnerabilities
- Provides secure coding suggestions

## Security Checklist

### 🔐 Authentication & Authorization
- [ ] Use secure authentication libraries
- [ ] Implement proper role-based access control
- [ ] Validate user permissions on all endpoints
- [ ] Use JWT tokens with proper expiration

### 🛡️ Input Validation
- [ ] Validate all user inputs
- [ ] Sanitize data before processing
- [ ] Use parameterized queries for databases
- [ ] Implement Content Security Policy

### 🔒 Data Protection
- [ ] Never hardcode secrets in code
- [ ] Use environment variables for sensitive data
- [ ] Encrypt sensitive data at rest
- [ ] Use HTTPS for all communications

### 🚨 Error Handling
- [ ] Don't expose sensitive information in errors
- [ ] Log security events appropriately
- [ ] Implement proper exception handling
- [ ] Use generic error messages for users

## Best Practices

1. **Principle of Least Privilege**: Grant minimum necessary permissions
2. **Defense in Depth**: Implement multiple security layers
3. **Regular Updates**: Keep dependencies and systems updated
4. **Security Testing**: Regular security audits and penetration testing
5. **Incident Response**: Have a plan for security incidents 