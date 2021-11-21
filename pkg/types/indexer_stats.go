package types

type IndexerStats struct {
	IndexerShards int

	TotalContentSize   uint64
	TotalDocumentCount uint64

	TotalIndexSize     uint64
	IndexSortTriggered uint64

	FailedAddingSymbol uint64
}
