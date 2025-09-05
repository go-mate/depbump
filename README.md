# depbump

Check and upgrade outdated dependencies in Go modules, with version bumping.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Key Features

🔄 **Smart Dependency Upgrades**: Auto detect and upgrade outdated Go module dependencies  
⚡ **Multiple Update Strategies**: Support direct, indirect, and all dependency updates  
🎯 **Version Control Integration**: Git tag synchronization for consistent dependency versions  
🌍 **Source Filtering**: Selective updates for GitHub/GitLab sources  
📋 **Workspace Support**: Go workspace multi-module batch dependency management

## Install

```bash
go install github.com/go-mate/depbump/cmd/depbump@latest
```

## Usage

### Basic Usage

```bash
# Basic module update (updates go.mod dependencies)
cd project-path && depbump

# Update module dependencies (same as above, explicit)
cd project-path && depbump module

# Update direct dependencies only
cd project-path && depbump direct

# Update direct dependencies to latest versions
cd project-path && depbump direct latest

# Update every dependency
cd project-path && depbump everyone

# Update every dependency to latest versions  
cd project-path && depbump everyone latest
```

### Advanced Usage

```bash
# Update only GitHub dependencies
depbump direct --github-only

# Skip GitLab dependencies
depbump direct --skip-gitlab

# Update only GitLab dependencies
depbump direct --gitlab-only

# Skip GitHub dependencies
depbump direct --skip-github

# Sync workspace dependencies
depbump sync

# Sync dependencies to Git tag versions
depbump sync tags

# Sync dependencies, use latest when tags missing
depbump sync subs
```

### Dependency Categories

- **module**: Update module dependencies using `go get -u ./...`
- **direct**: Update only direct (explicit) dependencies declared in go.mod  
- **everyone**: Update every dependency - aliases: `require`, `requires`
- **latest**: Get latest available versions (may have breaking changes)
- **update**: Get compatible updates (respects semantic versioning)

### Source Filtering Options

- `--github-only`: Update only GitHub-hosted dependencies
- `--skip-github`: Skip GitHub-hosted dependencies
- `--gitlab-only`: Update only GitLab-hosted dependencies
- `--skip-gitlab`: Skip GitLab-hosted dependencies

## Features

### Smart Dependency Management

depbump provides intelligent dependency management that can:
- Auto parse dependency information from `go.mod` files
- Detect available upgrade versions
- Handle version compatibility issues
- Support Go toolchain version management

### Workspace Integration

Supports Go 1.18+ workspace features:
- Auto discover all modules in workspace
- Batch process dependency updates for multiple modules
- Maintain consistency across workspace dependencies
- Auto execute `go work sync`

### Git Tag Synchronization

Provides Git tag integration functionality:
- Sync dependency versions to corresponding Git tags
- Support tag version verification
- Handle missing tag scenarios

## Command Reference

### Update Commands

```bash
# Update module dependencies (default action)
depbump

# Update module dependencies (explicit)
depbump module

# Update direct dependencies with compatible versions
depbump direct

# Update direct dependencies to latest versions
depbump direct latest

# Update all dependencies including indirect ones
depbump everyone

# Update all dependencies to latest versions
depbump everyone latest
```

### Sync Commands

```bash
# Execute go work sync for workspace
depbump sync

# Sync dependencies to their Git tag versions
depbump sync tags

# Sync dependencies with latest fallback
depbump sync subs
```

### Filtering Examples

```bash
# GitHub/GitLab specific updates
depbump direct --github-only      # Only GitHub dependencies
depbump direct --skip-github      # Skip GitHub dependencies
depbump direct --gitlab-only      # Only GitLab dependencies
depbump direct --skip-gitlab      # Skip GitLab dependencies

# Combine with latest mode
depbump direct latest --github-only
depbump everyone latest --skip-gitlab
```

## Troubleshooting

### Common Issues

1. **Toolchain Version Mismatch**
   - depbump automatically manages Go toolchain versions
   - Uses project's Go version from go.mod to ensure compatibility
   - Set GOTOOLCHAIN environment variable if needed

2. **Dependency Conflicts**
   - Run `go mod tidy -e` after updates to clean up
   - Use `depbump direct` instead of `depbump everyone` for safer updates
   - Check go.mod for incompatible version constraints

3. **Workspace Issues**
   - Ensure go.work file exists for workspace commands
   - Run `depbump sync` to synchronize workspace dependencies
   - Check that all modules are properly listed in go.work

## Tips and Best Practices

- **Start with direct dependencies**: Use `depbump direct` for safer updates
- **Test after updates**: Always run tests after dependency updates
- **Use version control**: Commit go.mod/go.sum before major updates
- **Incremental updates**: Update dependencies gradually, not all at once
- **Monitor breaking changes**: Use `depbump direct` (compatible) before `depbump direct latest`
- **Workspace consistency**: Run `depbump sync` after module updates in workspaces

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-08-28 08:33:43.829511 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a bug?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share your use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize by reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo for new releases and features
- 🌟 **Success stories?** Share how this package improved your workflow
- 💬 **General feedback?** All suggestions and comments are welcome

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage interface).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement your changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation for user-facing changes and use meaningful commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a pull request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Happy Coding with this package!** 🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/depbump.svg?variant=adaptive)](https://starchart.cc/go-mate/depbump)
