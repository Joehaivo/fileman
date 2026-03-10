package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
)

// OverlayTopRight 将 overlay 叠加在 base 内容的右上角
func OverlayTopRight(base, overlay string, baseWidth int) string {
	if overlay == "" {
		return base
	}

	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	// 移除末尾空行
	if len(baseLines) > 0 && baseLines[len(baseLines)-1] == "" {
		baseLines = baseLines[:len(baseLines)-1]
	}
	if len(overlayLines) > 0 && overlayLines[len(overlayLines)-1] == "" {
		overlayLines = overlayLines[:len(overlayLines)-1]
	}

	if len(baseLines) == 0 {
		return overlay
	}

	// 计算 overlay 的最大宽度
	overlayMaxWidth := 0
	for _, line := range overlayLines {
		w := lipgloss.Width(line)
		if w > overlayMaxWidth {
			overlayMaxWidth = w
		}
	}

	// 计算叠加的起始列位置（右上角）
	startCol := baseWidth - overlayMaxWidth
	if startCol < 0 {
		startCol = 0
	}

	// 叠加 overlay 到 base 的右上角
	result := make([]string, len(baseLines))
	for i := range baseLines {
		baseLine := baseLines[i]

		if i < len(overlayLines) {
			overlayLine := overlayLines[i]
			overlayWidth := lipgloss.Width(overlayLine)

			// 计算这一行 overlay 的起始位置
			lineStartCol := baseWidth - overlayWidth
			if lineStartCol < 0 {
				lineStartCol = 0
			}

			// 截断 base 行到起始位置
			truncatedBase := truncateStringByWidth(baseLine, lineStartCol)

			// 直接拼接
			result[i] = truncatedBase + overlayLine
		} else {
			result[i] = baseLine
		}
	}

	return strings.Join(result, "\n")
}

// truncateStringByWidth 按显示宽度截断字符串（保留 ANSI 转义序列）
func truncateStringByWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	// 使用 lipgloss.Width 来测量，但需要保留 ANSI 序列
	// 简单方法：逐字符处理，跳过 ANSI 序列
	var result strings.Builder
	currentWidth := 0
	inEscape := false

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// 检测 ANSI 转义序列
		if r == '\x1b' {
			inEscape = true
			result.WriteRune(r)
			continue
		}

		if inEscape {
			result.WriteRune(r)
			// ANSI 序列以字母结束
			if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' {
				inEscape = false
			}
			continue
		}

		// 计算显示宽度
		rw := runewidth.RuneWidth(r)
		if currentWidth+rw > maxWidth {
			break
		}
		result.WriteRune(r)
		currentWidth += rw
	}

	return result.String()
}