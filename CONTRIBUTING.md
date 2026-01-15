# Contributing to Orb

Thank you for your interest in contributing to Orb! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help maintain a positive community

## Security First

Orb is designed with security as the top priority. Any contribution that weakens encryption, privacy, or isolation will not be accepted. If you discover a security vulnerability, please email security@orb.example.com instead of creating a public issue.

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report:

- Check the [existing issues](https://github.com/Zayan-Mohamed/orb/issues)
- Try the latest version from the main branch
- Collect relevant information (version, OS, logs)

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.yml) when creating an issue.

### Suggesting Features

Feature suggestions are welcome! However, please remember:

- Features must align with Orb's security-first philosophy
- Consider the complexity vs. benefit trade-off
- Provide clear use cases and examples

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.yml).

### Improving Documentation

Documentation improvements are always appreciated:

- Fix typos and grammar
- Add examples and clarifications
- Update outdated information
- Create tutorials or guides

### Contributing Code

#### Getting Started

1. **Fork the repository**

   ```bash
   gh repo fork Zayan-Mohamed/orb --clone
   cd orb
   ```

2. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

3. **Set up development environment**

   ```bash
   # Install Go 1.22 or later
   go version

   # Install dependencies
   go mod download

   # Run tests to verify setup
   make test
   ```

#### Development Workflow

1. **Make your changes**

   - Write clear, commented code
   - Follow Go best practices
   - Add tests for new functionality
   - Update documentation as needed

2. **Test your changes**

   ```bash
   # Run unit tests
   make test

   # Run linter
   golangci-lint run

   # Build locally
   make build-local

   # Test manually
   ./orb share ./test-folder
   ```

3. **Commit your changes**

   ```bash
   # Use conventional commit messages
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve issue with..."
   git commit -m "docs: update README"
   ```

   Commit message prefixes:

   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `test:` - Test additions or changes
   - `refactor:` - Code refactoring
   - `perf:` - Performance improvements
   - `chore:` - Build process or auxiliary tool changes

4. **Push to your fork**

   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request**
   - Use the [Pull Request template](.github/PULL_REQUEST_TEMPLATE.md)
   - Provide a clear description of changes
   - Link related issues
   - Ensure CI checks pass

#### Code Style

- Follow standard Go formatting (`gofmt`, `goimports`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small
- Avoid premature optimization

Example:

```go
// Good
func validateSessionID(id string) error {
    if len(id) != SessionIDLength {
        return ErrInvalidSessionID
    }
    return nil
}

// Not recommended
func validate(s string) error {
    if len(s) != 6 { return errors.New("invalid") }
    return nil
}
```

#### Testing Guidelines

- Write tests for all new functionality
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

Example:

```go
func TestValidateSessionID(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        wantErr bool
    }{
        {"valid ID", "ABC123", false},
        {"too short", "ABC", true},
        {"too long", "ABC1234567", true},
        {"empty", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateSessionID(tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateSessionID() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Project Structure

```
orb/
├── cmd/              # CLI commands
├── internal/         # Internal packages
│   ├── crypto/      # Cryptography
│   ├── filesystem/  # File system operations
│   ├── relay/       # Relay server
│   ├── session/     # Session management
│   ├── tui/         # Terminal UI
│   └── tunnel/      # Network tunnel
├── pkg/             # Public packages
├── docs/            # Documentation
└── build/           # Build artifacts
```

## Areas for Contribution

### High Priority

- [ ] Performance optimizations
- [ ] Cross-platform testing
- [ ] Documentation improvements
- [ ] Bug fixes

### Medium Priority

- [ ] Additional tests
- [ ] Code refactoring
- [ ] CI/CD improvements
- [ ] Example scripts

### Low Priority

- [ ] Optional features (that don't compromise security)
- [ ] Developer tooling
- [ ] Community resources

## Review Process

1. **Automated Checks**

   - All tests must pass
   - Linting must pass
   - Security scans must pass
   - Code coverage should not decrease

2. **Code Review**

   - At least one maintainer approval required
   - Address all review comments
   - Maintain security standards

3. **Testing**

   - Changes tested on multiple platforms
   - Manual testing completed
   - No regressions introduced

4. **Documentation**
   - User-facing changes documented
   - Code comments updated
   - CHANGELOG.md updated if needed

## Release Process

Releases are managed by maintainers:

1. Version bump in code
2. Update CHANGELOG.md
3. Create git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions automatically builds and publishes release

## Communication

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: Questions, ideas, general discussion
- **Pull Requests**: Code contributions
- **Email**: security@orb.example.com (security issues only)

## Recognition

Contributors are recognized in:

- GitHub contributors page
- Release notes (for significant contributions)
- Documentation acknowledgments

## Questions?

If you have questions:

1. Check the [documentation](https://zayan-mohamed.github.io/orb/)
2. Search [existing issues](https://github.com/Zayan-Mohamed/orb/issues)
3. Ask in [Discussions](https://github.com/Zayan-Mohamed/orb/discussions)
4. Create a new issue if needed

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Orb! Together, we can build secure and user-friendly tools.
