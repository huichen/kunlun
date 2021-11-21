package searcher

import (
	"regexp"

	"kunlun/internal/query"
	"kunlun/pkg/types"
)

func ParseQuery(pattern string) (*SearchQuery, error) {
	originalQuery, err := query.Parse(pattern)
	if err != nil {
		return nil, err
	}

	if originalQuery == nil {
		return nil, nil
	}

	// query 校验
	err = validateQuery(originalQuery)
	if err != nil {
		return nil, err
	}

	// 解析出来的大小写类型也用于修饰词
	queryCase := getCase(originalQuery)

	searchQ := query.Copy(originalQuery)
	searchQ = TrimQuery(searchQ)

	fileQ, err := collectModifier(originalQuery, query.FileQuery)
	if err != nil {
		return nil, err
	}
	if fileQ != nil {
		fileQ.Case = queryCase
	}

	repoQ, err := collectModifier(originalQuery, query.RepoQuery)
	if err != nil {
		return nil, err
	}
	if repoQ != nil {
		repoQ.Case = queryCase
	}

	langQ, err := collectModifier(originalQuery, query.LanguageQuery)
	if err != nil {
		return nil, err
	}

	numNodes := 0
	if searchQ != nil {
		numNodes = searchQ.NumNodes
	}
	retQuery := &SearchQuery{
		OriginalQuery: originalQuery,
		TrimmedQuery:  searchQ,
		QueryResults:  make([]*[]types.DocumentWithSections, numNodes),

		Case: queryCase,

		// 修饰词
		LanguageQuery: langQ,
		LangRe:        compileRE(langQ),
		LanguageNames: getNames(langQ),
		RepoQuery:     repoQ,
		RepoRe:        compileRE(repoQ),
		RepoNames:     getNames(repoQ),
		FileQuery:     fileQ,
		FileRe:        compileRE(fileQ),
	}

	return retQuery, nil
}

// 大小写不敏感匹配
func compileRE(q *query.Query) []*regexp.Regexp {
	if q == nil {
		return nil
	}

	ret := make([]*regexp.Regexp, q.NumNodes)

	internalCompileRE(q, &ret)

	return ret
}

func internalCompileRE(q *query.Query, res *[]*regexp.Regexp) {
	switch q.Type {
	case query.LanguageQuery:
		re, err := regexp.Compile("(?i)" + q.Language)
		if err == nil {
			(*res)[q.ID] = re
		}
	case query.RepoQuery:
		re, err := regexp.Compile("(?i)" + q.RepoRegexString)
		if err == nil {
			(*res)[q.ID] = re
		}

	case query.FileQuery:
		re, err := regexp.Compile("(?i)" + q.FileRegexString)
		if err == nil {
			(*res)[q.ID] = re
		}

	case query.TreeQuery:
		for _, sq := range q.SubQueries {
			internalCompileRE(sq, res)
		}
	}
}

// 大小写不敏感匹配
func getNames(q *query.Query) []string {
	if q == nil {
		return nil
	}

	ret := make([]string, q.NumNodes)

	internalGetNames(q, &ret)

	return ret
}

func internalGetNames(q *query.Query, res *[]string) {
	switch q.Type {
	case query.LanguageQuery:
		(*res)[q.ID] = q.Language
	case query.RepoQuery:
		(*res)[q.ID] = q.RepoRegexString
	case query.TreeQuery:
		for _, sq := range q.SubQueries {
			internalGetNames(sq, res)
		}
	}
}
