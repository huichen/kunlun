package types

import (
	"runtime"
	"strings"
)

var (
	defaultMaxFileSize  = 1 << 17 // 128 KB
	defaultMaxFileLines = 5000
)

// 遍历器创建参数
type IndexWalkerOptions struct {
	// 只对该 map 中 extension （比如 "go", "java"，包含 .）的文件做检索
	// 如果为空或者 nil 则不做过滤
	AllowedFileExtensions    map[string]bool
	DisallowedFileExtensions map[string]bool // 优先级高于白名单

	// 只对该 map 中编程语言的文件做检索
	// 如果为空或者 nil 则不做过滤
	AllowedCodeLanguages map[string]bool
	// 是否允许未检测到语言的文件进入索引
	AllowUnknownLanguage bool
	// 不索引这些语言的文件
	DisallowedCodeLanguages map[string]bool // 优先级高于白名单

	// 过滤超过这么大尺寸的文件
	// 如果不设置则用 defaultMaxFileSize
	MaxFileSize int

	// 过滤超过这么多行的文件
	// 如果不设置则用 defaultMaxFileSize
	MaxFileLines int

	// 过滤这些文件夹，文件夹名取 base path，比如 /xxx/yyy/zzz，则取 zzz
	IgnoreDirs map[string]bool

	// 过滤 . 开头的文件和目录
	FilterDotPrefix bool

	// 只对这个白名单中的仓库做索引，如果不设置则不过滤
	AllowedRepoRemoteURLs map[string]bool

	// 使用多少个线程做文件处理
	NumFileProcessors int

	// Ctags 选项
	CTagsParserOptions *CTagsParserOptions
}

func NewIndexWalkerOptions() *IndexWalkerOptions {
	return &IndexWalkerOptions{
		AllowedFileExtensions:    make(map[string]bool),
		DisallowedFileExtensions: make(map[string]bool),
		AllowedCodeLanguages:     make(map[string]bool),
		DisallowedCodeLanguages:  make(map[string]bool),
		MaxFileSize:              defaultMaxFileSize,
		MaxFileLines:             defaultMaxFileLines,
		IgnoreDirs:               make(map[string]bool),
		FilterDotPrefix:          true,
		NumFileProcessors:        runtime.NumCPU() * 2,
	}
}

// 半角逗号分隔的文件后缀列表
// 不设置意味着不过滤
func (options *IndexWalkerOptions) SetAllowedFileExtensions(extensions string) *IndexWalkerOptions {
	fields := strings.Split(extensions, ",")

	for _, f := range fields {
		if f != "" {
			options.AllowedFileExtensions[f] = true
		}
	}

	return options
}
func (options *IndexWalkerOptions) SetDisallowedFileExtensions(extensions string) *IndexWalkerOptions {
	fields := strings.Split(extensions, ",")

	for _, f := range fields {
		if f != "" {
			options.DisallowedFileExtensions[f] = true
		}
	}

	return options
}

// 半角逗号分隔的编程语言列表
// 不设置意味着不过滤
func (options *IndexWalkerOptions) SetAllowedCodeLanguages(languages string) *IndexWalkerOptions {
	fields := strings.Split(languages, ",")

	for _, f := range fields {
		if f != "" {
			options.AllowedCodeLanguages[f] = true
		}
	}

	return options
}
func (options *IndexWalkerOptions) SetDisallowedCodeLanguages(languages string) *IndexWalkerOptions {
	fields := strings.Split(languages, ",")

	for _, f := range fields {
		if f != "" {
			options.DisallowedCodeLanguages[f] = true
		}
	}

	return options
}

// 设置为 true 则允许未检测到语言的文件被索引到
func (options *IndexWalkerOptions) SetAllowUnknownLanuage(allow bool) *IndexWalkerOptions {
	options.AllowUnknownLanguage = allow

	return options
}

// 过滤超过这么大尺寸的文件
func (options *IndexWalkerOptions) SetMaxFileSize(size int) *IndexWalkerOptions {
	if size > 0 {
		options.MaxFileSize = size
	}
	return options
}

// 过滤超过这么多行的文件
func (options *IndexWalkerOptions) SetMaxFileLines(lines int) *IndexWalkerOptions {
	if lines > 0 {
		options.MaxFileLines = lines
	}
	return options
}

// 半角逗号分隔的文件夹名列表
// 不设置意味着不过滤
func (options *IndexWalkerOptions) SetIgnoreDirs(dirs string) *IndexWalkerOptions {
	fields := strings.Split(dirs, ",")

	for _, f := range fields {
		if f != "" {
			options.IgnoreDirs[f] = true
		}
	}

	return options
}

func (options *IndexWalkerOptions) SetFilterDotPrefix(filter bool) *IndexWalkerOptions {
	options.FilterDotPrefix = filter
	return options
}

func (options *IndexWalkerOptions) SetNumProcessors(num int) *IndexWalkerOptions {
	if num > 0 {
		options.NumFileProcessors = num
	}
	return options
}

func (options *IndexWalkerOptions) SetCTagsParserOptions(opt *CTagsParserOptions) *IndexWalkerOptions {
	if opt.BinaryPath != "" {
		options.CTagsParserOptions = opt
	}
	return options
}
