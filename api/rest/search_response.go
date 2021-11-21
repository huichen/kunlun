package rest

type SearchResponse struct {
	ResponseBase

	Repos                        []Repo `json:"repos"`
	NumRepos                     int    `json:"numRepos"`
	NumDocuments                 int    `json:"numDocuments"`
	NumLines                     int    `json:"numLines"`
	NumSections                  int    `json:"numSections"`
	NumRegexMatches              int    `json:"numRegexMatches"`
	SearchDurationInMicroSeconds int64  `json:"searchDurationInMicroSeconds"`
	RecallDurationInMicroSeconds int64  `json:"recallDurationInMicroSeconds"`
}

type Repo struct {
	RepoID uint64 `json:"repoID"`

	LocalPath string `json:"localPath,omitempty"`

	RemoteURL string `json:"remoteURL,omitempty"`

	NumDocumentsInRepo int `json:"numDocumentsInRepo"`

	NumLinesInRepo int `json:"numLinesInRepo"`

	Documents []Document `json:"documents"`
}

type Document struct {
	DocumentID uint64 `json:"documentID"`

	Language string `json:"language,omitempty"`

	Filename string `json:"filename,omitempty"`

	NumLinesInDocument int `json:"numLinesInDocument"`

	Lines []Line `json:"lines,omitempty"`
}

type Line struct {
	LineNumber uint32 `json:"lineNumber"`

	Content string `json:"content,omitempty"`

	Highlights []Section `json:"highlights,omitempty"`
}

type Section struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
