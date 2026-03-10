package types

import "time"

// FileType 文件类型
type FileType int

const (
	FileTypeRegular   FileType = iota // 普通文件
	FileTypeDirectory                 // 目录
	FileTypeSymlink                   // 符号链接
	FileTypeOther                     // 其他类型
)

// FileEntry 表示目录中的一个文件或目录条目
type FileEntry struct {
	Name    string      // 文件名
	Size    int64       // 文件大小（字节）
	ModTime time.Time   // 修改时间
	Mode    string      // 权限字符串，如 "-rw-r--r--"
	Type    FileType    // 文件类型
	IsDir   bool        // 是否为目录
	Ext     string      // 扩展名（小写，含点，如 ".go"）
	Path    string      // 完整路径
}

// SelectionSet 多选文件集合，使用 map 实现 O(1) 查找
type SelectionSet map[string]struct{}

// Add 添加文件路径到选择集
func (s SelectionSet) Add(path string) {
	s[path] = struct{}{}
}

// Remove 从选择集移除文件路径
func (s SelectionSet) Remove(path string) {
	delete(s, path)
}

// Toggle 切换文件路径的选择状态
func (s SelectionSet) Toggle(path string) {
	if _, ok := s[path]; ok {
		delete(s, path)
	} else {
		s[path] = struct{}{}
	}
}

// Has 检查文件路径是否在选择集中
func (s SelectionSet) Has(path string) bool {
	_, ok := s[path]
	return ok
}

// Clear 清空选择集
func (s SelectionSet) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Len 返回选择集大小
func (s SelectionSet) Len() int {
	return len(s)
}

// FocusTarget 焦点目标
type FocusTarget int

const (
	FocusPanelA   FocusTarget = iota // 上方文件面板
	FocusPanelB                      // 下方文件面板
	FocusPreview                     // 预览/编辑区
)

// ModalType 模态弹窗类型
type ModalType int

const (
	ModalNone     ModalType = iota // 无弹窗
	ModalNewDir                    // 新建目录
	ModalNewFile                   // 新建文件
	ModalRename                    // 重命名
	ModalDelete                    // 删除确认
	ModalProgress                  // 复制/移动进度
	ModalError                     // 错误提示
	ModalSettings                  // 设置弹窗
)

// Settings 应用设置
type Settings struct {
	ShowDate   bool // 是否显示修改时间
	ShowHidden bool // 是否显示隐藏文件
}

// ProgressInfo 复制/移动进度信息
type ProgressInfo struct {
	Total     int64   // 总字节数
	Done      int64   // 已完成字节数
	Percent   float64 // 进度百分比 0-1
	FileName  string  // 当前处理文件名
	IsFinish  bool    // 是否完成
}

// FileOpResult 单个文件操作结果
type FileOpResult struct {
	SrcPath string
	DstPath string
	Err     error
}

// FloatingProgress 悬浮进度窗口状态
type FloatingProgress struct {
	OpType     string         // "复制" 或 "移动"
	Total      int            // 总文件数
	Done       int            // 已完成数
	Results    []FileOpResult // 操作结果列表
	IsComplete bool           // 是否完成
}
