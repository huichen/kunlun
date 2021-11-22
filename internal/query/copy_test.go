package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	flag.Parse()

	q, _ := Parse("a b sym:ddd repo:fff file:ggg case:yes lang:java (e or f and g \"h\")")
	assert.Equal(t, "(case:yes@[0,1,0] AND lang:java@[1,1,0] AND repo:fff@[2,1,0] AND file:ggg@[3,1,0] AND sym:ddd@[4,1,0] AND a@[5,1,0] AND b@[6,1,0] AND (e@[7,1,0] OR (f@[8,1,0] AND g@[9,1,0] AND h@[10,1,0])@[11,4,1])@[12,6,2])@[13,14,3]", q.fullDebugString(false, true, false))

	dq := Copy(q)
	assert.Equal(t, "(case:yes@[0,1,0] AND lang:java@[1,1,0] AND repo:fff@[2,1,0] AND file:ggg@[3,1,0] AND sym:ddd@[4,1,0] AND a@[5,1,0] AND b@[6,1,0] AND (e@[7,1,0] OR (f@[8,1,0] AND g@[9,1,0] AND h@[10,1,0])@[11,4,1])@[12,6,2])@[13,14,3]", dq.fullDebugString(false, true, false))
}
