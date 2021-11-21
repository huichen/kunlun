package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompactString(t *testing.T) {
	flag.Parse()

	q, _ := Parse(" a b c(d or (e f)) ")
	assert.Equal(t, "(a AND b AND c AND (d OR (e AND f)))", q.String())

}
