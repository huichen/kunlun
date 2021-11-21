package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	flag.Parse()

	_, err := ParseQuery(" a or -b  ")
	logger.Error(err)
	assert.NotNil(t, err)

	_, err = ParseQuery(" a and -(b and c)  ")
	logger.Error(err)
	assert.NotNil(t, err)

	_, err = ParseQuery(" repo:a OR (repo:b repo:c)  ")
	logger.Error(err)
	assert.NotNil(t, err)

	_, err = ParseQuery(" a OR (repo:b repo:c)  ")
	logger.Error(err)
	assert.Nil(t, err)

	_, err = ParseQuery(" a OR file:b")
	logger.Error(err)
	assert.Nil(t, err)

	_, err = ParseQuery(" a OR (c file:b)")
	logger.Error(err)
	assert.Nil(t, err)

}
