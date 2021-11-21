package ranker

import (
	"sort"
	"strings"

	"kunlun/pkg/types"
)

// 默认的排序器，先按照仓库名（local path 或者 remote URL）升序，然后按照文件名升序
type DefaultRanker struct {
}

func (ranker DefaultRanker) Rank(response *types.SearchResponse) {
	// 对 repo 排序
	sort.Slice(response.Repos, func(i, j int) bool {
		return compareRepos(response.Repos[i], response.Repos[j])
	})

	// 然后对 repo 内的文档排序

	for _, repo := range response.Repos {
		sort.Slice(repo.Documents, func(i, j int) bool {
			return compareDocuments(repo.Documents[i], repo.Documents[j])
		})
	}
}

// 先按照 repo 匹配的文件数倒序，再按照 repo 名字排序
// id = 0 的空 repo 始终排在最后面
func compareRepos(repo1 *types.SearchedRepo, repo2 *types.SearchedRepo) bool {
	if repo1.RepoID == 0 || repo2.RepoID == 0 {
		return repo1.RepoID > repo2.RepoID
	}

	if repo1.NumDocumentsInRepo != repo2.NumDocumentsInRepo {
		return repo1.NumDocumentsInRepo > repo2.NumDocumentsInRepo
	}

	repoPath1 := repo1.LocalPath
	if repoPath1 == "" {
		repoPath1 = repo1.RemoteURL
	}
	repoPath2 := repo2.LocalPath
	if repoPath2 == "" {
		repoPath2 = repo2.RemoteURL
	}

	return strings.Compare(repoPath1, repoPath2) < 0
}

// 先按照文档匹配的section数倒排序，如果相同按照文件名正序
func compareDocuments(doc1 types.SearchedDocument, doc2 types.SearchedDocument) bool {
	if doc1.NumSectionsInDocument != doc2.NumSectionsInDocument {
		return doc1.NumSectionsInDocument > doc2.NumSectionsInDocument
	}

	return strings.Compare(doc1.Filename, doc2.Filename) < 0
}
