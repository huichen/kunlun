package query

// 深度拷贝所有子节点
func Copy(query *Query) *Query {
	if query == nil {
		return nil
	}

	retQuery := &Query{}
	CopyTo(query, retQuery)

	if query.SubQueries == nil {
		retQuery.SubQueries = nil
	} else {
		retQuery.SubQueries = make([]*Query, len(query.SubQueries))
		for id, sq := range query.SubQueries {
			retQuery.SubQueries[id] = Copy(sq)
		}
	}

	return retQuery
}

func CopyTo(src *Query, dst *Query) {
	if src == nil || dst == nil {
		return
	}

	dst.Type = src.Type
	dst.ID = src.ID
	dst.Token = src.Token
	dst.Negate = src.Negate
	dst.Or = src.Or

	dst.RegexString = src.RegexString
	dst.IsSymbol = src.IsSymbol
	dst.FileRegexString = src.FileRegexString
	dst.RepoRegexString = src.RepoRegexString
	dst.Case = src.Case
	dst.Language = src.Language
	dst.NumNodes = src.NumNodes
	dst.MaxDepth = src.MaxDepth
	dst.RootDistance = src.RootDistance
	dst.stringRepresentation = src.stringRepresentation

	tokens := make([]string, len(src.RegexTokens))
	copy(tokens, src.RegexTokens)
	dst.RegexTokens = tokens

	if src.SubQueries == nil {
		dst.SubQueries = nil
	} else {
		dst.SubQueries = make([]*Query, len(src.SubQueries))
		for id, sq := range src.SubQueries {
			dst.SubQueries[id] = sq
		}
	}
}
