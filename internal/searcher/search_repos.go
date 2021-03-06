package searcher

import (
	"errors"
	"time"

	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/pkg/types"
)

// 使用表达式对仓库做检索
// 表达式为 repo:xxx 的形式，且不做内容或文件名匹配，只搜索匹配的仓库名
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
