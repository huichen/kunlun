package query

// 深度优先给 query 逐级更新 ID 和统计信息等
func UpdateStats(query *Query) {
	if query == nil {
		return
	}

	internalUpdateStats(query, 0, 0)

	// 充填 DebugString 信息
	query.setCompactString()
}

// 返回 (下一个 ID，子节点数目，子节点最大深度)
func internalUpdateStats(query *Query, id int, rootDistance int) (int, int, int) {
	numChildren := 0
	maxDepth := 0
	if query.Type == TreeQuery {
		for _, sq := range query.SubQueries {
			if sq == nil {
				continue
			}
			var children, depth int
			id, children, depth = internalUpdateStats(sq, id, rootDistance+1)
			if maxDepth < depth+1 {
				maxDepth = depth + 1
			}
			numChildren += children
		}
	}
	query.NumNodes = numChildren + 1
	query.MaxDepth = maxDepth
	query.ID = id
	query.RootDistance = rootDistance
	id++
	return id, query.NumNodes, query.MaxDepth
}
