package ui

import (
	"charm.land/lipgloss/v2"
)

// Toast Toast 通知组件
type Toast struct {
	Message string
	Width   int
}

// NewToast 创建新的 Toast 组件
func NewToast(message string) *Toast {
	return &Toast{
		Message: message,
	}
}

// SetWidth 设置宽度
func (t *Toast) SetWidth(width int) {
	t.Width = width
}

// Render 渲染 Toast 组件
func (t *Toast) Render() string {
	if t.Message == "" {
		return ""
	}

	content := t.Message

	// Toast 样式：成功绿色背景，圆角边框
	// 使用 lipgloss.Renderer 确保边框正确渲染
	style := lipgloss.NewStyle().
		Foreground(ColorBackground).
		Background(ColorSuccess).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSuccess).
		BorderBackground(ColorSuccess)

	return style.Render(content)
}

// RenderToast 渲染 Toast 的便捷函数
func RenderToast(message string, maxWidth int) string {
	toast := NewToast(message)
	if maxWidth > 0 {
		toast.SetWidth(maxWidth)
	}
	return toast.Render()
}

// RenderSuccessToast 渲染成功 Toast
func RenderSuccessToast(message string, maxWidth int) string {
	return RenderToast("✓ "+message, maxWidth)
}

// RenderErrorToast 渲染错误 Toast
func RenderErrorToast(message string, maxWidth int) string {
	// 错误 Toast 使用红色背景
	style := lipgloss.NewStyle().
		Foreground(ColorBackground).
		Background(ColorError).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorError).
		BorderBackground(ColorError)

	content := "✗ " + message
	return style.Render(content)
}