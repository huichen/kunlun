package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"kunlun/internal/ngram_index"
	"kunlun/pkg/types"
)

func TestSearchUtils(t *testing.T) {

	docs := []ngram_index.DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{1, 3, 4},
		},
		{
			DocumentID:     2,
			StartLocations: []uint32{11, 13, 24},
		},
	}

	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{1, 4}, {3, 6}, {4, 7}},
		},
		{
			DocumentID: 2,
			Sections:   []types.Section{{11, 14}, {13, 16}, {24, 27}},
		},
	}, DocLocationsToSections(docs, 3))
}
