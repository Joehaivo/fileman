package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Footer 底部快捷键提示组件（两行）
type Footer struct {
	Width      int  // 可用宽度
	IsSearching bool // 是否处于搜索模式
}

// NewFooter 创建 Footer 组件
func NewFooter() *Footer {
	return &Footer{}
}

// SetSize 设置宽度
func (f *Footer) SetSize(width int) {
	f.Width = width
}

// Render 渲染 Footer（两行）
func (f *Footer) Render() string {
	if f.Width <= 0 {
		return ""
	}

	var line1, line2 string

	if f.IsSearching {
		line1 = f.renderKeys([]keyHint{
			{"enter", "确认"},
			{"esc", "取消搜索"},
			{"↑↓", "选择"},
		})
		line2 = f.renderKeys([]keyHint{
			{"backspace", "删除字符"},
		})
	} else {
		line1 = f.renderKeys([]keyHint{
			{"←", "上一级"},
			{"→", "进入"},
			{"enter", "打开"},
			{"space", "多选"},
			{"tab", "切换窗口"},
			{"f3", "搜索"},
			{"f9", "全选"},
		})
		line2 = f.renderKeys([]keyHint{
			{"f1", "删除"},
			{"f2", "重命名"},
			{"f4", "编辑"},
			{"f5", "复制"},
			{"f6", "移动"},
			{"f7", "新建目录"},
			{"f8", "设置"},
			{"f10", "退出"},
		})
	}

	return line1 + "\n" + line2
}

// keyHint 快捷键提示结构
type keyHint struct {
	Key  string
	Desc string
}

// renderKeys 渲染一行快捷键提示
func (f *Footer) renderKeys(hints []keyHint) string {
	var parts []string
	for _, h := range hints {
		key := DefaultTheme.KeyHighlight.Render(h.Key)
		desc := DefaultTheme.KeyHintStyle.Render(" " + h.Desc)
		parts = append(parts, key+desc)
	}

	sep := DefaultTheme.SubduedStyle.Render("  ")
	line := strings.Join(parts, sep)

	lineLen := lipgloss.Width(line)
	if lineLen > f.Width {
		// 超出宽度时截断
		line = line[:f.Width]
	}

	return lipgloss.NewStyle().Width(f.Width).Render(line)
}
