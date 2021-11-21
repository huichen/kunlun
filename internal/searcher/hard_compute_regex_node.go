package searcher

import (
	"github.com/huichen/kunlun/internal/query"
)

func (schr *Searcher) hardComputeOneRegexNode(context *Context, q *query.Query) error {
	_, err := schr.internalHardComputeOneRegexNode(context, q)
	if err != nil {
		return err
	}
	return nil
}

func (schr *Searcher) internalHardComputeOneRegexNode(context *Context, q *query.Query) (bool, error) {
	if q == nil {
		return false, nil
	}

	// 跳过已经计算过的节点
	r := context.query.QueryResults[q.ID]
	if r != nil {
		return false, nil
	}

	// 如果是正则表达式，计算
	if q.Type == query.RegexQuery {
		err := schr.computeRegexNode(context, q, nil)
		if err != nil {
			return false, err
		}

		if err := context.checkTimeout(); err != nil {
			return false, err
		}
		return true, nil
	}

	// 否则对树子节点遍历
	if q.Type == query.TreeQuery {
		for _, sq := range q.SubQueries {
			done, err := schr.internalHardComputeOneRegexNode(context, sq)
			if err != nil {
				return false, err
			}
			if done {
				// 只要有一个子节点完成计算，退出
				return true, nil
			}
		}
	}

	return false, nil
}
