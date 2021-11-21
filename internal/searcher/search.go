package searcher

import (
	"errors"
	"time"

	"kunlun/internal/indexer"
	"kunlun/pkg/log"
	"kunlun/pkg/types"
)

var (
	logger = log.GetLogger()
)

// 做一次搜索
func (schr *Searcher) Search(idxr *indexer.Indexer, request types.SearchRequest) (*types.SearchResponse, error) {
	startTime := time.Now()

	query, err := ParseQuery(request.Query)

	if err != nil {
		return nil, err
	}

	if query == nil {
		return nil, errors.New("query 不能为空")
	}

	docFilter := NewDocFilter(query, idxr, request.ShouldRecallRepo, request.DocumentIDs)

	context := &Context{
		searchStartTime: &startTime,
		idxr:            idxr,
		query:           query,
		docFilter:       docFilter,
		request:         &request,
		searcherOptions: schr.options,
	}
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	if query.TrimmedQuery == nil {
		if query.FileQuery != nil {
			return schr.searchFiles(context, idxr, request)
		}

		if query.RepoQuery != nil {
			return schr.searchRepos(context)
		}

		return nil, errors.New("搜索表达式为空的情况下 file: 和 repo: 不能都为空")
	}

	if query.TrimmedQuery.Negate {
		return nil, errors.New("搜索表达式不能为纯非操作")
	}

	// 第一步：计算 query 中的单串值
	q := query.TrimmedQuery
	err = schr.searchTokenQuery(context, q)
	if err != nil {
		return nil, err
	}

	// 第二步：合并单串树
	err = schr.mergeTreeNodes(context, q)
	if err != nil {
		return nil, err
	}

	// 软硬优化迭代：
	// 1、如果能做软优化，则合并结果，并循环；
	// 2、如果已经无法做软优化了，做一次硬优化并合并结果，然后重复第 1 步
	for !context.query.DoneCompute() {
		for {
			softComputed, err := schr.softComputeOneRegexNode(context, q)
			if err != nil {
				return nil, err
			}
			if !softComputed {
				break
			}

			err = schr.mergeTreeNodes(context, q)
			if err != nil {
				return nil, err
			}
		}

		err = schr.hardComputeOneRegexNode(context, q)
		if err != nil {
			return nil, err
		}

		err = schr.mergeTreeNodes(context, q)
		if err != nil {
			return nil, err
		}
	}
	now := time.Now()
	context.recallEndTime = &now

	if q.Negate {
		return nil, errors.New("搜索表达式不能为纯非操作")
	}

	// 整理格式输出
	return annotateResponse(context, idxr, request)
}
