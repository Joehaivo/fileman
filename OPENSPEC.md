# FileMan TUI 开发进度追踪

## 项目信息
ddddasd收拾收拾
- **二进制名**: `fm`
- **模块路径**: `github.com/haivo/fileman`
- **Go 版本**: 1.22+
- **TUI 框架**: Bubble Tea v1.3.10 + Lip Gloss v1.1.0 + Bubbles v1.0.0

---

## Phase 进度

### ✅ Phase 1: 项目脚手架与基础框架
- [x] 初始化 Go module (`github.com/haivo/fileman`)
- [x] 安装依赖：bubbletea v1.3.10, lipgloss v1.1.0, bubbles v1.0.0
- [x] 搭建 Bubble Tea 基础骨架 (Model/Init/Update/View)
- [x] 实现 AltScreen 全屏模式
- [x] 实现 `tea.WindowSizeMsg` 自适应布局
- [x] 实现 Header/Content/Footer 三段布局框架
- [x] 定义主题样式 (`ui/theme.go`)：颜色常量、边框样式、焦点/非焦点样式
- [x] 创建 OPENSPEC.md

### ✅ Phase 2: 文件面板核心
- [x] 实现 `FileEntry` 类型 (`types/types.go`)
- [x] 实现目录扫描 (`fileops/scanner.go`)：读取目录、排序（目录优先、按名称）
- [x] 实现 Nerd Font 图标映射 (`fileops/icons.go`)：50+ 扩展名映射
- [x] 实现单个文件面板组件 (`ui/panel.go`)：
  - [x] 文件列表渲染：图标 + 文件名 + 大小 + 日期
  - [x] 第一行固定 `..`（上级目录）
  - [x] 光标移动（上/下/PageUp/PageDown/Home/End）
  - [x] 列表滚动（超出可视区域时）
  - [x] Enter 打开目录 / `..` 返回上级

### ✅ Phase 3: 双面板 + 焦点系统
- [x] Content 区域左侧放置 PanelA（上）+ PanelB（下）
- [x] 实现 Tab 焦点切换
- [x] 焦点面板/非焦点面板文字样式区分
- [x] 两个面板独立维护各自路径和光标状态

### ✅ Phase 4: 文件预览
- [x] 实现文件预览读取 (`fileops/preview.go`)：限制最大 1MB
- [x] 预览区显示：文件名（标题）+ 带行号内容 + 底部文件信息
- [x] 文件信息：类型、大小、修改时间、权限、行数/滚动进度
- [x] 对非文本/二进制文件显示基本信息
- [x] 支持预览区上下滚动

### ✅ Phase 5: 多选系统
- [x] 实现 `SelectionSet` 数据结构（map-based O(1) 查找）
- [x] Space 键 toggle 选择当前项
- [x] 已选文件高亮显示（橙色）
- [x] Ctrl+A 全选当前目录
- [x] Header 实时显示：`已选: N 个 (size)`
- [x] 大小自动格式化（B/KB/MB/GB

### ✅ Phase 6: 搜索模式
- [x] `/` 进入搜索模式
- [x] Header 切换显示：`搜索: keyword`
- [x] 实时过滤当前面板文件列表（模糊匹配，不区分大小写）
- [x] Esc 退出搜索模式，恢复完整列表
- [x] 搜索模式下方向键仍可移动光标
- [x] 搜索模式下 Enter 确认并打开选中项
- [x] Footer 搜索模式提示

### ✅ Phase 7: 模态弹窗系统
- [x] 实现统一弹窗组件 (`ui/modal.go`)：居中渲染、圆角边框
- [x] **输入型**：新建目录（Ctrl+N）、重命名（F2）
- [x] **确认型**：删除确认（Del）
- [x] **进度型**：复制/移动进度显示
- [x] **错误型**：操作失败提示
- [x] Enter 确认，Esc 取消

### ✅ Phase 8: 文件操作
- [x] **删除**：支持单个/批量，弹窗确认
- [x] **重命名**：F2 弹出输入框，预填原文件名
- [x] **新建目录**：Ctrl+N 弹出输入框
- [x] **复制**（F5）：从焦点面板复制到另一面板路径
- [x] **移动**（F6）：从焦点面板移动到另一面板路径
- [x] **编辑**（Ctrl+E）：调用 `$EDITOR` 打开文件（`tea.ExecProcess`）
- [x] 操作完成后刷新两个面板并清空选择集

### 🔲 Phase 9: 进度条与动画
- [x] 复制/移动时显示进度弹窗
- [ ] 实时进度百分比更新（`tea.Tick` + goroutine channel）
- [ ] 光标移动 ease-out 视觉效果

### 🔲 Phase 10: 性能优化与收尾
- [x] 列表渲染仅渲染可视区域（O(n) 虚拟滚动）
- [x] 样式对象复用（`DefaultTheme` 单例）
- [ ] 大目录延迟加载提示
- [ ] 预览大文件分块读取改进
- [x] Ctrl+Q 优雅退出

---

## 文件结构

```
fileman/
├── main.go                        ✅ 入口，初始化 tea.Program（AltScreen）
├── fm                             ✅ 编译后的二进制文件
├── go.mod / go.sum                ✅
├── OPENSPEC.md                    ✅ 本文件
├── internal/
│   ├── app/
│   │   ├── model.go               ✅ 主 Model 定义 + 组件初始化
│   │   ├── update.go              ✅ Update 逻辑 + 按键分发
│   │   ├── view.go                ✅ View：边框布局 + 左右栏组合
│   │   └── keys.go                ✅ 快捷键判断函数
│   ├── ui/
│   │   ├── theme.go               ✅ 主题颜色常量 + 样式缓存
│   │   ├── header.go              ✅ Header 组件（标题/版本/选择统计/搜索）
│   │   ├── footer.go              ✅ Footer 组件（两行快捷键提示）
│   │   ├── panel.go               ✅ 文件面板组件（列表/光标/选择/搜索）
│   │   ├── preview.go             ✅ 右侧预览组件（行号/内容/文件信息）
│   │   └── modal.go               ✅ 模态弹窗系统
│   ├── fileops/
│   │   ├── scanner.go             ✅ 目录扫描 + 排序 + 格式化工具
│   │   ├── operations.go          ✅ 复制/移动/删除/重命名/新建
│   │   ├── preview.go             ✅ 文件预览读取（限制 1MB）
│   │   └── icons.go               ✅ Nerd Font 图标映射（50+ 扩展名）
│   └── types/
│       └── types.go               ✅ 共享类型：FileEntry/SelectionSet/FocusTarget/ModalType
```

---

## 快捷键速查

| 按键 | 功能 |
|------|------|
| ↑ / k | 光标上移 |
| ↓ / j | 光标下移 |
| PgUp / Ctrl+U | 上翻页 |
| PgDn / Ctrl+D | 下翻页 |
| Home / g | 移至顶部 |
| End / G | 移至底部 |
| Enter | 打开目录 / 确认操作 |
| Tab | 切换焦点面板 |
| Space | 多选 toggle |
| Ctrl+A | 全选当前目录 |
| / | 进入搜索模式 |
| Esc | 退出搜索 / 取消弹窗 |
| Del | 删除（需确认） |
| F2 | 重命名 |
| Ctrl+N | 新建目录 |
| F5 | 复制到另一面板 |
| F6 | 移动到另一面板 |
| Ctrl+E | 用 $EDITOR 打开 |
| Ctrl+Q | 退出程序 |