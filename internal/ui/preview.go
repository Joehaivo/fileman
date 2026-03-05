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
	Entry    *types.FileEntry       // 当前预览的文件条目
	Result   *fileops.PreviewResult // 预览读取结果
	Width    int                    // 预览区宽度（不含边框）
	Height   int                    // 预览区高度（不含边框）
	ScrollY  int                    // 垂直滚动偏移（行数）
	IsEdit   bool                   // 是否处于编辑模式
	Lines    []string               // 编辑内容（按行存储）
	CursorY  int                    // 光标行号（0-based）
	CursorX  int                    // 光标列号（字符位置，0-based）
	ScrollX  int                    // 水平滚动偏移（字符数）
	Modified bool                   // 是否已修改
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
	pv.IsEdit = false
	pv.Lines = nil
	pv.CursorY = 0
	pv.CursorX = 0
	pv.ScrollX = 0
	pv.Modified = false

	if entry == nil || entry.IsDir {
		return
	}

	pv.Result = fileops.ReadPreview(*entry)
	// 为编辑准备内容
	if pv.Result != nil && !pv.Result.IsBinary && !pv.Result.IsTooLarge && pv.Result.Error == "" {
		pv.Lines = make([]string, len(pv.Result.Lines))
		copy(pv.Lines, pv.Result.Lines)
	}
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

// IsEditable 返回当前文件是否可编辑（文本文件且有预览内容）
func (pv *PreviewPane) IsEditable() bool {
	return pv.Entry != nil &&
		!pv.Entry.IsDir &&
		pv.Result != nil &&
		!pv.Result.IsBinary &&
		!pv.Result.IsTooLarge &&
		pv.Result.Error == "" &&
		len(pv.Lines) > 0
}

// EnterEdit 进入编辑模式
func (pv *PreviewPane) EnterEdit() {
	if !pv.IsEditable() {
		return
	}
	pv.IsEdit = true
	pv.CursorY = 0
	pv.CursorX = 0
	pv.ScrollY = 0
	pv.ScrollX = 0
	pv.clampEditCursor()
}

// ExitEdit 退出编辑模式
func (pv *PreviewPane) ExitEdit() {
	pv.IsEdit = false
}

// GetContent 获取编辑内容（用于保存）
func (pv *PreviewPane) GetContent() string {
	return strings.Join(pv.Lines, "\n")
}

// clampEditCursor 确保光标在有效范围内
func (pv *PreviewPane) clampEditCursor() {
	if len(pv.Lines) == 0 {
		pv.CursorY = 0
		pv.CursorX = 0
		return
	}

	// 限制行号
	if pv.CursorY < 0 {
		pv.CursorY = 0
	}
	if pv.CursorY >= len(pv.Lines) {
		pv.CursorY = len(pv.Lines) - 1
	}

	// 限制列号（允许在行尾+1位置）
	lineLen := len([]rune(pv.Lines[pv.CursorY]))
	if pv.CursorX < 0 {
		pv.CursorX = 0
	}
	if pv.CursorX > lineLen {
		pv.CursorX = lineLen
	}
}

// clampEditScroll 确保光标在可视区域内
func (pv *PreviewPane) clampEditScroll() {
	contentH := pv.contentHeight()
	lineNumWidth := pv.lineNumWidth()
	contentW := pv.Width - lineNumWidth - 2
	if contentW < 1 {
		contentW = 1
	}

	// 垂直滚动
	if pv.CursorY < pv.ScrollY {
		pv.ScrollY = pv.CursorY
	}
	if pv.CursorY >= pv.ScrollY+contentH {
		pv.ScrollY = pv.CursorY - contentH + 1
	}

	// 水平滚动（处理超长行）
	cursorVisualX := pv.CursorX // 字符位置
	if cursorVisualX < pv.ScrollX {
		pv.ScrollX = cursorVisualX
	}
	if cursorVisualX >= pv.ScrollX+contentW {
		pv.ScrollX = cursorVisualX - contentW + 1
	}
}

// lineNumWidth 计算行号宽度
func (pv *PreviewPane) lineNumWidth() int {
	total := len(pv.Lines)
	if total < 1 {
		total = 1
	}
	width := len(fmt.Sprintf("%d", total))
	if width < 3 {
		width = 3
	}
	return width
}

// MoveCursorUp 光标上移
func (pv *PreviewPane) MoveCursorUp() {
	if pv.CursorY > 0 {
		pv.CursorY--
		pv.clampEditCursor()
		pv.clampEditScroll()
	}
}

// MoveCursorDown 光标下移
func (pv *PreviewPane) MoveCursorDown() {
	if pv.CursorY < len(pv.Lines)-1 {
		pv.CursorY++
		pv.clampEditCursor()
		pv.clampEditScroll()
	}
}

// MoveCursorLeft 光标左移
func (pv *PreviewPane) MoveCursorLeft() {
	if pv.CursorX > 0 {
		pv.CursorX--
		pv.clampEditScroll()
	} else if pv.CursorY > 0 {
		// 移到上一行末尾
		pv.CursorY--
		pv.CursorX = len([]rune(pv.Lines[pv.CursorY]))
		pv.clampEditScroll()
	}
}

// MoveCursorRight 光标右移
func (pv *PreviewPane) MoveCursorRight() {
	lineLen := len([]rune(pv.Lines[pv.CursorY]))
	if pv.CursorX < lineLen {
		pv.CursorX++
		pv.clampEditScroll()
	} else if pv.CursorY < len(pv.Lines)-1 {
		// 移到下一行开头
		pv.CursorY++
		pv.CursorX = 0
		pv.clampEditScroll()
	}
}

// InsertChar 在光标位置插入字符
func (pv *PreviewPane) InsertChar(ch rune) {
	line := pv.Lines[pv.CursorY]
	runes := []rune(line)

	// 在 CursorX 位置插入
	newRunes := make([]rune, 0, len(runes)+1)
	newRunes = append(newRunes, runes[:pv.CursorX]...)
	newRunes = append(newRunes, ch)
	newRunes = append(newRunes, runes[pv.CursorX:]...)

	pv.Lines[pv.CursorY] = string(newRunes)
	pv.CursorX++
	pv.Modified = true
	pv.clampEditScroll()
}

// InsertNewline 在光标位置插入换行
func (pv *PreviewPane) InsertNewline() {
	line := pv.Lines[pv.CursorY]
	runes := []rune(line)

	// 分割当前行
	before := string(runes[:pv.CursorX])
	after := string(runes[pv.CursorX:])

	pv.Lines[pv.CursorY] = before
	// 在下一行插入
	newLines := make([]string, 0, len(pv.Lines)+1)
	newLines = append(newLines, pv.Lines[:pv.CursorY+1]...)
	newLines = append(newLines, after)
	newLines = append(newLines, pv.Lines[pv.CursorY+1:]...)
	pv.Lines = newLines

	pv.CursorY++
	pv.CursorX = 0
	pv.Modified = true
	pv.clampEditScroll()
}

// DeleteChar 删除光标位置的字符
func (pv *PreviewPane) DeleteChar() {
	line := pv.Lines[pv.CursorY]
	runes := []rune(line)

	if pv.CursorX < len(runes) {
		// 删除当前字符
		newRunes := make([]rune, 0, len(runes)-1)
		newRunes = append(newRunes, runes[:pv.CursorX]...)
		newRunes = append(newRunes, runes[pv.CursorX+1:]...)
		pv.Lines[pv.CursorY] = string(newRunes)
		pv.Modified = true
	}
}

// Backspace 退格删除
func (pv *PreviewPane) Backspace() {
	if pv.CursorX > 0 {
		// 删除前一个字符
		pv.CursorX--
		pv.DeleteChar()
	} else if pv.CursorY > 0 {
		// 合并到上一行
		prevLine := pv.Lines[pv.CursorY-1]
		currLine := pv.Lines[pv.CursorY]
		prevLen := len([]rune(prevLine))

		pv.Lines[pv.CursorY-1] = prevLine + currLine
		// 删除当前行
		pv.Lines = append(pv.Lines[:pv.CursorY], pv.Lines[pv.CursorY+1:]...)

		pv.CursorY--
		pv.CursorX = prevLen
		pv.Modified = true
		pv.clampEditScroll()
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

	// 编辑模式使用编辑渲染
	if pv.IsEdit {
		return pv.renderEditContent(height)
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

// renderEditContent 渲染编辑模式内容（带光标）
func (pv *PreviewPane) renderEditContent(height int) string {
	total := len(pv.Lines)
	lineNumWidth := pv.lineNumWidth()

	// 内容可用宽度（减去行号区域）
	contentWidth := pv.Width - lineNumWidth - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	var sb strings.Builder
	rendered := 0

	for i := pv.ScrollY; i < total && rendered < height; i++ {
		lineNum := fmt.Sprintf("%*d", lineNumWidth, i+1)
		lineNumStr := DefaultTheme.SubduedStyle.Render(lineNum + " ")

		line := pv.Lines[i]
		runes := []rune(line)

		// 处理超长行：水平滚动
		var displayRunes []rune
		if len(runes) > contentWidth {
			end := pv.ScrollX + contentWidth
			if end > len(runes) {
				end = len(runes)
			}
			displayRunes = runes[pv.ScrollX:end]
		} else {
			displayRunes = runes
		}

		// 当前行的光标位置（屏幕坐标）
		cursorScreenX := pv.CursorX - pv.ScrollX

		// 判断是否是光标所在行
		if i == pv.CursorY && cursorScreenX >= 0 && cursorScreenX <= len(displayRunes) {
			// 渲染带光标的行
			before := string(displayRunes[:cursorScreenX])
			after := ""
			if cursorScreenX < len(displayRunes) {
				// 光标在字符上，高亮该字符
				after = string(displayRunes[cursorScreenX:])
			}
			// 光标使用反色显示
			cursorChar := " " // 光标在行尾时显示空格
			if cursorScreenX < len(displayRunes) {
				cursorChar = string(displayRunes[cursorScreenX])
				after = string(displayRunes[cursorScreenX+1:])
			}
			cursorStyle := lipgloss.NewStyle().Reverse(true)
			content := before + cursorStyle.Render(cursorChar) + after
			// 补全宽度
			if len([]rune(content)) < contentWidth {
				content += strings.Repeat(" ", contentWidth-len([]rune(content)))
			}
			sb.WriteString(lineNumStr)
			sb.WriteString(content)
		} else {
			// 普通行
			content := string(displayRunes)
			// 补全宽度
			if len([]rune(content)) < contentWidth {
				content += strings.Repeat(" ", contentWidth-len([]rune(content)))
			}
			sb.WriteString(lineNumStr)
			sb.WriteString(content)
		}

		sb.WriteByte('\n')
		rendered++
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
