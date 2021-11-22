package indexer

import "github.com/huichen/kunlun/pkg/types"

// 返回索引统计信息
func (indexer *Indexer) GetStats() types.IndexerStats {
	indexer.indexerLock.RLock()
	defer indexer.indexerLock.RUnlock()

	stats := types.IndexerStats{
		IndexerShards:      indexer.numIndexerShards,
		TotalContentSize:   indexer.totalContentSize,
		TotalDocumentCount: indexer.totalDocumentCount,
		FailedAddingSymbol: indexer.failedDocs,
	}

	for i := 0; i < indexer.numIndexerShards; i++ {
		stats.IndexSortTriggered += indexer.contentNgramIndices[i].GetSortTriggered()
	}

	for i := 0; i < indexer.numIndexerShards; i++ {
		stats.TotalIndexSize += indexer.contentNgramIndices[i].GetTotalIndexSize()
	}

	return stats
}
