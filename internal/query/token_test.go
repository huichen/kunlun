package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokens(t *testing.T) {
	flag.Parse()

	q, _ := Parse("a b -(c OR d OR a OR .*ab.*d OR -(e f))")
	tokens := q.GetTokens()
	assert.Equal(t, []string{"a", "b", "e", "f", "c", "d", "ab"}, tokens)

}
