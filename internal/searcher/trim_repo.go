package searcher

import "kunlun/pkg/types"

// 当不止一个 repo 的时候，单个 repo 中的文档数不能超过 MaxDocumentsPerRepo
func trimRepo(context *Context, response *types.SearchResponse) {
	maxDocs := 0
	if len(response.Repos) == 1 {
		if context.request.MaxDocumentsInSingleRepoReturn == 0 {
			return
		} else {
			maxDocs = context.request.MaxDocumentsInSingleRepoReturn
		}
	} else {
		if context.request.MaxDocumentsPerRepo <= 0 {
			return
		} else {
			maxDocs = context.request.MaxDocumentsPerRepo
		}
	}

	for _, repo := range response.Repos {
		if len(repo.Documents) <= maxDocs {
			continue
		}
		repo.Documents = repo.Documents[:maxDocs]
	}
}
