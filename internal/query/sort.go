package query

import (
	"sort"
	"strings"
)

var (
	QueryTypeOrder = map[QueryType]int{
		CaseQuery:     1,
		LanguageQuery: 2,
		RepoQuery:     3,
		FileQuery:     4,
		TokenQuery:    6,
		RegexQuery:    7,
		TreeQuery:     8,
	}
)

func Sort(q *Query) *Query {
	if q == nil {
		return nil
	}

	for {
		oldString := q.String()
		internalSort(q)
		q.setCompactString()
		newString := q.String()
		if oldString == newString {
			// 循环直到排序结果稳定为止
			break
		}
	}

	// 顺序发生变化后需要更新 ID
	UpdateStats(q)

	return q
}

func internalSort(q *Query) {
	if q == nil || len(q.SubQueries) == 0 {
		return
	}

	sort.SliceStable(q.SubQueries, func(i, j int) bool {
		o1 := QueryTypeOrder[q.SubQueries[i].Type]
		o2 := QueryTypeOrder[q.SubQueries[j].Type]
		if o1 == o2 {
			q1 := q.SubQueries[i]
			q2 := q.SubQueries[j]

			// 类型相同的情况下，先把 negate 放后面
			if q1.Negate != q2.Negate {
				return !q1.Negate
			}

			// 再把 sym: 放前面
			if q1.IsSymbol != q2.IsSymbol {
				return q1.IsSymbol
			}

			// 然后比较内容按照字符串字母顺序排列
			return strings.Compare(q.SubQueries[i].String(), q.SubQueries[j].String()) < 0
		}
		return o1 < o2
	})

	for _, sq := range q.SubQueries {
		internalSort(sq)
	}
}
