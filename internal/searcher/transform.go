package searcher

import (
	"kunlun/internal/indexer"
	"kunlun/pkg/types"
)

func transformSearchedDocsToResponse(context *Context, idxr *indexer.Indexer, docs []types.SearchedDocument) (*types.SearchResponse, error) {
	outputRepos := []*types.SearchedRepo{}
	searchedRepoMap := make(map[uint64]*types.SearchedRepo)
	for _, doc := range docs {
		meta := idxr.GetMeta(doc.DocumentID)

		// 识别文档的 repo
		var repo *types.SearchedRepo
		var repoID uint64
		var repoLocalPath string
		var repoRemoteURL string
		if meta.Repo != nil {
			repoID = meta.Repo.ID
			repoLocalPath = meta.Repo.LocalPath
			repoRemoteURL = meta.Repo.RemoteURL
		}

		var ok bool
		if repo, ok = searchedRepoMap[repoID]; !ok {
			repo = &types.SearchedRepo{
				RepoID:    repoID,
				LocalPath: repoLocalPath,
				RemoteURL: repoRemoteURL,
				Documents: []types.SearchedDocument{doc},
			}
			outputRepos = append(outputRepos, repo)
			searchedRepoMap[repoID] = repo
		} else {
			if len(repo.Documents) == 0 {
				repo.Documents = []types.SearchedDocument{doc}
			} else {
				repo.Documents = append(repo.Documents, doc)
			}
		}
	}

	totalSections := 0
	for id := range outputRepos {
		outputRepos[id].NumDocumentsInRepo = len(outputRepos[id].Documents)

		numSections := 0
		for _, doc := range outputRepos[id].Documents {
			numSections += doc.NumSectionsInDocument
		}
		outputRepos[id].NumSectionsInRepo = numSections
		totalSections += numSections
	}

	response := types.SearchResponse{
		Repos:           outputRepos,
		NumRepos:        len(outputRepos),
		NumDocuments:    len(docs),
		NumSections:     totalSections,
		NumRegexMatches: context.regexSearchTimes,
	}

	return &response, nil
}
