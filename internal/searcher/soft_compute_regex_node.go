package searcher

import (
	"errors"

	"github.com/huichen/kunlun/internal/query"
)

// 试图对树中的一个正则表达式做软计算优化
//
// 所谓软计算，指的是通过已经计算好的中间结果 a 对形如下面的表达式做搜索匹配
// 		a AND regex
// 其中 regex 是某个待搜索的正则表达式
//
// 软计算的好处是可以通过 a 中的 docIDs 极大减少做正则匹配的文档范围，从而达到提升计算速度的目的
// 另一个可以优化的软计算是如下形式：
//		a AND (b OR regex)
// 我们通过如下变换来做优化
//		a AND (b OR regex) => a AND (b OR (a AND regex))
//
// 函数试图在树中找到任何一个这样的优化（只找到一个即可），计算存放在 context 中并返回 true
// 没找到的话返回 false
func (schr *Searcher) softComputeOneRegexNode(context *Context, q *query.Query) (bool, error) {
	return schr.internalSoftComputeOneRegexNode(context, q, nil)
}

func (schr *Searcher) internalSoftComputeOneRegexNode(context *Context, q *query.Query, parentComputedQuery *query.Query) (bool, error) {
	if q == nil {
		return false, nil
	}

	// 略过非树节点
	if q.Type != query.TreeQuery {
		return false, nil
	}

	// 如果该树已经计算过，返回
	r := context.query.QueryResults[q.ID]
	if r != nil {
		return false, nil
	}

	// 合法性校验，negate 树是不允许的
	if q.Negate {
		return false, errors.New("树不能为 negate，没有归一化？")
	}

	// 收集已经计算的中间结果，和是否有正则节点信息
	var computedTreeNode *query.Query
	hasRegexNode := false
	for _, sq := range q.SubQueries {
		r = context.query.QueryResults[sq.ID]
		if r != nil {
			if computedTreeNode != nil {
				return false, errors.New("树下包含两个已计算节点，先调用 mergeTreeNodes 函数？")
			}
			computedTreeNode = sq
		}
		if sq.Type == query.RegexQuery {
			hasRegexNode = true
		}
	}

	// 没有正则表达式，则进入子节点
	if !hasRegexNode {
		return schr.softComputeSubTrees(context, q, computedTreeNode)
	}

	// 如果是 OR 树，当有父输入时计算正则，否则进入子节点
	if q.Or {
		if parentComputedQuery != nil {
			return schr.softComputeFirstSubRegexNode(context, q, parentComputedQuery)
		}
		return schr.softComputeSubTrees(context, q, computedTreeNode)
	}

	// 如果没有中间结果，进入子节点
	if computedTreeNode == nil {
		return schr.softComputeSubTrees(context, q, nil)
	}

	// 最后一种情况：有 computedTreeNode 且为 AND 树，则对第一个正则表达式做计算
	_, err := schr.softComputeFirstSubRegexNode(context, q, computedTreeNode)
	if err != nil {
		return false, err
	}

	return true, nil
}

// 对 q 的子树做处理
func (schr *Searcher) softComputeSubTrees(context *Context, q *query.Query, computedTreeNode *query.Query) (bool, error) {
	for _, sq := range q.SubQueries {
		if sq.Type != query.TreeQuery {
			continue
		}

		var err error
		var done bool
		if q.Or {
			// OR 树的情况，直接计算子节点
			done, err = schr.internalSoftComputeOneRegexNode(context, sq, nil)
		} else {
			// AND 树的情况，利用已经计算的中间结果做简化
			done, err = schr.internalSoftComputeOneRegexNode(context, sq, computedTreeNode)
		}
		if err != nil {
			return false, err
		}

		if done {
			// 有任何简化的话不再继续
			return true, nil
		}
	}

	return false, nil
}

// 对 q 的第一个 regex 子节点，利用 computedTreeNode 做检索
func (schr *Searcher) softComputeFirstSubRegexNode(context *Context, q *query.Query, computedTreeNode *query.Query) (bool, error) {
	done := false
	for _, sq := range q.SubQueries {
		if sq.Type != query.RegexQuery {
			continue
		}

		err := schr.computeRegexNode(context, sq, computedTreeNode)
		if err != nil {
			return false, err
		}
		if err := context.checkTimeout(); err != nil {
			return false, err
		}
		done = true
		break
	}

	if !done {
		return false, errors.New("没有找到 regex 节点")
	}

	return true, nil
}

func (schr *Searcher) computeRegexNode(context *Context, q *query.Query, computedTreeNode *query.Query) error {
	request, err := context.getSearchRegexRequest(q, computedTreeNode, q.IsSymbol)
	if err != nil {
		return err
	}

	resp, err := context.idxr.SearchRegex(*request)
	if err != nil {
		return err
	}

	// 计数器更新
	context.regexSearchTimes += resp.RegexSearchTimes

	// 归并结果
	q.Negate = resp.Negate
	context.query.QueryResults[q.ID] = &resp.Documents

	return nil
}
