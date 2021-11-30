package ngram_index

// 第二个参数返回 key 是否存在于索引中
func (index *NgramIndex) GetKeyFrequency(key IndexKey) (uint64, bool) {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	if v, ok := index.keyFrequencies[key]; ok {
		return v, ok
	}

	return 0, false
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
