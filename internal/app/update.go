package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/Joehaivo/fileman/internal/config"
	"github.com/Joehaivo/fileman/internal/fileops"
	"github.com/Joehaivo/fileman/internal/i18n"
	"github.com/Joehaivo/fileman/internal/types"
	"github.com/Joehaivo/fileman/internal/ui"
	"github.com/atotto/clipboard"
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
			m.modal.ShowError(msg.err.Error(), m.msg)
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

	case toastMsg:
		// Toast 自动消失
		m.toastMessage = ""
		return m, nil

	case fileOpResultMsg:
		return m.handleFileOpResult(msg)

	case tea.KeyPressMsg:
		return m.handleKey(msg)

	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)

	case tea.MouseWheelMsg:
		return m.handleMouseWheel(msg)
	}

	return m, nil
}

// handleKey 处理键盘事件，根据当前状态分发
func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// 悬浮进度窗口完成后的按键处理
	if m.floatingProgress != nil && m.floatingProgress.IsComplete {
		if isEscape(msg) || isEnter(msg) {
			m.floatingProgress = nil
			m.selection.Clear()
			return m, tea.Batch(
				m.loadPanel(m.panelA),
				m.loadPanel(m.panelB),
			)
		}
		return m, nil
	}

	// 弹窗模式下的按键处理
	if m.modal.IsVisible() {
		return m.handleModalKey(msg)
	}

	// 编辑模式下的按键处理
	if m.isEditing {
		return m.handleEditKey(msg)
	}

	// 搜索模式下的按键处理
	if m.isSearching {
		return m.handleSearchKey(msg)
	}

	// 普通模式按键处理
	return m.handleNormalKey(msg)
}

// handleModalKey 处理弹窗模式下的按键
func (m Model) handleModalKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch m.modal.Type {
	case types.ModalDelete:
		if isEnter(msg) {
			return m.executeDelete()
		}
		if isEscape(msg) {
			m.modal.Hide()
		}

	case types.ModalNewDir, types.ModalNewFile, types.ModalRename:
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

	case types.ModalSettings:
		if isEnter(msg) {
			if m.modal.Settings != nil {
				m.settings = *m.modal.Settings
				m.panelA.ShowDate = m.settings.ShowDate
				m.panelB.ShowDate = m.settings.ShowDate
				m.msg = i18n.GetMessages(m.settings.UseEnglish)
				_ = config.SaveConfig(&config.Config{UseEnglish: m.settings.UseEnglish})
			}
			m.modal.Hide()
			return m, tea.Batch(
				m.loadPanel(m.panelA),
				m.loadPanel(m.panelB),
			)
		}
		if isEscape(msg) {
			m.modal.Hide()
			return m, nil
		}

		switch msg.String() {
		case "up":
			if m.modal.SettingsIdx > 0 {
				m.modal.SettingsIdx--
			}
		case "down":
			if m.modal.SettingsIdx < 2 {
				m.modal.SettingsIdx++
			}
		case "space":
			if m.modal.Settings != nil {
				switch m.modal.SettingsIdx {
				case 0:
					m.modal.Settings.ShowDate = !m.modal.Settings.ShowDate
				case 1:
					m.modal.Settings.ShowHidden = !m.modal.Settings.ShowHidden
				case 2:
					m.modal.Settings.UseEnglish = !m.modal.Settings.UseEnglish
				}
			}
		}
		return m, nil

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
func (m Model) handleSearchKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if isEscape(msg) {
		// 退出搜索，恢复完整列表
		m.isSearching = false
		m.searchQuery = ""
		panel := m.activePanel()
		// 保存当前选中项，退出搜索后重新定位
		var currentName string
		if entry := panel.CurrentEntry(); entry != nil {
			currentName = entry.Name
		}
		panel.IsSearching = false
		panel.SetSearch("")
		if currentName != "" {
			panel.SetCursorByName(currentName)
		}
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

	if isLeft(msg) {
		m.isSearching = false
		panel := m.activePanel()
		panel.IsSearching = false
		m.header.IsSearching = false
		m.footer.IsSearching = false
		return m.handleGoUp()
	}

	if isRight(msg) {
		m.isSearching = false
		panel := m.activePanel()
		panel.IsSearching = false
		m.header.IsSearching = false
		m.footer.IsSearching = false
		return m.handleEnter()
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
	// v2: msg.Runes -> msg.Text
	if msg.Text != "" {
		m.searchQuery += msg.Text
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

// handleEditKey 处理编辑模式下的按键
func (m Model) handleEditKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// F1 保存并退出
	if isSave(msg) {
		return m.saveEdit()
	}

	// F2 放弃更改并退出
	if isExitEdit(msg) {
		return m.cancelEdit()
	}

	// F3 复制当前行到剪贴板
	if isCopyLine(msg) {
		line := m.preview.GetCurrentLine()
		if line != "" {
			_ = clipboard.WriteAll(line)
			m.showCopiedToast("当前行")
			return m, m.startToastTimer()
		}
		return m, nil
	}

	// Ctrl+A 复制全部内容到剪贴板
	if isCopyAll(msg) {
		content := m.preview.GetAllContent()
		if content != "" {
			_ = clipboard.WriteAll(content)
			m.showCopiedToast("全部内容")
			return m, m.startToastTimer()
		}
		return m, nil
	}

	// 其他按键交给 textarea 处理
	cmd := m.preview.UpdateEditor(msg)
	return m, cmd
}

// handleNormalKey 处理普通模式按键
func (m Model) handleNormalKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
		// 如果当前文件可编辑，进入编辑模式
		if m.preview.IsEditable() {
			return m.enterEditMode()
		}
		return m, nil
	}

	if isLeft(msg) {
		// 左方向键：返回上一级目录
		return m.handleGoUp()
	}

	if isRight(msg) {
		// 右方向键：进入选中的目录
		return m.handleEnter()
	}

	// 暂时禁用多选功能
	// if isSpace(msg) {
	// 	m.activePanel().ToggleSelection()
	// 	m.selectionTotalSize = m.computeSelectionSize()
	// 	m.activePanel().MoveCursorDown()
	// 	m.updatePreview()
	// 	return m, nil
	// }

	if isRename(msg) {
		return m.showRenameModal()
	}

	if isNewDir(msg) {
		m.modal.ShowNewDir(m.msg)
		return m, nil
	}

	if isNewFile(msg) {
		m.modal.ShowNewFile(m.msg)
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

	if isSettings(msg) {
		m.modal.ShowSettings(m.settings, m.msg)
		return m, nil
	}

	if isToggleHidden(msg) {
		m.settings.ShowHidden = !m.settings.ShowHidden
		// 重新加载两个面板
		return m, tea.Batch(
			m.loadPanel(m.panelA),
			m.loadPanel(m.panelB),
		)
	}

	if isDelete(msg) {
		return m.showDeleteConfirm()
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

// enterEditMode 进入编辑模式
func (m Model) enterEditMode() (tea.Model, tea.Cmd) {
	if !m.preview.IsEditable() {
		return m, nil
	}
	m.isEditing = true
	m.preview.EnterEdit()
	m.footer.IsEditing = true
	return m, nil
}

// saveEdit 保存编辑内容并退出编辑模式
func (m Model) saveEdit() (tea.Model, tea.Cmd) {
	if m.preview.Entry == nil {
		return m, nil
	}

	content := m.preview.GetContent()
	filePath := m.preview.Entry.Path

	// 保存文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		m.modal.ShowError(err.Error(), m.msg)
		return m, nil
	}

	// 退出编辑模式，返回文件面板
	m.isEditing = false
	m.preview.ExitEdit()
	m.footer.IsEditing = false
	m.footer.CanEdit = m.preview.IsEditable()

	// 刷新预览
	m.updatePreview()

	return m, nil
}

// cancelEdit 放弃更改并退出编辑模式
func (m Model) cancelEdit() (tea.Model, tea.Cmd) {
	m.isEditing = false
	m.preview.ExitEdit()
	m.footer.IsEditing = false
	m.footer.CanEdit = m.preview.IsEditable()

	// 刷新预览（恢复原始内容）
	m.updatePreview()

	return m, nil
}

// handleGoUp 处理左方向键（返回上一级目录）
func (m Model) handleGoUp() (tea.Model, tea.Cmd) {
	panel := m.activePanel()
	currentDirName := filepath.Base(panel.Path)
	parentPath := filepath.Dir(panel.Path)

	if parentPath != panel.Path {
		// 设置待选中的文件名，以便加载后自动选中之前的目录
		panel.PendingSelect = currentDirName
		return m, m.navigateTo(parentPath)
	}
	return m, nil
}

// handleEnter 处理 Enter 键（打开目录/文件）
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	panel := m.activePanel()

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
		m.modal.ShowDelete("", m.selection.Len(), m.msg)
		return m, nil
	}

	entry := m.activePanel().CurrentEntry()
	if entry == nil {
		return m, nil
	}

	m.modal.ShowDelete(entry.Name, 1, m.msg)
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
	m.modal.ShowRename(entry.Name, m.msg)
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
		case types.ModalNewFile:
			if err := fileops.CreateFile(panel.Path, value); err != nil {
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

	// 单文件操作：使用 Toast 显示结果
	if len(entries) == 1 {
		return m.startSingleFileCopy(entries[0], dstDir)
	}

	// 多文件操作：使用悬浮进度窗口
	m.floatingProgress = &types.FloatingProgress{
		OpType: "复制",
		Total:  len(entries),
		Done:   0,
	}

	return m, m.executeMultiFileCopy(entries, dstDir)
}

// startSingleFileCopy 开始单文件复制操作
func (m Model) startSingleFileCopy(entry types.FileEntry, dstDir string) (tea.Model, tea.Cmd) {
	dst := filepath.Join(dstDir, entry.Name)

	return m, func() tea.Msg {
		var err error
		if entry.IsDir {
			err = fileops.CopyDir(entry.Path, dst, nil)
		} else {
			err = fileops.CopyFileProgress(entry.Path, dst, nil)
		}

		return fileOpResultMsg{
			opType:  "copy",
			srcPath: entry.Path,
			dstPath: dst,
			err:     err,
		}
	}
}

// executeMultiFileCopy 执行多文件复制操作
func (m *Model) executeMultiFileCopy(entries []types.FileEntry, dstDir string) tea.Cmd {
	return func() tea.Msg {
		results := make([]types.FileOpResult, 0, len(entries))
		successCount := 0

		for _, entry := range entries {
			dst := filepath.Join(dstDir, entry.Name)
			var err error
			if entry.IsDir {
				err = fileops.CopyDir(entry.Path, dst, nil)
			} else {
				err = fileops.CopyFileProgress(entry.Path, dst, nil)
			}

			result := types.FileOpResult{
				SrcPath: entry.Path,
				DstPath: dst,
				Err:     err,
			}
			results = append(results, result)

			if err == nil {
				successCount++
			}
		}

		return fileOpResultMsg{
			opType:       "copy",
			totalCount:   len(entries),
			successCount: successCount,
			results:      results,
		}
	}
}

// startMoveOperation 开始移动操作（F6）
func (m Model) startMoveOperation() (tea.Model, tea.Cmd) {
	entries, ok := m.getSelectedOrCurrent()
	if !ok {
		return m, nil
	}

	dstDir := m.otherPanelPath()

	// 单文件操作：使用 Toast 显示结果
	if len(entries) == 1 {
		return m.startSingleFileMove(entries[0], dstDir)
	}

	// 多文件操作：使用悬浮进度窗口
	m.floatingProgress = &types.FloatingProgress{
		OpType: "移动",
		Total:  len(entries),
		Done:   0,
	}

	return m, m.executeMultiFileMove(entries, dstDir)
}

// startSingleFileMove 开始单文件移动操作
func (m Model) startSingleFileMove(entry types.FileEntry, dstDir string) (tea.Model, tea.Cmd) {
	dst := filepath.Join(dstDir, entry.Name)

	return m, func() tea.Msg {
		err := fileops.MoveEntry(entry.Path, dst)

		return fileOpResultMsg{
			opType:  "move",
			srcPath: entry.Path,
			dstPath: dst,
			err:     err,
		}
	}
}

// executeMultiFileMove 执行多文件移动操作
func (m *Model) executeMultiFileMove(entries []types.FileEntry, dstDir string) tea.Cmd {
	return func() tea.Msg {
		results := make([]types.FileOpResult, 0, len(entries))
		successCount := 0

		for _, entry := range entries {
			dst := filepath.Join(dstDir, entry.Name)
			err := fileops.MoveEntry(entry.Path, dst)

			result := types.FileOpResult{
				SrcPath: entry.Path,
				DstPath: dst,
				Err:     err,
			}
			results = append(results, result)

			if err == nil {
				successCount++
			}
		}

		return fileOpResultMsg{
			opType:       "move",
			totalCount:   len(entries),
			successCount: successCount,
			results:      results,
		}
	}
}

// handleFileOpResult 处理文件操作结果消息
func (m Model) handleFileOpResult(msg fileOpResultMsg) (tea.Model, tea.Cmd) {
	// 多文件操作完成
	if msg.totalCount > 0 {
		if m.floatingProgress != nil {
			m.floatingProgress.Done = msg.successCount
			m.floatingProgress.Results = msg.results
			m.floatingProgress.IsComplete = true
		}
		// 刷新两个面板
		return m, tea.Batch(
			m.loadPanel(m.panelA),
			m.loadPanel(m.panelB),
		)
	}

	// 单文件操作完成
	if msg.err != nil {
		m.modal.ShowError(msg.err.Error(), m.msg)
		return m, nil
	}

	// 单文件操作成功，显示 Toast
	// 格式: .../dir/ff.mp3 -> ~/dir2/ff.mp3 复制成功
	opName := m.msg.ToastCopySuccess
	if msg.opType == "move" {
		opName = m.msg.ToastMoveSuccess
	}

	// 简化源路径（省略头部）
	srcDisplay := ui.SimplifyPath(msg.srcPath)
	srcDisplay = ui.TruncatePathHead(srcDisplay, 50)

	// 简化目标路径（用 ~ 替换 home）
	dstDisplay := ui.SimplifyPath(msg.dstPath)

	m.toastMessage = srcDisplay + " → " + dstDisplay + " " + opName

	// 刷新两个面板并启动 3 秒定时器自动消失
	return m, tea.Batch(
		m.loadPanel(m.panelA),
		m.loadPanel(m.panelB),
		tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return toastMsg{}
		}),
	)
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

// handleMouseClick 处理鼠标点击事件
func (m Model) handleMouseClick(msg tea.MouseClickMsg) (tea.Model, tea.Cmd) {
	mouse := msg.Mouse()
	x, y := mouse.X, mouse.Y

	if m.modal.IsVisible() {
		return m, nil
	}

	if m.isEditing {
		return m.handlePreviewMouseClick(x, y)
	}

	region := m.getClickRegion(x, y)
	now := time.Now()

	switch region {
	case 1: // panelA 路径行
		m.focus = types.FocusPanelA
		m.panelA.IsFocused = true
		m.panelB.IsFocused = false
		m.copyPathToClipboard(m.panelA.Path)
		return m, m.startToastTimer()

	case 2: // panelA 列表区
		m.focus = types.FocusPanelA
		m.panelA.IsFocused = true
		m.panelB.IsFocused = false
		relY := m.getPanelListY(y, true)
		m.panelA.SetCursorByY(relY)
		m.updatePreview()

		if m.isDoubleClick(now, y, 1) {
			return m.handleDoubleClick()
		}
		m.lastMouseClick = now
		m.lastClickY = y
		m.lastClickPanel = 1
		return m, nil

	case 3: // panelB 路径行
		m.focus = types.FocusPanelB
		m.panelA.IsFocused = false
		m.panelB.IsFocused = true
		m.copyPathToClipboard(m.panelB.Path)
		return m, m.startToastTimer()

	case 4: // panelB 列表区
		m.focus = types.FocusPanelB
		m.panelA.IsFocused = false
		m.panelB.IsFocused = true
		relY := m.getPanelListY(y, false)
		m.panelB.SetCursorByY(relY)
		m.updatePreview()

		if m.isDoubleClick(now, y, 2) {
			return m.handleDoubleClick()
		}
		m.lastMouseClick = now
		m.lastClickY = y
		m.lastClickPanel = 2
		return m, nil

	case 5: // 预览区
		return m.handlePreviewMouseClick(x, y)
	}

	return m, nil
}

// handleMouseWheel 处理鼠标滚轮事件
func (m Model) handleMouseWheel(msg tea.MouseWheelMsg) (tea.Model, tea.Cmd) {
	mouse := msg.Mouse()
	x, y := mouse.X, mouse.Y

	if m.modal.IsVisible() {
		return m, nil
	}

	region := m.getClickRegion(x, y)

	switch region {
	case 2: // panelA 列表区
		if m.focus != types.FocusPanelA {
			m.focus = types.FocusPanelA
			m.panelA.IsFocused = true
			m.panelB.IsFocused = false
		}
		if mouse.Button == tea.MouseWheelUp {
			m.panelA.MoveCursorUp()
		} else {
			m.panelA.MoveCursorDown()
		}
		m.updatePreview()
		return m, nil

	case 4: // panelB 列表区
		if m.focus != types.FocusPanelB {
			m.focus = types.FocusPanelB
			m.panelA.IsFocused = false
			m.panelB.IsFocused = true
		}
		if mouse.Button == tea.MouseWheelUp {
			m.panelB.MoveCursorUp()
		} else {
			m.panelB.MoveCursorDown()
		}
		m.updatePreview()
		return m, nil

	case 5: // 预览区
		if mouse.Button == tea.MouseWheelUp {
			m.preview.ScrollUp()
		} else {
			m.preview.ScrollDown()
		}
		return m, nil
	}

	return m, nil
}

// isDoubleClick 检测是否为双击
func (m *Model) isDoubleClick(now time.Time, y int, panel int) bool {
	return now.Sub(m.lastMouseClick) < 500*time.Millisecond &&
		y == m.lastClickY &&
		panel == m.lastClickPanel
}

// handleDoubleClick 处理双击事件
func (m Model) handleDoubleClick() (tea.Model, tea.Cmd) {
	panel := m.activePanel()
	entry := panel.CurrentEntry()

	if entry == nil {
		return m, nil
	}

	if entry.IsDir {
		return m, m.navigateTo(entry.Path)
	}

	if m.preview.IsEditable() {
		return m.enterEditMode()
	}

	return m, nil
}

// handlePreviewMouseClick 处理预览区的鼠标点击
func (m Model) handlePreviewMouseClick(x, y int) (tea.Model, tea.Cmd) {
	if !m.isEditing {
		return m, nil
	}

	mouseClick := tea.MouseClickMsg{
		X: x,
		Y: y,
	}
	var cmd tea.Cmd
	m.preview.Editor, cmd = m.preview.Editor.Update(mouseClick)
	return m, cmd
}

// copyPathToClipboard 复制路径到剪贴板
func (m *Model) copyPathToClipboard(path string) {
	_ = clipboard.WriteAll(path)
	m.showCopiedToast(path)
}

// startToastTimer 启动 Toast 自动消失定时器
func (m Model) startToastTimer() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return toastMsg{}
	})
}
