package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/huichen/kunlun/internal/query"
)

func TestTrim(t *testing.T) {
	flag.Parse()

	q, _ := query.Parse("repo:repo_a (aa or bb)")
	TrimQuery(q)

	assert.Equal(t, "(aa OR bb)", q.String())
	assert.Equal(t, 3, q.NumNodes)
	assert.Equal(t, 1, q.MaxDepth)
}
