# 🚀 Quick Start Guide to Contributing to Gin

Welcome! This guide will help you get started contributing to the Gin Web Framework.

## ✅ Environment Setup Complete

Your development environment is ready:
- **Go Version**: 1.25.0 ✓ (requires 1.24+)
- **Dependencies**: All modules verified ✓
- **Dev Tools**: golint and misspell installed ✓
- **Tests**: All passing (99.1% coverage) ✓

## 📚 Understanding the Project Structure

### Core Files
- **`gin.go`** - Main engine and router
- **`context.go`** - Request/response context handling
- **`routergroup.go`** - Route grouping and middleware
- **`tree.go`** - High-performance routing tree (radix tree)
- **`logger.go`** - Logging middleware
- **`recovery.go`** - Panic recovery middleware

### Important Directories
- **`binding/`** - Request data binding (JSON, XML, YAML, forms, etc.)
- **`render/`** - Response rendering (JSON, XML, HTML, etc.)
- **`internal/`** - Internal utilities (not public API)
- **`examples/`** - Usage examples

## 🎯 How to Contribute

### 1. Find Something to Work On

**Good First Issues:**
- Check [GitHub Issues](https://github.com/gin-gonic/gin/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- Look for issues labeled "good first issue" or "help wanted"
- Browse [Discussions](https://github.com/gin-gonic/gin/discussions) for feature ideas

**Areas You Can Help:**
- 🐛 Fix bugs
- 📝 Improve documentation in `docs/doc.md`
- ✨ Add new features
- 🧪 Improve test coverage
- ⚡ Performance optimizations
- 🌐 Translate documentation

### 2. Development Workflow

**Before Starting:**
```bash
# Create a new branch for your work
git checkout -b feature/my-contribution
```

**Make Your Changes:**
```bash
# Format your code (required)
make fmt

# Run tests
make test

# Check for common issues
make vet

# Check formatting (CI will run this)
make fmt-check

# Fix misspellings
make misspell
```

**Run Specific Tests:**
```bash
# Test a specific package
go test -v ./binding/

# Test a specific function
go test -v -run TestContextQuery

# Run with race detector
go test -race ./...
```

### 3. Code Standards

**Important Rules:**
- ✅ All tests must pass
- ✅ Code must be formatted with `make fmt`
- ✅ Add tests for new features
- ✅ Update `docs/doc.md` for new features (NOT README.md)
- ✅ Keep commits focused and atomic
- ✅ Squash to max 2 commits before submitting PR

**Example Test Pattern:**
```go
func TestMyNewFeature(t *testing.T) {
    // Setup
    router := gin.New()
    
    // Test your feature
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, 200, w.Code)
}
```

### 4. Submit Your PR

**Checklist:**
- [ ] Tests pass locally (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] No linting errors (`make vet`)
- [ ] Documentation updated (if needed)
- [ ] Commits squashed to ≤2 commits
- [ ] PR targets `master` branch
- [ ] PR description explains what and why

**PR Template Will Ask:**
- What does this PR do?
- Why is this change needed?
- How has this been tested?
- What issues does it fix? (use `Fixes #123`)

## 🛠️ Common Development Tasks

### Running the Examples
```bash
cd examples/
go run <example-name>/main.go
```

### Checking Code Coverage
```bash
make test
# coverage.out file will be generated
go tool cover -html=coverage.out
```

### Testing Against Different Go Versions
```bash
# The project supports Go 1.24+
# Make sure your changes work with that version
```

### Debugging Tips
```bash
# Enable debug mode
export GIN_MODE=debug

# Run with verbose output
go test -v -run TestName
```

## 📖 Key Resources

- **API Docs**: https://pkg.go.dev/github.com/gin-gonic/gin
- **User Guides**: https://gin-gonic.com/
- **Discussions**: https://github.com/gin-gonic/gin/discussions
- **Examples**: https://github.com/gin-gonic/examples

## 🤔 Getting Help

**Questions?**
- Post in [GitHub Discussions](https://github.com/gin-gonic/gin/discussions)
- Use English for all communications
- Search existing issues/discussions first

**Security Issues?**
- Email: appleboy.tw@gmail.com
- Do NOT post publicly

## 💡 Contribution Ideas for Beginners

1. **Documentation Improvements**
   - Fix typos or unclear explanations
   - Add code examples
   - Improve error messages

2. **Test Coverage**
   - Add tests for edge cases
   - Improve test clarity

3. **Performance Benchmarks**
   - Add benchmarks for new features
   - Run: `go test -bench=. -benchmem`

4. **Code Quality**
   - Simplify complex functions
   - Add helpful comments
   - Improve error handling

## 🎉 Ready to Contribute!

**Quick Command Reference:**
```bash
# Check everything before submitting
make fmt         # Format code
make test        # Run tests
make vet         # Check for issues
make misspell    # Fix spelling

# Or check formatting without changing files
make fmt-check
make misspell-check
```

**Next Steps:**
1. Browse issues or discussions for ideas
2. Comment on an issue to claim it
3. Create your branch and start coding
4. Run the checks above
5. Submit your PR!

---

Thank you for contributing to Gin! Every contribution, no matter how small, helps make Gin better for the entire community. 🙌
