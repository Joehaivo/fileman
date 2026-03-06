package ui

import (
	"strings"

	"image/color"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/haivo/fileman/internal/types"
)

// Modal 模态弹窗组件
type Modal struct {
	Type         types.ModalType     // 弹窗类型
	Title        string              // 弹窗标题
	Message      string              // 提示消息（确认型）
	Input        textinput.Model     // 输入框（输入型）
	HasInput     bool                // 是否有输入框
	ScreenWidth  int                 // 屏幕宽度（用于居中）
	ScreenHeight int                 // 屏幕高度（用于居中）
	Progress     *types.ProgressInfo // 进度信息（进度型）
	Settings     *types.Settings     // 设置信息（设置型，临时状态）
	SettingsIdx  int                 // 设置项当前索引
}

// NewModal 创建新的模态弹窗
func NewModal() *Modal {
	ti := textinput.New()
	ti.CharLimit = 255
	ti.SetWidth(30)

	return &Modal{
		Type:  types.ModalNone,
		Input: ti,
	}
}

// SetSize 设置屏幕尺寸（用于居中计算）
func (m *Modal) SetSize(width, height int) {
	m.ScreenWidth = width
	m.ScreenHeight = height
	m.Input.SetWidth(width/3 - 4)
	if m.Input.Width() < 20 {
		m.Input.SetWidth(20)
	}
}

// IsVisible 返回弹窗是否可见
func (m *Modal) IsVisible() bool {
	return m.Type != types.ModalNone
}

// ShowNewDir 显示新建目录弹窗
func (m *Modal) ShowNewDir() {
	m.Type = types.ModalNewDir
	m.Title = "新建目录"
	m.Message = "请输入目录名称："
	m.HasInput = true
	m.Input.Reset()
	m.Input.SetValue("")
	m.Input.Placeholder = "目录名称"
	m.Input.Focus()
}

// ShowNewFile 显示新建文件弹窗
func (m *Modal) ShowNewFile() {
	m.Type = types.ModalNewFile
	m.Title = "新建文件"
	m.Message = "请输入文件名称："
	m.HasInput = true
	m.Input.Reset()
	m.Input.SetValue("")
	m.Input.Placeholder = "文件名称"
	m.Input.Focus()
}

// ShowRename 显示重命名弹窗
// currentName: 当前文件名（预填）
func (m *Modal) ShowRename(currentName string) {
	m.Type = types.ModalRename
	m.Title = "重命名"
	m.Message = "请输入新名称："
	m.HasInput = true
	m.Input.Reset()
	m.Input.SetValue(currentName)
	m.Input.Placeholder = "新名称"
	m.Input.Focus()
}

// ShowDelete 显示删除确认弹窗
// name: 要删除的文件/目录名
// count: 如果是批量删除，大于 1
func (m *Modal) ShowDelete(name string, count int) {
	m.Type = types.ModalDelete
	m.Title = "确认删除"
	m.HasInput = false
	if count > 1 {
		m.Message = "确定要删除选中的 " + itoa(count) + " 个文件吗？"
	} else {
		m.Message = "确定要删除 \"" + name + "\" 吗？"
	}
}

// ShowError 显示错误弹窗
func (m *Modal) ShowError(msg string) {
	m.Type = types.ModalError
	m.Title = "错误"
	m.Message = msg
	m.HasInput = false
}

// ShowProgress 显示进度弹窗
// title: 操作标题（如 "正在复制..."）
func (m *Modal) ShowProgress(title string, info *types.ProgressInfo) {
	m.Type = types.ModalProgress
	m.Title = title
	m.HasInput = false
	m.Progress = info
}

// Hide 隐藏弹窗
func (m *Modal) Hide() {
	m.Type = types.ModalNone
	m.Progress = nil
}

// GetInputValue 获取输入框内容
func (m *Modal) GetInputValue() string {
	return m.Input.Value()
}

// Render 渲染模态弹窗，返回覆盖在主界面上的字符串
func (m *Modal) Render() string {
	if !m.IsVisible() {
		return ""
	}

	// 弹窗内容宽度
	boxWidth := m.ScreenWidth / 3
	if boxWidth < 40 {
		boxWidth = 40
	}
	if boxWidth > 60 {
		boxWidth = 60
	}

	// 构建弹窗内容
	var content strings.Builder

	// 标题
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorTitle).
		Width(boxWidth - 4).
		Align(lipgloss.Center)
	content.WriteString(titleStyle.Render(m.Title))
	content.WriteByte('\n')
	content.WriteByte('\n')

	// 消息内容
	if m.Message != "" {
		msgLines := strings.Split(m.Message, "\n")
		for _, line := range msgLines {
			msgStyle := lipgloss.NewStyle().
				Foreground(ColorForeground).
				Width(boxWidth - 4)
			content.WriteString(msgStyle.Render(line))
			content.WriteByte('\n')
		}
		content.WriteByte('\n')
	}

	// 输入框
	if m.HasInput {
		content.WriteString(m.Input.View())
		content.WriteByte('\n')
		content.WriteByte('\n')
	}

	// 进度条
	if m.Type == types.ModalProgress && m.Progress != nil {
		content.WriteString(m.renderProgressBar(boxWidth - 4))
		content.WriteByte('\n')
		content.WriteByte('\n')
	}

	// 设置列表
	if m.Type == types.ModalSettings {
		content.WriteString(m.renderSettingsList(boxWidth - 4))
		content.WriteByte('\n')
		content.WriteByte('\n')
	}

	// 操作提示
	hints := m.renderHints(boxWidth - 4)
	content.WriteString(hints)

	// 弹窗边框样式 - 使用艳丽配色
	var borderColor color.Color
	switch m.Type {
	case types.ModalDelete:
		borderColor = ColorError // 红色
	case types.ModalError:
		borderColor = ColorError // 红色
	case types.ModalProgress:
		borderColor = ColorSelected // 粉色，更醒目
	case types.ModalSettings:
		borderColor = ColorBorderFocus // 设置也是强调色
	default:
		borderColor = ColorBorderFocus // 紫色，强调色
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(2, 3). // 增加 padding，让内容与边框有更多间距
		Width(boxWidth)

	box := boxStyle.Render(content.String())

	// 居中渲染
	return lipgloss.Place(
		m.ScreenWidth,
		m.ScreenHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

// renderProgressBar 渲染进度条
func (m *Modal) renderProgressBar(width int) string {
	if m.Progress == nil {
		return ""
	}

	percent := m.Progress.Percent
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	filled := int(float64(width) * percent)
	empty := width - filled

	bar := DefaultTheme.SuccessStyle.Render(strings.Repeat("█", filled)) +
		DefaultTheme.SubduedStyle.Render(strings.Repeat("░", empty))

	percentStr := lipgloss.NewStyle().Foreground(ColorForeground).
		Render(" " + itoa(int(percent*100)) + "%")

	return bar + percentStr
}

// renderHints 渲染操作提示
func (m *Modal) renderHints(width int) string {
	var hints string
	switch m.Type {
	case types.ModalDelete:
		hints = DefaultTheme.KeyHighlight.Render("enter") +
			DefaultTheme.KeyHintStyle.Render(" 确认删除  ") +
			DefaultTheme.KeyHighlight.Render("esc") +
			DefaultTheme.KeyHintStyle.Render(" 取消")
	case types.ModalError:
		hints = DefaultTheme.KeyHighlight.Render("enter/esc") +
			DefaultTheme.KeyHintStyle.Render(" 关闭")
	case types.ModalProgress:
		if m.Progress != nil && m.Progress.IsFinish {
			hints = DefaultTheme.KeyHighlight.Render("enter/esc") +
				DefaultTheme.KeyHintStyle.Render(" 关闭")
		} else {
			hints = DefaultTheme.SubduedStyle.Render("操作进行中...")
		}
	case types.ModalSettings:
		hints = DefaultTheme.KeyHighlight.Render("enter") +
			DefaultTheme.KeyHintStyle.Render(" 确认  ") +
			DefaultTheme.KeyHighlight.Render("esc") +
			DefaultTheme.KeyHintStyle.Render(" 取消  ") +
			DefaultTheme.KeyHighlight.Render("space") +
			DefaultTheme.KeyHintStyle.Render(" 切换  ") +
			DefaultTheme.KeyHighlight.Render("↑↓") +
			DefaultTheme.KeyHintStyle.Render(" 选择")
	default:
		hints = DefaultTheme.KeyHighlight.Render("enter") +
			DefaultTheme.KeyHintStyle.Render(" 确认  ") +
			DefaultTheme.KeyHighlight.Render("esc") +
			DefaultTheme.KeyHintStyle.Render(" 取消")
	}

	return lipgloss.NewStyle().Width(width).Render(hints)
}

// ShowSettings 显示设置弹窗
func (m *Modal) ShowSettings(currentSettings types.Settings) {
	m.Type = types.ModalSettings
	m.Title = "设置"
	m.HasInput = false
	// 复制设置到临时状态
	s := currentSettings
	m.Settings = &s
	m.SettingsIdx = 0
}

// renderSettingsList 渲染设置列表
func (m *Modal) renderSettingsList(width int) string {
	if m.Settings == nil {
		return ""
	}

	var sb strings.Builder

	// 设置项1：展示修改时间
	label1 := "展示修改时间"
	status1 := "[ ] "
	if m.Settings.ShowDate {
		status1 = "[x] "
	}

	style1 := lipgloss.NewStyle().Foreground(ColorForeground)
	cursor1 := "  "

	if m.SettingsIdx == 0 {
		style1 = lipgloss.NewStyle().Foreground(ColorSelected).Bold(true)
		cursor1 = "> "
	}

	line1 := style1.Render(cursor1 + status1 + label1)
	sb.WriteString(line1)
	sb.WriteByte('\n')

	// 设置项2：显示隐藏文件
	label2 := "显示隐藏文件"
	status2 := "[ ] "
	if m.Settings.ShowHidden {
		status2 = "[x] "
	}

	style2 := lipgloss.NewStyle().Foreground(ColorForeground)
	cursor2 := "  "

	if m.SettingsIdx == 1 {
		style2 = lipgloss.NewStyle().Foreground(ColorSelected).Bold(true)
		cursor2 = "> "
	}

	line2 := style2.Render(cursor2 + status2 + label2)
	sb.WriteString(line2)

	return sb.String()
}

// itoa 简单整数转字符串
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	buf := make([]byte, 20)
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
