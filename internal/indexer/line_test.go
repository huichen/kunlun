package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"kunlun/pkg/types"
)

func TestLine(t *testing.T) {
	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)
	idxr := NewIndexer(options)

	idxr.IndexFile([]byte("aaaa\nbbbb\ncccc\ndddd\neeee\nffff\ngggg"), types.IndexFileInfo{Path: "this.doc.txt"})
	idxr.Finish()

	docs, _, err := idxr.GetLinesFromSections(1, []types.Section{
		{0, 1}, {3, 4}, {6, 7}, {11, 17},
	}, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, []types.Line{
		{
			LineNumber: 0,
			Highlights: []types.Section{{0, 1}, {3, 4}},
		},
		{
			LineNumber: 1,
			Highlights: []types.Section{{1, 2}},
		},
		{
			LineNumber: 2,
			Highlights: []types.Section{{1, 5}},
		},
		{
			LineNumber: 3,
			Highlights: []types.Section{{0, 2}},
		},
	}, docs)

	_, _, err = idxr.GetLinesFromSections(1, []types.Section{
		{0, 100},
	}, 0, 0)
	logger.Error(err)
	assert.NotNil(t, err)

	_, _, err = idxr.GetLinesFromSections(1, []types.Section{
		{10, 1},
	}, 0, 0)
	logger.Error(err)
	assert.NotNil(t, err)

	docs, _, err = idxr.GetLinesFromSections(1, []types.Section{
		{3, 4}, {11, 17},
	}, 1, 0)
	assert.Nil(t, err)
	assert.Equal(t, []types.Line{
		{
			LineNumber: 0,
			Highlights: []types.Section{{3, 4}},
		},
		{
			LineNumber: 1,
		},
		{
			LineNumber: 2,
			Highlights: []types.Section{{1, 5}},
		},
		{
			LineNumber: 3,
			Highlights: []types.Section{{0, 2}},
		},
		{
			LineNumber: 4,
			Highlights: nil,
		},
	}, docs)

	docs, _, err = idxr.GetLinesFromSections(1, []types.Section{
		{3, 4}, {11, 17},
	}, 1, 0)
	assert.Nil(t, err)
	assert.Equal(t, []types.Line{
		{
			LineNumber: 0,
			Highlights: []types.Section{{3, 4}},
		},
		{
			LineNumber: 1,
		},
		{
			LineNumber: 2,
			Highlights: []types.Section{{1, 5}},
		},
		{
			LineNumber: 3,
			Highlights: []types.Section{{0, 2}},
		},
		{
			LineNumber: 4,
			Highlights: nil,
		},
	}, docs)

}
