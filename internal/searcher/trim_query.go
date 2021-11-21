package searcher

import "github.com/huichen/kunlun/internal/query"

// 对 q 做精简，只保留 token/regex/tree 三种类型的节点
func TrimQuery(q *query.Query) *query.Query {
	q = internalTrimQuery(q)

	// 扁平化之后重排序
	q = query.Flatten(q)
	q = query.Sort(q)

	return q
}

func internalTrimQuery(q *query.Query) *query.Query {
	if q == nil {
		return nil
	}

	if q.Type != query.TokenQuery &&
		q.Type != query.RegexQuery &&
		q.Type != query.TreeQuery {
		return nil
	}

	if q.Type == query.TreeQuery {
		sqs := []*query.Query{}
		for _, sq := range q.SubQueries {
			newSq := internalTrimQuery(sq)
			if newSq != nil {
				sqs = append(sqs, newSq)
			}
		}
		if len(sqs) != 0 {
			q.SubQueries = sqs
		} else {
			return nil
		}
	}

	return q
}
