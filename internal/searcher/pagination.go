package searcher

import "kunlun/pkg/types"

// 根据输入参数对返回的仓库做分页截断
func paginateRepos(context *Context, response *types.SearchResponse) {
	if context.request.PageSize > 0 {
		start := context.request.PageSize * context.request.PageNum
		if start > len(response.Repos) {
			response.Repos = []*types.SearchedRepo{}
		} else {
			end := start + context.request.PageSize
			if end > len(response.Repos) {
				end = len(response.Repos)
			}
			response.Repos = response.Repos[start:end]
		}
	}
}

// 分页截断文档，仅当只有一个仓库的时候做此操作
func paginateDocuments(context *Context, response *types.SearchResponse) {
	if len(response.Repos) != 1 {
		return
	}
	repo := response.Repos[0]

	if context.request.PageSize > 0 {
		start := context.request.PageSize * context.request.PageNum
		if start > len(repo.Documents) {
			repo.Documents = []types.SearchedDocument{}
		} else {
			end := start + context.request.PageSize
			if end > len(response.Repos) {
				end = len(response.Repos)
			}
			response.Repos = response.Repos[start:end]
		}
	}
}
