package searcher

import (
	"errors"
	"time"

	"kunlun/internal/indexer"
	"kunlun/pkg/types"
)

func (schr *Searcher) searchRepos(context *Context) (*types.SearchResponse, error) {
	repoQuery := context.query.RepoQuery

	if repoQuery == nil {
		return nil, errors.New("repo: 不能为空")
	}

	request := indexer.SearchRepoRequest{
		RepoFilter: context.docFilter.ShouldRecallRepo,
	}
	resp := context.idxr.SearchRepos(&request)
	if err := context.checkTimeout(); err != nil {
		return nil, err
	}

	outputRepos := []*types.SearchedRepo{}
	for _, repo := range resp.Repos {
		outputRepos = append(outputRepos, &types.SearchedRepo{
			RepoID:    repo.ID,
			LocalPath: repo.LocalPath,
			RemoteURL: repo.RemoteURL,
		})
	}

	searchResponse := types.SearchResponse{
		Repos:                        outputRepos,
		NumRepos:                     len(outputRepos),
		SearchDurationInMicroSeconds: time.Since(*context.searchStartTime).Microseconds(),
		ResponseType:                 "repos",
	}

	// 对结果排序
	rkr := getRanker(*context.request)
	rkr.Rank(&searchResponse)

	return &searchResponse, nil
}
