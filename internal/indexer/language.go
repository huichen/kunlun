package indexer

type Language struct {
	ID uint64

	Name string
}

// 将代码仓库添加到索引
func (indexer *Indexer) IndexLanguage(lang string) *Language {
	indexer.indexerLock.Lock()
	defer indexer.indexerLock.Unlock()

	if l, ok := indexer.langNameToIDMap[lang]; ok {
		return l
	}

	// 更新 repo 计数
	indexer.langCounter++

	l := &Language{
		ID: indexer.langCounter,
	}

	indexer.langNameToIDMap[lang] = l

	return l
}
