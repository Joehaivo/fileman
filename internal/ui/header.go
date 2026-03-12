package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/Joehaivo/fileman/internal/fileops"
	"github.com/Joehaivo/fileman/internal/i18n"
	"github.com/Joehaivo/fileman/internal/types"
)

// Header 顶部标题组件
type Header struct {
	Width       int                // 可用宽度
	Selection   types.SelectionSet // 当前选择集
	IsSearching bool               // 是否处于搜索模式
	SearchQuery string             // 搜索关键词
	Msg         *i18n.Messages     // 国际化文本
}

// NewHeader 创建 Header 组件
func NewHeader(selection types.SelectionSet) *Header {
	return &Header{
		Selection: selection,
	}
}

// SetSize 设置宽度
func (h *Header) SetSize(width int) {
	h.Width = width
}

// Render 渲染 Header（单行）
func (h *Header) Render() string {
	if h.Width <= 0 || h.Msg == nil {
		return ""
	}

	// 左侧：应用名称 + 版本
	leftStr := ""
	if h.IsSearching {
		leftStr = DefaultTheme.SearchStyle.Render(h.Msg.HeaderSearchLabel) +
			DefaultTheme.TitleStyle.Render(h.SearchQuery)
	} else {
		leftStr = DefaultTheme.TitleStyle.Render(h.Msg.AppName) +
			DefaultTheme.SubduedStyle.Render(" "+h.Msg.Version)
	}

	// 右侧：选择统计
	rightStr := h.renderSelectionInfo()

	// 组合左右两侧
	leftPlain := lipgloss.NewStyle().Render(leftStr)
	rightPlain := lipgloss.NewStyle().Render(rightStr)

	leftLen := lipgloss.Width(leftPlain)
	rightLen := lipgloss.Width(rightPlain)

	padding := h.Width - leftLen - rightLen
	if padding < 1 {
		padding = 1
	}

	return leftStr + strings.Repeat(" ", padding) + rightStr
}

// renderSelectionInfo 渲染右侧选择统计信息
func (h *Header) renderSelectionInfo() string {
	if h.Selection.Len() == 0 || h.Msg == nil {
		return ""
	}

	count := h.Selection.Len()
	return DefaultTheme.SelectionStyle.Render(
		fmt.Sprintf(h.Msg.HeaderSelectedCount, count),
	)
}

// RenderWithSize 渲染包含选中文件大小的 Header
// totalSize: 所有选中文件的总字节数
func (h *Header) RenderWithSize(totalSize int64) string {
	if h.Width <= 0 || h.Msg == nil {
		return ""
	}

	leftStr := ""
	if h.IsSearching {
		leftStr = DefaultTheme.SearchStyle.Render(h.Msg.HeaderSearchLabel) +
			DefaultTheme.TitleStyle.Render(h.SearchQuery)
	} else {
		leftStr = DefaultTheme.TitleStyle.Render(h.Msg.AppName) +
			DefaultTheme.SubduedStyle.Render(" "+h.Msg.Version)
	}

	rightStr := ""
	if h.Selection.Len() > 0 {
		sizeStr := fileops.FormatSizeTotal(totalSize)
		rightStr = DefaultTheme.SelectionStyle.Render(
			fmt.Sprintf(h.Msg.HeaderSelectedCount, h.Selection.Len()) + " (" + sizeStr + ")",
		)
	}

	leftLen := lipgloss.Width(leftStr)
	rightLen := lipgloss.Width(rightStr)

	padding := h.Width - leftLen - rightLen
	if padding < 1 {
		padding = 1
	}

	return leftStr + strings.Repeat(" ", padding) + rightStr
}
