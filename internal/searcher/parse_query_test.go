package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	flag.Parse()

	q, _ := ParseQuery("a b sym:ddd -repo:fff -repo:ddd file:ggg -file:kkk case:yes lang:java (e or f and g)")
	assert.Equal(t, "(case:yes AND lang:java AND -repo:ddd AND -repo:fff AND file:ggg AND -file:kkk AND sym:ddd AND a AND b AND (e OR (f AND g)))", q.OriginalQuery.String())

	assert.Equal(t, "(sym:ddd AND a AND b AND (e OR (f AND g)))", q.TrimmedQuery.String())
	assert.Equal(t, "lang:java", q.LanguageQuery.String())
	assert.Equal(t, "(-repo:ddd AND -repo:fff)", q.RepoQuery.String())
	assert.Equal(t, "(file:ggg AND -file:kkk)", q.FileQuery.String())
}

func TestDeepModifier(t *testing.T) {
	flag.Parse()

	q, _ := ParseQuery("a b file:a (repo:a or repo:b)")
	assert.Equal(t, "(repo:a OR repo:b)", q.RepoQuery.String())

	var err error
	q, err = ParseQuery("a b repo:c (repo:a or repo:b)")
	logger.Error(err)
	assert.Nil(t, q)

	q, err = ParseQuery("a or (repo:a or repo:b)")
	logger.Error(err)
	assert.Equal(t, "(repo:a OR repo:b)", q.RepoQuery.String())

}
