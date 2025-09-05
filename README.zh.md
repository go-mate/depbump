# depbump

检查并升级 Go 模块中的过时依赖，支持版本升级功能。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🔄 **智能依赖升级**: 自动检测和升级过时的 Go 模块依赖  
⚡ **多种更新策略**: 支持直接依赖、间接依赖和全部依赖更新  
🎯 **版本控制集成**: 集成 Git 标签同步，确保依赖版本一致性  
🌍 **源过滤支持**: 支持 GitHub/GitLab 源的选择性更新  
📋 **工作区支持**: 支持 Go workspace 跨模块批量依赖管理

## 安装

```bash
go install github.com/go-mate/depbump/cmd/depbump@latest
```

## 使用方法

### 基础用法

```bash
# 基本模块更新（更新 go.mod 依赖）
cd project-path && depbump

# 更新模块依赖（同上，显式指定）
cd project-path && depbump module

# 仅更新直接依赖
cd project-path && depbump direct

# 更新直接依赖到最新版本
cd project-path && depbump direct latest

# 更新每个依赖
cd project-path && depbump everyone

# 更新每个依赖到最新版本
cd project-path && depbump everyone latest
```

### 高级用法

```bash
# 仅更新 GitHub 依赖
depbump direct --github-only

# 跳过 GitLab 依赖
depbump direct --skip-gitlab

# 仅更新 GitLab 依赖
depbump direct --gitlab-only

# 跳过 GitHub 依赖
depbump direct --skip-github

# 同步工作区依赖
depbump sync

# 同步依赖到 Git 标签版本
depbump sync tags

# 同步依赖，缺失标签时使用最新版本
depbump sync subs
```

### 依赖类别

- **module**: 使用 `go get -u ./...` 更新模块依赖
- **direct**: 仅更新 go.mod 中直接声明的依赖  
- **everyone**: 更新每个依赖 - 别名：`require`, `requires`
- **latest**: 获取最新可用版本（可能有破坏性更改）
- **update**: 获取兼容的更新版本（遵循语义化版本）

### 源过滤选项

- `--github-only`: 仅更新 GitHub 上的依赖
- `--skip-github`: 跳过 GitHub 上的依赖
- `--gitlab-only`: 仅更新 GitLab 上的依赖
- `--skip-gitlab`: 跳过 GitLab 上的依赖

## 功能说明

### 智能依赖管理

depbump 提供了智能的依赖管理功能，能够：
- 自动解析 `go.mod` 文件中的依赖信息
- 检测可用的升级版本
- 处理版本兼容性问题
- 支持 Go toolchain 版本管理

### 工作区集成

支持 Go 1.18+ 的工作区功能：
- 自动发现工作区中的所有模块
- 批量处理多个模块的依赖更新
- 保持工作区依赖的一致性
- 自动执行 `go work sync`

### Git 标签同步

提供与 Git 标签的集成功能：
- 同步依赖版本到对应的 Git 标签
- 支持标签版本验证
- 处理缺失标签的情况

## 命令参考

### 更新命令

```bash
# 更新模块依赖（默认操作）
depbump

# 更新模块依赖（显式指定）
depbump module

# 更新直接依赖到兼容版本
depbump direct

# 更新直接依赖到最新版本
depbump direct latest

# 更新每个依赖包括间接依赖
depbump everyone

# 更新每个依赖到最新版本
depbump everyone latest
```

### 同步命令

```bash
# 执行 go work sync 同步工作区
depbump sync

# 同步依赖到其 Git 标签版本
depbump sync tags

# 同步依赖，带最新版本回退
depbump sync subs
```

### 过滤示例

```bash
# GitHub/GitLab 特定更新
depbump direct --github-only      # 仅更新 GitHub 依赖
depbump direct --skip-github      # 跳过 GitHub 依赖
depbump direct --gitlab-only      # 仅更新 GitLab 依赖
depbump direct --skip-gitlab      # 跳过 GitLab 依赖

# 与 latest 模式结合
depbump direct latest --github-only
depbump everyone latest --skip-gitlab
```

## 故障排除

### 常见问题

1. **工具链版本不匹配**
   - depbump 自动管理 Go 工具链版本
   - 使用项目 go.mod 中的 Go 版本确保兼容性
   - 如需要可设置 GOTOOLCHAIN 环境变量

2. **依赖冲突**
   - 更新后运行 `go mod tidy -e` 进行清理
   - 使用 `depbump direct` 而非 `depbump everyone` 以获得更安全的更新
   - 检查 go.mod 中的不兼容版本约束

3. **工作区问题**
   - 确保 go.work 文件存在以使用工作区命令
   - 运行 `depbump sync` 同步工作区依赖
   - 检查所有模块是否正确列在 go.work 中

## 技巧和最佳实践

- **从直接依赖开始**: 使用 `depbump direct` 进行更安全的更新
- **更新后测试**: 依赖更新后务必运行测试
- **使用版本控制**: 大型更新前提交 go.mod/go.sum
- **渐进式更新**: 逐步更新依赖，不要一次全部更新
- **监控破坏性变更**: 先使用 `depbump direct`（兼容）再使用 `depbump direct latest`
- **工作区一致性**: 在工作区中更新模块后运行 `depbump sync`

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-08-28 08:33:43.829511 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **意见反馈？** 欢迎所有建议和宝贵意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Pull Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Pull Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**使用这个包快乐编程！** 🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-mate/depbump.svg?variant=adaptive)](https://starchart.cc/go-mate/depbump)