package engine

func (engine *KunlunEngine) GetContent(docID uint64) []byte {
	return engine.indexer.GetContent(docID)
}
