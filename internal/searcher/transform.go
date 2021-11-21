package searcher

import (
	"sort"
	"strings"

	"github.com/huichen/kunlun/internal/indexer"

	"github.com/huichen/kunlun/pkg/types"
)

func transformSearchedDocsToResponse(context *Context, idxr *indexer.Indexer, docs []types.SearchedDocument) (*types.SearchResponse, error) {
	outputRepos := []*types.SearchedRepo{}
	searchedRepoMap := make(map[uint64]*types.SearchedRepo)
	langMap := make(map[uint64]*types.SearchedLanguage)
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

		// 识别文档语言
		var langID uint64
		var langName string
		if meta.Language != nil {
			langID = meta.Language.ID
			langName = meta.Language.Name

			if l, ok := langMap[langID]; ok {
				l.NumDocumentsInLanguage = l.NumDocumentsInLanguage + 1
				l.NumSectionsInLanguage += doc.NumSectionsInDocument
			} else {
				lang := types.SearchedLanguage{
					LanguageID:             langID,
					Name:                   langName,
					NumDocumentsInLanguage: 1,
					NumSectionsInLanguage:  doc.NumSectionsInDocument,
				}
				langMap[langID] = &lang
			}
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

	outputLangs := []types.SearchedLanguage{}
	for _, l := range langMap {
		outputLangs = append(outputLangs, *l)
	}

	sort.Slice(outputLangs, func(i, j int) bool {
		if outputLangs[i].NumDocumentsInLanguage != outputLangs[j].NumDocumentsInLanguage {
			return outputLangs[i].NumDocumentsInLanguage > outputLangs[j].NumDocumentsInLanguage
		}

		if outputLangs[i].NumSectionsInLanguage != outputLangs[j].NumSectionsInLanguage {
			return outputLangs[i].NumSectionsInLanguage > outputLangs[j].NumSectionsInLanguage
		}

		return strings.Compare(outputLangs[i].Name, outputLangs[j].Name) < 0
	})

	response := types.SearchResponse{
		Repos:           outputRepos,
		Languages:       outputLangs,
		NumRepos:        len(outputRepos),
		NumDocuments:    len(docs),
		NumSections:     totalSections,
		NumRegexMatches: context.regexSearchTimes,
	}

	return &response, nil
}
