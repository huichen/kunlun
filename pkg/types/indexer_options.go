package types

import "runtime"

// 索引器生成选项，请使用下面的方式设置
type IndexerOptions struct {
	NumIndexerShards int

	MaxDocsPerShard int
}

func NewIndexerOptions() *IndexerOptions {
	return &IndexerOptions{
		NumIndexerShards: runtime.NumCPU(),
	}
}

func (opts *IndexerOptions) SetNumIndexerShards(num int) *IndexerOptions {
	if num > 0 {
		opts.NumIndexerShards = num
	}
	return opts
}

func (opts *IndexerOptions) SetMaxDocsPerShard(num int) *IndexerOptions {
	if num > 0 {
		opts.MaxDocsPerShard = num
	}
	return opts
}
