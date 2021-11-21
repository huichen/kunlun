package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegateOr(t *testing.T) {
	flag.Parse()

	q, _ := Parse(" a or -b  ")
	assert.Equal(t, "(a OR -b)", q.debugString())

	q, _ = Parse(" a and -(b and c)  ")
	assert.Equal(t, "(a AND (-b OR -c))", q.debugString())
}
