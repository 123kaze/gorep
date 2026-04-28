# gorep - A Fast File Search CLI in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`gorep` is a high-performance command-line tool for searching files, built with Go. It aims to provide a fast and efficient alternative to traditional tools like `grep` and `ripgrep`, with a focus on concurrent file processing and user-friendly output.

## Features

- **Concurrent Search**: Utilizes Go's goroutines for parallel file processing with a worker pool model.
- **Colored Output**: Highlights matches and file paths for better readability (can be disabled with `--no-color`).
- **Context Lines**: Supports showing lines before (`-B`), after (`-A`), or around (`-C`) matches, with deduplicated context blocks similar to `ripgrep`.
- **Binary File Detection**: Automatically skips binary files using extension-based filtering and null byte detection.
- **Glob Filtering**: Include or exclude files based on glob patterns (`--include` and `--exclude`).
- **Gitignore Support**: Respects `.gitignore` rules by default (can be disabled with `--no-gitignore`).
- **Fixed String and Regex Modes**: Search with literal strings (`-F`) or regular expressions.

## Installation

### From Source

```bash
go install github.com/yourusername/gorep@latest
```

(Replace `yourusername` with your GitHub username once uploaded.)

### Using Pre-built Binary

Once releases are available on GitHub, you can download pre-built binaries from the [releases page](https://github.com/yourusername/gorep/releases).

## Usage

```bash
# Basic search
gorep "pattern" /path/to/search

# Fixed string search (literal match)
gorep -F "function name" .

# Case-insensitive search
gorep -i "function" .

# Show 2 lines of context before and after matches
gorep -C 2 "function" .

# Include only specific file types
gorep --include "*.go" --include "*.py" "function" .

# Exclude certain directories or file types
gorep --exclude "vendor/*" --exclude "*.log" "function" .

# Ignore .gitignore rules
gorep --no-gitignore "function" .

# Disable colored output
gorep --no-color "function" .
```

## Performance

Benchmark results on a directory with ~21,667 files (`/home/123kaze/Project/claude-code-source-code-main`):

| Case | rg median | grep median | gorep median | Conclusion |
|------|-----------|-------------|--------------|------------|
| Fixed String: `function` | 0.0436s | 0.3404s | 0.1957s | `rg` fastest |
| Fixed String: `TODO` | 0.0314s | 0.2400s | 0.1461s | `rg` fastest |
| Fixed String: No match | 0.0324s | 0.2426s | 0.1387s | `rg` fastest |
| Regex: `function|class|const` | 0.0509s | 0.4846s | 0.9682s | `rg` fastest |
| Regex: `(function|class|const)\s+Identifier` | 0.0781s | 0.7447s | 1.3255s | `rg` fastest |
| Fixed String + TS/TSX only: `import` | 0.0249s | 0.9552s | 0.1339s | `rg` fastest |

**Key Takeaways**:
- `gorep` outperforms `grep` in fixed-string searches.
- `gorep`'s regex performance needs optimization compared to `rg` and sometimes `grep`.
- Full report: [benchmark_report.md](benchmark_report.md)

## Development

```bash
git clone https://github.com/yourusername/gorep.git
cd gorep
go build .
```

## Roadmap

- [x] Concurrent file walking and searching
- [x] Colored terminal output with match highlighting
- [x] Context line display (`-B`, `-A`, `-C`)
- [x] Binary file detection and skipping
- [x] Glob pattern include/exclude
- [x] `.gitignore` support
- [ ] Search statistics (`--stats`)
- [ ] JSON output format (`--json`)
- [ ] Performance optimizations (buffer tuning, byte-level matching)
- [ ] Benchmark comparisons with `grep`/`ripgrep`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
