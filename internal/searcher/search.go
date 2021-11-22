package searcher

import (
	"errors"
	"time"

	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/pkg/log"
	"github.com/huichen/kunlun/pkg/types"
)

var (
	logger = log.GetLogger()
)

// 做一次搜索
func (schr *Searcher) Search(idxr *indexer.Indexer, request types.SearchRequest) (*types.SearchResponse, error) {
	// 记录起始时间
	startTime := time.Now()

	// 解析搜索表达式
	query, err := ParseQuery(request.Query)
	if err != nil {
		return nil, err
	}
	if query == nil {
		return nil, errors.New("query 不能为空")
	}

	// 生成过滤的钩子函数
	docFilter := NewDocFilter(query, idxr, request.ShouldRecallRepo, request.DocumentIDs)

	// 封装上下文
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

	// 当内容搜索为空的时候，看看是不是为只搜索文件或者仓库的情况
	if query.TrimmedQuery == nil {
		if query.FileQuery != nil {
			return schr.searchFiles(context, idxr, request)
		}

		if query.RepoQuery != nil {
			return schr.searchRepos(context)
		}

		return nil, errors.New("搜索表达式为空的情况下 file: 和 repo: 不能都为空")
	}

	// 合法性校验
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

	// 我们采取软硬迭代的方式优化正则表达式的计算过程：
	// 1、如果能做正则表达式软优化，则先做软优化，然后合并树节点，并循环；
	// 2、如果已经无法做软优化了，做一次正则表达式硬优化再合并树节点，然后重复第 1 步
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

	// 记录召回时间
	now := time.Now()
	context.recallEndTime = &now

	// 合法性校验
	if q.Negate {
		return nil, errors.New("搜索表达式不能为纯非操作")
	}

	// 整理格式输出
	return annotateResponse(context, idxr, request)
}
