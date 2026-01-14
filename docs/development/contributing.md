# Contributing to Orb

Thank you for your interest in contributing to Orb!

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Submit a pull request

## Development Setup

```bash
git clone https://github.com/yourusername/orb.git
cd orb
go mod download
make build
```

## Code Guidelines

### Style

- Follow Go conventions
- Use `gofmt`
- Run `go vet`
- Pass `golangci-lint`

### Testing

- Write unit tests
- Maintain >80% coverage
- Test edge cases
- Include examples

### Documentation

- Comment exported functions
- Update README if needed
- Add examples
- Document breaking changes

## Pull Request Process

1. **Create Issue First**

   - Describe the problem
   - Propose solution
   - Get feedback

2. **Write Code**

   - Follow style guide
   - Add tests
   - Update docs

3. **Submit PR**

   - Clear description
   - Link to issue
   - Pass CI checks

4. **Code Review**

   - Address feedback
   - Make requested changes
   - Maintain civility

5. **Merge**
   - Squash commits
   - Update changelog
   - Close issue

## Testing

```bash
# Run all tests
make test

# With coverage
go test -cover ./...

# Specific package
go test ./internal/crypto
```

## Commit Messages

```
type(scope): subject

body

footer
```

**Types:**

- feat: New feature
- fix: Bug fix
- docs: Documentation
- test: Tests
- refactor: Refactoring
- chore: Maintenance

**Example:**

```
feat(crypto): add support for XChaCha20

Implement XChaCha20-Poly1305 for extended nonce space.

Closes #123
```

## Code Review Checklist

- [ ] Tests pass
- [ ] Code formatted
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
- [ ] Security implications considered
- [ ] Performance impact assessed

## Areas for Contribution

### High Priority

- Security audits
- Performance optimization
- Documentation improvements
- Test coverage
- Bug fixes

### Features

- Post-quantum cryptography
- Resume capability
- Progress bars
- Bulk operations
- GUI client

### Infrastructure

- CI/CD improvements
- Docker images
- Package managers
- Benchmarks
- Fuzzing

## Community

- GitHub Issues: Bug reports
- GitHub Discussions: Questions
- Pull Requests: Contributions
- Security: security@orb.example.com

## Code of Conduct

Be respectful, inclusive, and professional.

## License

By contributing, you agree your code is licensed under MIT.

## Next Steps

- [Building from Source](building.md)
- [API Reference](api.md)
