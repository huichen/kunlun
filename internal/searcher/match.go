package searcher

import (
	"regexp"

	"kunlun/internal/query"
)

// 在 q 的树状结构中使用正则表达式匹配 pattern
// 树叶子节点的正则表达式从 res 中读取，res 的 index 和数的 ID 匹配
// 如果 res 中没有找到响应的正则表达式，则返回 false
func matchRegexpQueries(pattern string, q *query.Query, res *[]*regexp.Regexp, names []string) bool {
	if q == nil {
		return false
	}

	negate := q.Negate

	if q.Type != query.TreeQuery {
		// 先检查是否完全匹配
		if names != nil && names[q.ID] == pattern {
			return negate != true
		}

		r := (*res)[q.ID]
		if r == nil {
			// 正则匹配失败的情况
			return negate != false
		}

		return negate != r.MatchString(pattern)
	}

	if len(q.SubQueries) == 0 {
		return false
	}

	subMatch := true
	if q.Or {
		subMatch = false
	}
	for _, sq := range q.SubQueries {
		m := matchRegexpQueries(pattern, sq, res, names)
		if q.Or {
			subMatch = subMatch || m
		} else {
			subMatch = subMatch && m
		}
	}
	return negate != subMatch
}
