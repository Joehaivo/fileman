package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	width := 10
	longWord := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	
	style := lipgloss.NewStyle().Width(width)
	rendered := style.Render(longWord)
	
	fmt.Println("--- Start ---")
	fmt.Println(rendered)
	fmt.Println("--- End ---")

	lines := strings.Split(rendered, "\n")
	for i, line := range lines {
		fmt.Printf("Line %d len: %d\n", i, len(line))
	}
}
