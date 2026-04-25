# Go CLI 并发文件搜索工具 · 开发路线

---

## 一、为什么做这个项目？

| 维度 | 说明 |
|------|------|
| **定位** | 简历第三个项目，展示 Go CLI 开发 + 系统编程能力，与 go-rag 形成互补 |
| **实用性** | 自己日常开发就能用，替代 grep/find 的场景 |
| **技术差异化** | go-rag 偏 AI 工程，这个偏**系统编程 + 并发模型 + IO 优化**，面试能聊不同方向 |
| **开发周期** | 1 周可完成核心功能 |

### 和其他项目的配合

| 项目 | 展示能力 |
|------|----------|
| **go-rag** | AI 工程、Pipeline 并发、向量检索、流式响应 |
| **go-search（本项目）** | 系统编程、文件 IO、worker pool、CLI 框架、性能优化 |
| 两者互补 | 一个偏业务+AI，一个偏底层+工具链，覆盖面广 |

---

## 二、项目概述

### 项目名：`go-search`

一个用 Go 编写的**高性能并发文件内容搜索 CLI 工具**。支持正则表达式、glob 过滤、递归目录遍历、自动忽略 `.gitignore` 规则、彩色高亮输出。核心设计是 **worker pool + channel pipeline** 并发模型，充分利用多核 CPU。

### 核心功能

```
基础功能（Day 1-3）：
├── 递归目录遍历 + 文件内容搜索（正则/固定字符串）
├── 并发 worker pool（goroutine + channel）
├── 彩色终端输出（匹配行高亮、文件名着色）
├── glob 模式过滤文件（--include / --exclude）
└── 基本 CLI 参数解析（Cobra）

进阶功能（Day 4-5）：
├── .gitignore 规则解析与跳过
├── 二进制文件自动检测与跳过
├── 上下文行显示（-B before / -A after / -C context）
├── 搜索结果统计（匹配文件数、匹配行数、耗时）
└── 大文件流式读取（bufio.Scanner，避免一次性加载）

加分功能（Day 6-7，时间充裕再做）：
├── 多模式搜索（同时搜多个关键词）
├── 替换模式（--replace，类似 sed）
├── JSON 输出格式（--json，方便管道处理）
├── 搜索结果缓存（对同一目录的重复搜索加速）
└── 性能 benchmark 对比 grep/ripgrep
```

### 技术栈

| 组件 | 选型 | 理由 |
|------|------|------|
| **CLI 框架** | Cobra | Go CLI 标准库，支持子命令、flag、自动补全 |
| **正则引擎** | regexp / regexp2 | 标准库 regexp 基于 RE2（线性时间），性能可控 |
| **文件遍历** | filepath.WalkDir | Go 1.16+ 引入，比 Walk 更高效（不调用 Stat） |
| **并发控制** | goroutine + channel + sync.WaitGroup | 原生并发原语，worker pool 模式 |
| **终端着色** | fatih/color 或 ANSI 转义码 | 彩色高亮输出 |
| **gitignore 解析** | go-git/go-git 或自研简易解析 | 尊重 .gitignore 规则 |
| **构建发布** | goreleaser | 自动交叉编译 + GitHub Release |

---

## 三、系统架构

```
                        ┌─────────────────────────────────────────────┐
  $ go-search "pattern" │              go-search CLI                  │
  --include="*.go"      │                                             │
  ./project/            │  ┌──────────┐   ┌────────────────────────┐ │
                        │  │  Cobra   │──►│   Search Engine        │ │
                        │  │ (参数解析)│   │                        │ │
                        │  └──────────┘   │  ┌──────┐  ┌────────┐ │ │
                        │                 │  │Walker│  │Workers │ │ │
                        │                 │  │(遍历)│─►│(搜索)  │ │ │
                        │                 │  └──────┘  └───┬────┘ │ │
                        │                 │       channel   │      │ │
                        │                 │  ┌──────────────▼────┐ │ │
                        │                 │  │   Printer         │ │ │
                        │                 │  │ (彩色格式化输出)    │ │ │
                        │                 │  └───────────────────┘ │ │
                        │                 └────────────────────────┘ │
                        └─────────────────────────────────────────────┘
```

### 核心并发模型（Fan-out / Fan-in）

```
                    ┌──────────┐
                    │  Walker  │  1 个 goroutine 遍历目录
                    │ (生产者)  │
                    └────┬─────┘
                         │ fileCh (chan string)
              ┌──────────┼──────────┐
              ▼          ▼          ▼
        ┌──────────┐┌──────────┐┌──────────┐
        │ Worker 1 ││ Worker 2 ││ Worker N │  N 个 goroutine 并发搜索
        │ (消费者)  ││ (消费者)  ││ (消费者)  │
        └────┬─────┘└────┬─────┘└────┬─────┘
             │           │           │
             └─────┬─────┘───────────┘
                   │ resultCh (chan SearchResult)
                   ▼
             ┌──────────┐
             │ Printer  │  1 个 goroutine 收集并输出
             │ (汇总)    │
             └──────────┘
```

**三阶段 Pipeline：**
1. **Walker**（1 个 goroutine）：遍历目录，过滤文件（glob/gitignore/二进制检测），将文件路径发送到 `fileCh`
2. **Workers**（N 个 goroutine）：从 `fileCh` 读取文件路径，打开文件逐行搜索，匹配结果发送到 `resultCh`
3. **Printer**（1 个 goroutine）：从 `resultCh` 读取结果，格式化彩色输出到终端

---

## 四、项目结构

```
go-search/
├── cmd/
│   └── root.go                    # Cobra 根命令 + flag 定义
├── internal/
│   ├── engine/                    # ★ 搜索引擎核心
│   │   ├── engine.go              # 搜索引擎入口（编排 Walker/Worker/Printer）
│   │   ├── walker.go              # 目录遍历器（生产者）
│   │   ├── worker.go              # 文件搜索器（消费者，worker pool）
│   │   └── matcher.go             # 匹配器（正则/固定字符串/多模式）
│   ├── filter/                    # 文件过滤
│   │   ├── gitignore.go           # .gitignore 规则解析
│   │   ├── glob.go                # glob 模式匹配（--include/--exclude）
│   │   └── binary.go              # 二进制文件检测（读前 512 字节）
│   ├── printer/                   # 输出格式化
│   │   ├── printer.go             # Printer 接口
│   │   ├── terminal.go            # 彩色终端输出
│   │   └── json.go                # JSON 输出（--json）
│   ├── model/                     # 数据结构
│   │   └── result.go              # SearchResult、FileMatch 等
│   └── config/
│       └── config.go              # 搜索配置（并发数、上下文行数等）
├── main.go                        # 入口
├── go.mod
├── Makefile
├── .goreleaser.yaml               # 交叉编译发布配置
└── README.md
```

---

## 五、核心代码设计

### 1. 搜索引擎 — Worker Pool + Channel Pipeline（核心亮点）

```go
// internal/engine/engine.go

type Engine struct {
    config  *config.SearchConfig
    matcher *matcher.Matcher
    filter  *filter.Filter
    printer printer.Printer
    logger  *log.Logger
}

type SearchConfig struct {
    Pattern     string   // 搜索模式
    RootDir     string   // 搜索根目录
    Workers     int      // worker 数量（默认 runtime.NumCPU()）
    Includes    []string // glob 包含
    Excludes    []string // glob 排除
    IgnoreCase  bool     // 忽略大小写
    FixedString bool     // 固定字符串模式（非正则）
    BeforeCtx   int      // 前 N 行上下文
    AfterCtx    int      // 后 N 行上下文
}

func (e *Engine) Run(ctx context.Context) (*SearchStats, error) {
    fileCh := make(chan string, 128)        // Walker → Workers
    resultCh := make(chan *model.FileMatch, 64) // Workers → Printer

    var wg sync.WaitGroup
    stats := &SearchStats{}

    // 阶段 1：Walker（1 个 goroutine 遍历目录）
    go func() {
        defer close(fileCh)
        e.walk(ctx, e.config.RootDir, fileCh)
    }()

    // 阶段 2：Workers（N 个 goroutine 并发搜索）
    for i := 0; i < e.config.Workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for filePath := range fileCh {
                match, err := e.searchFile(ctx, filePath)
                if err != nil {
                    continue
                }
                if match != nil {
                    resultCh <- match
                }
            }
        }()
    }

    // Workers 全部完成后关闭 resultCh
    go func() {
        wg.Wait()
        close(resultCh)
    }()

    // 阶段 3：Printer（主 goroutine 收集输出）
    for match := range resultCh {
        stats.MatchedFiles++
        stats.MatchedLines += len(match.Lines)
        e.printer.Print(match)
    }

    return stats, nil
}
```

### 2. Walker — 目录遍历 + 过滤

```go
// internal/engine/walker.go

func (e *Engine) walk(ctx context.Context, root string, fileCh chan<- string) {
    filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil // 跳过无权限的目录
        }

        select {
        case <-ctx.Done():
            return ctx.Err() // 支持取消
        default:
        }

        // 跳过隐藏目录（.git 等）
        if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
            return filepath.SkipDir
        }

        // .gitignore 过滤
        if d.IsDir() && e.filter.IsGitIgnored(path) {
            return filepath.SkipDir
        }

        if d.IsDir() {
            return nil
        }

        // glob 过滤
        if !e.filter.MatchGlob(path) {
            return nil
        }

        // 二进制文件检测
        if filter.IsBinary(path) {
            return nil
        }

        fileCh <- path
        return nil
    })
}
```

### 3. Worker — 文件搜索 + 上下文行

```go
// internal/engine/worker.go

func (e *Engine) searchFile(ctx context.Context, filePath string) (*model.FileMatch, error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var matches []model.LineMatch
    scanner := bufio.NewScanner(f)
    // 设置更大的 buffer 处理长行
    scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

    lineNum := 0
    // 环形缓冲区存储前 N 行（用于 -B 上下文）
    beforeBuf := make([]string, 0, e.config.BeforeCtx)

    for scanner.Scan() {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }

        lineNum++
        line := scanner.Text()

        if e.matcher.Match(line) {
            matches = append(matches, model.LineMatch{
                LineNum:    lineNum,
                Content:    line,
                BeforeCtx:  copySlice(beforeBuf),
            })
        }

        // 维护环形缓冲区
        if len(beforeBuf) >= e.config.BeforeCtx {
            beforeBuf = beforeBuf[1:]
        }
        beforeBuf = append(beforeBuf, line)
    }

    if len(matches) == 0 {
        return nil, nil
    }

    return &model.FileMatch{
        FilePath: filePath,
        Lines:    matches,
    }, nil
}
```

### 4. Matcher — 正则 / 固定字符串匹配

```go
// internal/engine/matcher.go

type Matcher struct {
    pattern     string
    regex       *regexp.Regexp
    fixedString bool
    ignoreCase  bool
}

func NewMatcher(pattern string, fixedString, ignoreCase bool) (*Matcher, error) {
    m := &Matcher{
        pattern:     pattern,
        fixedString: fixedString,
        ignoreCase:  ignoreCase,
    }

    if fixedString {
        if ignoreCase {
            m.pattern = strings.ToLower(pattern)
        }
        return m, nil
    }

    // 正则模式
    regexPattern := pattern
    if ignoreCase {
        regexPattern = "(?i)" + pattern
    }
    re, err := regexp.Compile(regexPattern)
    if err != nil {
        return nil, fmt.Errorf("无效的正则表达式: %w", err)
    }
    m.regex = re
    return m, nil
}

func (m *Matcher) Match(line string) bool {
    if m.fixedString {
        if m.ignoreCase {
            return strings.Contains(strings.ToLower(line), m.pattern)
        }
        return strings.Contains(line, m.pattern)
    }
    return m.regex.MatchString(line)
}

// FindHighlightRanges 返回匹配位置，用于彩色高亮
func (m *Matcher) FindHighlightRanges(line string) [][2]int {
    if m.fixedString {
        return findAllFixed(line, m.pattern, m.ignoreCase)
    }
    locs := m.regex.FindAllStringIndex(line, -1)
    ranges := make([][2]int, len(locs))
    for i, loc := range locs {
        ranges[i] = [2]int{loc[0], loc[1]}
    }
    return ranges
}
```

### 5. Printer — 彩色终端输出

```go
// internal/printer/terminal.go

type TerminalPrinter struct {
    fileColor    *color.Color // 紫色
    lineNumColor *color.Color // 绿色
    matchColor   *color.Color // 红色加粗
    sepColor     *color.Color // 灰色
    matcher      *matcher.Matcher
    mu           sync.Mutex   // 防止多 goroutine 输出交错
}

func (p *TerminalPrinter) Print(match *model.FileMatch) {
    p.mu.Lock()
    defer p.mu.Unlock()

    // 文件名（紫色）
    p.fileColor.Println(match.FilePath)

    for _, line := range match.Lines {
        // 上下文行
        for _, ctxLine := range line.BeforeCtx {
            p.sepColor.Printf("  %d-", line.LineNum)
            fmt.Println(ctxLine)
        }

        // 匹配行：行号（绿色）+ 内容（匹配部分红色高亮）
        p.lineNumColor.Printf("  %d:", line.LineNum)
        p.printHighlighted(line.Content)
        fmt.Println()
    }
    fmt.Println() // 文件间空行
}

func (p *TerminalPrinter) printHighlighted(line string) {
    ranges := p.matcher.FindHighlightRanges(line)
    if len(ranges) == 0 {
        fmt.Print(line)
        return
    }

    lastEnd := 0
    for _, r := range ranges {
        fmt.Print(line[lastEnd:r[0]])        // 非匹配部分
        p.matchColor.Print(line[r[0]:r[1]])  // 匹配部分（红色加粗）
        lastEnd = r[1]
    }
    fmt.Print(line[lastEnd:]) // 剩余部分
}
```

---

## 六、CLI 使用示例

```bash
# 基础搜索
go-search "TODO" ./project/

# 正则搜索 + 只搜 Go 文件
go-search "func\s+\w+Handler" --include="*.go" ./

# 固定字符串 + 忽略大小写
go-search -F -i "error" --include="*.log" /var/log/

# 显示上下文（前 2 行 后 2 行）
go-search -C 2 "panic" ./

# 排除 vendor 目录
go-search "import" --exclude="vendor/**" --include="*.go" ./

# JSON 输出（方便管道处理）
go-search --json "TODO" ./ | jq '.matches[] | .file'

# 指定 worker 数量
go-search -j 16 "pattern" ./large-project/
```

### CLI Flag 定义

| Flag | 短写 | 说明 | 默认值 |
|------|------|------|--------|
| `--include` | | glob 包含模式（可多次指定） | 所有文件 |
| `--exclude` | | glob 排除模式 | 无 |
| `--ignore-case` | `-i` | 忽略大小写 | false |
| `--fixed-string` | `-F` | 固定字符串（非正则） | false |
| `--before-context` | `-B` | 显示匹配前 N 行 | 0 |
| `--after-context` | `-A` | 显示匹配后 N 行 | 0 |
| `--context` | `-C` | 显示匹配前后 N 行 | 0 |
| `--workers` | `-j` | worker 数量 | runtime.NumCPU() |
| `--json` | | JSON 格式输出 | false |
| `--no-color` | | 禁用彩色输出 | false |
| `--stats` | | 显示搜索统计 | false |
| `--no-gitignore` | | 不尊重 .gitignore | false |

---

## 七、开发计划（1 周）

| 天 | 任务 | 产出 |
|----|------|------|
| Day 1 | 项目初始化 + Cobra CLI 骨架 + 基本参数解析 + 单文件搜索（regexp） | `go-search "pattern" file.go` 可用 |
| Day 2 | Walker 目录遍历 + Worker Pool（goroutine + channel）+ 基础并发搜索 | `go-search "pattern" ./dir/` 可用 |
| Day 3 | glob 过滤（--include/--exclude）+ .gitignore 解析 + 二进制检测 | 过滤功能完整 |
| Day 4 | 彩色终端输出（文件名/行号/匹配高亮）+ 上下文行（-B/-A/-C） | 输出美观可用 |
| Day 5 | 搜索统计 + JSON 输出 + 大文件流式处理优化 | 功能完善 |
| Day 6 | 性能优化：buffer 大小调优、内存池（sync.Pool）、benchmark 测试 | 性能可量化 |
| Day 7 | README（架构图、用法、性能对比）+ goreleaser 配置 + 代码整理 | 项目可展示 |

---

## 八、面试话术

### "介绍一下你的 CLI 搜索工具？"

> 这是一个用 Go 写的并发文件内容搜索工具，类似简化版的 ripgrep。
> 核心设计是 **Fan-out/Fan-in 的三阶段 Pipeline**：
> - **Walker** goroutine 遍历目录，通过 channel 将文件路径发给下游
> - **N 个 Worker** goroutine 从 channel 读取文件路径，逐行搜索匹配
> - **Printer** goroutine 收集结果并彩色格式化输出
>
> Worker 数量默认等于 CPU 核心数，通过 `runtime.NumCPU()` 自动设置。
> 支持正则搜索、glob 过滤、.gitignore 规则、上下文行显示等。

### "为什么用 channel 而不是 sync.Mutex？"

> 这个场景是典型的**生产者-消费者模式**，channel 是最自然的选择：
> - Walker 生产文件路径，Workers 消费——用 channel 传递数据
> - Workers 生产搜索结果，Printer 消费——又一个 channel
> - 关闭 channel 就是天然的"结束信号"，`range ch` 自动退出
>
> 如果用 Mutex，就需要手动维护共享队列 + 条件变量通知，代码复杂度高很多。
> Go 的设计哲学是 **"Don't communicate by sharing memory; share memory by communicating"**。

### "Worker 数量怎么定的？设多了会怎样？"

> 默认 `runtime.NumCPU()`，因为文件搜索主要瓶颈是 **磁盘 IO + CPU 正则匹配**。
> - SSD 上，IO 不是瓶颈，CPU 核心数个 Worker 刚好打满 CPU
> - HDD 上，IO 是瓶颈，多 Worker 反而增加随机读，可以适当减少
> - 设太多的问题：goroutine 本身开销小（2KB 栈），但大量并发打开文件会消耗文件描述符，
>   所以 `fileCh` 的 buffer size 也起到了背压（backpressure）的作用

### "怎么处理大文件？"

> 用 `bufio.Scanner` 逐行流式读取，不会一次性把整个文件加载到内存。
> Scanner 默认 buffer 64KB，对于超长行（如 minified JS），我设置了 1MB 上限。
> 如果一行超过 1MB 直接跳过——这种场景基本是二进制文件误判。
>
> 另外对于二进制文件，我读取文件前 512 字节，调用 `http.DetectContentType` 判断 MIME 类型，
> 非 `text/*` 的直接跳过，避免无意义搜索。

### "这个工具和 ripgrep 比差距在哪？"

> 主要差距：
> 1. **正则引擎**：ripgrep 用 Rust 的 `regex` crate，支持 SIMD 加速的字面量匹配；Go 的 regexp 基于 RE2，不支持 SIMD
> 2. **内存映射**：ripgrep 用 `mmap` 读取文件，避免内核/用户态数据拷贝；我用 `bufio.Scanner`，多一次拷贝
> 3. **Unicode 处理**：ripgrep 的 Unicode 支持更完善
>
> 但作为学习项目，核心并发架构是一样的（worker pool + channel pipeline），
> 关键是理解 **Fan-out/Fan-in 模式、背压控制、context 取消传播** 这些并发编程核心概念。

---

## 九、简历怎么写

```
go-search — 高性能并发文件搜索 CLI 工具                                Go / Cobra / goroutine
- 基于 Fan-out/Fan-in 三阶段 Pipeline：Walker 遍历 → N Workers 并发搜索 → Printer 汇总输出
- Worker Pool 通过 channel 传递文件路径，goroutine 数量自适应 CPU 核心数，搜索万级文件耗时 <1s
- 支持正则/固定字符串搜索、glob 过滤、.gitignore 规则解析、二进制文件自动跳过
- bufio.Scanner 流式读取大文件，避免内存暴涨；sync.Mutex 保证终端彩色输出不交错
```