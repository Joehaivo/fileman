package fileops

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/Joehaivo/fileman/internal/types"
)

// TestIsArchive 测试压缩格式检测
func TestIsArchive(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
		format   string
	}{
		{".zip", true, "ZIP"},
		{".tar", true, "TAR"},
		{".tar.gz", true, "TAR.GZ"},
		{".tgz", true, "TAR.GZ"},
		{".tar.bz2", true, "TAR.BZ2"},
		{".tbz2", true, "TAR.BZ2"},
		{".tar.xz", true, "TAR.XZ"},
		{".txz", true, "TAR.XZ"},
		{".gz", true, "GZIP"},
		{".bz2", true, "BZIP2"},
		{".xz", true, "XZ"},
		{".7z", true, "7Z"},
		{".rar", true, "RAR"},
		{".txt", false, ""},
		{".go", false, ""},
		{".md", false, ""},
		{"", false, ""},
	}

	for _, test := range tests {
		isArchive, format := IsArchive(test.ext)
		if isArchive != test.expected {
			t.Errorf("IsArchive(%q) = %v, want %v", test.ext, isArchive, test.expected)
		}
		if isArchive && format != test.format {
			t.Errorf("IsArchive(%q) format = %q, want %q", test.ext, format, test.format)
		}
	}
}

// TestIsArchiveEntry 测试 FileEntry 压缩格式检测
func TestIsArchiveEntry(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		isDir    bool
		expected bool
	}{
		{"test.zip", ".zip", false, true},
		{"test.tar.gz", ".tar.gz", false, true},
		{"test.txt", ".txt", false, false},
		{"folder", "", true, false},
	}

	for _, test := range tests {
		entry := types.FileEntry{
			Name:  test.name,
			Ext:   test.ext,
			IsDir: test.isDir,
		}
		isArchive, _ := IsArchiveEntry(entry)
		if isArchive != test.expected {
			t.Errorf("IsArchiveEntry(%+v) = %v, want %v", entry, isArchive, test.expected)
		}
	}
}

// TestReadArchivePreview_ZIP 测试 ZIP 文件预览
func TestReadArchivePreview_ZIP(t *testing.T) {
	// 创建临时测试 ZIP 文件
	tmpDir, err := os.MkdirTemp("", "archive_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "test.zip")

	// 创建测试 ZIP 文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("创建 ZIP 文件失败: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)

	// 添加文件
	files := []struct {
		name    string
		content string
	}{
		{"file1.txt", "Hello World"},
		{"dir/file2.txt", "Test Content"},
	}

	for _, f := range files {
		w, err := zipWriter.Create(f.name)
		if err != nil {
			t.Fatalf("创建 ZIP 条目失败: %v", err)
		}
		_, err = w.Write([]byte(f.content))
		if err != nil {
			t.Fatalf("写入 ZIP 条目失败: %v", err)
		}
	}

	zipWriter.Close()
	zipFile.Close()

	// 获取文件信息
	info, err := os.Stat(zipPath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	// 测试预览
	entry := types.FileEntry{
		Name:  "test.zip",
		Ext:   ".zip",
		Path:  zipPath,
		Size:  info.Size(),
		IsDir: false,
	}

	result := ReadArchivePreview(entry)

	// 验证结果
	if result.Error != "" {
		t.Errorf("ReadArchivePreview 返回错误: %s", result.Error)
	}

	if !result.IsArchive {
		t.Error("IsArchive 应该为 true")
	}

	if result.ArchiveFormat != "ZIP" {
		t.Errorf("ArchiveFormat = %q, want %q", result.ArchiveFormat, "ZIP")
	}

	if result.ArchiveCount != 2 {
		t.Errorf("ArchiveCount = %d, want 2", result.ArchiveCount)
	}

	if len(result.Lines) == 0 {
		t.Error("Lines 不应该为空")
	}
}

// TestReadArchivePreview_EmptyArchive 测试空压缩包
func TestReadArchivePreview_EmptyArchive(t *testing.T) {
	// 创建临时测试 ZIP 文件
	tmpDir, err := os.MkdirTemp("", "archive_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, "empty.zip")

	// 创建空 ZIP 文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("创建 ZIP 文件失败: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)
	zipWriter.Close()
	zipFile.Close()

	// 获取文件信息
	info, err := os.Stat(zipPath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	// 测试预览
	entry := types.FileEntry{
		Name:  "empty.zip",
		Ext:   ".zip",
		Path:  zipPath,
		Size:  info.Size(),
		IsDir: false,
	}

	result := ReadArchivePreview(entry)

	// 验证结果
	if result.Error != "" {
		t.Errorf("ReadArchivePreview 返回错误: %s", result.Error)
	}

	if result.ArchiveCount != 0 {
		t.Errorf("ArchiveCount = %d, want 0", result.ArchiveCount)
	}

	if len(result.Lines) == 0 {
		t.Error("空压缩包应该有提示行")
	}
}

// TestReadArchivePreview_NotArchive 测试非压缩文件
func TestReadArchivePreview_NotArchive(t *testing.T) {
	// 创建临时测试文件
	tmpDir, err := os.MkdirTemp("", "archive_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(filePath, []byte("Hello World"), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 获取文件信息
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	// 测试预览
	entry := types.FileEntry{
		Name:  "test.txt",
		Ext:   ".txt",
		Path:  filePath,
		Size:  info.Size(),
		IsDir: false,
	}

	result := ReadArchivePreview(entry)

	// 验证结果
	if result.Error == "" {
		t.Error("非压缩文件应该返回错误")
	}
}

// TestBuildArchiveTree 测试树形结构构建
func TestBuildArchiveTree(t *testing.T) {
	entries := []ArchiveEntry{
		{Name: "file1.txt", Size: 100, IsDir: false},
		{Name: "dir/file2.txt", Size: 200, IsDir: false},
		{Name: "dir/subdir/file3.txt", Size: 300, IsDir: false},
		{Name: "dir2/", IsDir: true},
	}

	lines := buildArchiveTree(entries)

	if len(lines) == 0 {
		t.Error("树形结构不应该为空")
	}

	// 检查是否包含文件名
	found := false
	for _, line := range lines {
		if line != "" {
			found = true
			break
		}
	}
	if !found {
		t.Error("树形结构应该包含文件名")
	}
}

// TestGetArchiveFormatDesc 测试格式描述
func TestGetArchiveFormatDesc(t *testing.T) {
	tests := []struct {
		format     string
		useEnglish bool
		expected   string
	}{
		{"ZIP", false, "ZIP 压缩包"},
		{"TAR.GZ", false, "TAR.GZ 压缩包"},
		{"ZIP", true, "ZIP Archive"},
		{"TAR.GZ", true, "TAR.GZ Archive"},
	}

	for _, test := range tests {
		result := GetArchiveFormatDesc(test.format, test.useEnglish)
		if result != test.expected {
			t.Errorf("GetArchiveFormatDesc(%q, %v) = %q, want %q",
				test.format, test.useEnglish, result, test.expected)
		}
	}
}

// TestGetArchiveExtensions 测试获取支持的扩展名
func TestGetArchiveExtensions(t *testing.T) {
	exts := GetArchiveExtensions()

	if len(exts) == 0 {
		t.Error("支持的扩展名列表不应该为空")
	}

	// 检查常见格式是否包含
	expectedExts := []string{".zip", ".tar", ".7z"}
	for _, expected := range expectedExts {
		found := false
		for _, ext := range exts {
			if ext == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("扩展名 %q 应该在支持列表中", expected)
		}
	}
}

// TestDetectArchiveFormat 测试从路径检测格式
func TestDetectArchiveFormat(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/file.zip", "ZIP"},
		{"/path/to/file.tar.gz", "TAR.GZ"},
		{"/path/to/file.7z", "7Z"},
		{"/path/to/file.txt", ""},
	}

	for _, test := range tests {
		result := DetectArchiveFormat(test.path)
		if result != test.expected {
			t.Errorf("DetectArchiveFormat(%q) = %q, want %q",
				test.path, result, test.expected)
		}
	}
}