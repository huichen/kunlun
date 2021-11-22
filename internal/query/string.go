package query

import (
	"fmt"
	"strings"
)

// 返回 query 的字符串表征
// 该字符串和 query 本身是一一对应的，也就是说
// 1、两个表征字符串相同的 query 是完全相同的
// 2、实际相同的两个 query 的表征字符串完全相等
// 这里所说的“实际相同”值得是做过归一化的相同，比如
//		(a AND b) == (b AND a)，表征字符串为 "(a AND b)"
//		(a AND (b AND c)) == (a AND b AND c)，表征字符串为 "(a AND b AND c)"
//
// 注意复杂运算下相同的两个逻辑相等的 query 可能有不同的字符串表征，比如
// 		(a AND (b OR c)) 表征字符串为 "(a AND (b OR c))"
// 		((a AND b) OR (a AND c)) 表征字符串为 "((a AND b) OR (a AND c))"
func (query *Query) String() string {
	if query == nil {
		return ""
	}

	// 如果预先有计算，则直接返回，否则重新计算
	if query.stringRepresentation != "" {
		return query.stringRepresentation
	}

	return query.compactString()
}

func (query *Query) compactString() string {
	retStr := query.fullDebugString(true, false, false)

	return retStr
}

func (query *Query) debugString() string {
	if query == nil {
		return ""
	}

	dbgString := query.fullDebugString(false, false, false)

	return dbgString
}

// 逐级设置 query 的 stringRepresentation 字段
func (query *Query) setCompactString() {
	if query == nil {
		return
	}

	query.fullDebugString(true, false, true)
}

func (query *Query) fullDebugString(compact bool, showStats bool, set bool) string {
	if query == nil {
		return ""
	}

	// 添加 NOT 前缀
	retStr := ""
	if query.Negate {
		retStr = "-"
	}

	// 添加 sym: 前缀
	if query.IsSymbol {
		retStr = retStr + "sym:"
	}

	switch query.Type {
	case TokenQuery:
		if query.Token != "" {
			token := escapeString(query.Token)
			retStr = retStr + token
		}
	case RegexQuery:
		if !compact {
			retStr = fmt.Sprintf("%s<%s>", retStr, escapeString(query.RegexString))
			retStr = fmt.Sprintf("%s(%s)", retStr, strings.Join(query.RegexTokens, ","))
		} else {
			retStr = fmt.Sprintf("%s%s", retStr, escapeString(query.RegexString))
		}
	case TreeQuery:
		var subStrs []string
		for _, q := range query.SubQueries {
			qStr := q.fullDebugString(compact, showStats, set)
			if qStr != "" {
				subStrs = append(subStrs, qStr)
			}
		}

		if len(subStrs) == 1 {
			retStr += subStrs[0]
		} else if len(subStrs) > 1 {
			op := " AND "
			if query.Or {
				op = " OR "
			}
			if compact && !query.Or {
				retStr = fmt.Sprintf("%s(%s)", retStr, strings.Join(subStrs, op))
			} else {
				retStr = fmt.Sprintf("%s(%s)", retStr, strings.Join(subStrs, op))
			}
		}
	case FileQuery:
		retStr = fmt.Sprintf("%sfile:%s", retStr, escapeString(query.FileRegexString))
	case RepoQuery:
		retStr = fmt.Sprintf("%srepo:%s", retStr, escapeString(query.RepoRegexString))
	case CaseQuery:
		retStr = fmt.Sprintf("%scase:yes", retStr)
	case LanguageQuery:
		retStr = fmt.Sprintf("%slang:%s", retStr, escapeString(query.Language))
	}

	if showStats {
		retStr = fmt.Sprintf("%s@[%d,%d,%d]", retStr, query.ID, query.NumNodes, query.MaxDepth)
	}

	if set {
		query.stringRepresentation = retStr
	}

	return retStr
}

func escapeString(token string) string {
	if strings.Contains(token, " ") || strings.Contains(token, "(") || strings.Contains(token, ")") {
		token = fmt.Sprintf("\"%s\"", token)
	}

	return token
}
