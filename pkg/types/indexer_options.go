package types

import "runtime"

// 索引器创建选项
type IndexerOptions struct {
	NumIndexerShards int

	MaxDocsPerShard int

	SkipIndexBigram  bool
	SkipIndexUnigram bool
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

func (opts *IndexerOptions) SetSkipIndexBigram(skip bool) *IndexerOptions {
	opts.SkipIndexBigram = skip
	return opts
}

func (opts *IndexerOptions) SetSkipIndexUnigram(skip bool) *IndexerOptions {
	opts.SkipIndexUnigram = skip
	return opts
}
