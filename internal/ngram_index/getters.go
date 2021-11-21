package ngram_index

func (index *NgramIndex) GetKeyFrequency(key IndexKey) uint64 {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	if v, ok := index.keyFrequencies[key]; ok {
		return v
	}

	return 0
}

func (index *NgramIndex) GetTotalIndexSize() uint64 {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	return index.totalIndexSize
}

func (index *NgramIndex) GetSortTriggered() uint64 {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	return index.sortTriggered
}
