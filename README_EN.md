# FileMan

<p align="center">
  <strong>A sleek, lightweight modern terminal file manager</strong>
</p>

<p align="center">
  Built with Go + Bubble Tea
</p>

<p align="center">
  <img src="docs/fileman-intro.gif" alt="FileMan Demo" width="600" />
</p>

<p align="center">
  <a href="README.md">简体中文</a> | English
</p>

---

## ✨ Features

- **Dual Panel Interface** — Dual panel design with Tab key to quickly switch focus
- **Real-time Preview** — Text file content preview with automatic file type detection
- **File Operations** — Quickly copy/move files between panels, plus delete, rename, create file/directory
- **Quick Search** — Real-time file filtering in current directory
- **Built-in Editor** — Text file editor included
- **Mouse Support** — Click to select, scroll to browse
- **Adaptive Layout** — Automatically adapts to terminal window size

## 📦 Installation

### One-line Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/Joehaivo/fileman/main/install.sh | bash
```

## 🚀 Usage

```bash
fm
```

Check version:

```bash
fm --version
```

## ⌨️ Keyboard Shortcuts

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move cursor up/down |
| `PgUp` / `PgDn` | Page up/down |
| `Home` / `End` | Jump to top/bottom |
| `←` | Go to parent directory |
| `→` / `Enter` | Enter directory or edit file |
| `Tab` | Switch between panels |

### File Operations

| Key | Action |
|-----|--------|
| `F1` | Rename |
| `F2` | Copy to other panel |
| `F3` | Move to other panel |
| `F4` | Create directory |
| `F5` | Create file |
| `F6` | Open in external editor |
| `F7` | Show/hide hidden files |
| `F8` | Settings |
| `F9` | Exit |
| `Del` | Delete |
| `/` | Search |
| `Esc` | Cancel search/modal |

### Edit Mode

| Key | Action |
|-----|--------|
| `↑` `↓` `←` `→` | Move cursor |
| `F1` | Save |
| `F2` | Exit edit mode |
| `Home` / `End` | Beginning/end of line |
| `PgUp` / `PgDn` | Page up/down |

### Build from Source

```bash
git clone https://github.com/Joehaivo/fileman.git
cd fileman
go build -ldflags "-s -w -X main.version=$(git describe --tags --always)" -o fm .
```

## 🛠️ Tech Stack

- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Style engine
- [Bubbles](https://github.com/charmbracelet/bubbles) — UI component library

## 📄 License

[MIT](LICENSE)