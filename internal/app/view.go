package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/haivo/fileman/internal/ui"
)

// View 实现 tea.Model 接口，渲染整个 TUI 界面
func (m Model) View() string {
	// 窗口过小时显示提示
	if m.width < minWidth || m.height < minHeight {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("窗口太小，请调整终端大小\n最小尺寸: 60×20")
	}

	// 构建完整界面
	frame := m.renderFrame()

	// 若弹窗可见，将弹窗覆盖在界面上
	if m.modal.IsVisible() {
		overlay := m.modal.Render()
		return overlay
	}

	return frame
}

// renderFrame 渲染主框架（边框 + Header + Content + Footer）
func (m Model) renderFrame() string {
	border := lipgloss.RoundedBorder()
	borderColor := ui.ColorBorderNormal

	// 计算内部宽度（去掉左右边框各1列）
	innerWidth := m.width - 2

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

	// 包裹外框边框
	outerStyle := lipgloss.NewStyle().
		Border(border).
		BorderForeground(borderColor).
		Width(innerWidth)

	return outerStyle.Render(sb.String())
}

// renderContent 渲染内容区域（左栏双面板 + 右栏预览）
func (m Model) renderContent() string {
	// 左栏：PanelA（上）+ 分隔线 + PanelB（下）
	leftContent := m.renderLeftColumn()
	// 右栏：预览区
	rightContent := m.renderRightColumn()

	// 计算左右栏宽度（含边框）
	leftWidth := m.panelA.Width + 2 // +2 for 左右边框
	rightWidth := m.preview.Width + 2

	// 垂直拼接左右两栏
	leftLines := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	// 确保两侧行数一致（取最多的，填充空行）
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for len(leftLines) < maxLines {
		leftLines = append(leftLines, strings.Repeat(" ", leftWidth))
	}
	for len(rightLines) < maxLines {
		rightLines = append(rightLines, strings.Repeat(" ", rightWidth))
	}

	var sb strings.Builder
	borderColor := ui.ColorBorderNormal

	for i, leftLine := range leftLines {
		rightLine := ""
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}
		// 中间垂直分隔符
		sep := lipgloss.NewStyle().Foreground(borderColor).Render("│")
		sb.WriteString(leftLine)
		sb.WriteString(sep)
		sb.WriteString(rightLine)
		if i < len(leftLines)-1 {
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
