package fileops

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Joehaivo/fileman/internal/types"
)

// ScanDir 扫描目录，返回排序后的文件条目列表（目录优先，按名称排序）
// path: 要扫描的目录路径
// showHidden: 是否显示隐藏文件（以 . 开头的文件）
func ScanDir(path string, showHidden bool) ([]types.FileEntry, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %w", err)
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	var files []types.FileEntry
	for _, entry := range entries {
		name := entry.Name()

		// 过滤隐藏文件（以 . 开头的文件）
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileEntry := buildFileEntry(entry, info, absPath)
		files = append(files, fileEntry)
	}

	// 排序：目录优先，然后按名称字母序
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return files, nil
}

// buildFileEntry 根据目录条目和文件信息构建 FileEntry
func buildFileEntry(entry fs.DirEntry, info fs.FileInfo, parentPath string) types.FileEntry {
	name := entry.Name()
	fullPath := filepath.Join(parentPath, name)

	fileType := types.FileTypeRegular
	isDir := entry.IsDir()

	if isDir {
		fileType = types.FileTypeDirectory
	} else if entry.Type()&fs.ModeSymlink != 0 {
		fileType = types.FileTypeSymlink
	} else if !entry.Type().IsRegular() {
		fileType = types.FileTypeOther
	}

	ext := ""
	if !isDir {
		ext = strings.ToLower(filepath.Ext(name))
	}

	return types.FileEntry{
		Name:    name,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Mode:    info.Mode().String(),
		Type:    fileType,
		IsDir:   isDir,
		Ext:     ext,
		Path:    fullPath,
	}
}

// FormatSize 格式化文件大小为人类可读格式（B/KB/MB/GB）
func FormatSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%dB", size)
	case size < MB:
		return fmt.Sprintf("%.1fK", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.1fM", float64(size)/MB)
	default:
		return fmt.Sprintf("%.1fG", float64(size)/GB)
	}
}

// FormatDate 格式化日期为简短格式
func FormatDate(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() {
		return t.Format("01-02 15:04")
	}
	return t.Format("2006-01-02")
}

// FormatSizeTotal 格式化总大小（用于多选统计）
func FormatSizeTotal(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.1f KB", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.1f MB", float64(size)/MB)
	default:
		return fmt.Sprintf("%.1f GB", float64(size)/GB)
	}
}

// GetFileTypeDesc 获取文件类型描述
// useEnglish: 是否使用英文描述
func GetFileTypeDesc(entry types.FileEntry, useEnglish bool) string {
	if entry.IsDir {
		if useEnglish {
			return "Directory"
		}
		return "目录"
	}

	if useEnglish {
		extDesc := map[string]string{
			".go":    "Go source",
			".js":    "JavaScript",
			".ts":    "TypeScript",
			".tsx":   "TypeScript JSX",
			".jsx":   "JavaScript JSX",
			".py":    "Python script",
			".rs":    "Rust source",
			".c":     "C source",
			".cpp":   "C++ source",
			".h":     "C/C++ header",
			".java":  "Java source",
			".kt":    "Kotlin source",
			".swift": "Swift source",
			".rb":    "Ruby script",
			".php":   "PHP file",
			".sh":    "Shell script",
			".bash":  "Bash script",
			".zsh":   "Zsh script",
			".json":  "JSON data",
			".yaml":  "YAML config",
			".yml":   "YAML config",
			".toml":  "TOML config",
			".ini":   "INI config",
			".conf":  "Config file",
			".md":    "Markdown",
			".txt":   "Text file",
			".pdf":   "PDF document",
			".png":   "PNG image",
			".jpg":   "JPEG image",
			".jpeg":  "JPEG image",
			".gif":   "GIF image",
			".svg":   "SVG image",
			".mp4":   "MP4 video",
			".mp3":   "MP3 audio",
			".zip":   "ZIP archive",
			".tar":   "TAR archive",
			".tar.gz": "TAR.GZ archive",
			".tgz":    "TAR.GZ archive",
			".tar.bz2": "TAR.BZ2 archive",
			".tbz2":    "TAR.BZ2 archive",
			".tar.xz":  "TAR.XZ archive",
			".txz":     "TAR.XZ archive",
			".gz":    "Gzip file",
			".bz2":   "BZIP2 archive",
			".xz":    "XZ archive",
			".7z":    "7-Zip archive",
			".rar":   "RAR archive",
			".sql":   "SQL script",
			".html":  "HTML file",
			".css":   "CSS stylesheet",
		}

		if desc, ok := extDesc[entry.Ext]; ok {
			return desc
		}

		if entry.Ext != "" {
			return strings.ToUpper(strings.TrimPrefix(entry.Ext, ".")) + " file"
		}

		return "File"
	}

	extDesc := map[string]string{
		".go":    "Go 源文件",
		".js":    "JavaScript 文件",
		".ts":    "TypeScript 文件",
		".tsx":   "TypeScript JSX 文件",
		".jsx":   "JavaScript JSX 文件",
		".py":    "Python 脚本",
		".rs":    "Rust 源文件",
		".c":     "C 源文件",
		".cpp":   "C++ 源文件",
		".h":     "C/C++ 头文件",
		".java":  "Java 源文件",
		".kt":    "Kotlin 源文件",
		".swift": "Swift 源文件",
		".rb":    "Ruby 脚本",
		".php":   "PHP 文件",
		".sh":    "Shell 脚本",
		".bash":  "Bash 脚本",
		".zsh":   "Zsh 脚本",
		".json":  "JSON 数据文件",
		".yaml":  "YAML 配置文件",
		".yml":   "YAML 配置文件",
		".toml":  "TOML 配置文件",
		".ini":   "INI 配置文件",
		".conf":  "配置文件",
		".md":    "Markdown 文档",
		".txt":   "文本文件",
		".pdf":   "PDF 文档",
		".png":   "PNG 图片",
		".jpg":   "JPEG 图片",
		".jpeg":  "JPEG 图片",
		".gif":   "GIF 图片",
		".svg":   "SVG 矢量图",
		".mp4":   "MP4 视频",
		".mp3":   "MP3 音频",
		".zip":   "ZIP 压缩包",
		".tar":   "TAR 归档",
		".tar.gz": "TAR.GZ 压缩包",
		".tgz":    "TAR.GZ 压缩包",
		".tar.bz2": "TAR.BZ2 压缩包",
		".tbz2":    "TAR.BZ2 压缩包",
		".tar.xz":  "TAR.XZ 压缩包",
		".txz":     "TAR.XZ 压缩包",
		".gz":    "Gzip 压缩文件",
		".bz2":   "BZIP2 压缩包",
		".xz":    "XZ 压缩包",
		".7z":    "7-Zip 压缩包",
		".rar":   "RAR 压缩包",
		".sql":   "SQL 脚本",
		".html":  "HTML 文件",
		".css":   "CSS 样式文件",
	}

	if desc, ok := extDesc[entry.Ext]; ok {
		return desc
	}

	if entry.Ext != "" {
		return strings.ToUpper(strings.TrimPrefix(entry.Ext, ".")) + " 文件"
	}

	return "文件"
}
