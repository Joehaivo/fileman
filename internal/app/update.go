package app

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"github.com/haivo/fileman/internal/fileops"
	"github.com/haivo/fileman/internal/types"
	"github.com/haivo/fileman/internal/ui"
)

// Update 实现 tea.Model 接口，处理所有消息
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.calcSizes()
		m.updatePreview()
		return m, nil

	case panelLoadMsg:
		if msg.panel == m.panelA {
			applyPanelLoad(m.panelA, msg)
		} else {
			applyPanelLoad(m.panelB, msg)
		}
		m.updatePreview()
		return m, nil

	case fileOpMsg:
		m.modal.Hide()
		if msg.err != nil {
			m.modal.ShowError(msg.err.Error())
			return m, nil
		}
		// 操作成功：刷新两个面板
		m.selection.Clear()
		return m, tea.Batch(
			m.loadPanel(m.panelA),
			m.loadPanel(m.panelB),
		)

	case progressMsg:
		if m.modal.Type == types.ModalProgress && m.modal.Progress != nil {
			*m.modal.Progress = msg.info
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

// handleKey 处理键盘事件，根据当前状态分发
func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 弹窗模式下的按键处理
	if m.modal.IsVisible() {
		return m.handleModalKey(msg)
	}

	// 搜索模式下的按键处理
	if m.isSearching {
		return m.handleSearchKey(msg)
	}

	// 普通模式按键处理
	return m.handleNormalKey(msg)
}

// handleModalKey 处理弹窗模式下的按键
func (m Model) handleModalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.modal.Type {
	case types.ModalDelete:
		if isEnter(msg) {
			return m.executeDelete()
		}
		if isEscape(msg) {
			m.modal.Hide()
		}

	case types.ModalNewDir, types.ModalRename:
		if isEnter(msg) {
			return m.executeInputModal()
		}
		if isEscape(msg) {
			m.modal.Hide()
			return m, nil
		}
		// 转发按键给输入框
		var cmd tea.Cmd
		m.modal.Input, cmd = m.modal.Input.Update(msg)
		return m, cmd

	case types.ModalError, types.ModalProgress:
		if isEnter(msg) || isEscape(msg) {
			if m.modal.Type == types.ModalProgress {
				if m.modal.Progress == nil || m.modal.Progress.IsFinish {
					m.modal.Hide()
				}
			} else {
				m.modal.Hide()
			}
		}
	}

	return m, nil
}

// handleSearchKey 处理搜索模式下的按键
func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if isEscape(msg) {
		// 退出搜索，恢复完整列表
		m.isSearching = false
		m.searchQuery = ""
		panel := m.activePanel()
		panel.IsSearching = false
		panel.SetSearch("")
		m.header.IsSearching = false
		m.footer.IsSearching = false
		m.updatePreview()
		return m, nil
	}

	if isEnter(msg) {
		// 搜索模式下 Enter = 打开当前选中项
		m.isSearching = false
		panel := m.activePanel()
		panel.IsSearching = false
		m.header.IsSearching = false
		m.footer.IsSearching = false
		return m.handleEnter()
	}

	if isUp(msg) {
		m.activePanel().MoveCursorUp()
		m.updatePreview()
		return m, nil
	}

	if isDown(msg) {
		m.activePanel().MoveCursorDown()
		m.updatePreview()
		return m, nil
	}

	// Backspace 删除最后一个搜索字符
	if msg.String() == "backspace" || msg.String() == "ctrl+h" {
		if len(m.searchQuery) > 0 {
			runes := []rune(m.searchQuery)
			m.searchQuery = string(runes[:len(runes)-1])
			m.applySearch()
		}
		return m, nil
	}

	// 其他字符追加到搜索词
	if len(msg.Runes) > 0 {
		m.searchQuery += string(msg.Runes)
		m.applySearch()
	}

	return m, nil
}

// applySearch 将搜索关键词应用到当前面板
func (m *Model) applySearch() {
	panel := m.activePanel()
	panel.SetSearch(m.searchQuery)
	m.header.SearchQuery = m.searchQuery
	m.updatePreview()
}

// handleNormalKey 处理普通模式按键
func (m Model) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if isQuit(msg) {
		return m, tea.Quit
	}

	if isTab(msg) {
		return m.switchFocus()
	}

	if isSearch(msg) {
		return m.enterSearchMode()
	}

	if isUp(msg) {
		m.activePanel().MoveCursorUp()
		m.updatePreview()
		return m, nil
	}

	if isDown(msg) {
		m.activePanel().MoveCursorDown()
		m.updatePreview()
		return m, nil
	}

	if isPageUp(msg) {
		m.activePanel().MoveCursorPageUp()
		m.updatePreview()
		return m, nil
	}

	if isPageDown(msg) {
		m.activePanel().MoveCursorPageDown()
		m.updatePreview()
		return m, nil
	}

	if isHome(msg) {
		m.activePanel().MoveCursorHome()
		m.updatePreview()
		return m, nil
	}

	if isEnd(msg) {
		m.activePanel().MoveCursorEnd()
		m.updatePreview()
		return m, nil
	}

	if isEnter(msg) {
		return m.handleEnter()
	}

	if isLeft(msg) {
		// 左方向键：返回上一级目录
		return m.handleGoUp()
	}

	if isRight(msg) {
		// 右方向键：进入选中的目录
		return m.handleEnter()
	}

	if isSpace(msg) {
		m.activePanel().ToggleSelection()
		m.selectionTotalSize = m.computeSelectionSize()
		m.activePanel().MoveCursorDown()
		m.updatePreview()
		return m, nil
	}

	if isSelectAll(msg) {
		m.activePanel().SelectAll()
		m.selectionTotalSize = m.computeSelectionSize()
		return m, nil
	}

	if isDelete(msg) {
		return m.showDeleteConfirm()
	}

	if isRename(msg) {
		return m.showRenameModal()
	}

	if isNewDir(msg) {
		m.modal.ShowNewDir()
		return m, nil
	}

	if isCopy(msg) {
		return m.startCopyOperation()
	}

	if isMove(msg) {
		return m.startMoveOperation()
	}

	if isEdit(msg) {
		return m.openInEditor()
	}

	return m, nil
}

// switchFocus 切换焦点面板
func (m Model) switchFocus() (tea.Model, tea.Cmd) {
	if m.focus == types.FocusPanelA {
		m.focus = types.FocusPanelB
	} else {
		m.focus = types.FocusPanelA
	}
	m.panelA.IsFocused = m.focus == types.FocusPanelA
	m.panelB.IsFocused = m.focus == types.FocusPanelB
	m.updatePreview()
	return m, nil
}

// enterSearchMode 进入搜索模式
func (m Model) enterSearchMode() (tea.Model, tea.Cmd) {
	m.isSearching = true
	m.searchQuery = ""
	panel := m.activePanel()
	panel.IsSearching = true
	panel.SetSearch("")
	m.header.IsSearching = true
	m.header.SearchQuery = ""
	m.footer.IsSearching = true
	return m, nil
}

// handleGoUp 处理左方向键（返回上一级目录）
func (m Model) handleGoUp() (tea.Model, tea.Cmd) {
	panel := m.activePanel()
	parentPath := filepath.Dir(panel.Path)
	if parentPath != panel.Path {
		return m, m.navigateTo(parentPath)
	}
	return m, nil
}

// handleEnter 处理 Enter 键（打开目录/文件）
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	panel := m.activePanel()

	if panel.Cursor == 0 {
		// 返回上级目录
		return m.handleGoUp()
	}

	entry := panel.CurrentEntry()
	if entry == nil {
		return m, nil
	}

	if entry.IsDir {
		return m, m.navigateTo(entry.Path)
	}

	return m, nil
}

// showDeleteConfirm 显示删除确认弹窗
func (m Model) showDeleteConfirm() (tea.Model, tea.Cmd) {
	if m.selection.Len() > 0 {
		m.modal.ShowDelete("", m.selection.Len())
		return m, nil
	}

	entry := m.activePanel().CurrentEntry()
	if entry == nil {
		return m, nil
	}

	m.modal.ShowDelete(entry.Name, 1)
	return m, nil
}

// executeDelete 执行删除操作
func (m Model) executeDelete() (tea.Model, tea.Cmd) {
	m.modal.Hide()

	entries, ok := m.getSelectedOrCurrent()
	if !ok {
		return m, nil
	}

	return m, func() tea.Msg {
		for _, e := range entries {
			if err := fileops.DeleteEntry(e.Path); err != nil {
				return fileOpMsg{err: err}
			}
		}
		return fileOpMsg{}
	}
}

// showRenameModal 显示重命名弹窗
func (m Model) showRenameModal() (tea.Model, tea.Cmd) {
	entry := m.activePanel().CurrentEntry()
	if entry == nil {
		return m, nil
	}
	m.modal.ShowRename(entry.Name)
	return m, nil
}

// executeInputModal 执行输入型弹窗的操作（新建目录 / 重命名）
func (m Model) executeInputModal() (tea.Model, tea.Cmd) {
	value := m.modal.GetInputValue()
	if value == "" {
		return m, nil
	}

	modalType := m.modal.Type
	m.modal.Hide()

	panel := m.activePanel()

	return m, func() tea.Msg {
		switch modalType {
		case types.ModalNewDir:
			if err := fileops.CreateDir(panel.Path, value); err != nil {
				return fileOpMsg{err: err}
			}
		case types.ModalRename:
			entry := panel.CurrentEntry()
			if entry != nil {
				if err := fileops.RenameEntry(entry.Path, value); err != nil {
					return fileOpMsg{err: err}
				}
			}
		}
		return fileOpMsg{}
	}
}

// startCopyOperation 开始复制操作（F5）
func (m Model) startCopyOperation() (tea.Model, tea.Cmd) {
	entries, ok := m.getSelectedOrCurrent()
	if !ok {
		return m, nil
	}

	dstDir := m.otherPanelPath()
	progressInfo := &types.ProgressInfo{}
	m.modal.ShowProgress("正在复制...", progressInfo)

	return m, func() tea.Msg {
		for _, e := range entries {
			dst := filepath.Join(dstDir, e.Name)
			if e.IsDir {
				if err := fileops.CopyDir(e.Path, dst, func(done, total int64) {
					pct := float64(done) / float64(total)
					_ = pct // 简化：不实时更新进度
				}); err != nil {
					return fileOpMsg{err: err}
				}
			} else {
				if err := fileops.CopyFileProgress(e.Path, dst, nil); err != nil {
					return fileOpMsg{err: err}
				}
			}
		}
		return fileOpMsg{}
	}
}

// startMoveOperation 开始移动操作（F6）
func (m Model) startMoveOperation() (tea.Model, tea.Cmd) {
	entries, ok := m.getSelectedOrCurrent()
	if !ok {
		return m, nil
	}

	dstDir := m.otherPanelPath()
	progressInfo := &types.ProgressInfo{}
	m.modal.ShowProgress("正在移动...", progressInfo)

	return m, func() tea.Msg {
		for _, e := range entries {
			dst := filepath.Join(dstDir, e.Name)
			if err := fileops.MoveEntry(e.Path, dst); err != nil {
				return fileOpMsg{err: err}
			}
		}
		return fileOpMsg{}
	}
}

// openInEditor 在编辑器中打开文件（Ctrl+E）
func (m Model) openInEditor() (tea.Model, tea.Cmd) {
	entry := m.activePanel().CurrentEntry()
	if entry == nil || entry.IsDir {
		return m, nil
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}

	filePath := entry.Path
	return m, tea.ExecProcess(exec.Command(editor, filePath), func(err error) tea.Msg {
		if err != nil {
			return fileOpMsg{err: err}
		}
		return fileOpMsg{}
	})
}

// updatePanelAfterOp 操作完成后刷新面板
func (m *Model) updatePanelAfterOp(p *ui.Panel) tea.Cmd {
	return m.loadPanel(p)
}
