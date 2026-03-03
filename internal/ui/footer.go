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
			{"Enter", "确认"},
			{"Esc", "取消搜索"},
			{"↑↓", "选择"},
		})
		line2 = f.renderKeys([]keyHint{
			{"Backspace", "删除字符"},
		})
	} else {
		line1 = f.renderKeys([]keyHint{
			{"←", "上一级"},
			{"→", "进入"},
			{"Enter", "打开"},
			{"Space", "多选"},
			{"Tab", "切换窗口"},
			{"/", "搜索"},
			{"A", "全选"},
		})
		line2 = f.renderKeys([]keyHint{
			{"Del", "删除"},
			{"F2", "重命名"},
			{"N", "新建目录"},
			{"F5", "复制"},
			{"F6", "移动"},
			{"E", "编辑"},
			{"Q", "退出"},
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
