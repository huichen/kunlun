package indexer

import "github.com/huichen/kunlun/pkg/types"

type SearchFileRequest struct {
	// 文档过滤器
	DocFilter func(docID uint64) bool
}

type SearchFileResponse struct {
	Documents []types.SearchedDocument
}

// 在索引中通过通过外部传入的 DocFilter 钩子函数搜索文件名
func (indexer *Indexer) SearchFile(request *SearchFileRequest) SearchFileResponse {
	retDocIDs := []types.SearchedDocument{}
	for _, meta := range indexer.documentIDToMetaMap {
		if request.DocFilter(meta.DocumentID) {
			// 语言
			lang := ""
			if meta.Language != nil {
				lang = meta.Language.Name
			}

			// 文件名优先使用仓库内路径
			filename := meta.PathInRepo
			if filename == "" {
				filename = meta.LocalPath
			}

			retDocIDs = append(retDocIDs, types.SearchedDocument{
				DocumentID: meta.DocumentID,
				Language:   lang,
				Filename:   filename,
			})
		}
	}
	return SearchFileResponse{
		Documents: retDocIDs,
	}
}
