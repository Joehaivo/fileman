package ui

import "github.com/charmbracelet/lipgloss"

// 颜色常量
const (
	ColorBackground    = lipgloss.Color("#1a1b26") // 深色背景
	ColorForeground    = lipgloss.Color("#c0caf5") // 主前景色
	ColorBorderFocus   = lipgloss.Color("#7aa2f7") // 焦点边框色（蓝色）
	ColorBorderNormal  = lipgloss.Color("#3b4261") // 普通边框色（暗灰）
	ColorHeaderBg      = lipgloss.Color("#16161e") // Header 背景色
	ColorFooterBg      = lipgloss.Color("#16161e") // Footer 背景色
	ColorSelected      = lipgloss.Color("#ff9e64") // 多选高亮色（橙色）
	ColorCursor        = lipgloss.Color("#2d3f76") // 光标背景色
	ColorCursorFg      = lipgloss.Color("#c0caf5") // 光标前景色
	ColorDirColor      = lipgloss.Color("#7aa2f7") // 目录颜色（蓝色）
	ColorSymlinkColor  = lipgloss.Color("#bb9af7") // 符号链接颜色（紫色）
	ColorExecColor     = lipgloss.Color("#9ece6a") // 可执行文件颜色（绿色）
	ColorArchiveColor  = lipgloss.Color("#e0af68") // 归档文件颜色（黄色）
	ColorImageColor    = lipgloss.Color("#ff9e64") // 图片颜色（橙色）
	ColorSubdued       = lipgloss.Color("#565f89") // 暗淡文字色
	ColorTitle         = lipgloss.Color("#7dcfff") // 标题色（浅蓝）
	ColorSearchActive  = lipgloss.Color("#f7768e") // 搜索激活色（红色）
	ColorSizeColor     = lipgloss.Color("#565f89") // 文件大小颜色
	ColorDateColor     = lipgloss.Color("#565f89") // 日期颜色
	ColorPreviewTitle  = lipgloss.Color("#7dcfff") // 预览标题颜色
	ColorInfoLabel     = lipgloss.Color("#565f89") // 信息标签颜色
	ColorInfoValue     = lipgloss.Color("#a9b1d6") // 信息值颜色
	ColorSelectionInfo = lipgloss.Color("#ff9e64") // 选择信息颜色（橙色）
	ColorKeyHint       = lipgloss.Color("#565f89") // 快捷键提示颜色
	ColorKeyHighlight  = lipgloss.Color("#7aa2f7") // 快捷键高亮颜色
	ColorError         = lipgloss.Color("#f7768e") // 错误颜色
	ColorSuccess       = lipgloss.Color("#9ece6a") // 成功颜色
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

// DefaultTheme 默认深色主题
var DefaultTheme = &Theme{
	BorderFocus:  lipgloss.NewStyle().Foreground(ColorBorderFocus),
	BorderNormal: lipgloss.NewStyle().Foreground(ColorBorderNormal),

	TitleStyle:     lipgloss.NewStyle().Foreground(ColorTitle).Bold(true),
	SubduedStyle:   lipgloss.NewStyle().Foreground(ColorSubdued),
	DirStyle:       lipgloss.NewStyle().Foreground(ColorDirColor).Bold(true),
	SymlinkStyle:   lipgloss.NewStyle().Foreground(ColorSymlinkColor),
	ExecStyle:      lipgloss.NewStyle().Foreground(ColorExecColor),
	ArchiveStyle:   lipgloss.NewStyle().Foreground(ColorArchiveColor),
	ImageStyle:     lipgloss.NewStyle().Foreground(ColorImageColor),
	SelectedStyle:  lipgloss.NewStyle().Foreground(ColorSelected).Bold(true),
	CursorStyle:    lipgloss.NewStyle().Background(ColorCursor).Foreground(ColorCursorFg),
	SizeStyle:      lipgloss.NewStyle().Foreground(ColorSizeColor),
	DateStyle:      lipgloss.NewStyle().Foreground(ColorDateColor),
	InfoLabelStyle: lipgloss.NewStyle().Foreground(ColorInfoLabel),
	InfoValueStyle: lipgloss.NewStyle().Foreground(ColorInfoValue),
	SearchStyle:    lipgloss.NewStyle().Foreground(ColorSearchActive).Bold(true),
	KeyHintStyle:   lipgloss.NewStyle().Foreground(ColorKeyHint),
	KeyHighlight:   lipgloss.NewStyle().Foreground(ColorKeyHighlight).Bold(true),
	ErrorStyle:     lipgloss.NewStyle().Foreground(ColorError).Bold(true),
	SuccessStyle:   lipgloss.NewStyle().Foreground(ColorSuccess),
	SelectionStyle: lipgloss.NewStyle().Foreground(ColorSelectionInfo).Bold(true),
	PreviewTitle:   lipgloss.NewStyle().Foreground(ColorPreviewTitle).Bold(true),
}
