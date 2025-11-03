# Demo 目录完整说明

## 📁 目录结构

```
demo/
├── README.md              # 详细使用文档
├── run_demos.sh          # 一键运行脚本
├── basic_demo.go         # 基础功能演示
├── web_app_demo.go       # Web 应用模拟
├── performance_demo.go   # 性能测试
└── demo.go              # 原始简单示例
```

## 🚀 快速开始

### 方式 1: 使用运行脚本（推荐）

```bash
# 显示帮助
./demo/run_demos.sh

# 运行基础演示
./demo/run_demos.sh basic

# 运行 Web 应用演示
./demo/run_demos.sh webapp

# 运行性能测试
./demo/run_demos.sh performance

# 运行所有演示
./demo/run_demos.sh all
```

### 方式 2: 直接运行

```bash
# 基础演示
go run demo/basic_demo.go

# Web 应用演示
go run demo/web_app_demo.go

# 性能测试
go run demo/performance_demo.go
```

### 方式 3: 使用 Makefile

```bash
# 运行默认 demo
make demo

# 构建 demo 可执行文件
make build-demo
./bin/demo
```

## 📚 Demo 详细说明

### 1. basic_demo.go - 基础功能演示

**包含 8 个核心场景：**

1. **基础日志** (`basicLogging`)
   - Info, Warn, Error 基本用法
   - 格式化日志输出
   
2. **字段日志** (`loggerWithFields`)
   - 单字段添加
   - 多字段链式调用
   - 错误字段

3. **Context 日志** (`contextLogging`)
   - 完整的请求追踪
   - TraceID, UserID, Pathname
   - 自定义字段累积
   - 跨函数日志传递

4. **错误日志** (`errorLogging`)
   - WithError 使用
   - 错误上下文添加

5. **命名日志器** (`namedLoggers`)
   - 组件级别日志隔离
   - Database, Cache, API 分离

6. **自定义配置** (`customConfiguration`)
   - 日志级别配置
   - 文件输出配置
   - Debug 级别启用

7. **并发日志** (`concurrentLogging`)
   - 5 个 goroutine 并发
   - 每个独立的 trace ID

8. **日志级别** (`logLevels`)
   - 所有级别演示
   - 动态级别切换

**适合场景：** 学习基础用法、理解核心概念

### 2. web_app_demo.go - Web 应用模拟

**模拟完整 Web 服务器生命周期：**

- **启动阶段**
  - 应用配置加载
  - 数据库连接初始化
  - Redis 缓存连接
  - API 服务器启动

- **请求处理**
  - HTTP 请求接收
  - 用户认证
  - 业务逻辑处理
  - 数据库查询
  - 缓存操作
  - 响应返回

- **关闭阶段**
  - 优雅关闭
  - 连接清理

**特色功能：**
- 每个请求独立的 trace ID
- 完整的请求生命周期日志
- 跨层级的日志传递
- Named logger 组件隔离

**输出文件：** `demo/web_app.log`

**适合场景：** 了解实际项目集成、学习最佳实践

### 3. performance_demo.go - 性能和压力测试

**包含 4 个性能测试场景：**

1. **高频日志测试**
   - 连续输出 10,000 条日志
   - 测试吞吐量

2. **并发压力测试**
   - 100 个 goroutine
   - 每个输出 100 条日志
   - 总计 10,000 条并发日志

3. **大字段测试**
   - 大数组
   - 大 Map
   - 长字符串

4. **Context 累积测试**
   - 累积 100 个字段
   - 测试字段管理性能

**输出文件：** `demo/performance.log`

**适合场景：** 性能评估、压力测试、优化参考

### 4. demo.go - 原始示例

**简单快速示例：**
- 基本日志输出
- 字段添加
- OTEL 集成

**适合场景：** 快速验证、入门体验

## 📊 输出示例

### Console 格式
```
[2025-11-03 14:52:45.741] [info] [] demo/basic_demo.go:46 []  Simple info message
[2025-11-03 14:52:45.743] [info] [database] demo/basic_demo.go:149 []  Connection pool initialized {"pool_size": 10}
[2025-11-03 14:52:45.743] [info] [] demo/basic_demo.go:93 []  Request started {
  "trace_id": "trace-abc-123",
  "user_id": 999,
  "pathname": "/api/orders"
}
```

### 文件输出
生成的日志文件：
- `demo/demo.log` - 基础演示日志
- `demo/web_app.log` - Web 应用日志
- `demo/performance.log` - 性能测试日志

## 🎯 学习路径

### 初学者路径
```bash
1. go run demo/demo.go           # 快速体验
2. go run demo/basic_demo.go     # 学习核心功能
3. 阅读代码注释                    # 理解实现
```

### 进阶路径
```bash
1. go run demo/basic_demo.go     # 复习基础
2. go run demo/web_app_demo.go   # 学习最佳实践
3. 修改代码实验                    # 动手实践
```

### 高级路径
```bash
1. go run demo/performance_demo.go  # 了解性能特性
2. 查看性能报告                       # 分析结果
3. 编写自己的 benchmark              # 深入优化
```

## 🔧 自定义 Demo

### 创建新的 Demo

```bash
# 1. 复制现有 demo
cp demo/basic_demo.go demo/my_demo.go

# 2. 修改 package main 和函数

# 3. 运行
go run demo/my_demo.go
```

### Demo 模板

```go
package main

import (
    "github.com/gw123/glog"
    "github.com/gw123/glog/common"
)

func main() {
    // 配置日志
    glog.SetDefaultLoggerConfig(
        common.Options{},
        common.WithLevel(common.DebugLevel),
        common.WithOutputPath("./demo/my_demo.log"),
    )
    
    // 你的代码
    glog.Info("Hello from my demo!")
}
```

## 📈 性能参考

基于 performance_demo.go 的结果：

| 测试项 | 数量 | 性能 |
|-------|------|------|
| 高频日志 | 10,000 | ~50,000 msg/s |
| 并发日志 | 10,000 (100x100) | ~40,000 msg/s |
| 大字段日志 | 1 | <5ms |
| 字段累积 | 100 fields | <1ms |

## 🧹 清理

```bash
# 清理所有日志文件
rm -f demo/*.log

# 或使用 make
make clean
```

## 💡 提示

1. **查看输出**：日志会同时输出到 console 和文件
2. **并发测试**：日志顺序可能不同，这是正常的
3. **性能测试**：会生成大量日志，注意磁盘空间
4. **实时查看**：可以使用 `tail -f demo/*.log` 实时查看日志

## 🐛 故障排除

### Go 命令未找到
```bash
# 设置 Go 路径
export PATH=$PATH:/usr/local/go/bin

# 或使用完整路径
/path/to/go run demo/basic_demo.go
```

### 权限错误
```bash
# 给脚本执行权限
chmod +x demo/run_demos.sh
```

### 文件写入错误
```bash
# 确保 demo 目录可写
chmod 755 demo/
```

## 📖 相关文档

- [README.md](../README.md) - 项目主文档
- [MAKEFILE_GUIDE.md](../MAKEFILE_GUIDE.md) - Makefile 使用指南
- [TEST_REPORT.md](../TEST_REPORT.md) - 测试报告
- [CLAUDE.md](../CLAUDE.md) - 架构说明

## 🤝 贡献

欢迎贡献新的 demo 示例！创建 PR 时请：
1. 添加清晰的注释
2. 更新 README.md
3. 提供使用示例
