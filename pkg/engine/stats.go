package engine

import "github.com/huichen/kunlun/pkg/types"

func (engine *KunlunEngine) GetIndexerStats() types.IndexerStats {
	return engine.indexer.GetStats()
}

func (engine *KunlunEngine) GetWalkerStats() types.IndexWalkerStats {
	return engine.walker.GetStats()
}
