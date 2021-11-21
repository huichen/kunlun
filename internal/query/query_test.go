package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	flag.Parse()

	_, err := Parse(" this is (a b (d d) ) query ")
	assert.Nil(t, err)

	_, err = Parse(" this is  \t\t  \"ad\\\"fs.*d dddfsdf\" query (asdf)")
	assert.Nil(t, err)

	_, err = Parse(" this is (a b  (d d) ) query ")
	assert.Nil(t, err)

	_, err = Parse(" this is (a b (d d) query ")
	logger.Error(err)
	assert.NotNil(t, err)

	_, err = Parse(" this is a b (d d)) query ")
	logger.Error(err)
	assert.NotNil(t, err)

	q, _ := Parse("a or b")
	assert.Equal(t, "(a OR b)", q.debugString())

	q, _ = Parse("a and b")
	assert.Equal(t, "(a AND b)", q.debugString())

	q, _ = Parse("a and b or c")
	assert.Equal(t, "(c OR (a AND b))", q.debugString())

	q, _ = Parse("(a or b) and c")
	assert.Equal(t, "(c AND (a OR b))", q.debugString())

	q, _ = Parse("a bb or     ccc ")
	assert.Equal(t, "(ccc OR (a AND bb))", q.debugString())

	q, _ = Parse("a bb or   \"  ccc\" ")
	assert.Equal(t, "(\"  ccc\" OR (a AND bb))", q.debugString())

	q, _ = Parse("a bb or(c and (d))")
	assert.Equal(t, "((a AND bb) OR (c AND d))", q.debugString())

	q, _ = Parse("a b or c and d or (e or f and g \"h\")")
	assert.Equal(t, "(e OR (a AND b) OR (c AND d) OR (f AND g AND h))", q.debugString())

	q, _ = Parse("*b \".*abc\"")
	assert.Equal(t, "(*b AND <.*abc>(abc))", q.debugString())

	q, _ = Parse("a a.*aaa.*[a-z]{3}\\d{3}[^b] a{3}")
	assert.Equal(t, "(a AND <a.*aaa.*[a-z]{3}\\d{3}[^b]>(a,aaa) AND <a{3}>(aaa))", q.debugString())

	q, _ = Parse("a b -c d")
	assert.Equal(t, "(a AND b AND d AND -c)", q.debugString())

	q, _ = Parse("- - a")
	assert.Equal(t, "", q.debugString())

	q, _ = Parse("-(-(-a))")
	assert.Equal(t, "-a", q.debugString())

	q, _ = Parse("-(c or -(e f))")
	assert.Equal(t, "(e AND f AND -c)", q.debugString())

	q, _ = Parse("-c (e f)")
	assert.Equal(t, "(e AND f AND -c)", q.debugString())

	q, _ = Parse("a b -(c or -(e f))")
	assert.Equal(t, "(a AND b AND e AND f AND -c)", q.debugString())

	q, _ = Parse("a b (c d)")
	assert.Equal(t, "(a AND b AND c AND d)", q.debugString())

	q, _ = Parse("---a")
	assert.Equal(t, "-a", q.debugString())

	q, _ = Parse("or or")
	assert.Nil(t, q)

	q, _ = Parse("or")
	assert.Nil(t, q)

	q, _ = Parse("a b \"c\" or d")
	assert.Equal(t, "(d OR (a AND b AND c))", q.debugString())

	q, _ = Parse("a *b*c *d e* *f*g*h")
	assert.Equal(t, "(*b*c AND *d AND *f*g*h AND a AND <e*>())", q.debugString())

	q, _ = Parse("a a* c+ ddd{3,6} e{3} f{0,7}")
	assert.Equal(t, "(a AND <a*>() AND <c+>(c) AND <ddd{3,6}>(dd,ddd) AND <e{3}>(eee) AND <f{0,7}>())", q.debugString())

	q, _ = Parse("a \"a*b \" (b or (c or d)) \"dd \"")
	assert.Equal(t, "(\"dd \" AND a AND \"a*b \" AND (b OR c OR d))", q.String())

	q, _ = Parse("a \"dd()\\\" \"")
	assert.Equal(t, "(a AND \"dd()\\\" \")", q.String())

	q, _ = Parse("a or sym:b")
	Sort(q)
	assert.Equal(t, "(sym:b OR a)", q.debugString())

}

func TestOnce(t *testing.T) {
	flag.Parse()

	q, _ := Parse("a b sym:ddd repo:fff file:ggg case:yes lang:java")
	assert.Equal(t, "(case:yes AND lang:java AND repo:fff AND file:ggg AND sym:ddd AND a AND b)", q.debugString())
}
