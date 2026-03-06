package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/haivo/fileman/internal/fileops"
	"github.com/haivo/fileman/internal/types"
)

const (
	appName    = "文件管家"
	appVersion = "v0.1.0"
)

// Header 顶部标题组件
type Header struct {
	Width       int                // 可用宽度
	Selection   types.SelectionSet // 当前选择集
	IsSearching bool               // 是否处于搜索模式
	SearchQuery string             // 搜索关键词
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
	if h.Width <= 0 {
		return ""
	}

	// 左侧：应用名称 + 版本
	leftStr := ""
	if h.IsSearching {
		leftStr = DefaultTheme.SearchStyle.Render("搜索: ") +
			DefaultTheme.TitleStyle.Render(h.SearchQuery)
	} else {
		leftStr = DefaultTheme.TitleStyle.Render(appName) +
			DefaultTheme.SubduedStyle.Render(" "+appVersion)
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
	if h.Selection.Len() == 0 {
		return ""
	}

	// 统计选中条目的总大小（此处无法直接获取文件大小，仅显示数量）
	// 在实际应用中，选择集应包含 FileEntry 以便计算大小
	count := h.Selection.Len()
	return DefaultTheme.SelectionStyle.Render(
		fmt.Sprintf("已选: %d 个", count),
	)
}

// RenderWithSize 渲染包含选中文件大小的 Header
// totalSize: 所有选中文件的总字节数
func (h *Header) RenderWithSize(totalSize int64) string {
	if h.Width <= 0 {
		return ""
	}

	leftStr := ""
	if h.IsSearching {
		leftStr = DefaultTheme.SearchStyle.Render("搜索: ") +
			DefaultTheme.TitleStyle.Render(h.SearchQuery)
	} else {
		leftStr = DefaultTheme.TitleStyle.Render(appName) +
			DefaultTheme.SubduedStyle.Render(" "+appVersion)
	}

	rightStr := ""
	if h.Selection.Len() > 0 {
		sizeStr := fileops.FormatSizeTotal(totalSize)
		rightStr = DefaultTheme.SelectionStyle.Render(
			fmt.Sprintf("已选: %d 个 (%s)", h.Selection.Len(), sizeStr),
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
