package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedup(t *testing.T) {
	flag.Parse()

	q, _ := Parse("a AND b AND a")
	Dedup(q)
	assert.Equal(t, "(a AND b)", q.String())

	q, _ = Parse("a AND (c OR b) AND d AND (c OR b)")
	Dedup(q)
	assert.Equal(t, "(a AND d AND (b OR c))", q.String())
}
