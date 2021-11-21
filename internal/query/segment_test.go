package query

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegment(t *testing.T) {
	flag.Parse()

	segs, _ := segment("or")
	assert.Equal(t, []string{"or"}, segs)

	_, err := segment("(adfdf")
	logger.Error(err)

	_, err = segment("(adfdf")
	logger.Error(err)

	segs, _ = segment("a -b -(c -d) -\"efg\"")
	assert.Equal(t, []string{"a", "-b", "-(c -d)", "-\"efg\""}, segs)

	segs, _ = segment("-a")
	assert.Equal(t, []string{"-a"}, segs)

	segs, _ = segment("-a -\"b\"")
	assert.Equal(t, []string{"-a", "-\"b\""}, segs)

	segs, _ = segment("a bb or(c and (d))")
	assert.Equal(t, []string{"a", "bb", "or", "(c and (d))"}, segs)

	segs, _ = segment("a bb or\"c and (d)\"")
	assert.Equal(t, []string{"a", "bb", "or", "\"c and (d)\""}, segs)

	segs, _ = segment("- - -a")
	assert.Equal(t, []string{"-", "-", "-a"}, segs)

	segs, _ = segment("---a")
	assert.Equal(t, []string{"---a"}, segs)

	segs, _ = segment("a or sym:b")
	assert.Equal(t, []string{"a", "or", "sym:b"}, segs)

	segs, _ = segment("a \"a*b \" (b or (c or d)) \"dd \"")
	assert.Equal(t, []string{"a", "\"a*b \"", "(b or (c or d))", "\"dd \""}, segs)

	segs, _ = segment("a or file:\"b b\"")
	assert.Equal(t, []string{"a", "or", "file:\"b b\""}, segs)

	segs, _ = segment("repo:docs (@todo or \" todo\" or \"todo \" or \"todo(\" ) -todouble")
	assert.Equal(t, []string{"repo:docs", "(@todo or \" todo\" or \"todo \" or \"todo(\" )", "-todouble"}, segs)

	segs, _ = segment("repo:docs (@todo or \" todo\" or \"\\\"todo \" or \"todo(\" ) -todouble")
	assert.Equal(t, []string{"repo:docs", "(@todo or \" todo\" or \"\\\"todo \" or \"todo(\" )", "-todouble"}, segs)

}
