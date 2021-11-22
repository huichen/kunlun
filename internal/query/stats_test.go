package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	flag.Parse()

	q, _ := Parse("a (b or (c or d)) e")
	assert.Equal(t, "(a@[0,1,0] AND e@[1,1,0] AND (b@[2,1,0] OR c@[3,1,0] OR d@[4,1,0])@[5,4,1])@[6,7,2]", q.fullDebugString(false, true, false))

}
