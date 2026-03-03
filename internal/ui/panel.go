package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/haivo/fileman/internal/fileops"
	"github.com/haivo/fileman/internal/types"
)

const (
	// 面板最小宽度
	panelMinWidth = 20
	// 图标宽度（含空格）
	iconWidth = 2
	// 大小列宽度
	sizeWidth = 6
	// 日期列宽度
	dateWidth = 11
)

// Panel 文件面板组件，管理单个目录的文件列表显示
type Panel struct {
	Path        string             // 当前目录路径
	Entries     []types.FileEntry  // 文件条目列表（不含 ".."）
	Filtered    []types.FileEntry  // 过滤后的文件条目列表
	Cursor      int                // 当前光标位置（0 = ".."）
	Offset      int                // 列表滚动偏移量
	Width       int                // 面板宽度（不含边框）
	Height      int                // 面板高度（不含边框）
	IsFocused   bool               // 是否为焦点面板
	Selection   types.SelectionSet // 选择集（共享）
	SearchQuery string             // 当前搜索关键词
	IsSearching bool               // 是否处于搜索模式
	Error       string             // 错误信息
}

// NewPanel 创建新的文件面板
func NewPanel(path string, selection types.SelectionSet) *Panel {
	return &Panel{
		Path:      path,
		Selection: selection,
		Cursor:    0,
	}
}

// SetSize 设置面板尺寸
func (p *Panel) SetSize(width, height int) {
	p.Width = width
	p.Height = height
	p.clampScrollOffset()
}

// visibleEntries 返回当前可见的文件条目列表（考虑搜索过滤）
// 索引 0 始终是 ".."（上级目录），返回列表从 1 开始是实际文件
func (p *Panel) visibleEntries() []types.FileEntry {
	if p.IsSearching && p.SearchQuery != "" {
		return p.Filtered
	}
	return p.Entries
}

// TotalItems 返回总条目数（含 ".." 占位）
func (p *Panel) TotalItems() int {
	return len(p.visibleEntries()) + 1 // +1 for ".."
}

// CurrentEntry 返回当前光标所在的 FileEntry，光标在 ".." 时返回 nil
func (p *Panel) CurrentEntry() *types.FileEntry {
	if p.Cursor == 0 {
		return nil
	}
	entries := p.visibleEntries()
	idx := p.Cursor - 1
	if idx < 0 || idx >= len(entries) {
		return nil
	}
	e := entries[idx]
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
	p.Cursor -= p.Height
	if p.Cursor < 0 {
		p.Cursor = 0
	}
	p.clampScrollOffset()
}

// MoveCursorPageDown 光标下翻页
func (p *Panel) MoveCursorPageDown() {
	p.Cursor += p.Height
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
	if p.Cursor >= p.Offset+p.Height {
		p.Offset = p.Cursor - p.Height + 1
	}

	// 边界检查
	total := p.TotalItems()
	maxOffset := total - p.Height
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
	total := len(entries) + 1 // +1 for ".."
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
// idx: 在总列表中的索引（0 = ".."）
func (p *Panel) renderLine(idx int, entries []types.FileEntry) string {
	isCursor := idx == p.Cursor

	var content string
	if idx == 0 {
		content = p.renderParentDirLine(isCursor)
	} else {
		entry := entries[idx-1]
		content = p.renderEntryLine(entry, isCursor)
	}

	return content
}

// renderParentDirLine 渲染 ".." 上级目录行
func (p *Panel) renderParentDirLine(isCursor bool) string {
	icon := " "
	text := fmt.Sprintf("%s ..", icon)

	if isCursor {
		return DefaultTheme.CursorStyle.Width(p.Width).Render(text)
	}
	return DefaultTheme.DirStyle.Width(p.Width).Render(text)
}

// renderEntryLine 渲染文件条目行
func (p *Panel) renderEntryLine(entry types.FileEntry, isCursor bool) string {
	icon := fileops.GetFileIcon(entry.Name, entry.IsDir)
	isSelected := p.Selection.Has(entry.Path)

	// 可用于文件名的宽度 = 总宽 - 图标(2) - 大小(6) - 日期(11) - 空格(3)
	nameWidth := p.Width - iconWidth - sizeWidth - dateWidth - 3
	if nameWidth < 1 {
		nameWidth = 1
	}

	// 截断过长的文件名
	name := entry.Name
	if len([]rune(name)) > nameWidth {
		runes := []rune(name)
		name = string(runes[:nameWidth-1]) + "…"
	}

	// 格式化大小
	sizeStr := ""
	if entry.IsDir {
		sizeStr = "<DIR>"
	} else {
		sizeStr = fileops.FormatSize(entry.Size)
	}

	// 格式化日期
	dateStr := fileops.FormatDate(entry.ModTime)

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

	iconPart := lipgloss.NewStyle().Width(iconWidth).Render(icon)
	namePart := nameStyle.Width(nameWidth).Render(name)
	sizePart := DefaultTheme.SizeStyle.Width(sizeWidth).Align(lipgloss.Right).Render(sizeStr)
	datePart := DefaultTheme.DateStyle.Width(dateWidth).Align(lipgloss.Right).Render(dateStr)

	line := fmt.Sprintf("%s%s %s %s", iconPart, namePart, sizePart, datePart)

	if isCursor {
		return DefaultTheme.CursorStyle.Width(p.Width).Render(line)
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
