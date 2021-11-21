package walker

import "github.com/huichen/kunlun/pkg/types"

func (walker *IndexWalker) GetStats() types.IndexWalkerStats {
	return walker.stats
}
