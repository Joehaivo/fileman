package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/Joehaivo/fileman/internal/fileops"
	"github.com/Joehaivo/fileman/internal/types"
)

func TestPreviewPane_Render_Wrap(t *testing.T) {
	pane := NewPreviewPane()
	pane.SetSize(20, 10) // Small width to force wrap

	entry := &types.FileEntry{
		Name:    "test.txt",
		Size:    100,
		ModTime: time.Now(),
		Mode:    "-rw-r--r--",
	}
	pane.Entry = entry

	// Manually set Result to avoid file IO and ensure content
	pane.Result = &fileops.PreviewResult{
		Lines: []string{
			"Line 1 short",
			"Line 2 is very very long and should wrap",
			"Line 3",
		},
		TotalLines: 3,
	}

	output := pane.Render()
	t.Logf("Output:\n%s", output)

	// Check if output contains expected content
	if !strings.Contains(output, "Line 1 short") {
		t.Errorf("Line 1 content missing")
	}
	if !strings.Contains(output, "Line 2 is very") {
		t.Errorf("Line 2 start missing")
	}
	if !strings.Contains(output, "Line 3") {
		t.Errorf("Line 3 content missing")
	}

	// Verify wrapping happened (output should have more lines than just title + 3 content + info)
	// Title: 1 line
	// Content: 3 lines logical.
	// Separator: 1 line
	// Info: 5 lines
	// Total base: 10 lines.
	// If wrapping happens, content takes more lines.
	// Count newlines.
	lineCount := strings.Count(output, "\n")
	if lineCount <= 10 {
		t.Errorf("Expected wrapping to increase line count, got %d", lineCount)
	}
}
