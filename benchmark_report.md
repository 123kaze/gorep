# gorep / ripgrep / grep 性能测试报告

- **测试时间**：2026-04-28 22:59:11
- **测试目录**：`/home/123kaze/Project/claude-code-source-code-main`
- **文件数量**：约 `21667` 个普通文件
- **测试轮数**：每个 case 每个工具 `9` 轮
- **统计口径**：报告 `median` 为主要结论，`avg/min/max` 辅助参考

## 工具版本

- **rg**：`ripgrep 14.1.1`
- **grep**：`grep (GNU grep) 3.12`
- **gorep**：当前目录编译产物 `./gorep`（`A fast file search CLI for Go`）

## 测试方法

- **避免终端渲染干扰**：所有命令 stdout 都重定向到临时文件。
- **关闭颜色**：`rg --color never`，`gorep --no-color`。
- **二进制跳过**：`grep` 使用 `-I`；`rg` 默认跳过二进制；`gorep` 使用项目内置后缀过滤 + null byte 检测。
- **注意**：三者默认忽略规则和输出格式不同，因此输出字节数不完全一致。特别是 `rg` 默认尊重 `.gitignore`，`grep` 不尊重 `.gitignore`，`gorep` 使用当前项目实现的简化 `.gitignore`。
- **注意**：`gorep` 现在按文件分组输出，`rg/grep` 默认每条命中都带文件路径，所以输出大小不能直接代表匹配数量。

## 汇总结果

| Case | rg median | grep median | gorep median | 结论 |
|---|---:|---:|---:|---|
| 固定字符串：常见词 `function` | 0.0436s | 0.3404s | 0.1957s | `rg` 最快 |
| 固定字符串：较低频 `TODO` | 0.0314s | 0.2400s | 0.1461s | `rg` 最快 |
| 固定字符串：无匹配 | 0.0324s | 0.2426s | 0.1387s | `rg` 最快 |
| 正则：`function|class|const` | 0.0509s | 0.4846s | 0.9682s | `rg` 最快 |
| 正则：声明形式 `(function|class|const)\s+Identifier` | 0.0781s | 0.7447s | 1.3255s | `rg` 最快 |
| 固定字符串 + 只搜 TS/TSX：`import` | 0.0249s | 0.9552s | 0.1339s | `rg` 最快 |

## 详细结果

### 固定字符串：常见词 `function`

| 工具 | median | avg | min | max | 输出行数 | 输出大小 | 返回码 |
|---|---:|---:|---:|---:|---:|---:|---|
| `rg` | 0.0436s | 0.0442s | 0.0392s | 0.0484s | 25884 | 4516982 | [0] |
| `grep` | 0.3404s | 0.3425s | 0.3278s | 0.3584s | 41241 | 7572252 | [0] |
| `gorep` | 0.1957s | 0.1968s | 0.1906s | 0.2021s | 29254 | 2034366 | [0] |

命令：
- **rg**：`rg --no-heading --line-number --color never -F function /home/123kaze/Project/claude-code-source-code-main`
- **grep**：`grep -RInI -F function /home/123kaze/Project/claude-code-source-code-main`
- **gorep**：`./gorep --no-color -F function /home/123kaze/Project/claude-code-source-code-main`

### 固定字符串：较低频 `TODO`

| 工具 | median | avg | min | max | 输出行数 | 输出大小 | 返回码 |
|---|---:|---:|---:|---:|---:|---:|---|
| `rg` | 0.0314s | 0.0318s | 0.0304s | 0.0338s | 1913 | 362389 | [0] |
| `grep` | 0.2400s | 0.2409s | 0.2382s | 0.2453s | 3690 | 707273 | [0] |
| `gorep` | 0.1461s | 0.1465s | 0.1426s | 0.1517s | 2728 | 225713 | [0] |

命令：
- **rg**：`rg --no-heading --line-number --color never -F TODO /home/123kaze/Project/claude-code-source-code-main`
- **grep**：`grep -RInI -F TODO /home/123kaze/Project/claude-code-source-code-main`
- **gorep**：`./gorep --no-color -F TODO /home/123kaze/Project/claude-code-source-code-main`

### 固定字符串：无匹配

| 工具 | median | avg | min | max | 输出行数 | 输出大小 | 返回码 |
|---|---:|---:|---:|---:|---:|---:|---|
| `rg` | 0.0324s | 0.0328s | 0.0281s | 0.0428s | 0 | 0 | [1] |
| `grep` | 0.2426s | 0.2434s | 0.2363s | 0.2558s | 0 | 0 | [1] |
| `gorep` | 0.1387s | 0.1386s | 0.1272s | 0.1521s | 38 | 5074 | [0] |

命令：
- **rg**：`rg --no-heading --line-number --color never -F THIS_PATTERN_SHOULD_NOT_EXIST_123456 /home/123kaze/Project/claude-code-source-code-main`
- **grep**：`grep -RInI -F THIS_PATTERN_SHOULD_NOT_EXIST_123456 /home/123kaze/Project/claude-code-source-code-main`
- **gorep**：`./gorep --no-color -F THIS_PATTERN_SHOULD_NOT_EXIST_123456 /home/123kaze/Project/claude-code-source-code-main`

### 正则：`function|class|const`

| 工具 | median | avg | min | max | 输出行数 | 输出大小 | 返回码 |
|---|---:|---:|---:|---:|---:|---:|---|
| `rg` | 0.0509s | 0.0519s | 0.0493s | 0.0557s | 123678 | 20839156 | [0] |
| `grep` | 0.4846s | 0.4842s | 0.4775s | 0.4907s | 194187 | 34338149 | [0] |
| `gorep` | 0.9682s | 1.0164s | 0.8627s | 1.2355s | 127200 | 7899189 | [0] |

命令：
- **rg**：`rg --no-heading --line-number --color never function|class|const /home/123kaze/Project/claude-code-source-code-main`
- **grep**：`grep -RInI -E function|class|const /home/123kaze/Project/claude-code-source-code-main`
- **gorep**：`./gorep --no-color function|class|const /home/123kaze/Project/claude-code-source-code-main`

### 正则：声明形式 `(function|class|const)\s+Identifier`

| 工具 | median | avg | min | max | 输出行数 | 输出大小 | 返回码 |
|---|---:|---:|---:|---:|---:|---:|---|
| `rg` | 0.0781s | 0.0768s | 0.0598s | 0.0880s | 72755 | 11347938 | [0] |
| `grep` | 0.7447s | 0.7462s | 0.7269s | 0.7808s | 99323 | 16246360 | [0] |
| `gorep` | 1.3255s | 1.3318s | 1.3109s | 1.3821s | 75577 | 4557386 | [0] |

命令：
- **rg**：`rg --no-heading --line-number --color never (function|class|const)\s+[A-Za-z_][A-Za-z0-9_]* /home/123kaze/Project/claude-code-source-code-main`
- **grep**：`grep -RInI -E (function|class|const)[[:space:]]+[A-Za-z_][A-Za-z0-9_]* /home/123kaze/Project/claude-code-source-code-main`
- **gorep**：`./gorep --no-color (function|class|const)\s+[A-Za-z_][A-Za-z0-9_]* /home/123kaze/Project/claude-code-source-code-main`

### 固定字符串 + 只搜 TS/TSX：`import`

| 工具 | median | avg | min | max | 输出行数 | 输出大小 | 返回码 |
|---|---:|---:|---:|---:|---:|---:|---|
| `rg` | 0.0249s | 0.0254s | 0.0215s | 0.0290s | 17802 | 2606141 | [0] |
| `grep` | 0.9552s | 0.9721s | 0.9006s | 1.2311s | 17802 | 2606141 | [0] |
| `gorep` | 0.1339s | 0.1330s | 0.1276s | 0.1382s | 17858 | 1087962 | [0] |

命令：
- **rg**：`rg --no-heading --line-number --color never -F -g *.ts -g *.tsx import /home/123kaze/Project/claude-code-source-code-main`
- **grep**：`bash -lc find /home/123kaze/Project/claude-code-source-code-main -type f \( -name '*.ts' -o -name '*.tsx' \) -print0 | xargs -0 grep -InI -F import`
- **gorep**：`./gorep --no-color --include *.ts --include *.tsx -F import /home/123kaze/Project/claude-code-source-code-main`

## 结论

1. **ripgrep (`rg`) 全面最快**。在所有测试中，`rg` 的中位耗时都显著低于 `grep` 和 `gorep`，尤其在正则和 glob 场景下优势明显。
2. **`gorep` 固定字符串搜索已经明显快于 `grep` 的全目录递归搜索**。例如 `function` case 中，`gorep` median 约 `0.19s`，`grep` 约 `0.33s`。
3. **`gorep` 正则模式明显慢于 `rg`，并且部分 case 慢于 `grep`**。主要原因是 Go 标准库 regexp、逐行 string 处理、结果结构构造和高亮 ranges 计算都会增加开销。
4. **`gorep` 的输出格式更紧凑但会影响横向比较**。它按文件分组输出，而 `rg/grep` 每条命中包含完整文件路径，因此输出字节数不同。
5. **`gorep` 当前最值得优化的是匹配层，而不是并发 walker**：固定字符串可增加 `scanner.Bytes()` 预筛/字节级匹配；正则可考虑避免重复分配、减少输出结构构造成本。

## 后续优化建议

- **固定字符串路径**：增加 `MatchBytes`，先用 `scanner.Bytes()` 做无分配预筛，命中后再转 string 计算高亮。
- **Scanner buffer**：设置 `scanner.Buffer(make([]byte, 0, 256*1024), 1024*1024)`，避免长行报错。
- **输出控制 benchmark**：增加 `--count` 或 `--quiet` 模式，隔离搜索耗时和打印耗时。
- **统计功能**：实现 `--stats`，直接输出扫描文件数、匹配文件数、匹配行数和耗时，方便未来持续 benchmark。
