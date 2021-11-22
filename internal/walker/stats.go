package walker

import "github.com/huichen/kunlun/pkg/types"

// 获得遍历器统计指标
func (walker *IndexWalker) GetStats() types.IndexWalkerStats {
	return walker.stats
}
