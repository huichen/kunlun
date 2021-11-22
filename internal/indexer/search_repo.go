package indexer

type SearchRepoRequest struct {
	// 文档过滤器
	RepoFilter func(repoID uint64) bool
}

type SearchRepoResponse struct {
	Repos []*CodeRepository
}

// 利用外部钩子函数 RepoFilter 搜索匹配的仓库
func (indexer *Indexer) SearchRepos(request *SearchRepoRequest) SearchRepoResponse {
	retRepoIDs := []*CodeRepository{}
	for repoID, repo := range indexer.idToRepoMap {
		if request.RepoFilter(repoID) {
			retRepoIDs = append(retRepoIDs, repo)
		}
	}
	return SearchRepoResponse{
		Repos: retRepoIDs,
	}
}
