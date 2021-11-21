package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	flag.Parse()

	q, _ := ParseQuery("repo:a repo:b")
	assert.True(t, matchRegexpQueries("abc", q.RepoQuery, &q.RepoRe, nil))
	assert.False(t, matchRegexpQueries("ac", q.RepoQuery, &q.RepoRe, nil))
	assert.False(t, matchRegexpQueries("a", q.RepoQuery, &q.RepoRe, nil))

	q, _ = ParseQuery("a (repo:repo_a or repo:repo_b)")
	assert.False(t, matchRegexpQueries("c", q.RepoQuery, &q.RepoRe, nil))
	assert.True(t, matchRegexpQueries("repo_ac", q.RepoQuery, &q.RepoRe, nil))
	assert.True(t, matchRegexpQueries("repo_a", q.RepoQuery, &q.RepoRe, nil))

	q, _ = ParseQuery("a -repo:a -repo:b")
	assert.True(t, matchRegexpQueries("c", q.RepoQuery, &q.RepoRe, nil))
	assert.False(t, matchRegexpQueries("ac", q.RepoQuery, &q.RepoRe, nil))
	assert.False(t, matchRegexpQueries("a", q.RepoQuery, &q.RepoRe, nil))

	q, _ = ParseQuery("a -repo:a repo:b")
	assert.False(t, matchRegexpQueries("c", q.RepoQuery, &q.RepoRe, nil))
	assert.False(t, matchRegexpQueries("ab", q.RepoQuery, &q.RepoRe, nil))
	assert.False(t, matchRegexpQueries("a", q.RepoQuery, &q.RepoRe, nil))
	assert.True(t, matchRegexpQueries("bcd", q.RepoQuery, &q.RepoRe, nil))

	q, _ = ParseQuery("a (file:a.*c or file:d+)")
	assert.False(t, matchRegexpQueries("c", q.FileQuery, &q.FileRe, nil))
	assert.True(t, matchRegexpQueries("ac", q.FileQuery, &q.FileRe, nil))
	assert.True(t, matchRegexpQueries("d", q.FileQuery, &q.FileRe, nil))

	q, _ = ParseQuery("a (lang:java or lang:cpp)")
	assert.True(t, matchRegexpQueries("java", q.LanguageQuery, &q.LangRe, nil))
	assert.True(t, matchRegexpQueries("cpp", q.LanguageQuery, &q.LangRe, nil))
	assert.False(t, matchRegexpQueries("python", q.LanguageQuery, &q.LangRe, nil))

	q, _ = ParseQuery("a (lang:java and lang:cpp)")
	assert.False(t, matchRegexpQueries("java", q.LanguageQuery, &q.LangRe, nil))
	assert.False(t, matchRegexpQueries("cpp", q.LanguageQuery, &q.LangRe, nil))
	assert.False(t, matchRegexpQueries("python", q.LanguageQuery, &q.LangRe, nil))

}
