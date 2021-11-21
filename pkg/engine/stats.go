package engine

import "kunlun/pkg/types"

func (engine *KunlunEngine) GetIndexerStats() types.IndexerStats {
	return engine.indexer.GetStats()
}

func (engine *KunlunEngine) GetWalkerStats() types.IndexWalkerStats {
	return engine.walker.GetStats()
}
