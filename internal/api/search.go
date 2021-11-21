package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"kunlun/api/rest"
	"kunlun/api/rest/http_error"
	"kunlun/internal/resource/engine"
	"kunlun/pkg/types"
)

var (
	defaultLineContext                 = 2
	defaultPageSize                    = 200
	defaultMaxLinesPerDoc              = 3
	defaultMaxDocsPerRepo              = 3
	defaultTimeout                     = 2000
	defaultMaxDocPerRepoInSingleReturn = 200
)

func Search(w http.ResponseWriter, req *http.Request) {
	// 得到搜索表达式
	query := req.URL.Query().Get("q")
	if query == "" {
		ReturnError(w, req, http_error.MissingParameter)
		return
	}

	// 单文档搜索参数
	var docID int
	var err error
	idStr := req.URL.Query().Get("id")
	if idStr != "" {
		docID, err = strconv.Atoi(idStr)
		if err != nil {
			ReturnError(w, req, http_error.GetError(err))
			return
		}
	}

	// 两个分页参数
	var pageSize int
	pageSizeStr := req.URL.Query().Get("pageSize")
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			ReturnError(w, req, http_error.GetError(err))
			return
		}
		if pageSize > defaultPageSize {
			pageSize = defaultPageSize
		}
	} else {
		pageSize = defaultPageSize
	}
	var page int
	pageStr := req.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			ReturnError(w, req, http_error.GetError(err))
			return
		}
	}

	// 引擎
	kgn := engine.GetEngine()

	// 生成查询请求
	var ids []uint64
	if docID != 0 {
		ids = []uint64{uint64(docID)}
	}
	lineContext := defaultLineContext
	if docID != 0 {
		lineContext = 10000
	}
	request := types.SearchRequest{
		Query:                          query,
		DocumentIDs:                    ids,
		MaxDocumentsPerRepo:            defaultMaxDocsPerRepo,
		MaxDocumentsInSingleRepoReturn: defaultMaxDocPerRepoInSingleReturn,
		ReturnLineContent:              true,
		NumContextLines:                lineContext,
		MaxLinesPerDocument:            defaultMaxLinesPerDoc,
		PageSize:                       pageSize,
		PageNum:                        page,
		HightlightStartTag:             "<b class=\"keywords\">",
		HightlightEndTag:               "</b>",
		TimeoutInMs:                    defaultTimeout,
	}

	// 调用引擎查询
	sr, err := kgn.Search(request)
	if err != nil {
		ReturnError(w, req, http_error.GetError(err))
		return
	}

	// 组织返回结构体
	response := rest.SearchResponse{
		NumRepos:                     sr.NumRepos,
		NumDocuments:                 sr.NumDocuments,
		NumLines:                     sr.NumLines,
		NumSections:                  sr.NumSections,
		NumRegexMatches:              sr.NumRegexMatches,
		SearchDurationInMicroSeconds: sr.SearchDurationInMicroSeconds,
		RecallDurationInMicroSeconds: sr.RecallDurationInMicroSeconds,
		Responsetype:                 sr.ResponseType,
	}
	// 添加 language
	languages := []rest.Language{}
	for _, lang := range sr.Languages {
		languages = append(languages, rest.Language{
			LanguageID:             lang.LanguageID,
			Name:                   lang.Name,
			NumDocumentsInLanguage: lang.NumDocumentsInLanguage,
			NumSectionsInLanguage:  lang.NumSectionsInLanguage,
		})
	}
	response.Languages = languages
	// 添加 repos 和 documents
	for _, repo := range sr.Repos {
		repoResponse := rest.Repo{
			RepoID:             repo.RepoID,
			RemoteURL:          repo.RemoteURL,
			LocalPath:          repo.LocalPath,
			NumDocumentsInRepo: repo.NumDocumentsInRepo,
			NumLinesInRepo:     repo.NumLinesInRepo,
			NumSectionsInRepo:  repo.NumSectionsInRepo,
		}
		for _, doc := range repo.Documents {
			lines := []rest.Line{}
			for _, l := range doc.Lines {
				highlights := []rest.Section{}
				for _, s := range l.Highlights {
					highlights = append(highlights, rest.Section{
						Start: int(s.Start),
						End:   int(s.End),
					})
				}
				lines = append(lines, rest.Line{
					LineNumber: l.LineNumber,
					Content:    string(l.Content),
					Highlights: highlights,
				})
			}

			repoResponse.Documents = append(repoResponse.Documents, rest.Document{
				DocumentID:            doc.DocumentID,
				Language:              doc.Language,
				Filename:              doc.Filename,
				Lines:                 lines,
				NumLinesInDocument:    doc.NumLinesInDocument,
				NumSectionsInDocument: doc.NumSectionsInDocument,
			})
		}
		response.Repos = append(response.Repos, repoResponse)
	}

	// 正常返回
	resp, _ := json.Marshal(&response)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(resp))
}
