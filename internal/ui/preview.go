package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/haivo/fileman/internal/fileops"
	"github.com/haivo/fileman/internal/types"
)

// PreviewPane 右侧预览组件
type PreviewPane struct {
	Entry   *types.FileEntry          // 当前预览的文件条目
	Result  *fileops.PreviewResult    // 预览读取结果
	Width   int                       // 预览区宽度（不含边框）
	Height  int                       // 预览区高度（不含边框）
	ScrollY int                       // 垂直滚动偏移（行数）
}

// NewPreviewPane 创建新的预览组件
func NewPreviewPane() *PreviewPane {
	return &PreviewPane{}
}

// SetSize 设置预览区尺寸
func (pv *PreviewPane) SetSize(width, height int) {
	pv.Width = width
	pv.Height = height
	pv.clampScroll()
}

// SetEntry 设置要预览的文件条目，并重新加载预览内容
func (pv *PreviewPane) SetEntry(entry *types.FileEntry) {
	pv.Entry = entry
	pv.ScrollY = 0
	pv.Result = nil

	if entry == nil || entry.IsDir {
		return
	}

	pv.Result = fileops.ReadPreview(*entry)
}

// ScrollUp 预览内容向上滚动
func (pv *PreviewPane) ScrollUp() {
	if pv.ScrollY > 0 {
		pv.ScrollY--
	}
}

// ScrollDown 预览内容向下滚动
func (pv *PreviewPane) ScrollDown() {
	pv.ScrollY++
	pv.clampScroll()
}

// ScrollPageUp 向上翻页
func (pv *PreviewPane) ScrollPageUp() {
	pv.ScrollY -= pv.contentHeight()
	if pv.ScrollY < 0 {
		pv.ScrollY = 0
	}
}

// ScrollPageDown 向下翻页
func (pv *PreviewPane) ScrollPageDown() {
	pv.ScrollY += pv.contentHeight()
	pv.clampScroll()
}

// contentHeight 计算内容区高度（减去标题行和信息区）
func (pv *PreviewPane) contentHeight() int {
	// 预览区高度 - 标题1行 - 分隔线1行 - 信息区5行
	h := pv.Height - 1 - 1 - 5
	if h < 1 {
		h = 1
	}
	return h
}

// clampScroll 限制滚动范围
func (pv *PreviewPane) clampScroll() {
	if pv.Result == nil || len(pv.Result.Lines) == 0 {
		pv.ScrollY = 0
		return
	}

	maxScroll := len(pv.Result.Lines) - pv.contentHeight()
	if maxScroll < 0 {
		maxScroll = 0
	}
	if pv.ScrollY > maxScroll {
		pv.ScrollY = maxScroll
	}
}

// Render 渲染预览内容
func (pv *PreviewPane) Render() string {
	if pv.Width <= 0 || pv.Height <= 0 {
		return ""
	}

	if pv.Entry == nil {
		return pv.renderEmpty()
	}

	var sb strings.Builder

	// 标题行：文件名
	titleLine := pv.renderTitle()
	sb.WriteString(titleLine)
	sb.WriteByte('\n')

	// 内容区
	contentH := pv.contentHeight()
	contentLines := pv.renderContent(contentH)
	sb.WriteString(contentLines)

	// 分隔线
	sep := strings.Repeat("─", pv.Width)
	sb.WriteString(DefaultTheme.SubduedStyle.Render(sep))
	sb.WriteByte('\n')

	// 文件信息区（5行）
	infoLines := pv.renderInfo()
	sb.WriteString(infoLines)

	return sb.String()
}

// renderEmpty 渲染空状态（无文件选中）
func (pv *PreviewPane) renderEmpty() string {
	msg := "选择文件以预览"
	style := DefaultTheme.SubduedStyle
	centered := lipgloss.NewStyle().
		Width(pv.Width).
		Height(pv.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(msg)
	return style.Render(centered)
}

// renderTitle 渲染标题行（文件名）
func (pv *PreviewPane) renderTitle() string {
	name := pv.Entry.Name
	if len(name) > pv.Width-2 {
		runes := []rune(name)
		if len(runes) > pv.Width-2 {
			name = string(runes[:pv.Width-3]) + "…"
		}
	}
	return DefaultTheme.PreviewTitle.Width(pv.Width).Render(name)
}

// renderContent 渲染文件内容区
func (pv *PreviewPane) renderContent(height int) string {
	if pv.Entry.IsDir {
		return pv.renderDirPreview(height)
	}

	if pv.Result == nil {
		return pv.renderPlaceholder(height, "加载中...")
	}

	if pv.Result.Error != "" && !pv.Result.IsBinary && !pv.Result.IsTooLarge {
		return pv.renderPlaceholder(height, "错误: "+pv.Result.Error)
	}

	if pv.Result.IsBinary {
		return pv.renderPlaceholder(height, "二进制文件，无法预览")
	}

	if pv.Result.IsTooLarge {
		return pv.renderPlaceholder(height, pv.Result.Error)
	}

	return pv.renderTextContent(height)
}

// renderPlaceholder 渲染占位符文字
func (pv *PreviewPane) renderPlaceholder(height int, msg string) string {
	lines := make([]string, height)
	lines[0] = DefaultTheme.SubduedStyle.Width(pv.Width).Render(msg)
	for i := 1; i < height; i++ {
		lines[i] = strings.Repeat(" ", pv.Width)
	}
	return strings.Join(lines, "\n") + "\n"
}

// renderDirPreview 渲染目录预览（显示目录信息）
func (pv *PreviewPane) renderDirPreview(height int) string {
	lines := make([]string, height)
	lines[0] = DefaultTheme.SubduedStyle.Width(pv.Width).Render("📁 目录")
	for i := 1; i < height; i++ {
		lines[i] = strings.Repeat(" ", pv.Width)
	}
	return strings.Join(lines, "\n") + "\n"
}

// renderTextContent 渲染文本内容
func (pv *PreviewPane) renderTextContent(height int) string {
	lines := pv.Result.Lines
	total := len(lines)

	var sb strings.Builder
	lineNumWidth := len(fmt.Sprintf("%d", total))
	if lineNumWidth < 3 {
		lineNumWidth = 3
	}

	// 内容可用宽度（减去行号区域）
	contentWidth := pv.Width - lineNumWidth - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	rendered := 0
	for i := pv.ScrollY; i < total && rendered < height; i++ {
		lineNum := fmt.Sprintf("%*d", lineNumWidth, i+1)
		lineNumStr := DefaultTheme.SubduedStyle.Render(lineNum + " ")

		content := lines[i]
		
		// 使用 lipgloss 自动折行
		wrappedContent := lipgloss.NewStyle().Width(contentWidth).Render(content)
		subLines := strings.Split(wrappedContent, "\n")

		for j, subLine := range subLines {
			if rendered >= height {
				break
			}

			if j == 0 {
				// 第一行显示行号
				sb.WriteString(lineNumStr)
			} else {
				// 后续行显示空白占位
				padding := strings.Repeat(" ", lineNumWidth+1) // +1 是为了匹配 lineNumStr 中的空格
				sb.WriteString(DefaultTheme.SubduedStyle.Render(padding))
			}

			sb.WriteString(subLine)
			sb.WriteByte('\n')
			rendered++
		}
	}

	// 填充空行
	for rendered < height {
		sb.WriteString(strings.Repeat(" ", pv.Width))
		sb.WriteByte('\n')
		rendered++
	}

	return sb.String()
}

// renderInfo 渲染文件信息区（5行）
func (pv *PreviewPane) renderInfo() string {
	if pv.Entry == nil {
		return strings.Repeat(strings.Repeat(" ", pv.Width)+"\n", 5)
	}

	entry := pv.Entry
	label := DefaultTheme.InfoLabelStyle
	value := DefaultTheme.InfoValueStyle

	typeDesc := fileops.GetFileTypeDesc(*entry)
	sizeStr := fileops.FormatSize(entry.Size)
	dateStr := entry.ModTime.Format("2006-01-02 15:04:05")
	modeStr := entry.Mode

	progressStr := ""
	if pv.Result != nil && !pv.Result.IsBinary && !pv.Result.IsTooLarge {
		total := pv.Result.TotalLines
		current := pv.ScrollY + 1
		if current > total {
			current = total
		}
		progressStr = fmt.Sprintf("%d / %d", current, total)
	}

	lines := []string{
		label.Render("类型: ") + value.Render(typeDesc),
		label.Render("大小: ") + value.Render(sizeStr),
		label.Render("修改: ") + value.Render(dateStr),
		label.Render("权限: ") + value.Render(modeStr),
		label.Render("行数: ") + value.Render(progressStr),
	}

	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(lipgloss.NewStyle().Width(pv.Width).Render(line))
		sb.WriteByte('\n')
	}

	return sb.String()
}
