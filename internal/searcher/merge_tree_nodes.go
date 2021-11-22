package searcher

import (
	"errors"

	"github.com/huichen/kunlun/internal/common_types"
	"github.com/huichen/kunlun/internal/query"
)

// 合并树节点
// 将所有同级别的已经检索过的节点做合并，比如
// 		a and (b or c) -> 合并计算 b or c 并作为单节点提升到上一级
//		a and b and (...) -> 合并计算 a and b
func (schr *Searcher) mergeTreeNodes(context *Context, q *query.Query) error {
	if q == nil {
		return nil
	}

	err := internalMergeTreeNodes(context, q)
	if err != nil {
		return err
	}

	return context.checkTimeout()
}

func internalMergeTreeNodes(context *Context, q *query.Query) error {
	if q == nil {
		return nil
	}

	// 遇到叶子节点，返回
	if q.Type != query.TreeQuery {
		return nil
	}

	r := context.query.QueryResults[q.ID]
	if r != nil {
		// 该树已经计算过，返回
		return nil
	}

	// 先对子树做归并
	for _, sq := range q.SubQueries {
		err := internalMergeTreeNodes(context, sq)
		if err != nil {
			return err
		}
		if err := context.checkTimeout(); err != nil {
			return err
		}
	}

	// 得到所有已经计算过的节点
	computedNodes := []*query.Query{}
	numNegateNodes := 0
	for _, sq := range q.SubQueries {
		if sq == nil {
			return errors.New("子 query 不允许为 nil")
		}
		r := context.query.QueryResults[sq.ID]
		if r != nil {
			computedNodes = append(computedNodes, sq)
			if sq.Negate {
				// 非法情况：OR 树 + negate 节点
				if q.Or {
					return errors.New("OR 运算符不能和 - 运算符并存")
				}

				numNegateNodes++
			}
		}
	}

	// 没有计算过的节点或者只有一个计算过的节点，无需归并，返回
	if len(computedNodes) <= 1 {
		return nil
	}

	// 合并计算过的节点
	newQIsNegate := false
	var newQ []common_types.DocumentWithSections
	var err error
	if numNegateNodes == len(computedNodes) {
		// 所有计算过的节点都是 negate 节点，这种情况只可能为 AND 树，我们做变换
		// 		-a AND -b -> -(a OR b)
		for _, n := range computedNodes {
			n.Negate = false
		}
		newQ, err = mergeQueries(context, computedNodes, true)
		newQIsNegate = true
	} else {
		newQ, err = mergeQueries(context, computedNodes, q.Or)
	}
	if err != nil {
		return err
	}

	// 树内所有节点都已经计算了，用归并结果替代树节点
	if len(computedNodes) == len(q.SubQueries) {
		context.query.QueryResults[q.ID] = &newQ
		q.Negate = (q.Negate != newQIsNegate)
		return nil
	}

	// 特殊情况：AND 树 + 空节点
	if !q.Or && len(newQ) == 0 {
		context.query.QueryResults[q.ID] = &newQ
		return nil
	}

	// 否则只保留第一个已经计算的节点并做替换，删除其他已经计算过的节点
	newSq := []*query.Query{}
	firstComputedNodeReplaced := false
	for _, sq := range q.SubQueries {
		r := context.query.QueryResults[sq.ID]
		if r != nil {
			if len(newQ) == 0 && q.Or {
				// 对 OR 树且归并计算为空的情况，不做添加
				continue
			}

			if !firstComputedNodeReplaced {
				// 第一个计算节点
				context.query.QueryResults[sq.ID] = &newQ
				sq.Negate = newQIsNegate
				newSq = append(newSq, sq)
				firstComputedNodeReplaced = true
			}
			// 其他计算节点忽略
		} else {
			newSq = append(newSq, sq)
		}
	}

	// 如果新的子节点只有一个，提升该子节点
	if len(q.SubQueries) == 1 {
		q.SubQueries = q.SubQueries[0].SubQueries
		q.Type = q.SubQueries[0].Type
		q.Negate = (q.SubQueries[0].Negate != q.Negate)
		context.query.QueryResults[q.ID] = context.query.QueryResults[q.SubQueries[0].ID]
		return nil
	}

	// 否则替换新的子节点
	q.SubQueries = newSq

	return nil
}
