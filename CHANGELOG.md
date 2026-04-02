# Changelog

All notable changes to this project will be documented in this file.

## [v1.2.0] - 2026-04-02

### Added

- Added support for the `--json` flag to the `fcm init` command to allow structured outputs in automation pipelines.
- Added `*firebase-adminsdk*.json` and `*.yaml` to `.gitignore` to prevent accidental leaks of service accounts and config files.

### Changed

- Reorganized `README.md` to explicitly list `fcm init` flags (`--file`, `--force`, `--json`) in their own structured table for better visibility.

### Fixed

- Stopped futile retries on HTTP 4xx client errors (e.g., `400 Bad Request`, `404 Not Found`), except for `429 Too Many Requests`. The CLI now immediately exits when a token is invalid instead of spamming duplicate failed requests.
- Fixed a silent failure issue where final errors (for single-send operations) were swallowed and not printed to the console when `--json` mode was disabled.

## [v1.1.2] - 2026-04-01

### Added

- Rich ANSI banner with gradient effects and project information.
- New `banner.go` file to handle CLI branding and usage formatting.

### Fixed

- Updated `.gitignore` to prevent recursive ignoring of the `cmd/fcm/` directory by removing the broad `fcm` pattern.

### Changed

- Refactored CLI entry point to suppress banner display during normal command execution. The banner now only appears on help commands or when no arguments are provided.
## [v1.1.1] - 2026-03-27

### Changed

- Updated GitHub Actions release workflow to support the new modular project structure.
- Updated `README.md` with corrected manual build instructions.
- Updated `.gitignore` to exclude build artifacts and `dist/` directory.

## [v1.0.1] - 2026-03-27

### Changed

- Refactored monolithic `fcm.go` into a modular project structure.
- Organized code into specialized internal packages: `auth`, `config`, `fcm`, `log`, `model`, and `util`.
- Moved CLI entry point to `cmd/fcm/main.go`.
- Improved maintainability and internal code organization.

## [v1.0.0] - 2026-03-27

### Added

- Ported FCM CLI tool from Node.js to Go.
- Full FCM v1 support: single token, multicast, topic, and condition messaging.
- Authentication via Google OAuth2 service account JSON keys.
- Automatic retry logic with exponential backoff.
- Parallel sending using Go routines for high performance.
- Interactive progress bar for multicast sending.
- CLI flags for all major options (key, token, notification, data, etc.).
- YAML configuration support via `fcm.yaml`.
- Profile-based configuration (`--profile <name>`).
- CLI flag overrides over YAML config.
- Support for config file loading via `--config` / `-f`.
- `fcm init` command to generate a starter `fcm.yaml`.
- Support for loading tokens from file via `--tokens-file`.
- `.env` support via `.env`, `FCM_ENV_FILE`, and `--env-file`.
- Environment-based defaults for `FCM_KEY`, `FCM_CONFIG`, and `FCM_LOG`.
- `--json` machine-readable output mode.
- Message ID output for successful single sends.
- Structured JSON result for multicast sends.
- Per-token success/error reporting for batch delivery.
- Custom usage help and version information.
- GitHub Actions for Continuous Integration (CI) and automated multi-platform releases.
- Quick installation script for easy deployment.
- Project licensing under **Mozilla Public License Version 2.0**.
- Initial documentation and README improvements.
- Better validation for message targets (token, tokens, topic, condition).
- More flexible configuration for production and CI workflows.
- Better local onboarding for first-time users.
- Easier batch sending workflows for QA and automation.
- Better CI integration.
- Easier programmatic handling of FCM responses.
