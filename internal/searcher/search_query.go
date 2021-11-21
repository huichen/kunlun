package searcher

import (
	"regexp"

	"kunlun/internal/query"
	"kunlun/pkg/types"
)

type SearchQuery struct {
	// 用户输入中解析出来的搜索表达式
	OriginalQuery *query.Query

	// 从 OriginalQuery 中提取的搜索表达式
	// 去掉了 file:/repo: 等 modifier 信息，并做了一些优化，方便查询使用
	TrimmedQuery *query.Query

	// 保存 TrimmedQuery 计算得到的 docID slice，长度和 TrimmedQuery 的节点数相同
	QueryResults []*[]types.DocumentWithSections

	// 从 OriginalQuery 中解析出的大小写
	Case bool

	// 从 OriginalQuery 中解析出的文件搜索条件
	LanguageQuery *query.Query
	LangRe        []*regexp.Regexp
	LanguageNames []string

	// 从 OriginalQuery 中解析出的代码仓库搜索条件
	RepoQuery *query.Query
	RepoRe    []*regexp.Regexp
	RepoNames []string

	// 从 OriginalQuery 中解析出的文件搜索条件
	FileQuery *query.Query
	FileRe    []*regexp.Regexp
}

// 判断是否计算完成，也就是 q.QueryResults 的最后一个元素不为 nil
func (q *SearchQuery) DoneCompute() bool {
	l := len(q.QueryResults)
	if l == 0 {
		return true
	}
	return q.QueryResults[l-1] != nil
}

func (q *SearchQuery) Result() *[]types.DocumentWithSections {
	l := len(q.QueryResults)
	if l == 0 {
		return nil
	}
	return q.QueryResults[l-1]
}
