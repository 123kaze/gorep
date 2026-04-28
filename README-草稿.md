# gorep — 高性能并发文件搜索 CLI

> **项目定位**: 用 Go 编写的并发文件内容搜索工具，类似简化版的 ripgrep/grep。  
> **核心能力**: Fan-out/Fan-in 三阶段 Pipeline 并发模型，万级文件搜索毫秒级响应。  
> **项目路径**: `～/GolandProjects/gorep`

---

## 快速开始

```bash
# 基础搜索
./gorep "main" .

# 正则搜索
./gorep "func\s+\w+Handler" ./src/

# 固定字符串（不当作正则）
./gorep -F "fmt.Println" .

# 忽略大小写
./gorep -i "error" .

# 显示帮助
./gorep --help
```

## 安装

```bash
# 从源码构建
cd ~/GolandProjects/gorep
go build -o gorep .

# 或直接使用预编译的二进制
cp gorep ~/local/bin/gorep
```

---

## CLI 用法

```
用法: gorep <pattern> <path> [flags]
```

### Flags

| Flag | 短写 | 说明 | 默认值 |
|------|------|------|--------|
| `--ignore-case` | `-i` | 忽略大小写 | `false` |
| `--fixed-string` | `-F` | 固定字符串模式（非正则） | `false` |
| `--workers` | `-j` | 并发 worker 数量 | `runtime.NumCPU()` |
| `--include` | | glob 包含模式（可多次指定） | 所有文件 |
| `--exclude` | | glob 排除模式 | 无 |
| `--all` | `-a` | 搜索隐藏文件/目录 | `false` |
| `--no-color` | `-n` | 禁用彩色输出 | `false` |

### 使用示例

```bash
# 只搜 .go 文件
./gorep "import" --include="*.go" .

# 排除 vendor 目录
./gorep "TODO" --exclude="vendor/**" .

# 指定 8 个 worker
./gorep -j 8 "pattern" ./large-project/

# 搜索任意文件（包括隐藏文件）
./gorep -a "password" .
```

---

## 系统架构

```
                        ┌───────────────────────────────────────────────────┐
  $ gorep "pattern"     │                   gorep CLI                       │
  --include="*.go"      │                                                   │
  ./project/            │  ┌──────────┐     ┌────────────────────────────┐  │
                        │  │  Cobra   │────►│       Search Engine        │  │
                        │  │ (参数解析) │    │                            │  │
                        │  └──────────┘    │  ┌───────┐   ┌──────────┐  │  │
                        │                  │  │Walker │   │ Workers  │  │  │
                        │                  │  │ (遍历) │──►│ (搜索)    │  │  │
                        │                  │  └───────┘   └────┬─────┘  │  │
                        │                  │      fileCh       │ resultCh    │
                        │                  │                   ▼         │  │
                        │                  │  ┌───────────────────────┐  │  │
                        │                  │  │       Printer         │  │  │
                        │                  │  │ (彩色格式化 + 输出)     │  │  │
                        │                  │  └───────────────────────┘  │  │
                        │                  └────────────────────────────┘  │
                        └───────────────────────────────────────────────────┘
```

### 核心并发模型

```
                    ┌──────────┐
                    │  Walker  │  1 goroutine 遍历目录
                    │ (生产者)  │
                    └────┬─────┘
                         │ fileCh
              ┌──────────┼──────────┐
              ▼          ▼          ▼
        ┌──────────┐┌──────────┐┌──────────┐
        │ Worker 1 ││ Worker 2 ││ Worker N │  N goroutine 并发搜索
        │ (消费者)  ││ (消费者)  ││ (消费者)  │
        └────┬─────┘└────┬─────┘└────┬─────┘
             │           │           │
             │  resultCh │           │
             ▼           ▼           ▼
          ┌──────────────────────────┐
          │        Printer           │  1 goroutine 收集 + 输出
          │       (汇总输出)          │
          └──────────────────────────┘
```

**三阶段 Pipeline：**

1. **Walker**：1 个 goroutine 遍历目录树，过滤文件，将路径送入 `fileCh`
2. **Workers**：N 个 goroutine 从 `fileCh` 消费文件，逐行搜索匹配，结果送入 `resultCh`
3. **Printer**：主 goroutine 从 `resultCh` 收集结果，彩色格式化输出

Worker 数量默认等于 `runtime.NumCPU()`，自动利用多核。

---

## 项目结构

```
gorep/
├── cmd/
│   └── root.go              # Cobra 根命令 + flag 定义
├── internal/
│   ├── config/
│   │   └── config.go        # 搜索配置结构体
│   ├── engine/
│   │   ├── engine.go        # 搜索引擎编排（Walker/Worker/Printer 调度）
│   │   ├── matcher.go       # 匹配器（正则/固定字符串 + 高亮范围计算）
│   │   └── worker.go        # 单文件搜索（逐行扫描 + 匹配）
│   ├── filter/
│   │   ├── binary.go        # 二进制文件检测（读取前 512 字节判 null byte）
│   │   ├── glob.go          # glob 模式匹配（--include / --exclude）
│   │   └── walker.go        # 目录遍历器（尚未填充）
│   ├── model/
│   │   └── result.go        # 数据结构（LineMatch / FileMatch）
│   └── printer/
│       ├── printer.go       # Printer 接口
│       └── terminal.go      # 彩色终端输出（文件名 / 行号 / 匹配高亮）
├── main.go                  # 程序入口
├── go.mod / go.sum          # Go 模块依赖
├── Makefile                 # 构建脚本
├── .goreleaser.yaml         # 交叉编译发布配置
├── GOcli.md                 # 开发文档与面试话术
└── README.md                # 本文件
```

---

## 依赖

| 包 | 用途 |
|---|------|
| `github.com/spf13/cobra` | CLI 框架 |
| `github.com/spf13/pflag` | 参数解析 |
| `github.com/fatih/color` | 终端彩色输出 |

---

## 性能基准

> 测试目录: `/home/123kaze/Project/claude-code-source-code-main`  
> 文件数: 21,604 文件 | 目录: 2,841 | 大小: 436MB  
> 测试工具对比: gorep vs ripgrep (rg 14.1.1) vs GNU grep (3.12)  
> 完整报告: [benchmark_report.md](./benchmark_report.md)

| 测试场景 | gorep | rg | grep | 结果 |
|---------|:----:|:--:|:----:|------|
| 普通搜索 `class` | 0.27s | **0.04s** | 0.23s | rg 快 7x |
| 中文搜索 `函数` | 0.22s | **0.03s** | 0.32s | rg 快 7x, gorep > grep |
| 正则搜索 IP 模式 | 0.58s | **0.04s** | 0.37s | rg 快 14x |
| 忽略大小写 | 0.23s | **0.03s** | 0.33s | rg 快 7x |
| 大目录(400MB) | 0.32s | **0.05s** | 0.22s | rg 快 6x |
| 不存在字符串 | 0.28s | **0.05s** | 0.37s | rg 快 5x |
| 固定字符串 | **0.002s** | 0.001s | **0.001s** | grep > rg > gorep |

**性能总结：**

| 排名 | 工具 | 说明 |
|:---:|------|------|
| 🥇 | **rg** | Rust 实现，SIMD 加速 + mmap 内存映射，全面领先 |
| 🥈 | **gorep** | Go 并发实现，简单搜索接近 rg 水平 |
| 🥉 | **grep** | 单线程 POSIX 标准，系统标配 |

### gorep 当前瓶颈

| 问题 | 影响 | 改进方向 |
|------|:----:|---------|
| 二进制检测 `IsBinary()` 对每个文件额外 open+read | 2x syscall 开销 | 改为扩展名过滤 + 文件头采样 |
| `bufio.Scanner` 逐行读取 + 字符串分配 | 内存分配频繁 | 使用 `mmap` 或 `bytes.Index` |
| `filepath.WalkDir` 单线程串行遍历 | 大目录遍历慢 | 实现并发目录遍历 |
| Go `regexp` 无 SIMD 加速 | 正则搜索慢 14x | rust regex 对比已超出 Go 范围 |

---

## 技术亮点

### 1. Fan-out/Fan-in 并发 Pipeline

```go
fileCh := make(chan string, workers*2)     // Walker → Workers
resultCh := make(chan *FileMatch, workers) // Workers → Printer

go walkFiles(root, fileCh)                  // 1 个 Walker
for i := 0; i < workers; i++ {
    go searchFile(fileCh, resultCh)         // N 个 Workers
}
for match := range resultCh {               // 主 goroutine 收集结果
    printer.Print(match)
}
```

### 2. 并发安全的高亮输出

终端输出通过 `sync.Mutex` 互斥，防止多 goroutine 输出交错。匹配部分使用红色加粗 ANSI 颜色标记，提升可读性。

### 3. 二进制文件智能跳过

读取前 512 字节判断是否包含 null byte，自动跳过 `.pyc` / `.o` / `.so` 等二进制文件。

### 4. 多层文件过滤体系

```
隐藏目录跳过 → 隐藏文件跳过 → glob 包含过滤 → glob 排除过滤 → 二进制检测 → 送入 Worker
```

---

## 开发路线

| 版本 | 状态 | 特性 |
|:----:|:----:|------|
| v0.1 | ✅ | Cobra CLI + 基本参数 + 单文件搜索 |
| v0.2 | ✅ | Walker 目录遍历 + Worker Pool 并发搜索 |
| v0.3 | ✅ | glob 过滤 + 二进制检测 + 彩色高亮输出 |
| v0.4 | ⬜ | `.gitignore` 规则解析 |
| v0.5 | ⬜ | 上下文行显示（-B / -A / -C） |
| v0.6 | ⬜ | 搜索统计 + 超时控制 |
| v0.7 | ⬜ | JSON 输出格式 |
| v0.8 | ⬜ | 性能优化（mmap / 内存池） |

**已完成清单：** `遇到的问题.md` 记录了每一步开发中的关键设计决策与 bug 排查过程。

---

## 应用场景

| 场景 | 推荐 |
|------|:----:|
| 日常开发搜索 | 使用 `rg`（最快）或 `gorep`（无依赖） |
| 远程/容器环境 | `gorep`（单二进制部署） |
| 系统脚本管道 | `grep`（POSIX 标准） |
| 跨平台分发工具 | `gorep`（Go 编译零依赖） |
| 定制化搜索需求 | `gorep`（源码可修改） |

---

## License

MIT

---

*更多开发细节和面试话术见 [GOcli.md](./GOcli.md)*
*开发过程中的问题和决策见 [遇到的问题.md](./遇到的问题.md)*
*性能基准测试详细数据见 [benchmark_report.md](./benchmark_report.md)*
