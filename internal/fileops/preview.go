package fileops

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/Joehaivo/fileman/internal/types"
)

const (
	// MaxPreviewSize 最大预览文件大小（1MB）
	MaxPreviewSize = 1 * 1024 * 1024
	// MaxPreviewLines 最大预览行数
	MaxPreviewLines = 10000
)

// PreviewResult 预览结果
type PreviewResult struct {
	Lines         []string // 文本行内容
	TotalLines    int      // 总行数
	IsBinary      bool     // 是否为二进制文件
	IsTooLarge    bool     // 是否超出大小限制
	FileSize      int64    // 文件大小（用于显示）
	Error         string   // 错误信息
	IsArchive     bool     // 是否为压缩文件
	ArchiveFormat string   // 压缩格式
	ArchiveCount  int      // 压缩包内文件数
	ArchiveSize   int64    // 压缩包解压后大小
}

// ReadPreview 读取文件预览内容，限制最大 1MB
// path: 文件路径
func ReadPreview(entry types.FileEntry) *PreviewResult {
	if entry.IsDir {
		return &PreviewResult{Error: "目录无法预览"}
	}

	// 检查是否为压缩文件
	if isArchive, _ := IsArchiveEntry(entry); isArchive {
		return ReadArchivePreview(entry)
	}

	if entry.Size > MaxPreviewSize {
		return &PreviewResult{
			IsTooLarge: true,
			FileSize:   entry.Size,
		}
	}

	f, err := os.Open(entry.Path)
	if err != nil {
		return &PreviewResult{Error: fmt.Sprintf("无法打开文件: %v", err)}
	}
	defer f.Close()

	// 读取前 8192 字节检测是否为文本文件（增加检测范围提高准确性）
	headerSize := 8192
	if entry.Size < int64(headerSize) {
		headerSize = int(entry.Size)
	}
	header := make([]byte, headerSize)
	n, err := f.Read(header)
	if err != nil && err != io.EOF {
		return &PreviewResult{Error: fmt.Sprintf("读取文件失败: %v", err)}
	}
	header = header[:n]

	if isBinary(header) {
		return &PreviewResult{IsBinary: true}
	}

	// 重新从头读取
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return &PreviewResult{Error: "重置文件位置失败"}
	}

	reader := io.LimitReader(f, MaxPreviewSize)
	scanner := bufio.NewScanner(reader)

	// 扩大 scanner buffer 处理长行
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1*1024*1024)

	var lines []string
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		if lineCount <= MaxPreviewLines {
			line := scanner.Text()
			// 将 Tab 替换为 4 个空格
			line = strings.ReplaceAll(line, "\t", "    ")
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		// 可能是 UTF-8 解码错误，标记为二进制
		if strings.Contains(err.Error(), "token too long") || !utf8.Valid([]byte(err.Error())) {
			return &PreviewResult{IsBinary: true}
		}
	}

	return &PreviewResult{
		Lines:      lines,
		TotalLines: lineCount,
	}
}

// isBinary 检测字节序列是否为二进制内容
// 通过检查是否存在 null 字节、大量非 UTF-8 字节或控制字符来判断
func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// 允许的文本控制字符（ASCII 0-31中常见的文本字符）
	allowedControlChars := map[byte]bool{
		0x09: true, // Tab
		0x0A: true, // LF (换行)
		0x0D: true, // CR (回车)
		0x1B: true, // ESC (ANSI转义序列)
		0x0C: true, // FF (换页)
	}

	nullCount := 0
	invalidControlCount := 0
	totalBytes := len(data)

	// 检查前1024字节（或全部，如果文件更小）
	checkSize := totalBytes
	if checkSize > 1024 {
		checkSize = 1024
	}

	for i := 0; i < checkSize; i++ {
		b := data[i]

		// null 字节通常是二进制文件的标志
		if b == 0 {
			nullCount++
			// 如果null字节超过1%，很可能是二进制文件
			if nullCount > checkSize/100 {
				return true
			}
		}

		// 检查控制字符（0-31，除了允许的）
		if b < 32 && !allowedControlChars[b] {
			invalidControlCount++
		}
	}

	// 如果无效控制字符超过5%，可能是二进制文件
	if invalidControlCount > checkSize/20 {
		return true
	}

	// 检查 UTF-8 有效性（检查前1024字节）
	checkData := data
	if len(checkData) > 1024 {
		checkData = data[:1024]
	}

	// 统计有效 UTF-8 序列的字节数
	validUTF8Bytes := 0
	i := 0
	for i < len(checkData) {
		r, size := utf8.DecodeRune(checkData[i:])
		if r == utf8.RuneError {
			// 如果是 RuneError，可能是因为 buffer 截断
			// 只有当不是因为截断导致的错误时，才认为是无效 UTF-8
			// RuneError width is 1
			if size == 1 {
				// 检查是否在末尾
				if i+size >= len(checkData) {
					// 在末尾截断，视为有效（防止因截断导致误判为二进制）
					validUTF8Bytes += size
					break
				}
				// 真正的无效 UTF-8 序列
				// 继续
			}
		} else {
			validUTF8Bytes += size
		}
		i += size
	}

	// 如果有效 UTF-8 字节比例低于 80%，且不是纯 ASCII
	if float64(validUTF8Bytes)/float64(len(checkData)) < 0.8 {
		// 统计可打印字符的比例
		printableCount := 0
		for _, b := range checkData {
			if (b >= 32 && b < 127) || allowedControlChars[b] {
				printableCount++
			}
		}
		// 如果可打印字符少于70%，可能是二进制文件
		if printableCount < len(checkData)*70/100 {
			return true
		}
	}

	return false
}
