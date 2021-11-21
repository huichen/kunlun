package walker

import "kunlun/pkg/types"

func (walker *IndexWalker) GetStats() types.IndexWalkerStats {
	return walker.stats
}
