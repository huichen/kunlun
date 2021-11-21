package searcher

import (
	"kunlun/internal/indexer"
	"kunlun/internal/ranker"
	"kunlun/pkg/types"
)

// 对搜索结果做注解，比如添加行信息、排序、过滤、分页等等
func annotateResponse(context *Context, idxr *indexer.Indexer, request types.SearchRequest) (*types.SearchResponse, error) {
	// 从搜索结果中收集文档
	docs, err := accumulateResultsToDocuments(context, idxr, request)
	if err != nil {
		return nil, err
	}
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	// 将文档通过 repo 组织为返回格式
	response, err := transformSearchedDocsToResponse(context, idxr, docs)
	if err != nil {
		return nil, err
	}
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
	appendLinesToResponse(context, idxr, response)
	appendContentToResponse(context, idxr, response)

	// 添加延时信息
	appendTimingInfo(context, response)

	return response, nil
}

func accumulateResultsToDocuments(
	context *Context, idxr *indexer.Indexer, request types.SearchRequest) ([]types.SearchedDocument, error) {
	results := context.query.Result()
	if results == nil {
		return []types.SearchedDocument{}, nil
	}

	contextLines := request.NumContextLines
	if contextLines < 0 {
		contextLines = 0
	}

	start := 0
	end := len(*results)
	retDocs := []types.SearchedDocument{}
	context.docIDToDocumentWithSectionsMap = make(map[uint64]*types.DocumentWithSections)
	for indexDoc := start; indexDoc < end; indexDoc++ {
		doc := (*results)[indexDoc]
		context.docIDToDocumentWithSectionsMap[doc.DocumentID] = &doc

		meta := idxr.GetMeta(doc.DocumentID)
		filename := ""
		if meta.PathInRepo != "" {
			filename = meta.PathInRepo
		} else {
			filename = meta.LocalPath
		}

		// 语言
		lang := ""
		if meta.Language != nil {
			lang = meta.Language.Name
		}

		retDocs = append(retDocs, types.SearchedDocument{
			DocumentID:            doc.DocumentID,
			Language:              lang,
			Filename:              filename,
			NumSectionsInDocument: len(doc.Sections),
		})

	}

	return retDocs, nil
}

func getRanker(request types.SearchRequest) types.Ranker {
	rkr := request.Ranker
	if rkr == nil {
		rkr = ranker.DefaultRanker{}
	}
	return rkr
}
