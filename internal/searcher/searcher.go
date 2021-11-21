package searcher

import "github.com/huichen/kunlun/pkg/types"

type Searcher struct {
	options *types.SearcherOptions
}

func NewSearcher(options *types.SearcherOptions) *Searcher {
	if options == nil {
		options = types.NewSearcherOptions()
	}

	s := &Searcher{
		options: options,
	}

	return s
}
