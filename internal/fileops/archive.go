package fileops

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Joehaivo/fileman/internal/types"
	"github.com/mholt/archives"
)

// ArchiveEntry 表示压缩包内的一个条目
type ArchiveEntry struct {
	Name    string // 文件名（相对路径）
	Size    int64  // 文件大小
	IsDir   bool   // 是否为目录
	ModTime string // 修改时间（可选）
}

// 支持的压缩格式映射
var archiveFormats = map[string]string{
	".zip":     "ZIP",
	".tar":     "TAR",
	".tar.gz":  "TAR.GZ",
	".tgz":     "TAR.GZ",
	".tar.bz2": "TAR.BZ2",
	".tbz2":    "TAR.BZ2",
	".tar.xz":  "TAR.XZ",
	".txz":     "TAR.XZ",
	".gz":      "GZIP",
	".bz2":     "BZIP2",
	".xz":      "XZ",
	".7z":      "7Z",
	".rar":     "RAR",
}

// IsArchive 检测文件扩展名是否为支持的压缩格式
// 返回: (是否为压缩文件, 格式名称)
func IsArchive(ext string) (bool, string) {
	// 先检查复合扩展名（如 .tar.gz）
	name := strings.ToLower(ext)
	if format, ok := archiveFormats[name]; ok {
		return true, format
	}
	return false, ""
}

// IsArchiveEntry 检测 FileEntry 是否为压缩文件
func IsArchiveEntry(entry types.FileEntry) (bool, string) {
	if entry.IsDir {
		return false, ""
	}
	return IsArchive(entry.Ext)
}

// ReadArchivePreview 读取压缩包内容，返回条目列表
func ReadArchivePreview(entry types.FileEntry) *PreviewResult {
	if entry.IsDir {
		return &PreviewResult{Error: "目录无法预览"}
	}

	isArchive, format := IsArchiveEntry(entry)
	if !isArchive {
		return &PreviewResult{Error: "不是支持的压缩格式"}
	}

	// 打开文件
	f, err := os.Open(entry.Path)
	if err != nil {
		return &PreviewResult{Error: fmt.Sprintf("无法打开文件: %v", err)}
	}
	defer f.Close()

	var entries []ArchiveEntry
	var totalSize int64

	// 使用 archives 库识别格式
	ctx := context.Background()
	formatObj, reader, err := archives.Identify(ctx, entry.Path, f)
	if err != nil {
		return &PreviewResult{Error: fmt.Sprintf("无法识别压缩格式: %v", err)}
	}

	// 根据格式类型提取文件列表
	switch ft := formatObj.(type) {
	case archives.Extractor:
		// 提取文件列表（使用 Identify 返回的 reader）
		err := ft.Extract(ctx, reader, func(ctx context.Context, fileInfo archives.FileInfo) error {
			entries = append(entries, ArchiveEntry{
				Name:    fileInfo.NameInArchive,
				Size:    fileInfo.Size(),
				IsDir:   fileInfo.IsDir(),
				ModTime: fileInfo.ModTime().Format("2006-01-02 15:04"),
			})
			if !fileInfo.IsDir() {
				totalSize += fileInfo.Size()
			}
			return nil
		})
		if err != nil {
			return &PreviewResult{Error: fmt.Sprintf("读取压缩包失败: %v", err)}
		}
	default:
		return &PreviewResult{Error: "不支持的压缩格式"}
	}

	// 排序：目录优先，然后按名称排序
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return entries[i].Name < entries[j].Name
	})

	// 构建树形显示
	lines := buildArchiveTree(entries)

	return &PreviewResult{
		Lines:         lines,
		TotalLines:    len(lines),
		IsArchive:     true,
		ArchiveFormat: format,
		ArchiveCount:  countFiles(entries),
		ArchiveSize:   totalSize,
	}
}

// countFiles 统计文件数量（不含目录）
func countFiles(entries []ArchiveEntry) int {
	count := 0
	for _, e := range entries {
		if !e.IsDir {
			count++
		}
	}
	return count
}

// buildArchiveTree 构建树形显示
func buildArchiveTree(entries []ArchiveEntry) []string {
	if len(entries) == 0 {
		return []string{"(空压缩包)"}
	}

	var lines []string

	// 构建树形结构
	tree := buildTree(entries)
	lines = renderTree(tree, "", true)

	return lines
}

// TreeNode 树形节点
type TreeNode struct {
	Name     string
	Size     int64
	IsDir    bool
	Children []*TreeNode
}

// buildTree 将扁平列表转换为树形结构
func buildTree(entries []ArchiveEntry) *TreeNode {
	root := &TreeNode{Name: "", IsDir: true, Children: []*TreeNode{}}

	for _, entry := range entries {
		parts := strings.Split(strings.Trim(entry.Name, "/"), "/")
		current := root

		for i, part := range parts {
			if part == "" {
				continue
			}

			isLast := i == len(parts)-1
			isDir := entry.IsDir || !isLast

			// 查找或创建子节点
			found := false
			for _, child := range current.Children {
				if child.Name == part {
					current = child
					found = true
					break
				}
			}

			if !found {
				node := &TreeNode{
					Name:  part,
					IsDir: isDir,
				}
				if isLast && !entry.IsDir {
					node.Size = entry.Size
				}
				current.Children = append(current.Children, node)
				current = node
			}
		}
	}

	// 对子节点排序
	sortChildren(root)

	return root
}

// sortChildren 递归排序子节点
func sortChildren(node *TreeNode) {
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	for _, child := range node.Children {
		sortChildren(child)
	}
}

// renderTree 渲染树形结构
func renderTree(node *TreeNode, prefix string, isLast bool) []string {
	var lines []string

	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		connector := "├── "
		newPrefix := prefix + "│   "
		if isLastChild {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		// 图标
		icon := ""
		if child.IsDir {
			icon = "" // 目录图标
		} else {
			icon = "" // 文件图标，使用 getFileIcon
		}

		// 文件名和大小
		name := child.Name
		if child.IsDir {
			name += "/"
		}
		sizeStr := ""
		if !child.IsDir {
			sizeStr = fmt.Sprintf(" (%s)", FormatSize(child.Size))
		}

		line := prefix + connector + icon + name + sizeStr
		lines = append(lines, line)

		// 递归渲染子节点
		if child.IsDir && len(child.Children) > 0 {
			lines = append(lines, renderTree(child, newPrefix, isLastChild)...)
		}
	}

	return lines
}

// GetArchiveFormatDesc 获取压缩格式的描述
func GetArchiveFormatDesc(format string, useEnglish bool) string {
	if useEnglish {
		return format + " Archive"
	}
	return format + " 压缩包"
}

// GetArchiveExtensions 返回支持的压缩格式扩展名列表
func GetArchiveExtensions() []string {
	exts := make([]string, 0, len(archiveFormats))
	for ext := range archiveFormats {
		exts = append(exts, ext)
	}
	sort.Strings(exts)
	return exts
}

// DetectArchiveFormat 从文件路径检测压缩格式
func DetectArchiveFormat(path string) string {
	name := strings.ToLower(filepath.Base(path))

	// 先检查复合扩展名（按长度降序排列，确保 .tar.gz 先于 .gz 匹配）
	compositeExts := []string{
		".tar.gz", ".tar.bz2", ".tar.xz",
		".tgz", ".tbz2", ".txz",
		".zip", ".tar", ".gz", ".bz2", ".xz",
		".7z", ".rar",
	}

	for _, ext := range compositeExts {
		if strings.HasSuffix(name, ext) {
			if format, ok := archiveFormats[ext]; ok {
				return format
			}
		}
	}

	return ""
}