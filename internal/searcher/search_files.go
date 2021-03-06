package searcher

import (
	"errors"
	"time"

	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/pkg/types"
)

// 使用表达式对文件做检索
// 表达式为 file:xxx 的形式，且不做内容匹配，只搜索匹配的文件名
func (schr *Searcher) searchFiles(context *Context, idxr *indexer.Indexer, request types.SearchRequest) (*types.SearchResponse, error) {
	fileQuery := context.query.FileQuery

	if fileQuery == nil {
		return nil, errors.New("file: 不能为空")
	}

	requestFileRequest := indexer.SearchFileRequest{
		DocFilter: context.docFilter.ShouldRecallDocument,
	}
	resp := context.idxr.SearchFile(&requestFileRequest)
	now := time.Now()
	context.recallEndTime = &now
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	outputDocs := resp.Documents

	// 将文档通过 repo 组织为返回格式
	response, err := transformSearchedDocsToResponse(context, idxr, outputDocs)
	if err != nil {
		return nil, err
	}
	response.ResponseType = "files"
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	// 排序
	rkr := getRanker(request)
	rkr.Rank(response)
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	// 对一个仓库最多能有多少文档做过滤
	trimRepo(context, response)
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	// 分页
	paginateRepos(context, response)
	paginateDocuments(context, response)
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	// 附加内容
	appendContentToResponse(context, idxr, response)
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	// 添加延时信息
	appendTimingInfo(context, response)

	return response, nil
}
