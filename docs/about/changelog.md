# Changelog

All notable changes to Orb will be documented here.

## [Unreleased]

### Added

- Complete MkDocs documentation
- Comprehensive security documentation
- User guides and examples

## [1.0.0] - 2024-01-15

### Added

- Initial release
- End-to-end encrypted file sharing
- Zero-trust relay server
- Argon2id key derivation
- Noise Protocol handshake
- ChaCha20-Poly1305 encryption
- TUI file browser
- CLI commands (share, connect, relay)
- Secure filesystem sandboxing
- Session management with expiration
- Cross-platform support (Linux, macOS, Windows)
- Build scripts and automation

### Security

- Path traversal protection
- Symlink escape prevention
- Rate limiting
- Session lockout after failed attempts
- Memory-hard key derivation

## Version History

- **1.0.0** - Initial stable release
- **0.9.0** - Beta release (testing)
- **0.1.0** - Alpha release (development)

## Upgrade Guide

### From 0.x to 1.0

No breaking changes. Binary compatible.

## Future Releases

### Planned for 1.1

- Progress bars for downloads
- Resume capability
- Bulk operations
- Performance improvements

### Planned for 2.0

- Post-quantum cryptography
- FUSE mounting support
- GUI client
- Enhanced metadata protection

## Release Process

1. Update version in code
2. Update CHANGELOG.md
3. Create git tag
4. Build binaries
5. Create GitHub release
6. Update documentation

## Support Policy

- **Latest release**: Full support
- **Previous minor**: Security fixes
- **Older versions**: Best effort

## Links

- [GitHub Releases](https://github.com/Zayan-Mohamed/orb/releases)
- [Security Advisories](https://github.com/Zayan-Mohamed/orb/security/advisories)
