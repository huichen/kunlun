package types

// 搜索返回结果
type SearchResponse struct {
	// 搜索得到的仓库
	Repos []*SearchedRepo

	// 总共检索到多少仓库，召回的总数量，包括因分页等截断的
	NumRepos int

	// 总共检索到多少文档，召回的总数量，包括因分页等截断的
	NumDocuments int

	// 检索到多少行，有分页截断
	NumLines int

	// 总共检索到多少片段，召回的总数量，包括因分页等截断的
	NumSections int

	// 做了多少次正则表达式匹配，一个（正则表达式，文档）算一次
	NumRegexMatches int

	// 搜索耗时，微秒
	SearchDurationInMicroSeconds int64

	// 召回耗时（不包含排序等），微秒
	RecallDurationInMicroSeconds int64
}

type SearchedRepo struct {
	// RepoID 是仓库的唯一标识，所有没有对应仓库的文档放在 RepoID 为 0 的“空仓库”中
	// 空仓库的 LocalPath、RemoteURL 均为空
	RepoID uint64

	// 仓库在文件系统中的路径，对没有本地拷贝的远程仓库来说，该值为空
	LocalPath string `json:"LocalPath,omitempty"`

	// 仓库的远程路径，对没有远程分支的仓库，该值为空
	RemoteURL string `json:"RemoteURL,omitempty"`

	// 该仓库内总共检索到多少文档，召回的总数量，包括因分页等截断的
	NumDocumentsInRepo int
	NumLinesInRepo     int
	NumSectionsInRepo  int

	// 搜索得到的文档
	Documents []SearchedDocument
}

// 仓库的文本表达，如果是远程仓库，优先用 RemoteURL，否则用 LocalPath
func (repo *SearchedRepo) String() string {
	if repo.RemoteURL != "" {
		return repo.RemoteURL
	}

	return repo.LocalPath
}

type SearchedDocument struct {
	// 在索引器中用来标识文档的 ID
	DocumentID uint64

	// 语言
	Language string `json:"Language,omitempty"`

	// 返回 repo path 或者文件系统路径
	Filename string `json:"Filename,omitempty"`

	// 文档打分
	Scores []float32 `json:"Scores,omitempty"`

	// 该文档内总共检索到多少行，召回的总数量（截断之前的）
	NumLinesInDocument int

	// 该文档内总共检索到多少匹配字段，召回的总数量（截断之前的）
	NumSectionsInDocument int

	// 返回匹配的行结果
	// 如果 request 里面设置了 MaxResultsPerDocument，那么最多返回这么多结果
	Lines []Line `json:"Lines,omitempty"`
}

type Line struct {
	// 行号
	LineNumber uint32

	// 行内容
	Content []byte `json:"Content,omitempty"`

	// 行高亮，标识了高亮段在行中的起止位置
	Highlights []Section `json:"Highlights,omitempty"`
}
