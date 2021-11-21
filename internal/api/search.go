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
	// 识别 search type
	query := req.URL.Query().Get("q")
	if query == "" {
		ReturnError(w, req, http_error.MissingParameter)
		return
	}

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

	var page int
	pageStr := req.URL.Query().Get("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			ReturnError(w, req, http_error.GetError(err))
			return
		}
	}

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

	sr, err := kgn.Search(request)
	if err != nil {
		ReturnError(w, req, http_error.GetError(err))
		return
	}

	response := rest.SearchResponse{
		NumRepos:                     sr.NumRepos,
		NumDocuments:                 sr.NumDocuments,
		NumLines:                     sr.NumLines,
		NumSections:                  sr.NumSections,
		NumRegexMatches:              sr.NumRegexMatches,
		SearchDurationInMicroSeconds: sr.SearchDurationInMicroSeconds,
		RecallDurationInMicroSeconds: sr.RecallDurationInMicroSeconds,
	}

	// 添加 repos 和 documents
	for _, repo := range sr.Repos {
		repoResponse := rest.Repo{
			RepoID:             repo.RepoID,
			RemoteURL:          repo.RemoteURL,
			LocalPath:          repo.LocalPath,
			NumDocumentsInRepo: repo.NumDocumentsInRepo,
			NumLinesInRepo:     repo.NumLinesInRepo,
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
				DocumentID:         doc.DocumentID,
				Language:           doc.Language,
				Filename:           doc.Filename,
				Lines:              lines,
				NumLinesInDocument: doc.NumLinesInDocument,
			})
		}
		response.Repos = append(response.Repos, repoResponse)
	}
	resp, _ := json.Marshal(&response)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(resp))
}

func trimOutput(response *types.SearchResponse) {
	for id := range response.Repos {
		response.Repos[id].LocalPath = ""
	}
}
