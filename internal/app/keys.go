package app

import tea "charm.land/bubbletea/v2"

// 按键常量（基于 bubbletea KeyPressMsg）
// 使用函数而非常量，方便匹配多种按键形式

// isQuit 检查是否为退出键
func isQuit(msg tea.KeyPressMsg) bool {
	return msg.String() == "f9"
}

// isUp 检查是否为向上移动键
func isUp(msg tea.KeyPressMsg) bool {
	return msg.String() == "up"
}

// isDown 检查是否为向下移动键
func isDown(msg tea.KeyPressMsg) bool {
	return msg.String() == "down"
}

// isPageUp 检查是否为向上翻页键
func isPageUp(msg tea.KeyPressMsg) bool {
	return msg.String() == "pgup"
}

// isPageDown 检查是否为向下翻页键
func isPageDown(msg tea.KeyPressMsg) bool {
	return msg.String() == "pgdown"
}

// isHome 检查是否为移至顶部键
func isHome(msg tea.KeyPressMsg) bool {
	return msg.String() == "home"
}

// isEnd 检查是否为移至底部键
func isEnd(msg tea.KeyPressMsg) bool {
	return msg.String() == "end"
}

// isEnter 检查是否为确认键
func isEnter(msg tea.KeyPressMsg) bool {
	return msg.String() == "enter"
}

// isEscape 检查是否为取消键
func isEscape(msg tea.KeyPressMsg) bool {
	return msg.String() == "esc"
}

// isTab 检查是否为 Tab（切换焦点）
func isTab(msg tea.KeyPressMsg) bool {
	return msg.String() == "tab"
}

// isSpace 检查是否为空格（多选）
// v2 中空格键返回 "space" 而非 " "
func isSpace(msg tea.KeyPressMsg) bool {
	return msg.String() == "space"
}

// isSearch 检查是否为搜索键
func isSearch(msg tea.KeyPressMsg) bool {
	return msg.String() == "/"
}

// isDelete 检查是否为删除键
func isDelete(msg tea.KeyPressMsg) bool {
	return msg.String() == "delete"
}

// isRename 检查是否为重命名键（F1）
func isRename(msg tea.KeyPressMsg) bool {
	return msg.String() == "f1"
}

// isNewDir 检查是否为新建目录键
func isNewDir(msg tea.KeyPressMsg) bool {
	return msg.String() == "f4"
}

// isNewFile 检查是否为新建文件键（F5）
func isNewFile(msg tea.KeyPressMsg) bool {
	return msg.String() == "f5"
}

// isCopy 检查是否为复制键（F2）
func isCopy(msg tea.KeyPressMsg) bool {
	return msg.String() == "f2"
}

// isMove 检查是否为移动键（F3）
func isMove(msg tea.KeyPressMsg) bool {
	return msg.String() == "f3"
}

// isEdit 检查是否为编辑键
func isEdit(msg tea.KeyPressMsg) bool {
	return msg.String() == "f6"
}

// isSettings 检查是否为设置键
func isSettings(msg tea.KeyPressMsg) bool {
	return msg.String() == "f8"
}

// isScrollUp 检查预览区是否为向上滚动（仅当焦点在预览区时）
func isScrollUp(msg tea.KeyPressMsg) bool {
	return msg.String() == "up"
}

// isScrollDown 检查预览区是否为向下滚动（仅当焦点在预览区时）
func isScrollDown(msg tea.KeyPressMsg) bool {
	return msg.String() == "down"
}

// isLeft 检查是否为向左键（返回上一级）
func isLeft(msg tea.KeyPressMsg) bool {
	return msg.String() == "left"
}

// isRight 检查是否为向右键（进入目录）
func isRight(msg tea.KeyPressMsg) bool {
	return msg.String() == "right"
}

// isSave 检查是否为保存键（Ctrl+S）
func isSave(msg tea.KeyPressMsg) bool {
	return msg.String() == "ctrl+s"
}

// isExitEdit 检查是否为退出编辑键（Esc）
func isExitEdit(msg tea.KeyPressMsg) bool {
	return msg.String() == "esc"
}

// isToggleHidden 检查是否为切换隐藏文件键（F7）
func isToggleHidden(msg tea.KeyPressMsg) bool {
	return msg.String() == "f7"
}