package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/haivo/fileman/internal/fileops"
	"github.com/haivo/fileman/internal/types"
)

const (
	// 面板最小宽度
	panelMinWidth = 20
	// 大小列宽度
	sizeWidth = 6
	// 日期列宽度
	dateWidth = 11
)

// Panel 文件面板组件，管理单个目录的文件列表显示
type Panel struct {
	Path          string             // 当前目录路径
	Entries       []types.FileEntry  // 文件条目列表
	Filtered      []types.FileEntry  // 过滤后的文件条目列表
	Cursor        int                // 当前光标位置
	Offset        int                // 列表滚动偏移量
	Width         int                // 面板宽度（不含边框）
	Height        int                // 面板高度（不含边框）
	IsFocused     bool               // 是否为焦点面板
	Selection     types.SelectionSet // 选择集（共享）
	SearchQuery   string             // 当前搜索关键词
	IsSearching   bool               // 是否处于搜索模式
	Error         string             // 错误信息
	ShowDate      bool               // 是否显示日期
	PendingSelect string             // 待选中的文件名（加载后自动选中）
}

// NewPanel 创建新的文件面板
func NewPanel(path string, selection types.SelectionSet) *Panel {
	return &Panel{
		Path:      path,
		Selection: selection,
		Cursor:    0,
		ShowDate:  false, // 默认不显示日期
	}
}

// SetSize 设置面板尺寸
func (p *Panel) SetSize(width, height int) {
	p.Width = width
	p.Height = height
	p.clampScrollOffset()
}

// visibleEntries 返回当前可见的文件条目列表（考虑搜索过滤）
func (p *Panel) visibleEntries() []types.FileEntry {
	if p.IsSearching && p.SearchQuery != "" {
		return p.Filtered
	}
	return p.Entries
}

// TotalItems 返回总条目数
func (p *Panel) TotalItems() int {
	return len(p.visibleEntries())
}

// CurrentEntry 返回当前光标所在的 FileEntry
func (p *Panel) CurrentEntry() *types.FileEntry {
	entries := p.visibleEntries()
	if p.Cursor < 0 || p.Cursor >= len(entries) {
		return nil
	}
	e := entries[p.Cursor]
	return &e
}

// MoveCursorUp 光标上移
func (p *Panel) MoveCursorUp() {
	if p.Cursor > 0 {
		p.Cursor--
		p.clampScrollOffset()
	}
}

// MoveCursorDown 光标下移
func (p *Panel) MoveCursorDown() {
	if p.Cursor < p.TotalItems()-1 {
		p.Cursor++
		p.clampScrollOffset()
	}
}

// MoveCursorPageUp 光标上翻页
func (p *Panel) MoveCursorPageUp() {
	visibleHeight := p.Height - 1
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	p.Cursor -= visibleHeight
	if p.Cursor < 0 {
		p.Cursor = 0
	}
	p.clampScrollOffset()
}

// MoveCursorPageDown 光标下翻页
func (p *Panel) MoveCursorPageDown() {
	visibleHeight := p.Height - 1
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	p.Cursor += visibleHeight
	total := p.TotalItems()
	if p.Cursor >= total {
		p.Cursor = total - 1
	}
	p.clampScrollOffset()
}

// MoveCursorHome 光标移至顶部
func (p *Panel) MoveCursorHome() {
	p.Cursor = 0
	p.Offset = 0
}

// MoveCursorEnd 光标移至底部
func (p *Panel) MoveCursorEnd() {
	p.Cursor = p.TotalItems() - 1
	p.clampScrollOffset()
}

// ToggleSelection 切换当前光标条目的选择状态
func (p *Panel) ToggleSelection() {
	entry := p.CurrentEntry()
	if entry == nil {
		return
	}
	p.Selection.Toggle(entry.Path)
}

// SelectAll 全选当前目录所有条目
func (p *Panel) SelectAll() {
	for _, e := range p.Entries {
		p.Selection.Add(e.Path)
	}
}

// SetSearch 设置搜索关键词并过滤列表
func (p *Panel) SetSearch(query string) {
	p.SearchQuery = query
	p.filterEntries()
	p.Cursor = 0
	p.Offset = 0
}

// SetCursorByName 根据文件名设置光标位置
func (p *Panel) SetCursorByName(name string) {
	if name == "" {
		return
	}

	// 查找匹配的文件名
	entries := p.visibleEntries()
	for i, e := range entries {
		if e.Name == name {
			p.Cursor = i
			p.clampScrollOffset()
			break
		}
	}
}

// filterEntries 根据搜索关键词过滤文件列表（模糊匹配）
func (p *Panel) filterEntries() {
	if p.SearchQuery == "" {
		p.Filtered = nil
		return
	}

	query := strings.ToLower(p.SearchQuery)
	p.Filtered = nil
	for _, e := range p.Entries {
		if strings.Contains(strings.ToLower(e.Name), query) {
			p.Filtered = append(p.Filtered, e)
		}
	}
}

// clampScrollOffset 调整滚动偏移量确保光标始终在可视区域内
func (p *Panel) clampScrollOffset() {
	if p.Height <= 0 {
		return
	}

	// 光标在偏移量之前：向上滚动
	if p.Cursor < p.Offset {
		p.Offset = p.Cursor
	}

	// 光标在可视区域之后：向下滚动
	// 注意：p.Height 包含了路径标题行，所以实际文件列表高度是 p.Height - 1
	// 当光标位于 (p.Offset + p.Height - 1) 时，已经是可视区域的最后一行
	// 所以当 Cursor >= Offset + (p.Height - 1) 时，需要向下滚动
	visibleHeight := p.Height - 1
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	if p.Cursor >= p.Offset+visibleHeight {
		// p.Cursor 是 0-based index
		// 如果 p.Cursor = 10, visibleHeight = 5
		// Offset 应该是 10 - 5 + 1 = 6
		// 显示范围: 6, 7, 8, 9, 10 (共5行)
		p.Offset = p.Cursor - visibleHeight + 1
	}

	// 边界检查
	total := p.TotalItems()
	maxOffset := total - visibleHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if p.Offset > maxOffset {
		p.Offset = maxOffset
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

// Render 渲染面板内容（不含边框）
// 始终返回固定 p.Height 行，确保布局稳定
func (p *Panel) Render() string {
	if p.Width < panelMinWidth || p.Height <= 0 {
		return ""
	}

	entries := p.visibleEntries()
	var sb strings.Builder

	// 渲染路径标题行（缩短路径）- 固定第1行
	pathLine := p.renderPathLine()
	sb.WriteString(pathLine)
	sb.WriteByte('\n')

	// 内容可视行数（减去路径标题行）
	visibleHeight := p.Height - 1
	if visibleHeight <= 0 {
		// 如果高度只有1行，只返回路径标题行
		return sb.String()
	}

	// 重新计算滚动范围（基于减去标题后的高度）
	total := len(entries)
	start := p.Offset
	end := start + visibleHeight
	if end > total {
		end = total
	}

	// 渲染文件列表行
	renderedLines := 0
	for i := start; i < end; i++ {
		line := p.renderLine(i, entries)
		sb.WriteString(line)
		sb.WriteByte('\n')
		renderedLines++
	}

	// 填充剩余空行，确保总行数 = p.Height（路径标题1行 + visibleHeight行）
	for renderedLines < visibleHeight {
		sb.WriteString(strings.Repeat(" ", p.Width))
		sb.WriteByte('\n')
		renderedLines++
	}

	// 移除最后一个换行符（因为整个字符串末尾不应该有换行）
	result := sb.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// renderPathLine 渲染路径标题行
func (p *Panel) renderPathLine() string {
	path := p.Path
	// 将 HOME 目录替换为 ~
	if home, err := filepath.Abs("~"); err == nil {
		path = strings.Replace(path, home, "~", 1)
	}

	maxWidth := p.Width - 2
	if maxWidth < 0 {
		maxWidth = 0
	}

	if len(path) > maxWidth {
		path = "…" + path[len(path)-maxWidth+1:]
	}

	style := DefaultTheme.SubduedStyle
	if p.IsFocused {
		style = DefaultTheme.TitleStyle
	}

	return style.Width(p.Width).Render(path)
}

// renderLine 渲染单行文件条目
func (p *Panel) renderLine(idx int, entries []types.FileEntry) string {
	isCursor := idx == p.Cursor

	entry := entries[idx]
	content := p.renderEntryLine(entry, isCursor)

	return content
}

// renderEntryLine 渲染文件条目行
func (p *Panel) renderEntryLine(entry types.FileEntry, isCursor bool) string {
	isSelected := p.Selection.Has(entry.Path)

	// 布局计算：
	// 总可用宽度（扣除左右 padding 各 1）
	contentWidth := p.Width - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	// 计算固定部分宽度：
	// Name(Flex) + Space(1) + Size(6) [+ Space(1) + Date(11)]
	fixedWidth := 1 + sizeWidth // 1 + 6 = 7
	if p.ShowDate {
		fixedWidth += 1 + dateWidth // + 1 + 11 = 12 -> Total 19
	}

	nameWidth := contentWidth - fixedWidth
	if nameWidth < 1 {
		nameWidth = 1
	}

	// 截断过长的文件名（使用 lipgloss.Width 计算实际显示宽度）
	name := entry.Name
	displayWidth := lipgloss.Width(name)
	if displayWidth > nameWidth {
		name = truncateString(name, nameWidth-1) + "…"
	}

	// 格式化大小
	sizeStr := ""
	if entry.IsDir {
		sizeStr = "<DIR>"
	} else {
		sizeStr = fileops.FormatSize(entry.Size)
	}

	// 格式化日期
	dateStr := ""
	if p.ShowDate {
		dateStr = fileops.FormatDate(entry.ModTime)
	}

	// 选择样式
	var nameStyle lipgloss.Style
	switch {
	case isSelected:
		nameStyle = DefaultTheme.SelectedStyle
	case entry.IsDir:
		nameStyle = DefaultTheme.DirStyle
	case entry.Type == types.FileTypeSymlink:
		nameStyle = DefaultTheme.SymlinkStyle
	default:
		// 根据扩展名使用不同颜色
		nameStyle = p.getFileStyle(entry)
	}

	// 如果面板未激活，使样式变淡
	if !p.IsFocused {
		nameStyle = nameStyle.Copy().Foreground(ColorSubdued)
	}

	namePart := nameStyle.Width(nameWidth).Render(name)
	sizeStyle := DefaultTheme.SizeStyle
	dateStyle := DefaultTheme.DateStyle
	if !p.IsFocused {
		sizeStyle = sizeStyle.Copy().Foreground(ColorSubdued)
		dateStyle = dateStyle.Copy().Foreground(ColorSubdued)
	}

	sizePart := sizeStyle.Width(sizeWidth).Align(lipgloss.Right).Render(sizeStr)

	var line string
	if p.ShowDate {
		datePart := dateStyle.Width(dateWidth).Align(lipgloss.Right).Render(dateStr)
		line = fmt.Sprintf("%s %s %s", namePart, sizePart, datePart)
	} else {
		line = fmt.Sprintf("%s %s", namePart, sizePart)
	}

	if isCursor {
		style := DefaultTheme.CursorStyle
		if !p.IsFocused {
			// 未激活面板的光标颜色变淡
			style = style.Copy().Background(ColorBorderNormal).Foreground(ColorSubdued)
		}
		// CursorStyle 包含 Padding(0, 1)
		return style.Width(p.Width).Render(line)
	}

	return lipgloss.NewStyle().Width(p.Width).Render(line)
}

// getFileStyle 根据文件类型返回合适的样式
func (p *Panel) getFileStyle(entry types.FileEntry) lipgloss.Style {
	archiveExts := map[string]bool{
		".zip": true, ".tar": true, ".gz": true,
		".bz2": true, ".xz": true, ".7z": true, ".rar": true,
	}
	imageExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".webp": true, ".svg": true, ".ico": true,
	}

	switch {
	case archiveExts[entry.Ext]:
		return DefaultTheme.ArchiveStyle
	case imageExts[entry.Ext]:
		return DefaultTheme.ImageStyle
	default:
		return lipgloss.NewStyle().Foreground(ColorForeground)
	}
}

// truncateString 截断字符串使其显示宽度不超过指定值
func truncateString(s string, maxDisplayWidth int) string {
	if maxDisplayWidth <= 0 {
		return ""
	}

	currentWidth := 0
	for i, r := range s {
		charWidth := lipgloss.Width(string(r))
		if currentWidth+charWidth > maxDisplayWidth {
			return s[:i]
		}
		currentWidth += charWidth
	}
	return s
}
