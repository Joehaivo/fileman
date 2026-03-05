# FileMan

FileMan 是一个基于 Go 语言和 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 框架构建的现代化终端文件管理器（TUI）。它提供了双面板界面、实时文件预览以及高效的键盘导航，专为追求效率的开发者和极客设计。

## ✨ 功能特性
- **双面板界面**：经典的左右双栏设计，支持 Tab 键快速切换焦点，文件管理更高效。
- **实时预览**：支持文本文件内容预览（带行号），自动识别并展示文件详细信息。
- **图标支持**：集成 Nerd Fonts 图标支持，让终端界面更加美观直观。
- **文件操作**：支持复制、移动、删除、重命名、新建目录等常用操作，均配有直观的弹窗确认。
- **快速搜索**：支持当前目录下的模糊搜索与实时过滤。
- **多选系统**：支持空格键多选文件，方便批量操作。
- **鼠标支持**：支持鼠标点击选择、滚动列表和预览内容。
- **现代化 UI**：基于 Lip Gloss 构建的精美界面，自适应终端大小。

## 🛠️ 安装说明

### 环境要求

 推荐使用支持 [Nerd Font](https://www.nerdfonts.com/) 的终端字体以获得最佳体验。

### 使用 Go 安装

```bash
go install github.com/haivo/fileman@latest
```

### 源码编译

```bash
git clone https://github.com/haivo/fileman.git
cd fileman
go build -o fm .
```

## 🚀 使用指南

运行程序：

```bash
./fm
```

### ⌨️ 快捷键列表

| 按键 | 功能 |
|------|------|
| `↑` / `k` | 光标上移 |
| `↓` / `j` | 光标下移 |
| `PgUp` / `Ctrl+U` | 上翻页 |
| `PgDn` / `Ctrl+D` | 下翻页 |
| `Home` / `g` | 跳转至顶部 |
| `End` / `G` | 跳转至底部 |
| `Enter` | 打开目录 / 确认操作 |
| `Tab` | 切换左右面板焦点 |
| `Space` | 选中/取消选中当前文件 |
| `Ctrl+A` | 全选当前目录所有文件 |
| `/` | 进入搜索模式 |
| `Esc` | 退出搜索 / 取消弹窗 |
| `Del` | 删除文件（需确认） |
| `F2` | 重命名 |
| `Ctrl+N` | 新建目录 |
| `F5` | 复制选中项到另一面板 |
| `F6` | 移动选中项到另一面板 |
| `Ctrl+E` | 使用 `$EDITOR` 编辑文件 |
| `Ctrl+Q` | 退出程序 |

## 🏗️ 技术栈

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 强大的 TUI 框架
dd - [Lip Gloss](https://github.com/charmbracelet/lipgloss) - 界面样式定义
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI 组件库

## 📄 许可证

本项目采用 [MIT](LICENSE) 许可证。