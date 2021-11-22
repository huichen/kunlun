package types

// 昆仑引擎创建选项
type EngineOptions struct {
	// 遍历器选项，如果后续不进行文件遍历，可以设为 nil
	WalkerOptions *IndexWalkerOptions

	// 索引器选项
	IndexerOptions *IndexerOptions

	// 搜索器选项
	SearcherOptions *SearcherOptions

	// 外部传入 repoID，当索引到对应的 remoteURL 时，使用外部 ID 标识仓库
	// 注意，repoID 不能为 0，否则会报错
	RepoRemoteURLToIDMap map[string]uint64
}

func NewEngineOptions() *EngineOptions {
	return &EngineOptions{}
}

func (options *EngineOptions) SetIndexerOptions(opt *IndexerOptions) *EngineOptions {
	options.IndexerOptions = opt
	return options
}

func (options *EngineOptions) SetWalkerOptions(opt *IndexWalkerOptions) *EngineOptions {
	options.WalkerOptions = opt
	return options
}

func (options *EngineOptions) SetSearcherOptions(opt *SearcherOptions) *EngineOptions {
	options.SearcherOptions = opt
	return options
}
