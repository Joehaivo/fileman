package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Joehaivo/fileman/internal/fileops"
	"github.com/Joehaivo/fileman/internal/i18n"
	"github.com/Joehaivo/fileman/internal/types"
)

// PreviewPane 右侧预览组件，使用 bubbles textarea 实现预览和编辑
type PreviewPane struct {
	Entry           *types.FileEntry       // 当前预览的文件条目
	Editor          textarea.Model         // textarea 编辑器（预览和编辑共用）
	Width           int                    // 预览区宽度（不含边框）
	Height          int                    // 预览区高度（不含边框）
	originalContent string                 // 原始内容（用于修改检测）
	loaded          bool                   // 是否已加载内容
	Result          *fileops.PreviewResult // 预览读取结果（用于显示信息）
	Msg             *i18n.Messages         // 国际化文本
}

// NewPreviewPane 创建新的预览组件
func NewPreviewPane() *PreviewPane {
	ta := textarea.New()
	ta.ShowLineNumbers = true
	ta.MaxHeight = 0 // 取消默认的 99 行限制，允许任意行数的文件编辑

	// 设置自定义样式
	styles := customTextareaStyles()
	ta.SetStyles(styles)

	// 默认处于预览模式（blur 状态）
	ta.Blur()

	return &PreviewPane{
		Editor: ta,
	}
}

// customTextareaStyles 自定义 textarea 样式以匹配 fileman 主题
func customTextareaStyles() textarea.Styles {
	// 使用暗色主题
	styles := textarea.DefaultDarkStyles()

	// 匹配 fileman 主题
	styles.Focused.LineNumber = DefaultTheme.SubduedStyle
	styles.Focused.CursorLine = lipgloss.NewStyle() // 无高亮背景
	styles.Focused.Placeholder = DefaultTheme.SubduedStyle

	styles.Blurred.LineNumber = DefaultTheme.SubduedStyle
	styles.Blurred.CursorLine = lipgloss.NewStyle()
	styles.Blurred.Placeholder = DefaultTheme.SubduedStyle

	return styles
}

// SetSize 设置预览区尺寸
func (pv *PreviewPane) SetSize(width, height int) {
	pv.Width = width
	pv.Height = height

	// textarea 高度 = 总高度 - 标题行 - 分隔线 - 信息区(5行)
	editorHeight := height - 1 - 1 - 5
	if editorHeight < 1 {
		editorHeight = 1
	}

	pv.Editor.SetWidth(width)
	pv.Editor.SetHeight(editorHeight)
}

// SetEntry 设置要预览的文件条目，并重新加载预览内容
func (pv *PreviewPane) SetEntry(entry *types.FileEntry) {
	pv.Entry = entry
	pv.loaded = false
	pv.originalContent = ""
	pv.Result = nil

	if entry == nil || entry.IsDir {
		pv.Editor.SetValue("")
		pv.Editor.ShowLineNumbers = true // 恢复默认设置
		pv.Editor.Blur()
		return
	}

	// 读取文件内容
	result := fileops.ReadPreview(*entry)
	pv.Result = result

	// 处理压缩文件
	if result.IsArchive {
		content := strings.Join(result.Lines, "\n")
		pv.Editor.SetValue(content)
		pv.Editor.ShowLineNumbers = false // 禁用行号避免错位
		pv.loaded = false                 // 压缩文件不可编辑
		pv.Editor.Blur()
		return
	}

	if result.Error != "" || result.IsBinary || result.IsTooLarge {
		// 显示错误信息
		pv.Editor.ShowLineNumbers = false // 错误信息不需要行号
		if result.Error != "" {
			pv.Editor.SetValue(result.Error)
		} else if result.IsBinary {
			if pv.Msg != nil {
				pv.Editor.SetValue(pv.Msg.PreviewBinary)
			} else {
				pv.Editor.SetValue("二进制文件，无法预览")
			}
		} else if result.IsTooLarge {
			if pv.Msg != nil {
				msg := fmt.Sprintf(pv.Msg.PreviewTooLargeFmt, fileops.FormatSize(result.FileSize))
				pv.Editor.SetValue(msg)
			} else {
				pv.Editor.SetValue(fmt.Sprintf("文件过大 (%s)，无法预览", fileops.FormatSize(result.FileSize)))
			}
		}
		pv.Editor.Blur()
		return
	}

	content := strings.Join(result.Lines, "\n")
	pv.Editor.SetValue(content)
	pv.Editor.ShowLineNumbers = true // 普通文本文件显示行号
	pv.originalContent = content
	pv.loaded = true
	pv.Editor.Blur() // 默认预览模式
}

// ScrollUp 预览内容向上滚动
func (pv *PreviewPane) ScrollUp() {
	// 直接调用 textarea 的光标移动方法，无需焦点
	pv.Editor.CursorUp()
}

// ScrollDown 预览内容向下滚动
func (pv *PreviewPane) ScrollDown() {
	pv.Editor.CursorDown()
}

// ScrollPageUp 向上翻页
func (pv *PreviewPane) ScrollPageUp() {
	pv.Editor.PageUp()
}

// ScrollPageDown 向下翻页
func (pv *PreviewPane) ScrollPageDown() {
	pv.Editor.PageDown()
}

// IsEditable 返回当前文件是否可编辑（文本文件且有预览内容）
func (pv *PreviewPane) IsEditable() bool {
	return pv.Entry != nil &&
		!pv.Entry.IsDir &&
		pv.loaded
}

// EnterEdit 进入编辑模式
func (pv *PreviewPane) EnterEdit() {
	if !pv.IsEditable() {
		return
	}
	// 将光标移到文本开头
	pv.Editor.MoveToBegin()
	pv.Editor.Focus()
}

// ExitEdit 退出编辑模式
func (pv *PreviewPane) ExitEdit() {
	pv.Editor.Blur()
}

// IsModified 检查内容是否已修改
func (pv *PreviewPane) IsModified() bool {
	return pv.Editor.Value() != pv.originalContent
}

// GetContent 获取编辑内容（用于保存）
func (pv *PreviewPane) GetContent() string {
	return pv.Editor.Value()
}

// GetCurrentLine 获取当前光标所在行的内容
func (pv *PreviewPane) GetCurrentLine() string {
	lineNum := pv.Editor.Line()
	value := pv.Editor.Value()
	lines := strings.Split(value, "\n")
	if lineNum >= 0 && lineNum < len(lines) {
		return lines[lineNum]
	}
	return ""
}

// GetAllContent 获取全部内容（等同于 GetContent）
func (pv *PreviewPane) GetAllContent() string {
	return pv.Editor.Value()
}

// ResetContent 重置内容为原始内容
func (pv *PreviewPane) ResetContent() {
	pv.Editor.SetValue(pv.originalContent)
}

// UpdateEditor 更新 textarea 状态（用于编辑模式下的按键处理）
func (pv *PreviewPane) UpdateEditor(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	pv.Editor, cmd = pv.Editor.Update(msg)
	return cmd
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

	// 内容区 (textarea)
	sb.WriteString(pv.Editor.View())
	sb.WriteByte('\n')

	// 分隔线
	sep := strings.Repeat("─", pv.Width)
	sb.WriteString(DefaultTheme.SubduedStyle.Render(sep))
	sb.WriteByte('\n')

	// 文件信息区（5行）
	sb.WriteString(pv.renderInfo())

	return sb.String()
}

// renderEmpty 渲染空状态（无文件选中）
func (pv *PreviewPane) renderEmpty() string {
	msg := "选择文件以预览"
	if pv.Msg != nil {
		msg = pv.Msg.PreviewSelectFile
	}
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
	if len(name) > pv.Width-3 { // -3: 左边距1 + 右边距2
		runes := []rune(name)
		if len(runes) > pv.Width-3 {
			name = string(runes[:pv.Width-4]) + "…"
		}
	}
	// 左边添加一个空格，使文件名不紧贴边缘
	return " " + DefaultTheme.PreviewTitle.Width(pv.Width-1).Render(name)
}

// renderInfo 渲染文件信息区（5行）
func (pv *PreviewPane) renderInfo() string {
	if pv.Entry == nil {
		return strings.Repeat(strings.Repeat(" ", pv.Width)+"\n", 5)
	}

	entry := pv.Entry
	label := DefaultTheme.InfoLabelStyle
	value := DefaultTheme.InfoValueStyle

	useEnglish := pv.Msg != nil && pv.Msg == i18n.English
	typeDesc := fileops.GetFileTypeDesc(*entry, useEnglish)
	sizeStr := fileops.FormatSize(entry.Size)
	dateStr := entry.ModTime.Format("2006-01-02 15:04:05")
	modeStr := entry.Mode

	// 使用国际化标签
	var labelType, labelSize, labelModified, labelMode, labelLines string
	var labelArchiveFiles, labelArchiveSize string
	if pv.Msg != nil {
		labelType = pv.Msg.InfoType
		labelSize = pv.Msg.InfoSize
		labelModified = pv.Msg.InfoModified
		labelMode = pv.Msg.InfoMode
		labelLines = pv.Msg.InfoLines
		labelArchiveFiles = pv.Msg.InfoArchiveFiles
		labelArchiveSize = pv.Msg.InfoArchiveSize
	} else {
		labelType = "类型: "
		labelSize = "大小: "
		labelModified = "修改: "
		labelMode = "权限: "
		labelLines = "行数: "
		labelArchiveFiles = "文件数: "
		labelArchiveSize = "解压大小: "
	}

	var lines []string

	// 压缩文件特殊处理
	if pv.Result != nil && pv.Result.IsArchive {
		archiveTypeDesc := fileops.GetArchiveFormatDesc(pv.Result.ArchiveFormat, useEnglish)
		archiveSizeStr := fileops.FormatSize(pv.Result.ArchiveSize)

		lines = []string{
			label.Render(labelType) + value.Render(archiveTypeDesc),
			label.Render(labelSize) + value.Render(sizeStr),
			label.Render(labelArchiveFiles) + value.Render(fmt.Sprintf("%d", pv.Result.ArchiveCount)),
			label.Render(labelArchiveSize) + value.Render(archiveSizeStr),
			label.Render(labelModified) + value.Render(dateStr),
		}
	} else {
		// 获取行数信息
		var progressStr string
		if pv.Result != nil && !pv.Result.IsBinary && !pv.Result.IsTooLarge && pv.Result.Error == "" {
			total := pv.Result.TotalLines
			progressStr = fmt.Sprintf("%d", total)
		}

		lines = []string{
			label.Render(labelType) + value.Render(typeDesc),
			label.Render(labelSize) + value.Render(sizeStr),
			label.Render(labelModified) + value.Render(dateStr),
			label.Render(labelMode) + value.Render(modeStr),
			label.Render(labelLines) + value.Render(progressStr),
		}
	}

	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(lipgloss.NewStyle().Width(pv.Width).Render(line))
		sb.WriteByte('\n')
	}

	return sb.String()
}
