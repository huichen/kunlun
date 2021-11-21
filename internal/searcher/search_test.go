package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/huichen/kunlun/internal/indexer"

	"github.com/huichen/kunlun/pkg/types"
)

func TestSearch1(t *testing.T) {
	flag.Parse()

	idxrOptions := types.NewIndexerOptions()
	idxrOptions.SetNumIndexerShards(1)
	idxr := indexer.NewIndexer(idxrOptions)

	idxr.IndexFile([]byte("aaaa"), types.IndexFileInfo{Path: "repo_a/file_a"})
	idxr.IndexFile([]byte("bbbb"), types.IndexFileInfo{Path: "repo_a/file_b"})
	idxr.IndexFile([]byte("cccc"), types.IndexFileInfo{Path: "repo_a/file_c"})
	idxr.IndexFile([]byte("bb\naa"), types.IndexFileInfo{Path: "repo_b/file_a"})
	idxr.IndexFile([]byte("dddd"), types.IndexFileInfo{Path: "repo_b/file_d"})
	idxr.IndexFile([]byte("dddd"), types.IndexFileInfo{Path: "repo_c/file_d"})
	idxr.IndexFile([]byte("bbbb"), types.IndexFileInfo{Path: "repo_c/file_b"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_a", "me@git.com:repo_a"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_b", "me@git.com:repo_b"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_c", "me@git.com:repo_c"})
	idxr.Finish()

	options := types.NewSearcherOptions()
	schr := NewSearcher(options)

	request := types.SearchRequest{
		Query: "aaaa",
	}
	resp, _ := schr.Search(idxr, request)

	testDocs(t,
		[]types.SearchedDocument{
			{
				DocumentID: 1,
				Filename:   "file_a",
				Lines: []types.Line{
					{
						LineNumber: 0,
						Highlights: []types.Section{{0, 4}},
					},
				},
				NumLinesInDocument: 1,
			},
		}, resp.Repos[0].Documents)

}

func testDocs(t *testing.T, expected, actual []types.SearchedDocument) {
	for i := range actual {
		for j := range actual[i].Lines {
			actual[i].Lines[j].Content = nil
		}
	}
	assert.Equal(t, expected, actual)
}

func TestSearch2(t *testing.T) {
	flag.Parse()

	idxrOptions := types.NewIndexerOptions()
	idxrOptions.SetNumIndexerShards(1)
	idxr := indexer.NewIndexer(idxrOptions)

	idxr.IndexFile([]byte("aaaa"), types.IndexFileInfo{Path: "repo_a/file_a"})
	idxr.IndexFile([]byte("bbbb"), types.IndexFileInfo{Path: "repo_a/file_b"})
	idxr.IndexFile([]byte("cccc"), types.IndexFileInfo{Path: "repo_a/file_c"})
	idxr.IndexFile([]byte("bb\naa"), types.IndexFileInfo{Path: "repo_b/file_a"})
	idxr.IndexFile([]byte("dddd"), types.IndexFileInfo{Path: "repo_b/file_d"})
	idxr.IndexFile([]byte("dddd"), types.IndexFileInfo{Path: "repo_c/file_d"})
	idxr.IndexFile([]byte("bbbb"), types.IndexFileInfo{Path: "repo_c/file_b"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_a", "me@git.com:repo_a"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_b", "me@git.com:repo_b"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_c", "me@git.com:repo_c"})
	idxr.Finish()

	options := types.NewSearcherOptions()
	schr := NewSearcher(options)

	request := types.SearchRequest{
		Query: "aa",
	}
	resp, _ := schr.Search(idxr, request)
	testDocs(t, []types.SearchedDocument{
		{
			DocumentID: 1,
			Filename:   "file_a",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 4}},
				},
			},
			NumLinesInDocument: 1,
		},
	}, resp.Repos[0].Documents)

	testDocs(t, []types.SearchedDocument{
		{
			DocumentID: 4,
			Filename:   "file_a",
			Lines: []types.Line{
				{
					LineNumber: 1,
					Highlights: []types.Section{{0, 2}},
				},
			},
			NumLinesInDocument: 1,
		},
	}, resp.Repos[1].Documents)

	request = types.SearchRequest{
		Query: "aa and bb",
	}
	resp, _ = schr.Search(idxr, request)
	testDocs(t, []types.SearchedDocument{
		{
			DocumentID: 4,
			Filename:   "file_a",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 2}},
				},
				{
					LineNumber: 1,
					Highlights: []types.Section{{0, 2}},
				},
			},
			NumLinesInDocument: 2,
		},
	}, resp.Repos[0].Documents)

	request = types.SearchRequest{
		Query: "aa or bb",
	}
	resp, _ = schr.Search(idxr, request)

	testDocs(t, []types.SearchedDocument{
		{
			DocumentID: 1,
			Filename:   "file_a",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 4}},
				},
			},
			NumLinesInDocument: 1,
		},
		{
			DocumentID: 2,
			Filename:   "file_b",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 4}},
				},
			},
			NumLinesInDocument: 1,
		},
	}, resp.Repos[0].Documents)

	testDocs(t, []types.SearchedDocument{
		{
			DocumentID: 4,
			Filename:   "file_a",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 2}},
				},
				{
					LineNumber: 1,
					Highlights: []types.Section{{0, 2}},
				},
			},
			NumLinesInDocument: 2,
		},
	}, resp.Repos[1].Documents)

	request = types.SearchRequest{
		Query: "repo:repo_a (aa or bb)",
	}
	resp, _ = schr.Search(idxr, request)
	testDocs(t, []types.SearchedDocument{
		{
			DocumentID: 1,
			Filename:   "file_a",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 4}},
				},
			},
			NumLinesInDocument: 1,
		},
		{
			DocumentID: 2,
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 4}},
				},
			},
			Filename:           "file_b",
			NumLinesInDocument: 1,
		},
	}, resp.Repos[0].Documents)
}

func TestSearchBuggy(t *testing.T) {
	flag.Parse()

	idxrOptions := types.NewIndexerOptions()
	idxrOptions.SetNumIndexerShards(1)
	idxr := indexer.NewIndexer(idxrOptions)

	idxr.IndexFile([]byte("aaaa"), types.IndexFileInfo{Path: "repo_a/file_a"})
	idxr.IndexFile([]byte("bbbb"), types.IndexFileInfo{Path: "repo_a/file_b"})
	idxr.IndexFile([]byte("cccc"), types.IndexFileInfo{Path: "repo_a/file_c"})
	idxr.IndexFile([]byte("bb\naa"), types.IndexFileInfo{Path: "repo_b/file_a"})
	idxr.IndexFile([]byte("dddd"), types.IndexFileInfo{Path: "repo_b/file_d"})
	idxr.IndexFile([]byte("dddd"), types.IndexFileInfo{Path: "repo_c/file_d"})
	idxr.IndexFile([]byte("bbbb"), types.IndexFileInfo{Path: "repo_c/file_b"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_a", "me@git.com:repo_a"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_b", "me@git.com:repo_b"})
	idxr.IndexRepo(types.IndexRepoInfo{"repo_c", "me@git.com:repo_c"})
	idxr.Finish()

	options := types.NewSearcherOptions()
	schr := NewSearcher(options)

	request := types.SearchRequest{
		Query: "repo:repo_a (a.*a or bbb)",
	}
	resp, err := schr.Search(idxr, request)
	logger.Info(err)
	testDocs(t, []types.SearchedDocument{
		{
			DocumentID: 1,
			Filename:   "file_a",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 4}},
				},
			},
			NumLinesInDocument: 1,
		},
		{
			DocumentID: 2,
			Filename:   "file_b",
			Lines: []types.Line{
				{
					LineNumber: 0,
					Highlights: []types.Section{{0, 4}},
				},
			},
			NumLinesInDocument: 1,
		},
	}, resp.Repos[0].Documents)
}
