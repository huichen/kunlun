package query

// 对 query 逐级去重，比如
// 		a AND b AND a -> a AND b
//		a AND (c OR b) AND d AND (b OR c) -> a AND d AND (b OR c)
func dedup(query *Query) bool {
	if query == nil {
		return false
	}

	// 去重
	return internalDedup(query)
}

// 返回是否做了去重操作
func internalDedup(query *Query) bool {
	if query == nil {
		return false
	}

	if query.Type != TreeQuery || len(query.SubQueries) == 0 {
		return false
	}

	done := false

	// 首先对所有子树做去重
	for _, sq := range query.SubQueries {
		d := internalDedup(sq)
		if d {
			done = true
			sq.setCompactString()
		}
	}

	// 然后检查所有子树是否重叠
	newSq := []*Query{}
	existingSq := make(map[string]bool)
	for _, sq := range query.SubQueries {
		str := sq.String()
		if _, ok := existingSq[str]; !ok {
			newSq = append(newSq, sq)
			existingSq[str] = true
		}
	}
	query.SubQueries = newSq

	return done
}
