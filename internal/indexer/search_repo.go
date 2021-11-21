package indexer

type SearchRepoRequest struct {
	// 文档过滤器
	RepoFilter func(repoID uint64) bool
}

type SearchRepoResponse struct {
	Repos []*CodeRepository
}

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
