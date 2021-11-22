package query

import (
	"errors"
	"regexp/syntax"
	"strings"
)

// 解析单 term
// term 可以是一个 token，或者正则表达式
func parseTerm(term string) (*Query, error) {
	if len(term) == 0 {
		return nil, nil
	}

	retQuery := &Query{}

	// 符号
	if strings.HasPrefix(term, "sym:") {
		q, err := parseTerm(strings.TrimPrefix(term, "sym:"))
		if err != nil {
			return nil, err
		}
		q.IsSymbol = true
		return q, nil
	}

	// 文件名
	if strings.HasPrefix(term, "file:") {
		retQuery.Type = FileQuery
		retQuery.FileRegexString = trimQuote(strings.TrimPrefix(term, "file:"))
		return retQuery, nil
	}

	// 仓库名
	if strings.HasPrefix(term, "repo:") {
		retQuery.Type = RepoQuery
		retQuery.RepoRegexString = trimQuote(strings.TrimPrefix(term, "repo:"))
		return retQuery, nil
	}

	// 大小写
	if strings.HasPrefix(term, "case:") {
		if strings.TrimPrefix(term, "case:") != "yes" {
			return nil, errors.New("case: 后只能为 yes")
		}
		retQuery.Type = CaseQuery
		retQuery.Case = true
		return retQuery, nil
	}

	// 语言
	if strings.HasPrefix(term, "lang:") {
		retQuery.Type = LanguageQuery
		retQuery.Language = trimQuote(strings.TrimPrefix(term, "lang:"))
		return retQuery, nil
	}

	r, err := syntax.Parse(term, syntax.ClassNL|syntax.PerlX|syntax.UnicodeGroups)
	if err != nil {
		// 正则表达式解析失败，当做文本查询处理
		retQuery.Type = TokenQuery
		retQuery.Token = term
		return retQuery, nil
	}

	if r.Op == syntax.OpLiteral {
		// 当做文本查询处理
		retQuery.Type = TokenQuery
		retQuery.Token = term
		return retQuery, nil
	}

	subTokens := parseRegexp(r)
	retQuery.Type = RegexQuery
	retQuery.RegexString = term
	retQuery.RegexTokens = subTokens
	return retQuery, nil
}

func trimQuote(str string) string {
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}

	return str
}
