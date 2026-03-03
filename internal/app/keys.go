package app

import "github.com/charmbracelet/bubbletea"

// 按键常量（基于 bubbletea KeyMsg）
// 使用函数而非常量，方便匹配多种按键形式

// isQuit 检查是否为退出键
func isQuit(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+q"
}

// isUp 检查是否为向上移动键
func isUp(msg tea.KeyMsg) bool {
	k := msg.String()
	return k == "up" || k == "k"
}

// isDown 检查是否为向下移动键
func isDown(msg tea.KeyMsg) bool {
	k := msg.String()
	return k == "down" || k == "j"
}

// isPageUp 检查是否为向上翻页键
func isPageUp(msg tea.KeyMsg) bool {
	k := msg.String()
	return k == "pgup" || k == "ctrl+u"
}

// isPageDown 检查是否为向下翻页键
func isPageDown(msg tea.KeyMsg) bool {
	k := msg.String()
	return k == "pgdown" || k == "ctrl+d"
}

// isHome 检查是否为移至顶部键
func isHome(msg tea.KeyMsg) bool {
	k := msg.String()
	return k == "home" || k == "g"
}

// isEnd 检查是否为移至底部键
func isEnd(msg tea.KeyMsg) bool {
	k := msg.String()
	return k == "end" || k == "G"
}

// isEnter 检查是否为确认键
func isEnter(msg tea.KeyMsg) bool {
	return msg.String() == "enter"
}

// isEscape 检查是否为取消键
func isEscape(msg tea.KeyMsg) bool {
	return msg.String() == "esc"
}

// isTab 检查是否为 Tab（切换焦点）
func isTab(msg tea.KeyMsg) bool {
	return msg.String() == "tab"
}

// isSpace 检查是否为空格（多选）
func isSpace(msg tea.KeyMsg) bool {
	return msg.String() == " "
}

// isSearch 检查是否为搜索键
func isSearch(msg tea.KeyMsg) bool {
	return msg.String() == "/"
}

// isDelete 检查是否为删除键
func isDelete(msg tea.KeyMsg) bool {
	k := msg.String()
	return k == "delete" || k == "backspace"
}

// isNewDir 检查是否为新建目录键
func isNewDir(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+n"
}

// isRename 检查是否为重命名键（F2）
func isRename(msg tea.KeyMsg) bool {
	return msg.String() == "f2"
}

// isCopy 检查是否为复制键（F5）
func isCopy(msg tea.KeyMsg) bool {
	return msg.String() == "f5"
}

// isMove 检查是否为移动键（F6）
func isMove(msg tea.KeyMsg) bool {
	return msg.String() == "f6"
}

// isEdit 检查是否为编辑键（Ctrl+E）
func isEdit(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+e"
}

// isSelectAll 检查是否为全选键（Ctrl+A）
func isSelectAll(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+a"
}

// isScrollUp 检查预览区是否为向上滚动（仅当焦点在预览区时）
func isScrollUp(msg tea.KeyMsg) bool {
	return msg.String() == "up" || msg.String() == "k"
}

// isScrollDown 检查预览区是否为向下滚动（仅当焦点在预览区时）
func isScrollDown(msg tea.KeyMsg) bool {
	return msg.String() == "down" || msg.String() == "j"
}
