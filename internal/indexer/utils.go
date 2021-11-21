package indexer

import (
	"reflect"
	"sort"

	"github.com/huichen/kunlun/internal/ngram_index"
	"github.com/huichen/kunlun/pkg/types"
)

func DocLocationsToSections(docs []ngram_index.DocumentWithLocations, length uint32) []types.DocumentWithSections {
	ret := make([]types.DocumentWithSections, len(docs))
	for idDoc, doc := range docs {
		sections := make([]types.Section, len(doc.StartLocations))
		for idLocations, loc := range doc.StartLocations {
			sections[idLocations] = types.Section{
				Start: loc,
				End:   loc + length,
			}
		}
		ret[idDoc] = types.DocumentWithSections{
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

func SortAndDedup(slicePtr interface{}, less func(i, j int) bool) {
	v := reflect.ValueOf(slicePtr).Elem()
	if v.Len() <= 1 {
		return
	}
	sort.Slice(v.Interface(), less)

	i := 0
	for j := 1; j < v.Len(); j++ {
		if !less(i, j) {
			continue
		}
		i++
		v.Index(i).Set(v.Index(j))
	}
	i++
	v.SetLen(i)
}
