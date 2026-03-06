# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 构建与测试命令

```bash
# 编译二进制文件
go build -o fm .

# 运行程序
./fm

# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./internal/fileops/...

# 运行单个测试
go test -run TestReadPreview_Chinese ./internal/fileops/
```

## 架构
得到的asd asda sd
FileMan 是一个终端文件管理器（TUI），基于 Bubble Tea 框架构建，采用 Elm 架构（Model-Update-View 模式）

### 核心模式

```
main.go → app.New() → tea.NewProgram(model).Run()
                     ↓
         Model.Init() → 初始命令（加载面板）
                     ↓
         Model.Update(msg) → (Model, Cmd)  // 处理所有事件
                     ↓
         Model.View() → string             // 渲染整个界面
```

### 包结构

```
internal/
├── app/           # 主应用逻辑
│   ├── model.go   # Model 结构体，包含所有状态、New()、Init()
│   ├── update.go  # Update() - 事件处理、按键分发
│   ├── view.go    # View() - 边框渲染、布局组合
│   └── keys.go    # 按键匹配函数（isUp、isDown 等）
│
├── ui/            # UI 组件（各有 Render 方法）
│   ├── panel.go   # 文件面板 - 光标、选择、搜索、滚动
│   ├── preview.go # 预览面板 - 带行号的文件内容
│   ├── modal.go   # 弹窗 - 输入/确认/进度/错误对话框
│   ├── header.go  # 顶栏 - 标题、选择计数、搜索
│   ├── footer.go  # 底栏 - 快捷键提示
│   └── theme.go   # 主题颜色和 lipgloss 样式
│
├── fileops/       # 文件操作
│   ├── scanner.go # ScanDir、文件排序、FormatSize/Date
│   ├── operations.go # 复制、移动、删除、重命名、创建
│   ├── preview.go # ReadPreview 及二进制检测
│   └── icons.go   # GetFileIcon（Nerd Font 映射）
│
└── types/         # 共享类型
    └── types.go   # FileEntry、SelectionSet、ModalType、FocusTarget
```

### 关键设计点

1. **双面板系统**：两个 `ui.Panel` 实例（panelA/panelB）共享一个 `SelectionSet`。Tab 键切换面板焦点。

2. **模式处理**（update.go）：应用有多种输入模式：
   - 普通模式：文件导航和操作
   - 搜索模式：`/` 触发，字符过滤列表
   - 编辑模式：内置 textarea 编辑文本文件
   - 弹窗模式：确认/输入/进度对话框

3. **消息类型**（model.go）：
   - `panelLoadMsg`：异步目录扫描结果
   - `fileOpMsg`：文件操作完成
   - `progressMsg`：复制/移动进度更新

4. **布局**（view.go）：固定布局，尺寸动态计算：
   - Header（1 行）+ 分隔线
   - 内容区：左栏（面板，40%）| 右栏（预览，60%）
   - 分隔线 + Footer（2 行）
   - 整体包裹圆角边框和内边距

5. **状态管理**：`Model` 持有所有状态；`activePanel()` 和 `inactivePanel()` 辅助函数返回焦点/非焦点面板。