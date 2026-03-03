package ui

import "github.com/charmbracelet/lipgloss"

// 颜色常量 - OpenCode 风格艳丽深色主题
const (
	ColorBackground    = lipgloss.Color("#0d1117") // 深色背景（接近黑色）
	ColorForeground    = lipgloss.Color("#e6edf3") // 主前景色（浅灰白）
	ColorBorderFocus   = lipgloss.Color("#c586c0") // 焦点边框色（艳丽紫色）
	ColorBorderNormal  = lipgloss.Color("#30363d") // 普通边框色（深灰）
	ColorHeaderBg      = lipgloss.Color("#161b22") // Header 背景色
	ColorFooterBg      = lipgloss.Color("#161b22") // Footer 背景色
	ColorSelected      = lipgloss.Color("#ff79c6") // 多选高亮色（艳丽粉色）
	ColorCursor        = lipgloss.Color("#6e40c9") // 光标背景色（深紫色）
	ColorCursorFg      = lipgloss.Color("#ffffff") // 光标前景色（纯白）
	ColorDirColor      = lipgloss.Color("#c586c0") // 目录颜色（紫色）
	ColorSymlinkColor  = lipgloss.Color("#d2a8ff") // 符号链接颜色（浅紫色）
	ColorExecColor     = lipgloss.Color("#7ee787") // 可执行文件颜色（亮绿色）
	ColorArchiveColor  = lipgloss.Color("#f0883e") // 归档文件颜色（橙色）
	ColorImageColor    = lipgloss.Color("#ff79c6") // 图片颜色（粉色）
	ColorSubdued       = lipgloss.Color("#8b949e") // 暗淡文字色（中灰）
	ColorTitle         = lipgloss.Color("#c586c0") // 标题色（紫色）
	ColorSearchActive  = lipgloss.Color("#ff79c6") // 搜索激活色（粉色）
	ColorSizeColor     = lipgloss.Color("#8b949e") // 文件大小颜色
	ColorDateColor     = lipgloss.Color("#8b949e") // 日期颜色
	ColorPreviewTitle  = lipgloss.Color("#c586c0") // 预览标题颜色（紫色）
	ColorInfoLabel     = lipgloss.Color("#8b949e") // 信息标签颜色
	ColorInfoValue     = lipgloss.Color("#e6edf3") // 信息值颜色（浅色）
	ColorSelectionInfo = lipgloss.Color("#ff79c6") // 选择信息颜色（粉色）
	ColorKeyHint       = lipgloss.Color("#8b949e") // 快捷键提示颜色
	ColorKeyHighlight  = lipgloss.Color("#c586c0") // 快捷键高亮颜色（紫色）
	ColorError         = lipgloss.Color("#f85149") // 错误颜色（红色）
	ColorSuccess       = lipgloss.Color("#7ee787") // 成功颜色（绿色）
)

// Theme 主题样式集合
type Theme struct {
	// 边框样式
	BorderFocus  lipgloss.Style
	BorderNormal lipgloss.Style

	// 文字样式
	TitleStyle      lipgloss.Style
	SubduedStyle    lipgloss.Style
	DirStyle        lipgloss.Style
	SymlinkStyle    lipgloss.Style
	ExecStyle       lipgloss.Style
	ArchiveStyle    lipgloss.Style
	ImageStyle      lipgloss.Style
	SelectedStyle   lipgloss.Style
	CursorStyle     lipgloss.Style
	SizeStyle       lipgloss.Style
	DateStyle       lipgloss.Style
	InfoLabelStyle  lipgloss.Style
	InfoValueStyle  lipgloss.Style
	SearchStyle     lipgloss.Style
	KeyHintStyle    lipgloss.Style
	KeyHighlight    lipgloss.Style
	ErrorStyle      lipgloss.Style
	SuccessStyle    lipgloss.Style
	SelectionStyle  lipgloss.Style
	PreviewTitle    lipgloss.Style
}

// DefaultTheme OpenCode 风格艳丽深色主题
var DefaultTheme = &Theme{
	BorderFocus:  lipgloss.NewStyle().Foreground(ColorBorderFocus).Bold(true),
	BorderNormal: lipgloss.NewStyle().Foreground(ColorBorderNormal),

	TitleStyle:     lipgloss.NewStyle().Foreground(ColorTitle).Bold(true),
	SubduedStyle:   lipgloss.NewStyle().Foreground(ColorSubdued),
	DirStyle:       lipgloss.NewStyle().Foreground(ColorDirColor).Bold(true),
	SymlinkStyle:   lipgloss.NewStyle().Foreground(ColorSymlinkColor),
	ExecStyle:      lipgloss.NewStyle().Foreground(ColorExecColor).Bold(true),
	ArchiveStyle:   lipgloss.NewStyle().Foreground(ColorArchiveColor),
	ImageStyle:     lipgloss.NewStyle().Foreground(ColorImageColor).Bold(true),
	SelectedStyle:  lipgloss.NewStyle().Foreground(ColorSelected).Bold(true),
	CursorStyle:    lipgloss.NewStyle().Background(ColorCursor).Foreground(ColorCursorFg).Bold(true),
	SizeStyle:      lipgloss.NewStyle().Foreground(ColorSizeColor),
	DateStyle:      lipgloss.NewStyle().Foreground(ColorDateColor),
	InfoLabelStyle: lipgloss.NewStyle().Foreground(ColorInfoLabel),
	InfoValueStyle: lipgloss.NewStyle().Foreground(ColorInfoValue),
	SearchStyle:    lipgloss.NewStyle().Foreground(ColorSearchActive).Bold(true),
	KeyHintStyle:   lipgloss.NewStyle().Foreground(ColorKeyHint),
	KeyHighlight:   lipgloss.NewStyle().Foreground(ColorKeyHighlight).Bold(true),
	ErrorStyle:     lipgloss.NewStyle().Foreground(ColorError).Bold(true),
	SuccessStyle:   lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true),
	SelectionStyle: lipgloss.NewStyle().Foreground(ColorSelectionInfo).Bold(true),
	PreviewTitle:   lipgloss.NewStyle().Foreground(ColorPreviewTitle).Bold(true),
}
