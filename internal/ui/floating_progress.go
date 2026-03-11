package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/Joehaivo/fileman/internal/types"
)

// FloatingProgressComponent 悬浮进度窗口组件
type FloatingProgressComponent struct {
	Progress *types.FloatingProgress
	Width    int
	Height   int
}

// NewFloatingProgress 创建新的悬浮进度窗口组件
func NewFloatingProgress(progress *types.FloatingProgress) *FloatingProgressComponent {
	return &FloatingProgressComponent{
		Progress: progress,
	}
}

// SetSize 设置尺寸
func (f *FloatingProgressComponent) SetSize(width, height int) {
	f.Width = width
	f.Height = height
}

// Render 渲染悬浮进度窗口
func (f *FloatingProgressComponent) Render() string {
	if f.Progress == nil {
		return ""
	}

	p := f.Progress

	// 计算进度百分比
	percent := 0.0
	if p.Total > 0 {
		percent = float64(p.Done) / float64(p.Total)
	}

	// 计算窗口内容宽度
	contentWidth := 40
	if f.Width > 0 && f.Width < contentWidth+4 {
		contentWidth = f.Width - 4
	}
	if contentWidth < 20 {
		contentWidth = 20
	}

	var content strings.Builder

	// 标题行
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorTitle).
		Width(contentWidth)

	title := p.OpType + "进度"
	if p.IsComplete {
		title = p.OpType + "完成"
	}
	titleText := titleStyle.Render(title)
	content.WriteString(titleText)
	content.WriteByte('\n')
	content.WriteByte('\n')

	// 进度行
	progressLine := formatProgressLine(percent, p.Done, p.Total, contentWidth)
	content.WriteString(progressLine)
	content.WriteByte('\n')
	content.WriteByte('\n')

	// 结果列表（最近5条）
	if len(p.Results) > 0 {
		// 限制显示条数
		maxResults := 5
		if f.Height > 0 && f.Height-8 < maxResults {
			maxResults = f.Height - 8
			if maxResults < 1 {
				maxResults = 1
			}
		}

		start := 0
		if len(p.Results) > maxResults {
			start = len(p.Results) - maxResults
		}

		for i := start; i < len(p.Results); i++ {
			r := p.Results[i]
			filename := TruncatePathHead(r.SrcPath, contentWidth-4)
			if r.Err != nil {
				line := DefaultTheme.ErrorStyle.Render("✗ ") + filename
				content.WriteString(lipgloss.NewStyle().Width(contentWidth).Render(line))
			} else {
				line := DefaultTheme.SuccessStyle.Render("✓ ") + filename
				content.WriteString(lipgloss.NewStyle().Width(contentWidth).Render(line))
			}
			if i < len(p.Results)-1 {
				content.WriteByte('\n')
			}
		}
	}

	// 边框样式
	borderColor := ColorBorderFocus
	if p.IsComplete {
		borderColor = ColorSuccess
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 1).
		Width(contentWidth + 2)

	return boxStyle.Render(content.String())
}

// formatProgressLine 格式化进度行
func formatProgressLine(percent float64, done, total, width int) string {
	// 进度条宽度
	barWidth := width - 15 // 为百分比留空间
	if barWidth < 10 {
		barWidth = 10
	}

	// 填充部分
	filled := int(float64(barWidth) * percent)
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := DefaultTheme.SuccessStyle.Render(strings.Repeat("█", filled)) +
		DefaultTheme.SubduedStyle.Render(strings.Repeat("░", empty))

	// 百分比和计数
	countStr := " " + itoa(done) + "/" + itoa(total)
	percentStr := itoa(int(percent*100)) + "%"

	return bar + " " + DefaultTheme.SubduedStyle.Render(countStr+" "+percentStr)
}

// RenderFloatingProgress 渲染悬浮进度窗口的便捷函数
func RenderFloatingProgress(progress *types.FloatingProgress, width, height int) string {
	f := NewFloatingProgress(progress)
	f.SetSize(width, height)
	return f.Render()
}