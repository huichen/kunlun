package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"kunlun/pkg/types"
)

func TestSearch(t *testing.T) {
	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)
	idxr := NewIndexer(options)

	idxr.IndexFile([]byte("this is a document"), types.IndexFileInfo{Path: "this.doc.txt"})
	idxr.Finish()

	docs, err := idxr.internalSearchToken(SearchTokenRequest{
		Token:         "this",
		CaseSensitive: false,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 4}},
		},
	}, docs)

	docs, err = idxr.internalSearchToken(SearchTokenRequest{
		Token:         "docu",
		CaseSensitive: false,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{10, 14}},
		},
	}, docs)

	docs, err = idxr.internalSearchToken(SearchTokenRequest{
		Token:         "doc",
		CaseSensitive: false,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{10, 13}},
		},
	}, docs)
}

func TestSearchCase(t *testing.T) {
	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)

	idxr := NewIndexer(options)

	idxr.IndexFile([]byte("KAFKA_LoG_PRODUCER_HELPER"), types.IndexFileInfo{Path: "this.doc.txt"})
	idxr.Finish()

	docs, err := idxr.internalSearchToken(SearchTokenRequest{
		Token:         "A_LoG",
		CaseSensitive: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{4, 9}},
		},
	}, docs)

	docs, err = idxr.internalSearchToken(SearchTokenRequest{
		Token:         "a_log",
		CaseSensitive: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{}, docs)

}

func TestSearchMore(t *testing.T) {
	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)
	idxr := NewIndexer(options)

	idxr.IndexFile([]byte("KAFKA_LOG_PRODUCER_HELPER"), types.IndexFileInfo{Path: "this.doc.txt"})
	idxr.Finish()

	docs, err := idxr.internalSearchToken(SearchTokenRequest{
		Token:         "KAFKA_LOG_PRODUCER_HELPER",
		CaseSensitive: false,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 25}},
		},
	}, docs)
}

func TestSearchToken(t *testing.T) {
	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)
	idxr := NewIndexer(options)

	idxr.IndexFile([]byte("KAFKA_LOG_PRODUCER_HELPER"), types.IndexFileInfo{Path: "this.doc.txt"})
	idxr.Finish()

	docs, err := idxr.SearchToken(SearchTokenRequest{
		Token:         "KAFKA_LOG_PRODUCER_HELPER",
		CaseSensitive: false,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 25}},
		},
	}, docs)
}

func TestSearchSymbol(t *testing.T) {
	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)
	idxr := NewIndexer(options)

	idxr.IndexFile([]byte("a1 a2 a3\nb1 b2 b3\nc1 c2 c3"),
		types.IndexFileInfo{
			Path: "this.doc.txt",
			CTagsEntries: []*types.CTagsEntry{
				{Sym: "a2", Line: 1},
				{Sym: "b1", Line: 2},
				{Sym: "c3", Line: 3},
			},
		})
	idxr.Finish()

	docs, err := idxr.SearchToken(SearchTokenRequest{
		Token:         "a3",
		CaseSensitive: false,
		IsSymbol:      true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{}, docs)

	docs, err = idxr.SearchToken(SearchTokenRequest{
		Token:         "b2",
		CaseSensitive: false,
		IsSymbol:      true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{}, docs)

	docs, err = idxr.SearchToken(SearchTokenRequest{
		Token:         "b2",
		CaseSensitive: false,
		IsSymbol:      false,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{DocumentID: 0x1, Sections: []types.Section{{Start: 12, End: 14}}},
	}, docs)

	docs, err = idxr.SearchToken(SearchTokenRequest{
		Token:         "b1",
		CaseSensitive: false,
		IsSymbol:      true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{DocumentID: 0x1, Sections: []types.Section{{Start: 9, End: 11}}},
	}, docs)

	docs, err = idxr.SearchToken(SearchTokenRequest{
		Token:         "a2",
		CaseSensitive: false,
		IsSymbol:      true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{DocumentID: 0x1, Sections: []types.Section{{Start: 3, End: 5}}},
	}, docs)

	docs, err = idxr.SearchToken(SearchTokenRequest{
		Token:         "c3",
		CaseSensitive: false,
		IsSymbol:      true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{DocumentID: 0x1, Sections: []types.Section{{Start: 24, End: 26}}},
	}, docs)

}
