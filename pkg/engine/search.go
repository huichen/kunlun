package engine

import (
	"errors"

	"github.com/huichen/kunlun/pkg/types"
)

func (engine *KunlunEngine) Search(request types.SearchRequest) (*types.SearchResponse, error) {
	if !engine.finished {
		return nil, errors.New("索引即将构建完成，请稍后再来搜索 ~")

	}
	return engine.searcher.Search(engine.indexer, request)
}
