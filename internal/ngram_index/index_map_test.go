package ngram_index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexMap(t *testing.T) {
	im := IndexMap{}
	im.insert(2, KeyedDocument{
		DocumentID:           2,
		SortedStartLocations: []uint32{3, 4},
	})

	idx, found := im.find(1)
	assert.Equal(t, idx, 0)
	assert.False(t, found)

	idx, found = im.find(2)
	assert.Equal(t, idx, 0)
	assert.True(t, found)

	idx, found = im.find(3)
	assert.Equal(t, idx, 1)
	assert.False(t, found)

	im.insert(2, KeyedDocument{
		DocumentID:           1,
		SortedStartLocations: []uint32{3, 4},
	})
	logger.Info(im.index[0].documents)

	im.insert(2, KeyedDocument{
		DocumentID:           3,
		SortedStartLocations: []uint32{3, 4},
	})
	logger.Info(im.index[0].documents)

	im.insert(1, KeyedDocument{
		DocumentID:           3,
		SortedStartLocations: []uint32{3, 4},
	})
	logger.Info(im.index[0].documents)
	logger.Info(im.index[1].documents)

	im.insert(3, KeyedDocument{
		DocumentID:           3,
		SortedStartLocations: []uint32{3, 4},
	})
	logger.Info(im.index[0].documents)
	logger.Info(im.index[1].documents)
	logger.Info(im.index[2].documents)
}

func TestSortedDocuments(t *testing.T) {
	documents := &SortedKeyedDocuments{}

	documents.insert(KeyedDocument{
		DocumentID:           2,
		SortedStartLocations: []uint32{3, 4},
	})
	logger.Info(documents)

	idx, found := documents.find(1)
	assert.Equal(t, idx, 0)
	assert.False(t, found)

	idx, found = documents.find(2)
	assert.Equal(t, idx, 0)
	assert.True(t, found)

	idx, found = documents.find(3)
	assert.Equal(t, idx, 1)
	assert.False(t, found)

	documents.insert(KeyedDocument{
		DocumentID:           1,
		SortedStartLocations: []uint32{3, 4},
	})
	logger.Info(documents)

	documents.insert(KeyedDocument{
		DocumentID:           3,
		SortedStartLocations: []uint32{3, 4},
	})
	logger.Info(documents)

}
