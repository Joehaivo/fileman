package ui

import (
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
)

// truncatePathHead 省略路径头部，保留尾部
// 例如: "/very/long/path/to/file.txt" -> "…path/to/file.txt"
func truncatePathHead(path string, maxDisplayWidth int) string {
	if path == "" {
		return ""
	}

	// 使用 lipgloss.Width 计算显示宽度
	if lipgloss.Width(path) <= maxDisplayWidth {
		return path
	}

	// 从头部开始省略
	// 先尝试只保留文件名
	filename := filepath.Base(path)
	ellipsis := "…"
	ellipsisWidth := lipgloss.Width(ellipsis)

	if lipgloss.Width(filename)+ellipsisWidth <= maxDisplayWidth {
		return ellipsis + filename
	}

	// 如果连文件名都太长，直接截断
	if lipgloss.Width(filename) > maxDisplayWidth-ellipsisWidth {
		// 截断文件名
		runes := []rune(filename)
		for len(runes) > 0 && lipgloss.Width(string(runes))+ellipsisWidth > maxDisplayWidth {
			runes = runes[1:]
		}
		if len(runes) > 0 {
			return ellipsis + string(runes)
		}
		return ellipsis
	}

	// 保留部分路径
	dir := filepath.Dir(path)
	parts := strings.Split(dir, string(filepath.Separator))

	// 从前往后逐步添加路径组件，直到达到最大宽度
	result := filename
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "" {
			continue
		}
		newResult := parts[i] + string(filepath.Separator) + result
		if lipgloss.Width(newResult)+ellipsisWidth <= maxDisplayWidth {
			result = newResult
		} else {
			break
		}
	}

	return ellipsis + result
}