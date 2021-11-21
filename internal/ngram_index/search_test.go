package ngram_index

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindStartLocationWithKeyDistance(t *testing.T) {
	locs1 := []uint32{1, 2, 32, 44}
	locs2 := []uint32{4, 11, 35, 45}

	res := findStartLocationWithKeyDistance(locs1, locs2, 1)
	assert.Equal(t, []uint32{44}, res)

	res = findStartLocationWithKeyDistance(locs1, locs2, 3)
	assert.Equal(t, []uint32{1, 32}, res)

	res = findStartLocationWithKeyDistance(locs1, locs2, 9)
	assert.Equal(t, []uint32{2}, res)

	res = findStartLocationWithKeyDistance(locs1, locs2, 19)
	assert.Equal(t, []uint32(nil), res)

	// distance = 0 则做并集
	res = findStartLocationWithKeyDistance(locs1, locs2, 0)
	assert.Equal(t, []uint32{1, 2, 4, 11, 32, 35, 44, 45}, res)
}

func TestSearchDocument(t *testing.T) {
	content := "this is a document"

	index := NewNgramIndex()

	err := index.IndexDocument(1, []byte(content), nil)

	if err != nil {
		logger.Fatal(err)
	}

	key1, err := StringToIndexKey("thi")
	assert.Nil(t, err)
	key2, err := StringToIndexKey("his")
	assert.Nil(t, err)
	key3, err := StringToIndexKey("doc")
	assert.Nil(t, err)
	key4, err := StringToIndexKey("xxx")
	assert.Nil(t, err)
	key5, err := StringToIndexKey("is ")
	assert.Nil(t, err)

	docs, err := index.SearchTwoKeys(key1, key2, 1, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{0},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key1, key3, 10, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{0},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key2, key4, 1, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations(nil), docs)

	docs, err = index.SearchTwoKeys(key5, key3, 5, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{5},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key5, key3, 8, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{2},
		},
	}, docs)
}

func TestSearchOneKey(t *testing.T) {
	content := "this is a document"

	index := NewNgramIndex()

	err := index.IndexDocument(1, []byte(content), nil)

	if err != nil {
		logger.Fatal(err)
	}

	key1, err := StringToIndexKey("thi")
	assert.Nil(t, err)
	key2, err := StringToIndexKey("his")
	assert.Nil(t, err)
	key3, err := StringToIndexKey("doc")
	assert.Nil(t, err)
	key4, err := StringToIndexKey("xxx")
	assert.Nil(t, err)
	key5, err := StringToIndexKey("is ")
	assert.Nil(t, err)

	docs, err := index.SearchOneKey(key1, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{0},
		},
	}, docs)

	docs, err = index.SearchOneKey(key2, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{1},
		},
	}, docs)

	docs, err = index.SearchOneKey(key3, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{10},
		},
	}, docs)

	docs, err = index.SearchOneKey(key4, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations(nil), docs)

	docs, err = index.SearchOneKey(key5, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{2, 5},
		},
	}, docs)
}

func TestSearchShortKey(t *testing.T) {
	content := "this is a document"

	index := NewNgramIndex()

	err := index.IndexDocument(1, []byte(content), nil)

	if err != nil {
		logger.Fatal(err)
	}

	key1, err := StringToIndexKey("th")
	assert.Nil(t, err)
	key2, err := StringToIndexKey("is")
	assert.Nil(t, err)
	key3, err := StringToIndexKey("doc")
	assert.Nil(t, err)
	key4, err := StringToIndexKey("xxx")
	assert.Nil(t, err)
	key5, err := StringToIndexKey("nt")
	assert.Nil(t, err)

	docs, err := index.SearchTwoKeys(key1, key2, 2, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{0},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key1, key3, 10, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{0},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key2, key3, 5, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{5},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key2, key3, 8, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{2},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key2, key4, 1, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations(nil), docs)

	docs, err = index.SearchTwoKeys(key3, key5, 6, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{10},
		},
	}, docs)

	docs, err = index.SearchTwoKeys(key5, key3, 8, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations(nil), docs)
}

func TestSearchUnderscore(t *testing.T) {
	content := "KAFKA_LOG_PRODUCER_HELPER_SUPPLIER"

	index := NewNgramIndex()

	err := index.IndexDocument(1, []byte(content), nil)

	if err != nil {
		logger.Fatal(err)
	}

	key1, err := StringToIndexKey("a_l")
	assert.Nil(t, err)

	docs, err := index.SearchOneKey(key1, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{4},
		},
	}, docs)
}

func TestSearchTwoKeyws(t *testing.T) {
	flag.Parse()
	content := "this is a document"

	index := NewNgramIndex()

	err := index.IndexDocument(1, []byte(content), nil)

	if err != nil {
		logger.Fatal(err)
	}

	key1, err := StringToIndexKey("doc")
	assert.Nil(t, err)
	key2, err := StringToIndexKey("thi")
	assert.Nil(t, err)

	docs, err := index.SearchTwoKeys(key1, key2, 0, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, []DocumentWithLocations{
		{
			DocumentID:     1,
			StartLocations: []uint32{0, 10},
		},
	}, docs)

}
