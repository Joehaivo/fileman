package app

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"github.com/haivo/fileman/internal/fileops"
	"github.com/haivo/fileman/internal/types"
	"github.com/haivo/fileman/internal/ui"
)

const (
	// 最小终端宽度
	minWidth = 60
	// 最小终端高度
	minHeight = 20
	// Header 固定高度（行数）
	headerHeight = 1
	// Footer 固定高度（行数）
	footerHeight = 2
	// 外框边框占用（上下各1）
	borderVertical = 2
	// 左右栏比例（左栏占 40%）
	leftRatio = 2
	leftDenom = 5
)

// fileOpMsg 文件操作完成消息
type fileOpMsg struct {
	err error
}

// progressMsg 进度更新消息
type progressMsg struct {
	info types.ProgressInfo
}

// Model 主应用 Model，包含所有子组件状态
type Model struct {
	// 终端尺寸
	width  int
	height int

	// 内容区域高度（用于固定布局）
	contentHeight int

	// 子组件
	header   *ui.Header
	footer   *ui.Footer
	panelA   *ui.Panel
	panelB   *ui.Panel
	preview  *ui.PreviewPane
	modal    *ui.Modal

	// 焦点状态
	focus types.FocusTarget

	// 全局选择集（两个面板共享）
	selection types.SelectionSet

	// 搜索输入（由 panel 内部管理，此处同步状态）
	isSearching bool
	searchQuery string

	// 选中文件总大小（Header 显示用）
	selectionTotalSize int64

	// 应用设置
	settings types.Settings

	// 待处理的复制/移动目标路径
	pendingOpSrc []string
	pendingOpDst string

	// 初始命令（Init 执行后清空）
	initCmd tea.Cmd
}

// New 创建并初始化 Model
func New() (Model, tea.Cmd) {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = cwd
	}

	selection := make(types.SelectionSet)
	settings := types.Settings{
		ShowDate: false,
	}

	header := ui.NewHeader(selection)
	footer := ui.NewFooter()
	panelA := ui.NewPanel(cwd, selection)
	panelA.ShowDate = settings.ShowDate
	panelB := ui.NewPanel(homeDir, selection)
	panelB.ShowDate = settings.ShowDate
	preview := ui.NewPreviewPane()
	modal := ui.NewModal()

	m := Model{
		header:    header,
		footer:    footer,
		panelA:    panelA,
		panelB:    panelB,
		preview:   preview,
		modal:     modal,
		focus:     types.FocusPanelA,
		selection: selection,
		settings:  settings,
	}

	// 初始加载两个面板的目录内容，保存为初始命令在 Init 中执行
	m.initCmd = tea.Batch(
		m.loadPanel(panelA),
		m.loadPanel(panelB),
	)

	return m, nil
}

// Init 实现 tea.Model 接口，返回初始命令（加载两个面板目录）
func (m Model) Init() tea.Cmd {
	return m.initCmd
}

// loadPanel 加载面板目录内容的命令
func (m *Model) loadPanel(p *ui.Panel) tea.Cmd {
	path := p.Path
	return func() tea.Msg {
		entries, err := fileops.ScanDir(path)
		if err != nil {
			return panelLoadMsg{panel: p, err: err}
		}
		return panelLoadMsg{panel: p, entries: entries}
	}
}

// panelLoadMsg 目录加载完成消息
type panelLoadMsg struct {
	panel   *ui.Panel
	entries []types.FileEntry
	err     error
}

// activePanel 返回当前焦点面板
func (m *Model) activePanel() *ui.Panel {
	if m.focus == types.FocusPanelA {
		return m.panelA
	}
	return m.panelB
}

// inactivePanel 返回非焦点面板
func (m *Model) inactivePanel() *ui.Panel {
	if m.focus == types.FocusPanelA {
		return m.panelB
	}
	return m.panelA
}

// applyPanelLoad 将加载结果应用到面板
func applyPanelLoad(p *ui.Panel, msg panelLoadMsg) {
	if msg.err != nil {
		p.Error = msg.err.Error()
		p.Entries = nil
		p.Filtered = nil
		return
	}
	p.Error = ""
	p.Entries = msg.entries
	p.Filtered = nil
	// 重置搜索过滤
	if p.IsSearching {
		p.SetSearch(p.SearchQuery)
	}

	// 检查是否有待选中的文件
	if p.PendingSelect != "" {
		p.SetCursorByName(p.PendingSelect)
		p.PendingSelect = ""
	}
}

// updatePreview 根据当前焦点面板更新预览
func (m *Model) updatePreview() {
	entry := m.activePanel().CurrentEntry()
	m.preview.SetEntry(entry)
}

// calcSizes 计算并应用所有组件的尺寸
func (m *Model) calcSizes() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	// 可用内容高度 = 终端高度 - 外框上下边框(2) - 上下padding(2) - 内部水平分隔线
	// 布局：上框(1) + 上padding(1) + header(1) + 水平分隔(1) + content + 水平分隔(1) + footer(2) + 下padding(1) + 下框(1) = 9行固定
	m.contentHeight = m.height - 9
	if m.contentHeight < 1 {
		m.contentHeight = 1
	}
	contentHeight := m.contentHeight

	// 计算可用内容宽度（去掉边框和padding）
	// 外框总宽度 = m.width，边框左右各1列，padding左右各1列
	contentWidth := m.width - 4 // -2 for 左右边框，-2 for 左右padding

	// 左栏宽度（40%），右栏宽度（60%）
	// 布局：leftWidth + 垂直分隔(1) + rightWidth = contentWidth
	leftWidth := contentWidth * leftRatio / leftDenom
	rightWidth := contentWidth - leftWidth - 1 // -1 for 垂直分隔符
	if rightWidth < 10 {
		rightWidth = 10
	}
	if leftWidth < 10 {
		leftWidth = 10
	}

	// 面板高度：contentHeight / 2，考虑中间水平分隔线
	panelHeight := (contentHeight - 1) / 2 // -1 for 中间水平分隔线

	// 设置各组件尺寸
	m.header.SetSize(contentWidth) // 使用内容宽度（已扣除边框和padding）
	m.footer.SetSize(contentWidth)
	m.panelA.SetSize(leftWidth, panelHeight)
	m.panelB.SetSize(leftWidth, panelHeight)
	m.preview.SetSize(rightWidth, contentHeight)
	m.modal.SetSize(m.width, m.height)

	// 设置面板焦点状态
	m.panelA.IsFocused = m.focus == types.FocusPanelA
	m.panelB.IsFocused = m.focus == types.FocusPanelB
}

// focusedPanelPath 返回当前焦点面板路径（用于文件操作目标）
func (m *Model) focusedPanelPath() string {
	return m.activePanel().Path
}

// otherPanelPath 返回非焦点面板路径（用于复制/移动的目标路径）
func (m *Model) otherPanelPath() string {
	return m.inactivePanel().Path
}

// getSelectedOrCurrent 获取选择集中的文件，若无选择则返回当前条目
func (m *Model) getSelectedOrCurrent() ([]types.FileEntry, bool) {
	panel := m.activePanel()

	if m.selection.Len() > 0 {
		// 收集选中的文件条目
		var entries []types.FileEntry
		for _, e := range panel.Entries {
			if m.selection.Has(e.Path) {
				entries = append(entries, e)
			}
		}
		return entries, len(entries) > 0
	}

	entry := panel.CurrentEntry()
	if entry == nil {
		return nil, false
	}
	return []types.FileEntry{*entry}, true
}

// computeSelectionSize 计算所有选中文件的总大小
func (m *Model) computeSelectionSize() int64 {
	var total int64
	panel := m.activePanel()
	for _, e := range panel.Entries {
		if m.selection.Has(e.Path) {
			total += e.Size
		}
	}
	return total
}

// navigateTo 在当前焦点面板中导航到指定路径
func (m *Model) navigateTo(path string) tea.Cmd {
	panel := m.activePanel()
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil
	}
	panel.Path = absPath
	panel.Cursor = 0
	panel.Offset = 0
	panel.SearchQuery = ""
	panel.IsSearching = false
	panel.Filtered = nil
	return m.loadPanel(panel)
}
