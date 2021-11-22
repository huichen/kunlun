package indexer

import (
	"github.com/huichen/kunlun/internal/common_types"
	"github.com/huichen/kunlun/internal/ngram_index"
	"github.com/huichen/kunlun/pkg/types"
)

// 将文档中的起始位置数组转化成分段区间数组，分段长度等于 length
func docLocationsToSections(docs []ngram_index.DocumentWithLocations, length uint32) []common_types.DocumentWithSections {
	ret := make([]common_types.DocumentWithSections, len(docs))
	for idDoc, doc := range docs {
		sections := make([]types.Section, len(doc.StartLocations))
		for idLocations, loc := range doc.StartLocations {
			sections[idLocations] = types.Section{
				Start: loc,
				End:   loc + length,
			}
		}
		ret[idDoc] = common_types.DocumentWithSections{
			DocumentID: doc.DocumentID,
			Sections:   sections,
		}
	}

	return ret
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}

	return b
}

func min(a, b uint32) uint32 {
	if a > b {
		return b
	}

	return a
}
