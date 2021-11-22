package engine

// 得到一个文档（docID 标识）的文本内容
func (engine *KunlunEngine) GetContent(docID uint64) []byte {
	return engine.indexer.GetContent(docID)
}
