package indexer

type DocumentMeta struct {
	DocumentID uint64
	Size       uint64

	// 文件路径和在 repo 中的相对路径
	LocalPath  string
	PathInRepo string

	// 文件属于哪个仓库
	Repo *CodeRepository

	// 文件的编程语言
	Language *Language

	// 行信息
	// 第 i 个元素保存了第 i 行的第一个位置
	// 如果该行为空（只有 '\n'）则指向 '\n' 否则指向该行第一个不为 '\n' 的字节
	LineStartLocations []uint32
}

func (indexer *Indexer) GetMeta(documentID uint64) *DocumentMeta {
	if !indexer.finished {
		logger.Fatal("必须先调用 Finish 函数")
	}

	if meta, ok := indexer.documentIDToMetaMap[documentID]; ok {
		return meta
	}

	return nil
}
