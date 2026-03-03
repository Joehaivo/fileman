package fileops

// 文件扩展名到 Nerd Font 图标的映射
// 使用 Nerd Font v3 字符
var extIcons = map[string]string{
	// Go
	".go": "󰟓",
	// JavaScript / TypeScript
	".js":  "",
	".jsx": "",
	".ts":  "",
	".tsx": "",
	// Web
	".html": "",
	".htm":  "",
	".css":  "",
	".scss": "",
	".sass": "",
	".less": "",
	// Python
	".py":  "",
	".pyc": "",
	// Rust
	".rs": "",
	// C/C++
	".c":   "",
	".cpp": "",
	".cc":  "",
	".h":   "",
	".hpp": "",
	// Java/Kotlin
	".java":   "",
	".kt":     "",
	".kts":    "",
	".class":  "",
	// Swift
	".swift": "",
	// Ruby
	".rb": "",
	// PHP
	".php": "",
	// Shell
	".sh":   "",
	".bash": "",
	".zsh":  "",
	".fish": "",
	// Config
	".json":   "",
	".yaml":   "",
	".yml":    "",
	".toml":   "",
	".ini":    "",
	".cfg":    "",
	".conf":   "",
	".env":    "",
	".editorconfig": "",
	// Markdown / Docs
	".md":   "",
	".mdx":  "",
	".txt":  "",
	".rst":  "",
	".pdf":  "",
	".doc":  "",
	".docx": "",
	".xls":  "󰈛",
	".xlsx": "󰈛",
	".ppt":  "󰈧",
	".pptx": "󰈧",
	// Images
	".png":  "",
	".jpg":  "",
	".jpeg": "",
	".gif":  "",
	".webp": "",
	".svg":  "",
	".ico":  "",
	".bmp":  "",
	".tiff": "",
	// Video
	".mp4":  "",
	".mov":  "",
	".avi":  "",
	".mkv":  "",
	".webm": "",
	// Audio
	".mp3":  "",
	".wav":  "",
	".flac": "",
	".ogg":  "",
	".m4a":  "",
	// Archives
	".zip": "",
	".tar": "",
	".gz":  "",
	".bz2": "",
	".xz":  "",
	".7z":  "",
	".rar": "",
	".dmg": "",
	// Fonts
	".ttf":   "",
	".otf":   "",
	".woff":  "",
	".woff2": "",
	// Database
	".sql":    "",
	".sqlite": "",
	".db":     "",
	// Git
	".gitignore":    "",
	".gitattributes": "",
	// Docker
	"dockerfile": "",
	// Makefile
	"makefile": "",
	// Lock files
	".lock": "󰌾",
	// Lua
	".lua": "",
	// Vim
	".vim":   "",
	".vimrc": "",
	// Dart
	".dart": "",
}

// specialNameIcons 特殊文件名图标映射（不依赖扩展名）
var specialNameIcons = map[string]string{
	"dockerfile":       "",
	"docker-compose.yml": "",
	"docker-compose.yaml": "",
	"makefile":         "",
	"rakefile":         "",
	".gitignore":       "",
	".gitconfig":       "",
	".gitattributes":   "",
	"go.mod":           "󰟓",
	"go.sum":           "󰟓",
	"package.json":     "",
	"package-lock.json": "",
	"yarn.lock":        "",
	"pnpm-lock.yaml":   "",
	"cargo.toml":       "",
	"cargo.lock":       "",
	"readme.md":        "",
	"license":          "",
	"licence":          "",
	".env":             "",
	".env.local":       "",
	".env.example":     "",
}

// GetFileIcon 根据文件名和类型返回对应的 Nerd Font 图标
func GetFileIcon(name string, isDir bool) string {
	if isDir {
		return ""
	}

	// 先检查特殊文件名（不区分大小写）
	lowerName := toLower(name)
	if icon, ok := specialNameIcons[lowerName]; ok {
		return icon
	}

	// 再检查扩展名
	ext := getExt(name)
	if ext != "" {
		if icon, ok := extIcons[ext]; ok {
			return icon
		}
	}

	// 默认文件图标
	return ""
}

// toLower 转换为小写（避免 strings 包依赖）
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result[i] = c
	}
	return string(result)
}

// getExt 获取文件扩展名（含点，小写）
func getExt(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			if i == 0 {
				return ""
			}
			ext := name[i:]
			return toLower(ext)
		}
	}
	return ""
}
