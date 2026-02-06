# Orb Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.3] 2026-02-07

### Fixes

* Harden SecureFilesystem.Read to prevent excessive memory allocation:
* Validates that requested read length is non-negative.
* Caps read length at 10MB and ensures it does not exceed file size.
* Safely converts length to int before allocation, returning an error if it would overflow.
* Prevents potential crashes or undefined behavior from malicious or invalid input.

## [1.1.2] 2026-02-06

### Changes

* Fixed path traversal issues in file downloads
* Strengthened filename validation
* Resolved gosec false positive (G304)
* Added gosec security scanning in CI
* Reduced TUI cyclomatic complexity
* Improved error handling in download flow
* Added progress UI for chunked downloads


### Security

- Fixed potential file inclusion vulnerability in file download (CVE-None)
- Strengthened filename validation with regex whitelist to prevent path traversal
- Added #nosec comment for gosec G304 false positive

### Fixed

- Resolved SARIF upload compatibility issues in GitHub security workflow

## [1.0.0] 2026-01-14

### Added

- GitHub Actions workflows for builds, releases, and documentation
- Comprehensive documentation with MkDocs
- Security scanning with gosec and Trivy
- Issue and PR templates
- Contributing guidelines
- GIF placeholders for screen recordings

### Changed

- Updated README with badges and improved developer experience
- Enhanced project structure documentation

### Fixed

- (Add bug fixes here)

### Security

- (Add security-related changes here)

## [1.0.0] - YYYY-MM-DD (Template)

### Added

- Initial release
- Zero-trust folder tunneling
- End-to-end encryption with Noise Protocol
- TUI file browser
- Cross-platform support (Linux, macOS, Windows)
- Relay server
- Secure session management

### Security

- Argon2id for key derivation
- ChaCha20-Poly1305 for transport encryption
- Path sanitization and symlink protection
- Replay protection
- Rate limiting

---

## Release Template

Use this template for new releases:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added

- New features

### Changed

- Changes to existing functionality

### Deprecated

- Features that will be removed in future versions

### Removed

- Features that were removed

### Fixed

- Bug fixes

### Security

- Security-related changes
```

[Unreleased]: https://github.com/Zayan-Mohamed/orb/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Zayan-Mohamed/orb/releases/tag/v1.0.0
