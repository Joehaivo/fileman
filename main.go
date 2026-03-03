package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/haivo/fileman/internal/app"
)

func main() {
	// 初始化 Model（初始命令通过 Init() 方法返回）
	model, _ := app.New()

	// 创建 Bubble Tea 程序，启用 AltScreen 全屏模式
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// 运行程序
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行错误: %v\n", err)
		os.Exit(1)
	}
}
