package searcher

import (
	"fmt"

	"kunlun/internal/query"
)

// 获得所有 queryType 类型，汇总在返回结果
func collectModifier(q *query.Query, queryType query.QueryType) (*query.Query, error) {
	if q == nil {
		return nil, nil
	}

	retQuery := internalCollectModifier(q, queryType)
	retQuery = query.Copy(retQuery)

	query.Flatten(retQuery)
	query.UpdateStats(retQuery)

	// modifier 校验
	err := validateModifiers(retQuery)
	if err != nil {
		return nil, err
	}

	if retQuery != nil {
		// 对修饰词使用的合法性做校验
		if retQuery.MaxDepth > 1 {
			return nil, fmt.Errorf("%s: 修饰词不能同时在多个级别使用", query.GetModifierName(queryType))
		}
	}

	return retQuery, nil
}

func internalCollectModifier(q *query.Query, queryType query.QueryType) *query.Query {
	if q == nil {
		return nil
	}

	if q.Type == queryType {
		return query.Copy(q)
	}

	retQuery := &query.Query{
		Type:       query.TreeQuery,
		SubQueries: []*query.Query{},
		Or:         q.Or,
	}

	for _, sq := range q.SubQueries {
		newSq := internalCollectModifier(sq, queryType)
		if newSq != nil {
			retQuery.SubQueries = append(retQuery.SubQueries, newSq)
		}
	}

	if len(retQuery.SubQueries) == 0 {
		return nil
	}

	return retQuery
}

func getCase(q *query.Query) bool {
	if q == nil {
		return false
	}

	if q.Type == query.CaseQuery {
		return true
	}

	ca := false
	for _, sq := range q.SubQueries {
		if getCase(sq) {
			ca = true
		}
	}

	return ca
}
