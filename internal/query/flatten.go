package query

// 对子查询和上级 AND/OR 一致的情况，做扁平化
// 比如
// 		(a AND (b AND c)) -> (a AND b AND c)
//		(a OR (b OR c)) -> (a OR b OR c)
// 另外，对只有一个子查询的情况，做级别提升
//      ((a)) -> (a)
func Flatten(query *Query) *Query {
	if query == nil {
		return nil
	}

	flattened := true
	dedupped := true
	for flattened || dedupped {
		dedupped = dedup(query)
		flattened = internalFlatten(query)
	}

	return query
}

func internalFlatten(query *Query) bool {
	if query == nil || len(query.SubQueries) == 0 {
		return false
	}

	parentOp := query.Or

	flattenedSubQueries := []*Query{}
	flattened := false

	if len(query.SubQueries) == 1 {
		// 只有一个子 query 的情况，提一级
		promoteSubQueryToParent(query.SubQueries[0], query)
		return true
	}

	for _, sq := range query.SubQueries {
		f := internalFlatten(sq)
		if f {
			flattened = true
		}

		if len(sq.SubQueries) > 0 && sq.Or == parentOp {
			flattenedSubQueries = append(flattenedSubQueries, sq.SubQueries...)
			flattened = true
		} else {
			flattenedSubQueries = append(flattenedSubQueries, sq)
		}
	}

	query.SubQueries = flattenedSubQueries

	return flattened
}

func promoteSubQueryToParent(child, parent *Query) {
	// XOR 操作
	parent.Negate = child.Negate != parent.Negate

	// 其他字段复制
	parent.Token = child.Token
	parent.Type = child.Type
	parent.RegexString = child.RegexString
	parent.RegexTokens = child.RegexTokens
	parent.SubQueries = child.SubQueries
	parent.Or = child.Or
	parent.IsSymbol = child.IsSymbol
	parent.FileRegexString = child.FileRegexString
	parent.RepoRegexString = child.RepoRegexString
	parent.Case = child.Case
	parent.Language = child.Language
	parent.stringRepresentation = child.stringRepresentation
}
