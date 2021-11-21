package query

// 从查询表达式树中获得所有的 tokens，包括了正则表达式中的 tokens
// 返回结果做了去重
func (query *Query) GetTokens() []string {
	tokens := internalGetTokens(query)

	finalTokens := []string{}
	tokenMap := make(map[string]bool)

	for _, t := range tokens {
		if _, ok := tokenMap[t]; !ok {
			finalTokens = append(finalTokens, t)
			tokenMap[t] = true
		}
	}
	return finalTokens
}

func internalGetTokens(query *Query) []string {
	var tokens []string

	if query.Type == RegexQuery {
		tokens = append(tokens, query.RegexTokens...)
		return tokens
	}

	// 包含子
	if len(query.SubQueries) > 0 {
		for _, sq := range query.SubQueries {
			sqTokens := internalGetTokens(sq)
			if len(sqTokens) > 0 {
				tokens = append(tokens, sqTokens...)
			}
		}
		return tokens
	}

	if query.Token != "" {
		return []string{query.Token}
	}

	return nil
}
