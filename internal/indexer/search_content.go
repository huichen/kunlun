package indexer

import (
	"sort"

	"github.com/huichen/kunlun/internal/ngram_index"

	"github.com/huichen/kunlun/pkg/types"
)

func (indexer *Indexer) searchContent(
	keyword []byte,
	offset uint32,
	key1 ngram_index.IndexKey,
	key2 ngram_index.IndexKey,
	distance uint32,
	caseSensitive bool,
	isSymbol bool,
	shouldDocBeRecalled func(uint64) bool,
) ([]types.DocumentWithSections, error) {
	responseChan := make(chan ContentSearchResponseInfo, indexer.numIndexerShards)

	// 启动搜索协程
	requestChan := make(chan ContentSearchRequestInfo, indexer.numIndexerShards)
	for i := 0; i < indexer.numIndexerShards; i++ {
		go indexer.contentSearchWorker(i, requestChan, responseChan)
	}

	// 发送搜索请求
	for i := 0; i < indexer.numIndexerShards; i++ {
		requestChan <- ContentSearchRequestInfo{
			Keyword:             keyword,
			Offset:              offset,
			Key1:                key1,
			Key2:                key2,
			Distance:            distance,
			CaseSensitive:       caseSensitive,
			IsSymbol:            isSymbol,
			ShouldDocBeRecalled: shouldDocBeRecalled,
		}
	}

	// 接受和合并搜索请求
	retDocs := []ngram_index.DocumentWithLocations{}
	var err error
	for i := 0; i < indexer.numIndexerShards; i++ {
		response := <-responseChan
		if response.Err != nil {
			err = response.Err
		} else {
			retDocs = append(retDocs, response.DocumentsWithLocations...)
		}
	}
	if err != nil {
		return nil, err
	}

	// 搜索结果为空，直接返回
	if len(retDocs) == 0 {
		return []types.DocumentWithSections{}, nil
	}

	sort.Sort(ngram_index.SortDocumentWithLocations(retDocs))

	return DocLocationsToSections(retDocs, uint32(len(keyword))), nil

}

type ContentSearchRequestInfo struct {
	Keyword             []byte
	Offset              uint32
	Key1                ngram_index.IndexKey
	Key2                ngram_index.IndexKey
	Distance            uint32
	CaseSensitive       bool
	IsSymbol            bool
	ShouldDocBeRecalled func(uint64) bool
}

type ContentSearchResponseInfo struct {
	Err                    error
	DocumentsWithLocations []ngram_index.DocumentWithLocations
}

func (indexer *Indexer) contentSearchWorker(
	shard int,
	requestChan chan ContentSearchRequestInfo,
	responseChan chan ContentSearchResponseInfo,
) {
	info := <-requestChan
	key1 := info.Key1
	key2 := info.Key2
	distance := info.Distance
	docFilterFunc := info.ShouldDocBeRecalled

	var contentMatchDocs []ngram_index.DocumentWithLocations
	var err error
	var docs []ngram_index.DocumentWithLocations
	if distance > 0 {
		docs, err = indexer.contentNgramIndices[shard].SearchTwoKeys(
			key1, key2, distance, docFilterFunc, info.IsSymbol)
	} else {
		// 单键的情况
		docs, err = indexer.contentNgramIndices[shard].SearchOneKey(key1, docFilterFunc, info.IsSymbol)
	}
	if err != nil {
		responseChan <- ContentSearchResponseInfo{
			Err: err,
		}
		return
	}
	if info.Keyword == nil {
		responseChan <- ContentSearchResponseInfo{
			Err:                    nil,
			DocumentsWithLocations: docs,
		}
		return
	}

	contentMatchDocs, err = indexer.filterDocumentsWithFullMatch(
		&(indexer.documentIDToContentMap), docs, info.Offset, info.Keyword, info.CaseSensitive)
	if err != nil {
		responseChan <- ContentSearchResponseInfo{
			Err: err,
		}
		return
	}

	responseChan <- ContentSearchResponseInfo{
		Err:                    nil,
		DocumentsWithLocations: contentMatchDocs,
	}
}
