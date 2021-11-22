package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/huichen/kunlun/internal/common_types"
	"github.com/huichen/kunlun/pkg/types"
)

func TestSearchRegex(t *testing.T) {
	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)
	idxr := NewIndexer(options)

	idxr.IndexFile([]byte("this is a document"), common_types.IndexFileInfo{Path: "this0.doc.txt"})
	idxr.IndexFile([]byte("this is a document"), common_types.IndexFileInfo{Path: "this1.doc.txt"})
	idxr.IndexFile([]byte("thids is a document"), common_types.IndexFileInfo{Path: "this2.doc.txt"})
	idxr.IndexFile([]byte("this is a document"), common_types.IndexFileInfo{Path: "this3.doc.txt"})
	idxr.IndexFile([]byte("this is a document"), common_types.IndexFileInfo{Path: "this4.doc.txt"})
	idxr.IndexFile([]byte("docs is good"), common_types.IndexFileInfo{Path: "this5.doc.txt"})
	idxr.Finish()

	docs, _, err := idxr.internalSearchRegex(SearchRegexRequest{
		Regex:  "this.*doc",
		Tokens: []string{"this", "doc"},
	})
	assert.Nil(t, err)
	assert.Equal(t, []common_types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 13}},
		},
		{
			DocumentID: 2,
			Sections:   []types.Section{{0, 13}},
		},
		{
			DocumentID: 4,
			Sections:   []types.Section{{0, 13}},
		},
		{
			DocumentID: 5,
			Sections:   []types.Section{{0, 13}},
		},
	}, docs)

	// -a b
	resp, err := idxr.SearchRegex(SearchRegexRequest{
		Regex:               "thisddd.*doc",
		Tokens:              []string{"this", "doc"},
		CandidateDocs:       &[]uint64{1},
		CandidateDocsNegate: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []common_types.DocumentWithSections{}, resp.Documents)

	// -a b
	resp, err = idxr.SearchRegex(SearchRegexRequest{
		Regex:               "this.*doc",
		Tokens:              []string{"this", "doc"},
		CandidateDocs:       &[]uint64{1},
		CandidateDocsNegate: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, []common_types.DocumentWithSections{
		{
			DocumentID: 2,
			Sections:   []types.Section{{0, 13}},
		},
		{
			DocumentID: 4,
			Sections:   []types.Section{{0, 13}},
		},
		{
			DocumentID: 5,
			Sections:   []types.Section{{0, 13}},
		},
	}, resp.Documents)

	// a - b
	resp, err = idxr.SearchRegex(SearchRegexRequest{
		Regex:         "this.*doc",
		Negate:        true,
		Tokens:        []string{"this", "doc"},
		CandidateDocs: &[]uint64{1, 2, 3},
	})
	assert.Nil(t, err)
	assert.Equal(t, []common_types.DocumentWithSections{
		{
			DocumentID: 3,
			Sections:   []types.Section(nil),
		},
	}, resp.Documents)

	// -a -b
	_, err = idxr.SearchRegex(SearchRegexRequest{
		Regex:               "this.*doc",
		Negate:              true,
		Tokens:              []string{"this", "doc"},
		CandidateDocs:       &[]uint64{1, 2, 3},
		CandidateDocsNegate: true,
	})
	assert.NotNil(t, err)

}
