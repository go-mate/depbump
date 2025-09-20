[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/depbump/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/depbump/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/depbump)](https://pkg.go.dev/github.com/go-mate/depbump)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/depbump/main.svg)](https://coveralls.io/github/go-mate/depbump?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.22+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-mate/depbump.svg)](https://github.com/go-mate/depbump/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mate/depbump)](https://goreportcard.com/report/github.com/go-mate/depbump)

# depbump

Check and upgrade outdated dependencies in Go modules, with version bumping.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🔄 **Smart Package Upgrades**: Auto detect and upgrade outdated Go module packages
⚡ **Multiple Update Strategies**: Support direct, indirect, and package updates
🧠 **Go Version Matching**: Intelligent analysis prevents toolchain contagion during upgrades
🎯 **Version Management Integration**: Git tag synchronization for consistent package versions
🌍 **Source Filtering**: Selective updates for GitHub/GitLab sources
📋 **Workspace Support**: Go workspace multi-module batch package management

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

# Update direct packages
cd project-path && depbump direct

# Update direct dependencies to latest versions
cd project-path && depbump direct latest

# Update each package
cd project-path && depbump everyone

# Update each package to latest versions  
cd project-path && depbump everyone latest
```

### Advanced Usage

```bash
# Update GitHub packages
depbump direct --github-only

# Skip GitLab dependencies
depbump direct --skip-gitlab

# Update GitLab packages
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

### Intelligent Package Management

```bash
# Smart Go version matching checks and upgrades
# Prevents Go toolchain contagion while upgrading dependencies
depbump bump

# Works in workspace environment (processes all modules)
cd workspace-root && depbump bump
```

**New `bump` Command Features:**
- 🧠 **Go Version Matching**: Analyzes each package's Go version requirements
- 🚫 **Toolchain Contagion Prevention**: Avoids upgrades that would force toolchain changes
- ⬆️ **Upgrade-First Approach**: Does not downgrade existing packages
- 📊 **Intelligent Analysis**: Shows version transitions with Go version requirements
- 🔄 **Workspace Integration**: Processes multiple Go modules well

### Package Categories

- **module**: Update module dependencies using `go get -u ./...`
- **direct**: Update direct (explicit) packages declared in go.mod  
- **everyone**: Update each package - aliases: `require`, `requires`
- **latest**: Get latest available versions (might have breaking changes)
- **update**: Get compatible updates (respects semantic versioning)

### Source Filtering Options

- `--github-only`: Update GitHub-hosted packages
- `--skip-github`: Skip GitHub-hosted dependencies
- `--gitlab-only`: Update GitLab-hosted packages
- `--skip-gitlab`: Skip GitLab-hosted dependencies

## Features

### Smart Package Management

depbump provides intelligent package management that can:
- Auto parse package information from `go.mod` files
- Detect available upgrade versions
- Handle version matching issues
- Support Go toolchain version management

### Workspace Integration

Supports Go 1.18+ workspace features:
- Auto find modules in workspace
- Batch process package updates for multiple modules
- Keep consistency across workspace packages
- Auto execute `go work sync`

### Git Tag Synchronization

Provides Git tag integration features:
- Sync package versions to corresponding Git tags
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
depbump direct --github-only      # GitHub packages
depbump direct --skip-github      # Skip GitHub dependencies
depbump direct --gitlab-only      # GitLab packages
depbump direct --skip-gitlab      # Skip GitLab dependencies

# Combine with latest mode
depbump direct latest --github-only
depbump everyone latest --skip-gitlab
```

## Troubleshooting

### Common Issues

1. **Toolchain Version Mismatch**
   - depbump manages Go toolchain versions
   - Uses project's Go version from go.mod to ensure matching
   - Set GOTOOLCHAIN environment variable if needed

2. **Package Conflicts**
   - Run `go mod tidy -e` following updates to clean up
   - Use `depbump direct` instead of `depbump everyone` for safe updates
   - Check go.mod for incompatible version constraints

3. **Workspace Issues**
   - Ensure go.work file exists for workspace commands
   - Run `depbump sync` to synchronize workspace dependencies
   - Check that modules are listed in go.work

## Tips and Best Practices

- **Start with direct packages**: Use `depbump direct` for safe updates
- **Test updates**: Run tests when updating packages
- **Use version management**: Commit go.mod/go.sum before big updates
- **Step-wise updates**: Update packages in steps, not at once
- **Watch breaking changes**: Use `depbump direct` (matching) before `depbump direct latest`
- **Workspace sync**: Run `depbump sync` when updating modules in workspaces

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-06 04:53:24.895249 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a bug?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a pull request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/depbump.svg?variant=adaptive)](https://starchart.cc/go-mate/depbump)
