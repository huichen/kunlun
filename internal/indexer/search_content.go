package indexer

import (
	"sort"

	"github.com/huichen/kunlun/internal/ngram_index"
	"github.com/huichen/kunlun/pkg/types"
)

// 使用 key1、key2 以及他们之间的距离 distance 在索引中查询满足条件的文档及分段的起始位置
// 分段需要匹配 keyword，同时 key1 在分段中的开始位置为 offset
// 参数：
//		key2：当为零时，只对 key1 做单键匹配
//		distance: 当且仅当 key2 为零时该值为零
//		caseSensitive：匹配时是否考虑大小写
//		isSymbol：如果为 true 则在符号索引中搜索，否则在全文索引搜索
//		shouldDocBeRecalled：外部钩子函数，用于判断某个文档 ID 是否应该被召回
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
	responseChan := make(chan contentSearchResponseInfo, indexer.numIndexerShards)

	// 启动搜索协程
	requestChan := make(chan contentSearchRequestInfo, indexer.numIndexerShards)
	for i := 0; i < indexer.numIndexerShards; i++ {
		go indexer.contentSearchWorker(i, requestChan, responseChan)
	}

	// 发送搜索请求
	for i := 0; i < indexer.numIndexerShards; i++ {
		requestChan <- contentSearchRequestInfo{
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

	return docLocationsToSections(retDocs, uint32(len(keyword))), nil

}

// 用于向内容查询线程发送请求
type contentSearchRequestInfo struct {
	Keyword             []byte
	Offset              uint32
	Key1                ngram_index.IndexKey
	Key2                ngram_index.IndexKey
	Distance            uint32
	CaseSensitive       bool
	IsSymbol            bool
	ShouldDocBeRecalled func(uint64) bool
}

// 用于从内容查询线程接受请求
type contentSearchResponseInfo struct {
	Err                    error
	DocumentsWithLocations []ngram_index.DocumentWithLocations
}

// 内容查询线程
func (indexer *Indexer) contentSearchWorker(
	shard int,
	requestChan chan contentSearchRequestInfo,
	responseChan chan contentSearchResponseInfo,
) {
	info := <-requestChan
	key1 := info.Key1
	key2 := info.Key2
	distance := info.Distance
	docFilterFunc := info.ShouldDocBeRecalled

	var contentMatchDocs []ngram_index.DocumentWithLocations
	var err error
	var docs []ngram_index.DocumentWithLocations
	if key2 != 0 {
		docs, err = indexer.contentNgramIndices[shard].SearchTwoKeys(
			key1, key2, distance, docFilterFunc, info.IsSymbol)
	} else {
		if distance != 0 {
			logger.Fatal("key2 为 0 时 distance 必须也为 0")
		}
		// 单键的情况
		docs, err = indexer.contentNgramIndices[shard].SearchOneKey(key1, docFilterFunc, info.IsSymbol)
	}
	if err != nil {
		responseChan <- contentSearchResponseInfo{
			Err: err,
		}
		return
	}
	if info.Keyword == nil {
		responseChan <- contentSearchResponseInfo{
			Err:                    nil,
			DocumentsWithLocations: docs,
		}
		return
	}

	contentMatchDocs, err = indexer.filterDocumentsWithFullMatch(
		&(indexer.documentIDToContentMap), docs, info.Offset, info.Keyword, info.CaseSensitive)
	if err != nil {
		responseChan <- contentSearchResponseInfo{
			Err: err,
		}
		return
	}

	responseChan <- contentSearchResponseInfo{
		Err:                    nil,
		DocumentsWithLocations: contentMatchDocs,
	}
}
