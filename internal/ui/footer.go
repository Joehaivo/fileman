package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// Footer 底部快捷键提示组件（两行）
type Footer struct {
	Width       int  // 可用宽度
	IsSearching bool // 是否处于搜索模式
	IsEditing   bool // 是否处于编辑模式
	CanEdit     bool // 当前是否可以进入编辑模式
}

// NewFooter 创建 Footer 组件
func NewFooter() *Footer {
	return &Footer{}
}

// SetSize 设置宽度
func (f *Footer) SetSize(width int) {
	f.Width = width
}

// keyHint 快捷键提示结构
type keyHint struct {
	Key  string
	Desc string
}

// Render 渲染 Footer（两行）
func (f *Footer) Render() string {
	if f.Width <= 0 {
		return ""
	}

	var hints []keyHint

	if f.IsEditing {
		// 编辑模式
		hints = []keyHint{
			{"↑↓←→", "移动光标"},
			{"F1", "保存"},
			{"F2", "退出"},
			{"F3", "复制行"},
			{"^A", "复制全部"},
			{"Home/End", "首/尾"},
			{"PgUp/PgDn", "翻页"},
		}
	} else if f.IsSearching {
		// 搜索模式
		hints = []keyHint{
			{"Enter", "确认"},
			{"Esc", "取消搜索"},
			{"↑↓", "选择"},
			{"Backspace", "删除字符"},
		}
	} else {
		// 普通模式
		hints = []keyHint{
			{"←", "上一级"},
			{"→", "下一级"},
			{"Tab", "切换面板"},
			{"/", "搜索"},
		}
		if f.CanEdit {
			hints = append(hints, keyHint{"Enter", "编辑"})
		}
		hints = append(hints,
			// keyHint{"Space", "多选"}, // 暂时禁用多选
			keyHint{"F1", "重命名"},
			keyHint{"F2", "复制"},
			keyHint{"F3", "移动"},
			keyHint{"F4", "新建目录"},
			keyHint{"F5", "新建文件"},
			keyHint{"F6", "外部编辑"},
			keyHint{"F7", "切换隐藏"},
			keyHint{"F8", "设置"},
			keyHint{"F9", "退出"},
			keyHint{"Del", "删除"},
		)
	}

	// 自动分行：从第一行开始摆放，放不下则放到第二行
	line1, line2 := f.layoutKeys(hints)

	return line1 + "\n" + line2
}

// layoutKeys 将快捷键自动分配到两行
func (f *Footer) layoutKeys(hints []keyHint) (string, string) {
	sep := "  " // 快捷键之间的分隔符
	sepWidth := 2

	var line1Parts, line2Parts []string
	line1Width := 0
	line2Width := 0

	for _, h := range hints {
		// 计算当前快捷键的宽度
		part := f.formatKeyHint(h)
		partWidth := lipgloss.Width(part)

		// 尝试放到第一行
		neededWidth := partWidth
		if len(line1Parts) > 0 {
			neededWidth += sepWidth
		}

		if line1Width+neededWidth <= f.Width {
			line1Parts = append(line1Parts, part)
			line1Width += neededWidth
		} else {
			// 尝试放到第二行
			neededWidth2 := partWidth
			if len(line2Parts) > 0 {
				neededWidth2 += sepWidth
			}
			// 只有第二行也放得下才添加
			if line2Width+neededWidth2 <= f.Width {
				line2Parts = append(line2Parts, part)
				line2Width += neededWidth2
			}
		}
	}

	// 构建最终的两行
	line1 := strings.Join(line1Parts, sep)
	line2 := strings.Join(line2Parts, sep)

	// 填充宽度
	if lipgloss.Width(line1) < f.Width {
		line1 += strings.Repeat(" ", f.Width-lipgloss.Width(line1))
	}
	if lipgloss.Width(line2) < f.Width {
		line2 += strings.Repeat(" ", f.Width-lipgloss.Width(line2))
	}

	return line1, line2
}

// formatKeyHint 格式化单个快捷键提示
func (f *Footer) formatKeyHint(h keyHint) string {
	key := DefaultTheme.KeyHighlight.Render(h.Key)
	desc := DefaultTheme.KeyHintStyle.Render(" " + h.Desc)
	return key + desc
}
