package indexer

import "github.com/huichen/kunlun/internal/ngram_index"

// 根据文档 ID 获得文档内容
func (indexer *Indexer) GetContent(docID uint64) []byte {
	indexer.indexerLock.RLock()
	defer indexer.indexerLock.RUnlock()

	contentPointer, ok := indexer.documentIDToContentMap[docID]
	if !ok {
		return nil
	}

	content := make([]byte, len(*contentPointer))
	copy(content, *contentPointer)
	return content
}

// 获得 ngram index key 在文档中出现的频率（文档个数）
func (indexer *Indexer) getKeyFrequency(key ngram_index.IndexKey) (uint64, bool) {
	freq := uint64(0)
	found := false
	for i := 0; i < indexer.numIndexerShards; i++ {
		if v, ok := indexer.contentNgramIndices[i].GetKeyFrequency(key); ok {
			found = true
			freq += v
		}
	}

	return freq, found
}
