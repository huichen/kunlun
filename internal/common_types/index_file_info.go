package common_types

// 用于添加索引时传入外部解析的文件信息
type IndexFileInfo struct {
	// 在操作系统中的路径
	Path string

	// 源文件语言
	Language string

	// 如果知道属于哪个仓库，传入仓库路径或者远程地址
	RepoLocalPath string
	RepoRemoteURL string

	// 当为某个仓库内代码时，传入仓库内的路径
	PathInRepo string

	// Ctags symbols
	CTagsEntries []*CTagsEntry
}
