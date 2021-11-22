package engine

import "github.com/huichen/kunlun/pkg/types"

// 得到引擎中索引器的统计指标
func (engine *KunlunEngine) GetIndexerStats() types.IndexerStats {
	return engine.indexer.GetStats()
}

// 得到引擎中遍历器的统计指标
func (engine *KunlunEngine) GetWalkerStats() types.IndexWalkerStats {
	return engine.walker.GetStats()
}
