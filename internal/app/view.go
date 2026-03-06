package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/haivo/fileman/internal/types"
	"github.com/haivo/fileman/internal/ui"
)

// View 实现 tea.Model 接口，渲染整个 TUI 界面
func (m Model) View() tea.View {
	var content string

	// 窗口过小时显示提示
	if m.width < minWidth || m.height < minHeight {
		content = lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("窗口太小，请调整终端大小\n最小尺寸: 60×20")
	} else {
		// 构建完整界面
		frame := m.renderFrame()

		// 若弹窗可见，将弹窗覆盖在界面上
		if m.modal.IsVisible() {
			content = m.modal.Render()
		} else {
			content = frame
		}
	}

	// 创建 tea.View 并设置 AltScreen 和 MouseMode
	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// renderFrame 渲染主框架（边框 + Header + Content + Footer）
func (m Model) renderFrame() string {
	border := lipgloss.RoundedBorder()
	borderColor := ui.ColorBorderNormal
	if m.focus == types.FocusPanelA || m.focus == types.FocusPanelB {
		borderColor = ui.ColorBorderFocus
	}

	// 计算内部宽度
	// lipgloss v2 中 Width 包含 Border 和 Padding
	// 内容区域宽度 = Width - border(2) - padding(2) = (m.width - 2) - 4 = m.width - 6
	innerWidth := m.width - 6

	// --- Header 行 ---
	headerContent := m.header.RenderWithSize(m.selectionTotalSize)
	headerLine := strings.Repeat("─", innerWidth)

	// --- Content 区域（左右两栏）---
	contentStr := m.renderContent()

	// --- Footer 行 ---
	footerLine := strings.Repeat("─", innerWidth)
	footerContent := m.footer.Render()

	// 组合内部内容
	// 结构：Header + 分隔线 + Content + 分隔线 + Footer
	var sb strings.Builder
	sb.WriteString(headerContent)
	sb.WriteByte('\n')
	sb.WriteString(lipgloss.NewStyle().Foreground(borderColor).Render(headerLine))
	sb.WriteByte('\n')
	sb.WriteString(contentStr)
	sb.WriteByte('\n')
	sb.WriteString(lipgloss.NewStyle().Foreground(borderColor).Render(footerLine))
	sb.WriteByte('\n')
	sb.WriteString(footerContent)

	// 包裹外框边框，增加 padding 使内容与边框有间距
	// lipgloss v2 中 Width 包含 Border 和 Padding，Margin 在 Width 之外额外添加
	// 总宽度 = Width + MarginRight(1)
	// 我们希望总宽度 = m.width - 1 (留出安全距离)
	// 所以 Width = m.width - 2
	outerStyle := lipgloss.NewStyle().
		Border(border).
		BorderForeground(borderColor).
		Padding(1, 1).     // 上下左右各1个字符的间距
		MarginRight(1).    // 右侧留出1个字符的 margin
		Width(m.width - 2) // Width 包含 border 和 padding

	return outerStyle.Render(sb.String())
}

// renderContent 渲染内容区域（左栏双面板 + 右栏预览）
// 强制返回固定 m.contentHeight 行，确保 Header 和 Footer 位置不变
func (m Model) renderContent() string {
	// 左栏：PanelA（上）+ 分隔线 + PanelB（下）
	leftContent := m.renderLeftColumn()
	// 右栏：预览区
	rightContent := m.renderRightColumn()

	// 计算左右栏宽度（不含边框，因为已经在内容区域内）
	leftWidth := m.panelA.Width
	rightWidth := m.preview.Width

	// 垂直拼接左右两栏
	leftLines := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	// 移除末尾空行（split 会在末尾产生空行）
	if len(leftLines) > 0 && leftLines[len(leftLines)-1] == "" {
		leftLines = leftLines[:len(leftLines)-1]
	}
	if len(rightLines) > 0 && rightLines[len(rightLines)-1] == "" {
		rightLines = rightLines[:len(rightLines)-1]
	}

	// 强制使用固定的 contentHeight，确保布局稳定
	fixedHeight := m.contentHeight
	if fixedHeight <= 0 {
		fixedHeight = 1
	}

	// 确保两侧行数 = fixedHeight
	for len(leftLines) < fixedHeight {
		leftLines = append(leftLines, strings.Repeat(" ", leftWidth))
	}
	for len(rightLines) < fixedHeight {
		rightLines = append(rightLines, strings.Repeat(" ", rightWidth))
	}
	// 如果超出，截断
	if len(leftLines) > fixedHeight {
		leftLines = leftLines[:fixedHeight]
	}
	if len(rightLines) > fixedHeight {
		rightLines = rightLines[:fixedHeight]
	}

	var sb strings.Builder
	borderColor := ui.ColorBorderNormal

	for i := 0; i < fixedHeight; i++ {
		leftLine := leftLines[i]
		rightLine := rightLines[i]
		// 中间垂直分隔符
		sep := lipgloss.NewStyle().Foreground(borderColor).Render("│")
		sb.WriteString(leftLine)
		sb.WriteString(sep)
		sb.WriteString(rightLine)
		if i < fixedHeight-1 {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}

// renderLeftColumn 渲染左栏（PanelA + 分隔线 + PanelB）
func (m Model) renderLeftColumn() string {
	borderColor := ui.ColorBorderNormal

	// 渲染 PanelA 内容
	panelAContent := m.renderPanel(m.panelA)
	// 中间水平分隔线
	sepLine := strings.Repeat("─", m.panelA.Width)
	// 渲染 PanelB 内容
	panelBContent := m.renderPanel(m.panelB)

	var sb strings.Builder
	sb.WriteString(panelAContent)
	sb.WriteByte('\n')
	sb.WriteString(lipgloss.NewStyle().Foreground(borderColor).Render(sepLine))
	sb.WriteByte('\n')
	sb.WriteString(panelBContent)

	return sb.String()
}

// renderPanel 渲染单个文件面板内容区
func (m Model) renderPanel(p *ui.Panel) string {
	if p.Error != "" {
		errStyle := lipgloss.NewStyle().
			Foreground(ui.ColorError).
			Width(p.Width).
			Height(p.Height)
		return errStyle.Render("错误: " + p.Error)
	}

	return p.Render()
}

// renderRightColumn 渲染右栏（预览区）
func (m Model) renderRightColumn() string {
	return m.preview.Render()
}