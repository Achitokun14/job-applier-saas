# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure
- Planning and architecture documentation

## [1.0.0] - 2026-03-27

### Added

#### Backend (Go + Chi)
- REST API server with Chi router
- JWT authentication (register, login)
- User profile management
- Job search and application endpoints
- Settings management
- SQLite support for local development
- PostgreSQL support for production
- GORM ORM integration
- CORS configuration
- Request logging and middleware

#### Frontend (SvelteKit + Bun)
- Landing page with feature highlights
- User authentication (login/register)
- Dashboard with application statistics
- Job search interface
- Application tracking page
- Profile editor
- Settings page with LLM configuration
- Responsive design with CSS styling
- Auth store for state management
- API client library

#### Python Service (FastAPI)
- Resume generation endpoint
- Cover letter generation endpoint
- Job description parsing endpoint
- Multiple resume styles (modern, classic, minimal)
- LLM integration (OpenAI, Anthropic, Google)
- HTML to PDF conversion
- Template-based generation (fallback)

#### TUI (Bubble Tea)
- Terminal user interface with Bubble Tea
- Login/Register forms
- Job search functionality
- Application management
- Profile viewing
- Settings management
- API client for backend communication

#### Infrastructure
- Docker configuration for all services
- Docker Compose for local development
- Docker Compose for production
- Caddy reverse proxy configuration
- GitHub Actions CI workflow
- Environment variable configuration

#### Documentation
- README.md
- LICENSE.md (Proprietary + AGPL-3.0)
- SECURITY.md
- CHANGELOG.md
- ARCHITECTURE.md
- GET_STARTED.md
- COMMANDS.md
- APIs.md
- Technical documentation for each service

### Changed
- N/A (Initial release)

### Deprecated
- N/A (Initial release)

### Removed
- N/A (Initial release)

### Fixed
- N/A (Initial release)

### Security
- JWT token authentication
- Password hashing with bcrypt
- CORS configuration
- Input validation

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 1.0.0 | 2026-03-27 | Initial release |

## Upgrade Guide

### From 0.x to 1.0.0

This is the initial release. No upgrade path needed.

## Roadmap

### v1.1.0 (Planned)
- [ ] LinkedIn integration for direct applications
- [ ] Indeed job scraping
- [ ] Application status tracking automation
- [ ] Email notifications

### v1.2.0 (Planned)
- [ ] Multi-user support with teams
- [ ] Role-based access control
- [ ] Audit logging
- [ ] API rate limiting

### v2.0.0 (Future)
- [ ] AI-powered job matching
- [ ] Interview preparation features
- [ ] Salary negotiation assistance
- [ ] Career path recommendations

## Contributors

- Achitokun14 - Project owner and lead developer

## Acknowledgments

- [Jobs_Applier_AI_Agent_AIHawk](https://github.com/feder-cr/Jobs_Applier_AI_Agent_AIHawk) - Original inspiration and resume generation logic
