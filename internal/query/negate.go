package query

// 将树状表达式中的 NOT 条件逐级下放到 term 级别
// 比如
// (not (a or b and c)) -> (not a and (not b or not c))
func populateDownNegate(query *Query, negate bool) {
	if query == nil {
		return
	}

	// 是否需要对子表达式取反
	negateChild := (negate != query.Negate)

	if len(query.SubQueries) == 0 {
		// 没有子表达式的情况
		query.Negate = (negate != query.Negate)
	} else {
		// 对每个子表达式做处理
		for _, subQuery := range query.SubQueries {
			populateDownNegate(subQuery, negateChild)
		}

		// 需要 or <-> and 切换
		if negateChild {
			query.Or = !query.Or
		}

		// 有子表达式，则自身不 negate
		query.Negate = false
	}
}
