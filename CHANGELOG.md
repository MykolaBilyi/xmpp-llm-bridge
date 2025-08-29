# Changelog
<!-- markdownlint-configure-file { "MD024": { "siblings_only": true } } -->

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.3] - 2025-08-29

### Changed

- Github Actions workflow for automated release builds and Docker image publishing, now uses CHANGELOG.md for release notes

## [0.1.2] - 2025-08-29

### Added

- Added project description in README.md
- Added CHANGELOG.md file
- devcontainer config
- Dependabot config

## [0.1.1] - 2025-08-19

### Changed

- Major refactoring of project structure to implement hexagonal architecture
  - Moved from `api/` structure to `internal/app/`, `internal/ports/`, and `internal/adapters/` pattern
  - Relocated handlers from `api/client/handlers/` to `internal/app/client/handlers/`
  - Moved server components from `api/server/` to `internal/app/server/`
  - Restructured main application logic into `internal/app/app.go` with proper dependency injection
  - Enhanced XMPP handling with new `pkg/xmpp/` utilities and message multiplexing
  - Improved configuration structure and expanded default configuration options
  - Added new port interfaces for better abstraction (`xmppSession.go`)
  - Refined entity models and adapter implementations for better separation of concerns

## [0.1.0] - 2025-08-19

### Added

- Initial release of XMPP-LLM Bridge service
- XMPP client implementation with SASL authentication using mellium.im/xmpp
- OpenAI LLM integration for processing chat messages
- Basic XMPP message handlers:
  - Debug handler for logging incoming stanzas
  - Echo message handler for testing
  - LLM message handler for AI-powered responses
- Configuration management with YAML config files and environment variable overrides
- Support for loading secrets from files using `_FILE` suffix environment variables
- Docker containerization with multi-stage builds
- HTTP server with health check endpoint on port 8080
- Development container configuration with VS Code devcontainer
- Basic project structure with hexagonal architecture foundations
- GitHub Actions workflow for automated release builds and Docker image publishing

[unreleased]: https://github.com/MykolaBilyi/xmpp-llm-bridge/compare/v0.1.3...HEAD
[0.1.3]: https://github.com/MykolaBilyi/xmpp-llm-bridge/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/MykolaBilyi/xmpp-llm-bridge/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/MykolaBilyi/xmpp-llm-bridge/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/MykolaBilyi/xmpp-llm-bridge/releases/tag/v0.1.0
