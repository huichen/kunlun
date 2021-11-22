package searcher

import (
	"github.com/huichen/kunlun/internal/query"
)

// “硬”计算第一个找到的（逐级遍历）一个正则表达式节点，计算时不使用表达式上下文信息
// 与之对应的，“软”计算会尝试使用表达式的上下文做联合优化
// 关于“软”计算的方式请见 soft_compute_regex_node.go 中注释
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
