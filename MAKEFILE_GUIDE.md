# Makefile 使用指南

## 快速开始

```bash
# 显示所有可用命令
make help

# 运行测试
make test

# 运行测试并显示覆盖率
make test-cover

# 运行所有质量检查
make quality
```

## 主要命令

### 测试相关

| 命令 | 说明 |
|------|------|
| `make test` | 运行所有测试（详细输出） |
| `make test-short` | 快速运行测试（无详细输出） |
| `make test-verbose` | 运行测试并显示详细输出 |
| `make test-cover` | 运行测试并生成覆盖率报告 |
| `make test-cover-html` | 生成 HTML 格式的覆盖率报告 |
| `make test-race` | 运行竞态检测测试 |
| `make test-all` | 运行所有测试套件 |
| `make quick` | 快速测试（短模式） |

### 性能测试

| 命令 | 说明 |
|------|------|
| `make bench` | 运行基准测试 |
| `make bench-cpu` | 运行基准测试并生成 CPU profile |
| `make bench-mem` | 运行基准测试并生成内存 profile |

### 代码质量

| 命令 | 说明 |
|------|------|
| `make fmt` | 格式化代码 |
| `make vet` | 运行 go vet 检查 |
| `make lint` | 运行 linter（需要 golangci-lint） |
| `make check` | 运行格式化和 vet 检查 |
| `make quality` | 运行所有质量检查 |

### 依赖管理

| 命令 | 说明 |
|------|------|
| `make tidy` | 整理 go modules |
| `make download` | 下载依赖 |
| `make verify` | 验证依赖 |

### 构建和运行

| 命令 | 说明 |
|------|------|
| `make demo` | 运行演示程序 |
| `make build-demo` | 构建演示程序 |

### CI/CD

| 命令 | 说明 |
|------|------|
| `make ci` | 本地运行 CI 流程 |
| `make all` | 运行检查和测试 |

### 其他

| 命令 | 说明 |
|------|------|
| `make clean` | 清理构建缓存和测试文件 |
| `make stats` | 显示代码统计信息 |
| `make install-tools` | 安装开发工具 |
| `make watch` | 监听文件变化并自动测试（需要 entr） |

## 使用示例

### 日常开发流程

```bash
# 1. 修改代码后，运行快速测试
make quick

# 2. 提交前运行完整检查
make quality

# 3. 本地模拟 CI 流程
make ci
```

### 测试覆盖率分析

```bash
# 生成覆盖率报告
make test-cover

# 生成 HTML 报告并在浏览器中查看
make test-cover-html
open coverage.html
```

### 性能分析

```bash
# 运行基准测试
make bench

# 生成 CPU profile 并分析
make bench-cpu
go tool pprof cpu.out

# 生成内存 profile 并分析
make bench-mem
go tool pprof mem.out
```

### 代码质量检查

```bash
# 格式化代码
make fmt

# 运行静态分析
make vet

# 运行 linter（首次使用需要先安装）
make install-tools
make lint
```

## 环境变量

如果系统中没有 `go` 命令，或需要使用特定版本的 Go，可以设置环境变量：

```bash
# 使用 gvm 管理的 Go 版本
export GOROOT=~/.gvm/gos/go1.21
make test

# 或者在命令中指定
GOROOT=~/.gvm/gos/go1.21 make test
```

## 持续集成

在 CI 环境中使用：

```bash
# GitHub Actions, GitLab CI 等
make ci
```

这将运行：
1. 清理环境
2. 整理依赖
3. 代码检查
4. 竞态测试
5. 覆盖率测试

## 注意事项

1. **首次使用**：建议先运行 `make install-tools` 安装必要的开发工具
2. **Go 版本**：项目要求 Go 1.17+
3. **竞态检测**：`make test-race` 会比常规测试慢，但能发现并发问题
4. **Watch 模式**：需要安装 `entr` 工具（`brew install entr` 或 `apt-get install entr`）

## 故障排除

### Go 命令未找到

```bash
# 设置 GOROOT 环境变量
export GOROOT=/path/to/your/go
make test
```

### golangci-lint 未安装

```bash
make install-tools
```

### 测试失败

```bash
# 查看详细输出
make test-verbose

# 运行特定包的测试
go test ./zap -v
```
