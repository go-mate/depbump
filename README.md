[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/depbump/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/depbump/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/depbump)](https://pkg.go.dev/github.com/go-mate/depbump)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/depbump/main.svg)](https://coveralls.io/github/go-mate/depbump?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
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
🎯 **Version Management Integration**: Git tag synchronization to maintain consistent package versions
🌍 **Source Filtering**: Selective updates targeting GitHub/GitLab sources
📋 **Workspace Support**: Go workspace multi-module batch package management

## Installation

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
cd project-path && depbump M        # Short alias // 短别名

# Update direct packages
cd project-path && depbump direct
cd project-path && depbump D        # Short alias // 短别名

# Update direct dependencies to latest versions
cd project-path && depbump direct latest
cd project-path && depbump D latest

# Update each package
cd project-path && depbump everyone
cd project-path && depbump E        # Short alias // 短别名
cd project-path && depbump each     # Alternative name // 别名

# Update each package to latest versions
cd project-path && depbump everyone latest
cd project-path && depbump E latest
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
# Smart Go version matching checks and upgrades (default: direct dependencies)
# Prevents Go toolchain contagion while upgrading dependencies
depbump bump

# Upgrade direct dependencies with Go version matching
depbump bump direct
depbump bump D              # Short alias // 短别名

# Upgrade direct dependencies to latest versions
depbump bump direct latest
depbump bump D L            # Short aliases // 短别名

# Upgrade direct dependencies across workspace modules
depbump bump direct recursive
depbump bump D R            # Short aliases // 短别名

# Upgrade each package (direct + indirect) with Go version matching
depbump bump everyone
depbump bump E              # Short alias // 短别名
depbump bump each           # Alternative name // 别名

# Upgrade each package across workspace modules
depbump bump everyone recursive
depbump bump E R            # Short aliases // 短别名
depbump bump -E -R          # Flag-based (equivalent) // 标志方式（等效）

# Works in workspace environment (processes each module)
cd workspace-root && depbump bump
```

**Convenient Flag-Based Usage:**
```bash
# Use flags instead of subcommands (more concise)
depbump bump -E             # everyone (each package)
depbump bump -R             # recursive (workspace modules)
depbump bump -E -R          # everyone + recursive
depbump bump -L -R          # latest + recursive

# Note: -E and -L cannot be used at the same time (exclusive flags)
```

**New `bump` Command Features:**
- 🧠 **Go Version Matching**: Analyzes each package's Go version requirements
- 🚫 **Toolchain Contagion Prevention**: Avoids upgrades that would force toolchain changes
- ⬆️ **Upgrade-First Method**: Does not downgrade existing packages
- 📊 **Intelligent Analysis**: Shows version transitions with Go version requirements
- 🔄 **Workspace Integration**: Processes multiple Go modules with ease

### Package Categories

- **module**: Update module dependencies using `go get -u ./...`
- **direct**: Update direct (explicit) packages declared in go.mod - aliases: `directs`
- **everyone**: Update each package - aliases: `require`, `requires`
- **bump**: Smart Go version matching upgrades (default: direct)
  - **bump direct**: Upgrade direct dependencies with version matching - aliases: `directs`
  - **bump everyone**: Upgrade each package with version matching - aliases: `require`, `requires`
- **latest**: Get latest available versions (might have breaking changes)
- **update**: Get compatible updates (respects semantic versioning)

### Source Filtering Options

- `--github-only`: Update packages hosted on GitHub
- `--skip-github`: Skip dependencies hosted on GitHub
- `--gitlab-only`: Update packages hosted on GitLab
- `--skip-gitlab`: Skip dependencies hosted on GitLab

## Features

### Smart Package Management

depbump provides intelligent package management that can:
- Auto parse package information from `go.mod` files
- Detect available upgrade versions
- Handle version matching issues
- Support Go toolchain version management

### Workspace Integration

Supports Go 1.18+ workspace features:
- Auto detect modules in workspace
- Batch process package updates across multiple modules
- Maintain coherence across workspace packages
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

# Update each package including indirect ones
depbump everyone

# Update each package to latest versions
depbump everyone latest
```

### Sync Commands

```bash
# Execute go work sync in workspace
depbump sync

# Sync dependencies to corresponding Git tag versions
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
   - Use `depbump direct` instead of `depbump everyone` to get safe updates
   - Check go.mod when encountering incompatible version constraints

3. **Workspace Issues**
   - Ensure go.work file exists when running workspace commands
   - Run `depbump sync` to synchronize workspace dependencies
   - Check that modules are listed in go.work

## Tips and Best Practices

- **Start with direct packages**: Use `depbump direct` to get safe updates
- **Test updates**: Run tests when updating packages
- **Use version management**: Commit go.mod/go.sum before big updates
- **Step-wise updates**: Update packages in steps, not at once
- **Watch breaking changes**: Use `depbump direct` (compatible) before `depbump direct latest`
- **Workspace sync**: Run `depbump sync` when updating modules in workspaces

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
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
7. **Documentation**: Update documentation to support client-facing changes
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/depbump.svg?variant=adaptive)](https://starchart.cc/go-mate/depbump)
