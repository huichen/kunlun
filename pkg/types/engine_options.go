package types

type EngineOptions struct {
	// 遍历器选项，如果后续不进行文件遍历，可以设为 nil
	WalkerOptions *IndexWalkerOptions

	// 索引器选项
	IndexerOptions *IndexerOptions

	// 搜索器选项
	SearcherOptions *SearcherOptions
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
