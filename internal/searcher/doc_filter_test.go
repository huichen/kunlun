package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/huichen/kunlun/internal/common_types"
	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/pkg/types"
)

func TestDocFilter(t *testing.T) {
	flag.Parse()

	options := types.NewIndexerOptions()
	options.SetNumIndexerShards(1)
	idxr := indexer.NewIndexer(options)

	idxr.IndexFile([]byte("aaaa"), common_types.IndexFileInfo{Path: "repo_a/file_a"})
	idxr.IndexFile([]byte("bbbb"), common_types.IndexFileInfo{Path: "repo_a/file_b"})
	idxr.IndexFile([]byte("cccc"), common_types.IndexFileInfo{Path: "repo_a/file_c"})
	idxr.IndexFile([]byte("bb\naa"), common_types.IndexFileInfo{Path: "repo_b/file_a"})
	idxr.IndexFile([]byte("dddd"), common_types.IndexFileInfo{Path: "repo_b/file_d"})
	idxr.IndexFile([]byte("dddd"), common_types.IndexFileInfo{Path: "repo_c/file_d"})
	idxr.IndexFile([]byte("bbbb"), common_types.IndexFileInfo{Path: "repo_c/file_b"})
	idxr.IndexRepo(common_types.IndexRepoInfo{RepoLocalPath: "repo_a", RepoRemoteURL: "me@git.com:repo_a"})
	idxr.IndexRepo(common_types.IndexRepoInfo{RepoLocalPath: "repo_b", RepoRemoteURL: "me@git.com:repo_b"})
	idxr.IndexRepo(common_types.IndexRepoInfo{RepoLocalPath: "repo_c", RepoRemoteURL: "me@git.com:repo_c"})
	idxr.Finish()

	q, _ := ParseQuery("repo:repo_a file:b")
	filter := NewDocFilter(q, idxr, nil, nil)
	assert.False(t, filter.ShouldRecallDocument(1))
	assert.True(t, filter.ShouldRecallDocument(2))
	assert.False(t, filter.ShouldRecallDocument(3))
	assert.False(t, filter.ShouldRecallDocument(4))

	q, _ = ParseQuery("repo:repo_a")
	filter = NewDocFilter(q, idxr, nil, nil)
	assert.True(t, filter.ShouldRecallDocument(1))
	assert.True(t, filter.ShouldRecallDocument(2))
	assert.True(t, filter.ShouldRecallDocument(3))
	assert.False(t, filter.ShouldRecallDocument(4))

	q, _ = ParseQuery("file:a")
	filter = NewDocFilter(q, idxr, nil, nil)
	assert.True(t, filter.ShouldRecallDocument(1))
	assert.False(t, filter.ShouldRecallDocument(2))
	assert.False(t, filter.ShouldRecallDocument(3))
	assert.True(t, filter.ShouldRecallDocument(4))
	assert.False(t, filter.ShouldRecallDocument(5))
	assert.False(t, filter.ShouldRecallDocument(6))
	assert.False(t, filter.ShouldRecallDocument(7))

	externalRepoFilter := func(repoID uint64) bool { return repoID == 2 }
	filter = NewDocFilter(q, idxr, externalRepoFilter, nil)
	assert.False(t, filter.ShouldRecallDocument(1))
	assert.False(t, filter.ShouldRecallDocument(2))
	assert.False(t, filter.ShouldRecallDocument(3))
	assert.True(t, filter.ShouldRecallDocument(4))
	assert.False(t, filter.ShouldRecallDocument(5))
	assert.False(t, filter.ShouldRecallDocument(6))
	assert.False(t, filter.ShouldRecallDocument(7))

	q, _ = ParseQuery("repo:repox a")
	filter = NewDocFilter(q, idxr, nil, nil)
	assert.False(t, filter.ShouldRecallDocument(1))
	assert.False(t, filter.ShouldRecallDocument(2))
	assert.False(t, filter.ShouldRecallDocument(3))
	assert.False(t, filter.ShouldRecallDocument(4))
	assert.False(t, filter.ShouldRecallDocument(5))
	assert.False(t, filter.ShouldRecallDocument(6))
	assert.False(t, filter.ShouldRecallDocument(7))

}
