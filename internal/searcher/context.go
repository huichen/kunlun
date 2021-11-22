package searcher

import (
	"errors"
	"fmt"
	"time"

	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/internal/query"
	"github.com/huichen/kunlun/pkg/types"
)

// 一次搜索请求的上下文信息
type Context struct {
	searchStartTime *time.Time
	recallEndTime   *time.Time

	idxr      *indexer.Indexer
	query     *SearchQuery
	docFilter *DocFilter
	request   *types.SearchRequest

	searcherOptions *types.SearcherOptions

	// 保存该次请求正则表达式搜索做了多少次，在一个文档做一次匹配则此计数器加一
	regexSearchTimes int

	docIDToDocumentWithSectionsMap map[uint64]*types.DocumentWithSections
}

func (c *Context) getSearchTokenRequest(token string, symbol bool) indexer.SearchTokenRequest {
	return indexer.SearchTokenRequest{
		Token:         token,
		IsSymbol:      symbol,
		CaseSensitive: c.query.Case,
		DocFilter:     c.docFilter.ShouldRecallDocument,
	}
}

func (c *Context) getSearchRegexRequest(q *query.Query, computedQuery *query.Query, symbol bool) (*indexer.SearchRegexRequest, error) {
	if q.Type != query.RegexQuery {
		return nil, errors.New("query 类型不为 RegexQuery")
	}

	var candidateDocs *[]uint64
	negate := false
	if computedQuery != nil {
		r := c.query.QueryResults[computedQuery.ID]
		if r != nil {
			if len(*r) == 0 {
				return nil, errors.New("computedQuery 不能为空")
			}
			docs := []uint64{}
			for _, doc := range *r {
				docs = append(docs, doc.DocumentID)
			}
			candidateDocs = &docs
			negate = computedQuery.Negate
		}
	}

	return &indexer.SearchRegexRequest{
		Regex:               q.RegexString,
		Negate:              q.Negate,
		CaseSensitive:       c.query.Case,
		Tokens:              q.RegexTokens,
		IsSymbol:            symbol,
		MaxResultsPerFile:   c.request.MaxLinesPerDocument,
		CandidateDocs:       candidateDocs,
		CandidateDocsNegate: negate,
		DocFilter:           c.docFilter.ShouldRecallDocument,
	}, nil
}

// 在一次搜索中检查是否超时，如果超时则返回 err，上游收到后会立刻返回或者做简化处理
func (context *Context) checkTimeout() error {
	if context.request != nil && context.request.TimeoutInMs > 0 && context.searchStartTime != nil {
		if time.Since(*context.searchStartTime).Milliseconds() > int64(context.request.TimeoutInMs) {
			return fmt.Errorf("请求超时（%d ms），请缩小搜索范围", context.request.TimeoutInMs)
		}
	}

	return nil
}
