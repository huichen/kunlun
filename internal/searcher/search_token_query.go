package searcher

import "github.com/huichen/kunlun/internal/query"

func (schr *Searcher) searchTokenQuery(context *Context, q *query.Query) error {
	if q == nil {
		return nil
	}
	return schr.internalSearchTokenQuery(context, q)
}

func (schr *Searcher) internalSearchTokenQuery(context *Context, q *query.Query) error {
	// 对 token 类型，直接计算
	if q.Type == query.TokenQuery {
		request := context.getSearchTokenRequest(q.Token, q.IsSymbol)
		resp, err := context.idxr.SearchToken(request)
		if err != nil {
			return err
		}
		if err := context.checkTimeout(); err != nil {
			return err
		}

		context.query.QueryResults[q.ID] = &resp
		return nil
	}

	// 对树类型，递归计算
	if q.Type == query.TreeQuery && len(q.SubQueries) > 0 {
		for _, sq := range q.SubQueries {
			err := schr.internalSearchTokenQuery(context, sq)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
