package fileops

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Joehaivo/fileman/internal/types"
)

func TestReadPreview_Chinese(t *testing.T) {
	// Create a temporary file with Chinese content
	content := []byte("# FileMan TUI 开发进度追踪\n\n## 项目信息\n\n- **二进制名**: `fm`\n- **模块路径**: `github.com/Joehaivo/fileman`")
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_chinese.md")
	err := os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	entry := types.FileEntry{
		Path: tmpFile,
		Size: int64(len(content)),
	}

	result := ReadPreview(entry)
	if result.IsBinary {
		t.Errorf("Expected text file, but detected as binary")
	}
	if len(result.Lines) == 0 {
		t.Errorf("Expected content, but got empty lines")
	}
}

func TestReadPreview_Binary(t *testing.T) {
	// Create a temporary binary file
	content := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_binary.bin")
	err := os.WriteFile(tmpFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	entry := types.FileEntry{
		Path: tmpFile,
		Size: int64(len(content)),
	}

	result := ReadPreview(entry)
	if !result.IsBinary {
		t.Errorf("Expected binary file, but detected as text")
	}
}

func TestIsBinary_TruncatedUTF8(t *testing.T) {
	// "中文" bytes
	content := []byte("中文")
	// "中" is e4 b8 ad
	// "文" is e6 96 87

	// Case 1: Full content
	if isBinary(content) {
		t.Errorf("Full Chinese content should be text")
	}

	// Case 2: Truncated "文" (first byte only) -> e6
	truncated := content[:4] // e4 b8 ad e6
	if isBinary(truncated) {
		t.Errorf("Truncated UTF-8 at end should be text")
	}

	// Case 3: Truncated "文" (first two bytes) -> e6 96
	truncated2 := content[:5] // e4 b8 ad e6 96
	if isBinary(truncated2) {
		t.Errorf("Truncated UTF-8 at end (2 bytes) should be text")
	}
}
