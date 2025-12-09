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

# Update module dependencies across workspace
cd project-path && depbump module -R

# Update direct packages (default, -D is optional)
cd project-path && depbump update
cd project-path && depbump update -D

# Update each package (direct + indirect)
cd project-path && depbump update -E

# Update to latest versions (including prerelease)
cd project-path && depbump update -L

# Update across workspace modules
cd project-path && depbump update -R

# Combine flags
cd project-path && depbump update -D -R    # direct + recursive
cd project-path && depbump update -DR      # same as above
cd project-path && depbump update -E -R    # each + recursive
cd project-path && depbump update -ER      # same as above
```

### Advanced Usage

```bash
# Update GitHub packages
depbump update --github-only

# Skip GitLab dependencies
depbump update --skip-gitlab

# Update GitLab packages
depbump update --gitlab-only

# Skip GitHub dependencies
depbump update --skip-github

# Combine flags
depbump update -E --github-only

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

# Upgrade each package (direct + indirect) with Go version matching
depbump bump -E

# Upgrade to latest versions (including prerelease)
depbump bump -L

# Upgrade across workspace modules
depbump bump -R

# Combine flags
depbump bump -D -R          # direct + recursive
depbump bump -DR            # same as above
depbump bump -E -R          # each + recursive
depbump bump -ER            # same as above

# Note: -D and -E are exclusive, -E and -L are exclusive
```

**`bump` Command Features:**
- 🧠 **Go Version Matching**: Analyzes each package's Go version requirements
- 🚫 **Toolchain Contagion Prevention**: Avoids upgrades that would force toolchain changes
- ⬆️ **Upgrade-First Method**: Does not downgrade existing packages
- 📊 **Intelligent Analysis**: Shows version transitions with Go version requirements
- 🔄 **Workspace Integration**: Processes multiple Go modules with ease

### Command Structure

- **depbump**: Default module update (same as `depbump module`)
- **module**: Update module dependencies using `go get -u ./...`
  - `-R`: Update across workspace modules
- **update**: Update dependencies with filtering options
  - `-D`: Update direct dependencies (default)
  - `-E`: Update each package (direct + indirect)
  - `-L`: Use latest versions (including prerelease)
  - `-R`: Update across workspace modules
  - `--github-only` / `--skip-github`: GitHub filtering
  - `--gitlab-only` / `--skip-gitlab`: GitLab filtering
  - Note: `-D` and `-E` are exclusive
- **bump**: Smart Go version matching upgrades
  - `-D`: Upgrade direct dependencies (default)
  - `-E`: Upgrade each package (direct + indirect)
  - `-L`: Use latest versions (including prerelease)
  - `-R`: Upgrade across workspace modules
  - Note: `-E` and `-L` are exclusive
- **sync**: Git tag synchronization
  - **tags**: Sync to Git tag versions
  - **subs**: Sync with latest fallback

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

# Update module dependencies across workspace
depbump module -R

# Update direct dependencies (default)
depbump update

# Update each package (direct + indirect)
depbump update -E

# Update to latest versions
depbump update -L

# Update across workspace
depbump update -R

# Combine flags
depbump update -DR
depbump update -ER
```

### Bump Commands

```bash
# Smart upgrade with Go version matching
depbump bump

# Upgrade each package
depbump bump -E

# Upgrade across workspace
depbump bump -R

# Combine flags
depbump bump -DR
depbump bump -ER
```

### Sync Commands

```bash
# Sync dependencies to corresponding Git tag versions
depbump sync tags

# Sync dependencies with latest fallback
depbump sync subs
```

### Filtering Examples

```bash
# GitHub/GitLab specific updates
depbump update --github-only        # GitHub packages
depbump update --skip-github        # Skip GitHub dependencies
depbump update --gitlab-only        # GitLab packages
depbump update --skip-gitlab        # Skip GitLab dependencies

# Combine flags
depbump update -E --github-only
depbump update -L --skip-gitlab
```

## Troubleshooting

### Common Issues

1. **Toolchain Version Mismatch**
   - depbump manages Go toolchain versions
   - Uses project's Go version from go.mod to ensure matching
   - Set GOTOOLCHAIN environment variable if needed

2. **Package Conflicts**
   - Run `go mod tidy -e` following updates to clean up
   - Use `depbump update D` instead of `depbump update E` to get safe updates
   - Check go.mod when encountering incompatible version constraints

3. **Workspace Issues**
   - Ensure go.work file exists when running workspace commands
   - Run `depbump sync` to synchronize workspace dependencies
   - Check that modules are listed in go.work

## Tips and Best Practices

- **Start with direct packages**: Use `depbump update` (default) to get safe updates
- **Test updates**: Run tests when updating packages
- **Use version management**: Commit go.mod/go.sum before big updates
- **Step-wise updates**: Update packages in steps, not at once
- **Watch breaking changes**: Use `depbump update` before `depbump update -L`
- **Use bump command**: Use `depbump bump` to prevent Go toolchain contagion

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
