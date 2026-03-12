package i18n

type Messages struct {
	AppName string
	Version string

	HeaderSearchLabel   string
	HeaderSelectedCount string

	FooterMoveCursor    string
	FooterSave          string
	FooterExit          string
	FooterCopyLine      string
	FooterCopyAll       string
	FooterHomeEnd       string
	FooterPageUpDown    string
	FooterConfirm       string
	FooterCancelSearch  string
	FooterSelect        string
	FooterDeleteChar    string
	FooterGoUp          string
	FooterGoDown        string
	FooterSwitchPanel   string
	FooterSearch        string
	FooterEdit          string
	FooterRename        string
	FooterCopy          string
	FooterMove          string
	FooterNewDir        string
	FooterNewFile       string
	FooterExternalEdit  string
	FooterToggleHidden  string
	FooterSettings      string
	FooterSettingsLabel string
	FooterQuit          string
	FooterDelete        string

	ModalTitleNewDir           string
	ModalTitleNewFile          string
	ModalTitleRename           string
	ModalTitleDelete           string
	ModalTitleError            string
	ModalTitleSettings         string
	ModalMsgDirName            string
	ModalMsgFileName           string
	ModalMsgNewName            string
	ModalMsgDirPlaceholder     string
	ModalMsgFilePlaceholder    string
	ModalMsgNewNamePlaceholder string
	ModalMsgDeleteSingle       string
	ModalMsgDeleteMulti        string
	ModalConfirmDelete         string
	ModalClose                 string
	ModalOperating             string
	ModalConfirm               string
	ModalCancel                string
	ModalToggle                string
	ModalSelect                string

	SettingShowDate   string
	SettingShowHidden string
	SettingLanguage   string

	PreviewBinary      string
	PreviewTooLarge    string
	PreviewTooLargeFmt string // "File too large (%s), cannot preview"
	PreviewSelectFile  string
	InfoType           string
	InfoSize           string
	InfoModified       string
	InfoMode           string
	InfoLines          string
	InfoArchiveFiles   string
	InfoArchiveSize    string
	ArchiveTitle       string
	ArchiveEmpty       string

	ToastCopied      string
	ToastCopySuccess string
	ToastMoveSuccess string
	ToastCurrentLine string
	ToastAllContent  string
}

var Chinese = &Messages{
	AppName: "文件管家",
	Version: "v0.1.0",

	HeaderSearchLabel:   "搜索: ",
	HeaderSelectedCount: "已选: %d 个",

	FooterMoveCursor:    "移动光标",
	FooterSave:          "保存",
	FooterExit:          "退出",
	FooterCopyLine:      "复制行",
	FooterCopyAll:       "复制全部",
	FooterHomeEnd:       "首/尾",
	FooterPageUpDown:    "翻页",
	FooterConfirm:       "确认",
	FooterCancelSearch:  "取消搜索",
	FooterSelect:        "选择",
	FooterDeleteChar:    "删除字符",
	FooterGoUp:          "上一级",
	FooterGoDown:        "下一级",
	FooterSwitchPanel:   "切换面板",
	FooterSearch:        "搜索",
	FooterEdit:          "编辑",
	FooterRename:        "重命名",
	FooterCopy:          "复制",
	FooterMove:          "移动",
	FooterNewDir:        "新建目录",
	FooterNewFile:       "新建文件",
	FooterExternalEdit:  "外部编辑",
	FooterToggleHidden:  "切换隐藏",
	FooterSettings:      "设置",
	FooterSettingsLabel: "设置",
	FooterQuit:          "退出",
	FooterDelete:        "删除",

	ModalTitleNewDir:           "新建目录",
	ModalTitleNewFile:          "新建文件",
	ModalTitleRename:           "重命名",
	ModalTitleDelete:           "确认删除",
	ModalTitleError:            "错误",
	ModalTitleSettings:         "设置",
	ModalMsgDirName:            "请输入目录名称：",
	ModalMsgFileName:           "请输入文件名称：",
	ModalMsgNewName:            "请输入新名称：",
	ModalMsgDirPlaceholder:     "目录名称",
	ModalMsgFilePlaceholder:    "文件名称",
	ModalMsgNewNamePlaceholder: "新名称",
	ModalMsgDeleteSingle:       "确定要删除 \"%s\" 吗？",
	ModalMsgDeleteMulti:        "确定要删除选中的 %d 个文件吗？",
	ModalConfirmDelete:         "确认删除",
	ModalClose:                 "关闭",
	ModalOperating:             "操作进行中...",
	ModalConfirm:               "确认",
	ModalCancel:                "取消",
	ModalToggle:                "切换",
	ModalSelect:                "选择",

	SettingShowDate:   "展示修改时间",
	SettingShowHidden: "显示隐藏文件",
	SettingLanguage:   "切换中文/English",

	PreviewBinary:      "二进制文件，无法预览",
	PreviewTooLarge:    "文件过大，无法预览",
	PreviewTooLargeFmt: "文件过大 (%s)，无法预览",
	PreviewSelectFile:  "选择文件以预览",
	InfoType:           "类型: ",
	InfoSize:           "大小: ",
	InfoModified:       "修改: ",
	InfoMode:           "权限: ",
	InfoLines:          "行数: ",
	InfoArchiveFiles:   "文件数: ",
	InfoArchiveSize:    "解压大小: ",
	ArchiveTitle:       "压缩包内容",
	ArchiveEmpty:       "(空压缩包)",

	ToastCopied:      "已复制: %s",
	ToastCopySuccess: "复制成功",
	ToastMoveSuccess: "移动成功",
	ToastCurrentLine: "当前行",
	ToastAllContent:  "全部内容",
}

var English = &Messages{
	AppName: "FileMan",
	Version: "v0.1.0",

	HeaderSearchLabel:   "Search: ",
	HeaderSelectedCount: "Selected: %d",

	FooterMoveCursor:    "Move",
	FooterSave:          "Save",
	FooterExit:          "Exit",
	FooterCopyLine:      "Copy Line",
	FooterCopyAll:       "Copy All",
	FooterHomeEnd:       "Home/End",
	FooterPageUpDown:    "PgUp/PgDn",
	FooterConfirm:       "Confirm",
	FooterCancelSearch:  "Cancel",
	FooterSelect:        "Select",
	FooterDeleteChar:    "Delete",
	FooterGoUp:          "Up",
	FooterGoDown:        "Down",
	FooterSwitchPanel:   "Switch",
	FooterSearch:        "Search",
	FooterEdit:          "Edit",
	FooterRename:        "Rename",
	FooterCopy:          "Copy",
	FooterMove:          "Move",
	FooterNewDir:        "New Dir",
	FooterNewFile:       "New File",
	FooterExternalEdit:  "Ext Edit",
	FooterToggleHidden:  "Hidden",
	FooterSettings:      "Settings",
	FooterSettingsLabel: "Settings",
	FooterQuit:          "Quit",
	FooterDelete:        "Delete",

	ModalTitleNewDir:           "New Directory",
	ModalTitleNewFile:          "New File",
	ModalTitleRename:           "Rename",
	ModalTitleDelete:           "Confirm Delete",
	ModalTitleError:            "Error",
	ModalTitleSettings:         "Settings",
	ModalMsgDirName:            "Enter directory name:",
	ModalMsgFileName:           "Enter file name:",
	ModalMsgNewName:            "Enter new name:",
	ModalMsgDirPlaceholder:     "Directory name",
	ModalMsgFilePlaceholder:    "File name",
	ModalMsgNewNamePlaceholder: "New name",
	ModalMsgDeleteSingle:       "Delete \"%s\"?",
	ModalMsgDeleteMulti:        "Delete %d selected files?",
	ModalConfirmDelete:         "Delete",
	ModalClose:                 "Close",
	ModalOperating:             "Operating...",
	ModalConfirm:               "Confirm",
	ModalCancel:                "Cancel",
	ModalToggle:                "Toggle",
	ModalSelect:                "Select",

	SettingShowDate:   "Show modification time",
	SettingShowHidden: "Show hidden files",
	SettingLanguage:   "中文/English",

	PreviewBinary:      "Binary file, cannot preview",
	PreviewTooLarge:    "File too large, cannot preview",
	PreviewTooLargeFmt: "File too large (%s), cannot preview",
	PreviewSelectFile:  "Select a file to preview",
	InfoType:           "Type: ",
	InfoSize:           "Size: ",
	InfoModified:       "Modified: ",
	InfoMode:           "Mode: ",
	InfoLines:          "Lines: ",
	InfoArchiveFiles:   "Files: ",
	InfoArchiveSize:    "Unpacked: ",
	ArchiveTitle:       "Archive Contents",
	ArchiveEmpty:       "(Empty archive)",

	ToastCopied:      "Copied: %s",
	ToastCopySuccess: "Copied",
	ToastMoveSuccess: "Moved",
	ToastCurrentLine: "Current line",
	ToastAllContent:  "All content",
}

func GetMessages(useEnglish bool) *Messages {
	if useEnglish {
		return English
	}
	return Chinese
}
