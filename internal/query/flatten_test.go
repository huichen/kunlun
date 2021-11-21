package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	flag.Parse()

	q, _ := Parse(" a b (c d) ")
	assert.Equal(t, "(a AND b AND c AND d)", q.debugString())

	q, _ = Parse(" a b (c (d) (e f)) ")
	assert.Equal(t, "(a AND b AND c AND d AND e AND f)", q.debugString())

	q, _ = Parse(" a b or (c or (d) or (e f)) ")
	assert.Equal(t, "(c OR d OR (a AND b) OR (e AND f))", q.debugString())

}
