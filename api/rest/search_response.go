package rest

// 搜索返回结果，按照下面的层级结构来组织
//		仓库（repos） -> 文档（documents） -> 行（lines） -> 片段（sections）
// 不属于一个仓库的文档也会被组织到一个空仓库（RepoID == 0）里方便展示
// “片段”指的是一个匹配的字符串，是搜索匹配的最小单元
type SearchResponse struct {
	ResponseBase

	Repos                        []Repo     `json:"repos,omitempty"`
	Languages                    []Language `json:"language,omitempty"`
	NumRepos                     int        `json:"numRepos"`
	NumDocuments                 int        `json:"numDocuments"`
	NumLines                     int        `json:"numLines"`
	NumSections                  int        `json:"numSections"`
	NumRegexMatches              int        `json:"numRegexMatches"`
	SearchDurationInMicroSeconds int64      `json:"searchDurationInMicroSeconds"`
	RecallDurationInMicroSeconds int64      `json:"recallDurationInMicroSeconds"`
	Responsetype                 string     `json:"responseType"`
}

type Language struct {
	LanguageID             uint64 `json:"languageID"`
	Name                   string `json:"name"`
	NumDocumentsInLanguage int    `json:"numDocumentsInLanguage"`
	NumSectionsInLanguage  int    `json:"numSectionsInLanguage"`
}

type Repo struct {
	RepoID             uint64     `json:"repoID"`
	LocalPath          string     `json:"localPath,omitempty"`
	RemoteURL          string     `json:"remoteURL,omitempty"`
	NumDocumentsInRepo int        `json:"numDocumentsInRepo"`
	NumLinesInRepo     int        `json:"numLinesInRepo"`
	NumSectionsInRepo  int        `json:"numSectionsInRepo"`
	Documents          []Document `json:"documents,omitempty"`
}

type Document struct {
	DocumentID            uint64 `json:"documentID"`
	Language              string `json:"language,omitempty"`
	Filename              string `json:"filename,omitempty"`
	NumLinesInDocument    int    `json:"numLinesInDocument"`
	NumSectionsInDocument int    `json:"numSectionsInDocument"`
	Lines                 []Line `json:"lines,omitempty"`
}

type Line struct {
	LineNumber uint32    `json:"lineNumber"`
	Content    string    `json:"content,omitempty"`
	Highlights []Section `json:"highlights,omitempty"`
}

type Section struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
