# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure with frontend, backend, and configuration files
- Docker Compose setup for development environment
- Basic authentication with JWT
- Log search functionality
- Log export functionality
- PDPA field masking in exports
- Comprehensive documentation including DEVELOPMENT.md, ARCHITECTURE.md, SECURITY.md, and CONTRIBUTING.md
- Health check endpoints
- CORS configuration
- Rate limiting middleware
- Input validation for search queries
- Automated testing setup

### Changed
- Updated README.md with project overview and setup instructions
- Improved error handling and logging
- Enhanced security measures including JWT authentication and CORS configuration
- Refactored code structure for better maintainability

### Fixed
- Resolved issues with Docker Compose configuration
- Fixed authentication flow issues
- Corrected search query parsing
- Resolved export functionality bugs
- Fixed CORS configuration issues

### Removed
- Removed backup and main diagram files
- Removed outdated backend README.md

## [0.1.0] - 2026-06-15

### Added
- Initial release of APP system log management dashboard
- Complete backend API with Gin framework
- Frontend interface with Astro framework
- QuickWit integration for log storage and search
- JWT-based authentication system
- Log export functionality with PDPA compliance
- Docker deployment configuration
- Comprehensive documentation
- Security implementation with CORS, rate limiting, and input validation

### Changed
- Updated project dependencies to latest versions
- Improved code documentation
- Enhanced error handling and logging
- Optimized search performance

### Fixed
- Resolved authentication token handling issues
- Fixed search query parameter parsing
- Corrected export file generation
- Fixed CORS configuration for production

## [0.0.1] - 2026-06-01

### Added
- Initial project scaffolding
- Basic backend structure with Go and Gin
- Frontend structure with Astro
- Docker configuration files
- Initial documentation files
- Basic authentication implementation
- Search functionality
- Export functionality

### Changed
- Initial project setup and configuration
- Basic UI components
- Initial API endpoints

### Fixed
- Initial build and deployment issues
- Basic configuration problems

## [0.0.0] - 2026-05-25

### Added
- Project initialization
- Repository structure
- Initial README.md
- Basic development environment setup
- Initial commit with core files

[Unreleased]: https://github.com/your-username/app/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/your-username/app/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/your-username/app/compare/v0.0.0...v0.0.1
[0.0.0]: https://github.com/your-username/app/commits/v0.0.0