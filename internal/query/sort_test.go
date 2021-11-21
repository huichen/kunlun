package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	flag.Parse()

	q, _ := Parse("sym:a file:.java case:yes lang:java repo:infra.*framework (a or b or sym:c ) k g h ")
	Sort(q)
	assert.Equal(t, "(case:yes AND lang:java AND repo:infra.*framework AND file:.java AND sym:a AND g AND h AND k AND (sym:c OR a OR b))", q.String())

	q, _ = Parse("a b sym:ddd repo:fff file:ggg case:yes lang:java")
	assert.Equal(t, "(case:yes AND lang:java AND repo:fff AND file:ggg AND sym:ddd AND a AND b)", q.String())

	q, _ = Parse("a b (c or d) (e or a)")
	assert.Equal(t, "(a AND b AND (a OR e) AND (c OR d))", q.String())

}
