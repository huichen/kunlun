package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModifier(t *testing.T) {
	flag.Parse()

	q, _ := Parse("sym:a file:.java case:yes lang:java repo:infra.*framework")
	assert.Equal(t, "(case:yes AND lang:java AND repo:infra.*framework AND file:.java AND sym:a)", q.debugString())

	q, _ = Parse("-sym:a -file:.java case:yes -lang:java -repo:infra.*framework")
	assert.Equal(t, "(case:yes AND -lang:java AND -repo:infra.*framework AND -file:.java AND -sym:a)", q.debugString())

	q, _ = Parse("sym:a (b or repo:c)")
	assert.Equal(t, "(sym:a AND (repo:c OR b))", q.debugString())

	q, _ = Parse("a b (repo:a or repo:b)")
	assert.Equal(t, "(a AND b AND (repo:a OR repo:b))", q.debugString())
}
