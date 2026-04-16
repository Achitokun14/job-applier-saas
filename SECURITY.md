# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of Auto Job Applier SaaS seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: **security@[your-domain].com**

Please include the following information in your report:

- Description of the vulnerability
- Steps to reproduce the issue
- Possible impact of the vulnerability
- Any potential fixes or mitigations you've identified

### What to Expect

1. **Acknowledgment**: You will receive an acknowledgment of your report within 48 hours.

2. **Assessment**: We will assess the vulnerability and determine its severity within 7 days.

3. **Updates**: We will keep you informed of our progress towards fixing the vulnerability.

4. **Resolution**: We aim to release a fix within 30 days of confirming the vulnerability.

5. **Credit**: With your permission, we will credit you in the security advisory.

### Severity Levels

| Level | Description | Response Time |
|-------|-------------|---------------|
| **Critical** | Remote code execution, SQL injection, authentication bypass | 24 hours |
| **High** | Privilege escalation, data exposure | 7 days |
| **Medium** | XSS, CSRF, limited data exposure | 14 days |
| **Low** | Information disclosure, minor issues | 30 days |

## Security Measures

### Authentication & Authorization

- JWT-based authentication with configurable expiry
- Password hashing using bcrypt with salt
- Role-based access control (RBAC)
- Session management with secure token storage

### Data Protection

- All API communications over HTTPS/TLS
- Sensitive data encrypted at rest
- Database credentials stored in environment variables
- No sensitive data in logs

### API Security

- CORS configuration with allowed origins
- Rate limiting on authentication endpoints
- Input validation and sanitization
- SQL injection prevention via GORM ORM
- XSS protection headers

### Infrastructure

- Docker containers run as non-root users
- Network isolation between services
- Regular security updates for base images
- Caddy reverse proxy with security headers

### Development Practices

- Dependencies regularly updated
- Security-focused code reviews
- Automated security scanning in CI/CD
- No secrets in version control

## Security Configuration

### Environment Variables

Never commit sensitive data. Use `.env` files for:

```env
JWT_SECRET=<strong-random-secret>
DB_PASSWORD=<strong-password>
LLM_API_KEY=<api-key>
```

### Recommended Settings

```env
# JWT Configuration
JWT_SECRET=<64-char-random-string>
JWT_EXPIRY=24h

# Database
DB_PASSWORD=<16-char-mixed-password>

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

### Generating Secure Secrets

```bash
# Generate JWT secret
openssl rand -base64 64

# Generate database password
openssl rand -base64 16
```

## Known Security Considerations

1. **LLM API Keys**: Store securely, never expose to frontend
2. **Selenium**: Running in headless mode reduces attack surface
3. **File Uploads**: Resume uploads should be validated
4. **Rate Limiting**: Implement on production deployments

## Security Checklist for Deployment

- [ ] Change default JWT_SECRET
- [ ] Use strong database passwords
- [ ] Enable HTTPS with valid certificates
- [ ] Configure firewall rules
- [ ] Set up monitoring and alerting
- [ ] Review CORS settings
- [ ] Enable security headers
- [ ] Regular backups configured
- [ ] Log rotation configured

## Contact

- Security Issues: security@[your-domain].com
- General Issues: GitHub Issues

## Acknowledgments

We thank the following security researchers for responsibly disclosing vulnerabilities:

- (None yet - be the first!)
