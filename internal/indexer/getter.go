package indexer

import "kunlun/internal/ngram_index"

func (indexer *Indexer) GetKeyFrequency(key ngram_index.IndexKey) uint64 {
	freq := uint64(0)
	for i := 0; i < indexer.numIndexerShards; i++ {
		freq += indexer.contentNgramIndices[i].GetKeyFrequency(key)
	}

	return freq
}

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
